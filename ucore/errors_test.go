package ucore

import (
	"errors"
	"fmt"

	//ge "github.com/WingGao/errors"
	"testing"
)

var errT = errors.New("test error")

func errorRecov(e error) {
	defer func() {
		if err := recover(); err != nil {
			//e2 := errors.Wrap(err, 1)
			e2 := err.(error)
			PrintError(e2)
		}
	}()
	errDep2(e)
}

func errDep2(e error){
	panic(e)
}

func TestPrintError(t *testing.T) {
	PrintError(errT)

	errorRecov(errors.New("hello"))
}


func errorRecovWe(e error) {
	defer func() {
		if err := recover(); err != nil {
			//e2 := errors.Wrap(err, 1)
			we := NewWError(err)
			we.Fmt()
			fmt.Println(we.ErrorStack())
		}
	}()
	errDep2(e)
}

func TestNewWError(t *testing.T) {
	err := errors.New("test-error")
	errorRecovWe(err)
}
