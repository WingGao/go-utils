package utils

import (
	"github.com/thoas/go-funk"
	"github.com/go-errors/errors"
)

//简单值检查
type ValueChecker struct {
	errs *ErrorList
}

func NewValueChecker() (v *ValueChecker) {
	v = &ValueChecker{}
	v.errs = NewErrorList()
	return
}

func (v *ValueChecker) NotEmpty(value interface{}, errMsg string) bool {
	if funk.IsEmpty(value) {
		v.errs.AppendE(errors.Wrap(errMsg, 1))
		return false
	}
	return true
}

func (v *ValueChecker) FirstError() error {
	return v.errs.FirstError()
}
