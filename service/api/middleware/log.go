/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/29
   Description :
-------------------------------------------------
*/

package middleware

import (
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/utils"
)

func LoggerMiddleware(app core.IApp) iris.Handler {
	return func(ctx iris.Context) {
		log := app.CreateLogger(ctx.Method(), ctx.Path())
		utils.Context.SaveLoggerToIrisContext(ctx, log)

		startTime := time.Now()
		log.Debug("api.request", zap.String("query", ctx.Request().URL.RawQuery))

		ctx.Next()

		fields := []interface{}{
			"api.response", zap.String("query", ctx.Request().URL.RawQuery),
			zap.String("latency", time.Since(startTime).String()),
			zap.String("ip", ctx.RemoteAddr()),
		}

		if err, ok := ctx.Values().Get("error").(error); ok {
			if err == nil {
				err = fmt.Errorf("err{nil}")
			}
			fields = append(fields, zap.Error(err))
			log.Warn(fields...)
		} else {
			fields = append(fields, zap.Any("result", ctx.Values().Get("result")))
			log.Debug(fields...)
		}
	}
}
