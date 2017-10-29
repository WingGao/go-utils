package utils

import (
	"github.com/go-errors/errors"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
)

var (
	ERR_REQUIRE_LOGIN = errors.New("require login")
	ERR_REQUIRE_ADMIN = errors.New("require admin")
	ERR_NO_ACCOUNT    = errors.New("no such account")
	ERR_PARAMS        = errors.New("params error")
	ERR_NO_PERMISSION = errors.New("no permission")
	ERR_CANNOT_MODIFY = errors.New("cannot modify")
)

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
