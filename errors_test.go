package utils

import (
	"testing"
	"github.com/go-errors/errors"
)

var errT = errors.New("test error")

func errorRecov() {
	defer func() {
		if err := recover(); err != nil {
			//e2 := errors.Wrap(err, 1)
			e2 := err.(error)
			PrintError(e2)
		}
	}()
	panic(errT)
}
func TestPrintError(t *testing.T) {
	PrintError(errT)

	errorRecov()
}
