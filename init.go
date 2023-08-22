package utils

import (
	"github.com/WingGao/go-utils/session"
	ucore "github.com/WingGao/go-utils/ucore"
	"wingao.net/utils/whandler"
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
	AfterHandler       = whandler.AfterHandler
	CancelAfterHandler = whandler.CancelAfterHandler
	//session
	XSESSION_KEY = session.XSESSION_KEY
)
