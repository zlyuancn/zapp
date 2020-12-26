/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/30
   Description :
-------------------------------------------------
*/

package api

import (
	"fmt"

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
func Wrap(handler Handler) iris.Handler {
	return func(irisCtx *iris_context.Context) {
		ctx := makeContext(irisCtx) // 构建上下文
		result := handler(ctx)      // 处理
		WriteToCtx(ctx, result)     // 写入结果
		ctx.StopExecution()         // 停止调用链
	}
}

// 包装中间件, 只有返回nil才能继续调用链, 非nil值表示拦截, 并将结果处理后返回给客户端
func WrapMiddleware(middleware Handler) iris.Handler {
	return func(irisCtx *iris_context.Context) {
		ctx := makeContext(irisCtx) // 构建上下文
		result := middleware(ctx)   // 处理
		if result == nil {          // 返回nil继续调用链
			ctx.Next()
			return
		}

		WriteToCtx(ctx, result) // 写入结果
		ctx.StopExecution()     // 停止调用链
	}
}

// 写入数据到ctx
//
// 如果返回bytes会直接返回给客户端
// 返回其它值会经过处理后再返回给客户端
func WriteToCtx(ctx *Context, result interface{}) {
	if err, ok := result.(error); ok {
		ctx.Values().Set("error", err)
		code, message := decodeErr(err)

		isProduction := !component.GlobalComponent().Config().Frame.Debug
		showDetailedErrorOfProduction := component.GlobalComponent().Config().Services.Api.ShowDetailedErrorOfProduction
		if !isProduction || showDetailedErrorOfProduction {
			message = err.Error()
		}
		_, _ = ctx.JSON(Response{
			ErrCode: code,
			ErrMsg:  message,
		})
		return
	}

	switch v := result.(type) {
	case []byte: // 直接写入
		ctx.Values().Set("result", fmt.Sprintf("bytes<len=%d>", len(v)))
		_, _ = ctx.Write(v)
	case *[]byte: // 直接写入
		ctx.Values().Set("result", fmt.Sprintf("bytes<len=%d>", len(*v)))
		_, _ = ctx.Write(*v)
	default:
		ctx.Values().Set("result", result)
		_, _ = ctx.JSON(Response{
			ErrCode: OK.Code,
			ErrMsg:  OK.Message,
			Data:    result,
		})
	}
}
