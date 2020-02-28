package ucore

import (
	"bytes"
	"fmt"
	ge "github.com/WingGao/errors"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/go-errors/errors"
	"regexp"
	"strconv"
	"strings"
)

type stackTracer interface {
	StackTrace() ge.StackTrace
}

var (
	ErrRequireAdmin = errors.New("require admin")
	ErrNoPermission = errors.New("no permission")
	ErrCannotModify = errors.New("cannot modify")
	ErrNoItem       = errors.New("no such item")
	ErrExisted      = errors.New("existed")
	ErrNotMatch     = errors.New("not match")
	ErrFormat       = errors.New("format error")

	UtilsErrList = []*errors.Error{ ErrRequireAdmin,
		ErrNoPermission, ErrCannotModify, ErrNoItem, ErrExisted, ErrNotMatch}
)

func NewErrNotFound() error {
	return ge.WrapSkip("不存在", 1)
}

func NewErrExisted() error {
	return ge.WrapSkip("已存在", 1)
}

func NewErrParams() error {
	return ge.WrapSkip("参数错误", 1)
}

func NewErrSystem() error {
	return ge.WrapSkip("系统错误", 1)
}

func NewErrNeedLogin() error {
	return ge.WrapSkip("require login", 1)
}

func NewErrPermission() error {
	return ge.WrapSkip("没有权限", 1)
	//return ge.New("没有权限")
}
func NewErrPassword() error {
	return ge.WrapSkip("密码错误", 1)
	//return ge.New("没有权限")
}

func NewErrCode() error {
	return ge.WrapSkip("验证码错误", 1)
}

func NewErrNoAccount() error {
	return ge.WrapSkip("账户不存在", 1)
}

func Nothing(...interface{}) {

}

// 获取第一个错误
func FirstError(es ...error) error {
	for _, e := range es {
		if e != nil {
			return e
		}
	}
	return nil
}

type ErrorList struct {
	list *sll.List
}

func NewErrorList() *ErrorList {
	l := &ErrorList{}
	l.list = sll.New()
	return l
}

//只添加非nil的error
func (l *ErrorList) AppendE(errs ...error) {
	for _, v := range errs {
		if v != nil {
			//l.List.Add(ge.WrapSkip(v, 1))
			l.list.Add(v)
		}
	}
}

func (l *ErrorList) AppendEWrap(err error, skip int) {
	if err != nil {
		l.list.Add(ge.WrapSkip(err, skip))
	}
}

func (l *ErrorList) FirstError() error {
	//_, err := l.list.Find(func(index int, value interface{}) bool {
	//	return value != nil
	//})
	err, ok := l.list.Get(0)
	if !ok {
		return nil
	}
	return err.(error)
}

func (l *ErrorList) Panic() {
	if err := l.FirstError(); err != nil {
		panic(err)
	}
}

// 没有错误的时候运行
func (l *ErrorList) Run(fo func() error) {
	if l.FirstError() == nil {
		l.AppendE(fo())
	}
}

func PrintError(err error) {
	//errN := ge.WrapSkip(err, 0)
	//fmt.Println(errN.ErrorStack())
}

func PanicIfErr(err error) {
	if err != nil {
		panic(errors.New(err))
	}
}

type WError struct {
	Err        *errors.Error
	Frames     []errors.StackFrame
	StackLines string
}

func NewWError(e interface{}) *WError {
	out := &WError{}
	if e2, ok := e.(*errors.Error); ok {
		out.Err = e2
		out.Frames = e2.StackFrames()
	} else if e3, ok3 := e.(stackTracer); ok3 {
		stacks := e3.StackTrace()
		out.StackLines = fmt.Sprintf("%+v", stacks)
		out.Frames = parseBlocks(out.StackLines)
	} else {
		out.Err = errors.Wrap(e, 1)
		out.Frames = out.Err.StackFrames()
	}
	return out
}

//我们只需要知道最短路径
func (e *WError) Fmt() {
	simpleFrames := e.Frames
	//fmt.Printf("%#v\n", simpleFrames)
	for i, frame := range e.Frames {
		if frame.Package == "wingao.net/webproj/core" && frame.Name == "handlers.IrisWrapper.func1" {
			// 跳过
			// /reflect.Value.call
			//	C:/Go/src/reflect/value.go:447
			///reflect.Value.Call
			//	C:/Go/src/reflect/value.go:308
			//wingao.net/webproj/core/handlers.IrisWrapper.func1
			//	D:/Projs/go-web/wingao.net/webproj/core/handlers/h.go:172
			simpleFrames = e.Frames[:i-2]
			break
		}
		if frame.Package == "wingao.net/webproj/mcmd/serv" || frame.Package == "github.com/kataras/iris/v12/middleware/logger" {
			simpleFrames = e.Frames[:i+1]
			break
		}
	}
	//fmt.Printf("%#v\n", simpleFrames)
	e.Frames = simpleFrames
}

func (e *WError) Stack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range e.Frames {
		if frame.Name != "" {
			buf.WriteString(frame.Package)
			buf.WriteString("/")
			buf.WriteString(frame.Name)
			buf.WriteString("\n\t")
			buf.WriteString(frame.File)
			buf.WriteString(":")
			buf.WriteString(strconv.Itoa(frame.LineNumber))
			buf.WriteString("\n")
		} else {
			buf.WriteString(frame.String())
		}
	}

	return buf.Bytes()
}

func (e *WError) ErrorStack() string {
	bs := ""
	if e.Err != nil {
		bs = e.Err.TypeName() + " " + e.Err.Error() + "\n"
	}
	return bs + string(e.Stack())
}

var stackLineR = regexp.MustCompile(`\t`) // 文件行

func parseBlocks(input string) ([]errors.StackFrame) {
	var blocks []errors.StackFrame

	frame := errors.StackFrame{}
	for _, l := range strings.Split(input, "\n") {
		isStackLine := stackLineR.MatchString(l)
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if isStackLine {
			fs := strings.Split(l, ":")
			if len(fs) > 2 {
				//windows
				frame.File = fs[0] + ":" + fs[1]
				frame.LineNumber, _ = strconv.Atoi(fs[2])
			} else {
				frame.File = fs[0]
				frame.LineNumber, _ = strconv.Atoi(fs[1])
			}
			blocks = append(blocks, frame)
			frame = errors.StackFrame{}
		} else {
			fns := strings.Split(l, "/")
			frame.Name = fns[len(fns)-1]
			frame.Package = strings.Join(fns[0:len(fns)-1], "/")
		}
	}

	return blocks
}

type CommError struct {
	ErrMsg    string `json:"err_msg"`
	ErrFormId string `json:"err_form_id,omitempty"`
}

func (m CommError) Error() string {
	return m.ErrMsg
}
