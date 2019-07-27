package ucore

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/thoas/go-funk"
	"reflect"
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

func (v *ValueChecker) check(val bool, errMsg interface{}) bool {
	//TODO 求值是否可以滞后
	if v.shouldSkip() {
		return false
	}

	if val {
		return true
	} else {
		return v.addErr(errMsg)
	}
}

func (v *ValueChecker) CheckBy(f func() bool, errMsg string) bool {
	return v.check(f(), errMsg)
}
func (v *ValueChecker) NotEmpty(value interface{}, errMsg string) bool {
	return v.check(!funk.IsEmpty(value), errMsg)
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

// 长度满足 len(val) >= minLength
// errMsg = "{name}长度少于{minLength}"
func (v *ValueChecker) LenLager(val interface{}, minLength int, name string) bool {
	if v.shouldSkip() {
		return false
	}
	exLen := 0
	vType := reflect.TypeOf(val)
	switch vType.Kind() {
	case reflect.String:
		exLen = len(val.(string))
	case reflect.Array, reflect.Slice, reflect.Map:
		exLen = vType.Len()
	}
	if exLen >= minLength {
		return true
	} else {
		return v.addErr(fmt.Sprintf("%s长度少于%d", name, minLength))
	}
}
func (v *ValueChecker) MustTrue(val bool, errMsg interface{}) bool {
	return v.check(val, errMsg)
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
