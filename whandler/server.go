package whandler

import (
	"github.com/kataras/iris/v12"
)

type IServer interface {
	RegisterIris(p iris.Party, prefix string)
}
