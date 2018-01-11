package utils

import (
	"github.com/go-errors/errors"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"bytes"
	"strings"
)

var (
	ERR_REQUIRE_LOGIN = errors.New("require login")
	ERR_REQUIRE_ADMIN = errors.New("require admin")
	ERR_NO_ACCOUNT    = errors.New("no such account")
	ERR_NO_PERMISSION = errors.New("no permission")
	ERR_CANNOT_MODIFY = errors.New("cannot modify")
	ErrNoItem         = errors.New("no such item")
	ErrExisted        = errors.New("existed")
	ErrNotMatch       = errors.New("not match")
	ErrFormat         = errors.New("format error")

	UtilsErrList = []*errors.Error{ERR_REQUIRE_LOGIN, ERR_REQUIRE_ADMIN, ERR_NO_ACCOUNT,
		ERR_NO_PERMISSION, ERR_CANNOT_MODIFY, ErrNoItem, ErrExisted, ErrNotMatch}
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

func NewErrPermission() *errors.Error {
	return errors.Wrap("没有权限", 1)
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
	List *sll.List
}

func NewErrorList() *ErrorList {
	l := &ErrorList{}
	l.List = sll.New()
	return l
}

//只添加非nil的error
func (l *ErrorList) AppendE(errs ...error) {
	for _, v := range errs {
		if v != nil {
			//l.List.Add(errors.Wrap(v, 1))
			l.List.Add(v)
		}
	}
}

func (l *ErrorList) FirstError() error {
	_, err := l.List.Find(func(index int, value interface{}) bool {
		return value != nil
	})
	if err == nil {
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
		if frame.Package == "main" || strings.HasSuffix(frame.File, "/mcmd/main.go") {
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
