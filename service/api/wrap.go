/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/30
   Description :
-------------------------------------------------
*/

package api

import (
	"github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"

	"github.com/zlyuancn/zapp/component"
)

// 处理程序
//
// 如果返回bytes会直接返回给客户端
// 返回其它值会经过处理后再返回给客户端
type Handler = func(ctx *Context) interface{}

type Response struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data,omitempty"`
}

// 包装处理程序
//
// 如果是中间件返回nil才能继续下一步, 非nil值表示拦截, 并将结果处理后返回给客户端
func Wrap(handler Handler, isMiddleware ...bool) iris.Handler {
	isProduction := !component.GlobalComponent().Config().Frame.Debug
	showDetailedErrorOfProduction := component.GlobalComponent().Config().ApiService.ShowDetailedErrorOfProduction
	return func(irisCtx *iris_context.Context) {
		ctx := makeContext(irisCtx)
		result := handler(ctx)

		if result == nil && len(isMiddleware) > 0 && isMiddleware[0] { // 返回了nil且为中间件
			ctx.Next()
			return
		}

		if err, ok := result.(error); ok {
			ctx.Values().Set("error", err)
			code, message := decodeErr(err)
			if !isProduction || showDetailedErrorOfProduction {
				message = err.Error()
			}
			_, _ = ctx.JSON(Response{
				ErrCode: code,
				ErrMsg:  message,
			})
			ctx.StopExecution()
			return
		}

		switch v := result.(type) {
		case []byte:
			ctx.Values().Set("result", "bytes")
			_, _ = ctx.Write(v)
		case *[]byte:
			ctx.Values().Set("result", "bytes")
			_, _ = ctx.Write(*v)
		default:
			ctx.Values().Set("result", result)
			_, _ = ctx.JSON(Response{
				ErrCode: 0,
				ErrMsg:  "ok",
				Data:    result,
			})
		}
		ctx.StopExecution()
	}
}
