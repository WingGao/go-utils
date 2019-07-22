package whandler

import (
	"github.com/kataras/iris"
)

type IServer interface {
	RegisterIris(p iris.Party, prefix string)
}
