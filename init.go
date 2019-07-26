package utils

import (
	ucore "github.com/WingGao/go-utils/ucore"
	"github.com/WingGao/go-utils/session"
	"github.com/WingGao/go-utils/whandler"
)

func init() {
	for _, v := range ucore.UtilsErrList {
		whandler.AddHandlerIgnoreErrors(v.Error())
	}
}

var (
	// error
	NewErrorList = ucore.NewErrorList
	NewWError    = ucore.NewWError
	PanicIfErr   = ucore.PanicIfErr
	FirstError   = ucore.FirstError
	//handler
	AfterHandler = whandler.AfterHandler
	CancelAfterHandler = whandler.CancelAfterHandler
	//session
	XSESSION_KEY = session.XSESSION_KEY
)
