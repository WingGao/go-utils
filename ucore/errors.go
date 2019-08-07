package ucore

import (
	"bytes"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/go-errors/errors"
)

var (
	ErrRequireLogin = errors.New("require login") // 这个错误不能改
	ErrRequireAdmin = errors.New("require admin")
	ErrNoPermission = errors.New("no permission")
	ErrCannotModify = errors.New("cannot modify")
	ErrNoItem       = errors.New("no such item")
	ErrExisted      = errors.New("existed")
	ErrNotMatch     = errors.New("not match")
	ErrFormat       = errors.New("format error")

	UtilsErrList = []*errors.Error{ErrRequireLogin, ErrRequireAdmin,
		ErrNoPermission, ErrCannotModify, ErrNoItem, ErrExisted, ErrNotMatch}
)

func NewErrNotFound() *errors.Error {
	return errors.Wrap("不存在", 1)
}

func NewErrExisted() *errors.Error {
	return errors.Wrap("已存在", 1)
}

func NewErrParams() *errors.Error {
	return errors.Wrap("参数错误", 1)
}

func NewErrNeedLogin() *errors.Error {
	return errors.Wrap(ErrRequireLogin, 1)
}

func NewErrPermission() *errors.Error {
	return errors.Wrap("没有权限", 1)
}

func NewErrCode() *errors.Error {
	return errors.Wrap("验证码错误", 1)
}

func NewErrNoAccount() *errors.Error {
	return errors.Wrap("账户不存在", 1)
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
			//l.List.Add(errors.Wrap(v, 1))
			l.list.Add(v)
		}
	}
}

func (l *ErrorList) AppendEWrap(err error, skip int) {
	if err != nil {
		l.list.Add(errors.Wrap(err, skip))
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
	//errN := errors.Wrap(err, 0)
	//fmt.Println(errN.ErrorStack())
}

func PanicIfErr(err error) {
	if err != nil {
		panic(errors.Wrap(err, 1))
	}
}

type WError struct {
	Err    *errors.Error
	Frames []errors.StackFrame
}

func NewWError(e interface{}) *WError {
	out := &WError{}
	if e2, ok := e.(*errors.Error); ok {
		out.Err = e2
		out.Frames = e2.StackFrames()
	} else {
		out.Err = errors.Wrap(e, 1)
		out.Frames = out.Err.StackFrames()
	}
	return out
}

//我们只需要知道最短路径
func (e *WError) Fmt() {
	for i, frame := range e.Frames {
		if frame.Package == "wingao.net/webproj/mcmd/serv	" || frame.Package == "github.com/kataras/iris/middleware/logger" {
			e.Frames = e.Frames[:i+1]
			break
		}
	}
}

func (e *WError) Stack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range e.Frames {
		buf.WriteString(frame.String())
	}

	return buf.Bytes()
}

func (e *WError) ErrorStack() string {
	return e.Err.TypeName() + " " + e.Err.Error() + "\n" + string(e.Stack())
}

type CommError struct {
	ErrMsg    string `json:"err_msg"`
	ErrFormId string `json:"err_form_id,omitempty"`
}

func (m CommError) Error() string {
	return m.ErrMsg
}
