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
// 应该只使用true情况，不要使用`!xxx`这种
type ValueChecker struct {
	errs        *ErrorList
	SkipOnError bool //遇到错误是否跳过
}

func NewValueChecker() (v *ValueChecker) {
	v = &ValueChecker{SkipOnError: true}
	v.errs = NewErrorList()
	return
}

// 不检查,直接返回错误
func (v *ValueChecker) shouldSkip() bool {
	return v.FirstError() != nil && v.SkipOnError
}

func (v *ValueChecker) CheckBy(f func() bool, errMsg string) bool {
	if v.shouldSkip() {
		return false
	}

	if f() {
		return true
	} else {
		return v.addErr(errMsg)
	}
}
func (v *ValueChecker) NotEmpty(value interface{}, errMsg string) bool {
	if v.shouldSkip() {
		return false
	}

	if funk.IsEmpty(value) {
		return v.addErr(errMsg)
	}
	return true
}

func (v *ValueChecker) NotError(val error, errMsg string) bool {
	if v.shouldSkip() {
		return false
	}

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

func (v *ValueChecker) Contains(val, items interface{}, errMsg string) bool {
	if v.shouldSkip() {
		return false
	}

	if funk.Contains(items, val) {
		return true
	} else {
		return v.addErr(errMsg)
	}
}

func (v *ValueChecker) PhoneCn(val string, errMsg string) bool {
	if v.shouldSkip() {
		return false
	}

	if !phoneReg.MatchString(val) {
		v.addErr(DefaultVal(errMsg, "电话错误"))
		return false
	}
	return true
}

func (v *ValueChecker) addErr(errMsg interface{}) bool {
	v.errs.AppendE(errors.Wrap(errMsg, 2))
	return false
}

func (v *ValueChecker) FirstError() error {
	return v.errs.FirstError()
}
