package wrpc

import (
	gcontext "context"
	"github.com/kataras/iris/v12/context"
)
func ToIrisContext(ctx gcontext.Context) context.Context  {
	context.NewContext(nil)
}
