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

type Handler = func(ctx *Context) interface{}

type Response struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data,omitempty"`
}

func Wrap(handler Handler) iris.Handler {
	isProduction := !component.GlobalComponent().Config().Frame.Debug
	showDetailedErrorOfProduction := component.GlobalComponent().Config().ApiService.ShowDetailedErrorOfProduction
	return func(irisCtx *iris_context.Context) {
		ctx := newContext(irisCtx)
		result := handler(ctx)

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
