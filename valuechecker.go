package utils

import (
	"github.com/thoas/go-funk"
	"github.com/go-errors/errors"
	"regexp"
)

var (
	phoneReg = regexp.MustCompile(`^1\d{10}$`)
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

func (v *ValueChecker) NotError(val error, errMsg string) bool {
	if val != nil {
		if errMsg == "" {
			v.errs.AppendE(errors.Wrap(val, 1))
		} else {
			v.errs.AppendE(errors.Wrap(errMsg, 1))
		}
		return false
	}
	return true
}

func (v *ValueChecker) PhoneCn(val string, errMsg string) bool {
	if !phoneReg.MatchString(val) {
		v.addErr(DefaultVal(errMsg, "电话错误"))
		return false
	}
	return true
}

func (v *ValueChecker) addErr(errMsg interface{}) {
	v.errs.AppendE(errors.Wrap(errMsg, 2))
}

func (v *ValueChecker) FirstError() error {
	return v.errs.FirstError()
}
