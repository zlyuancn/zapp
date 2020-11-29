/*
-------------------------------------------------
   Author :       Zhang Fan
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
		// todo 打印请求信息
		log.Debug("api.request")

		ctx.Next()
		if err, ok := ctx.Values().Get("error").(error); ok {
			if err == nil {
				err = fmt.Errorf("nil")
			}
			log.Warn("api.response", zap.String("spent_time", time.Since(startTime).String()), zap.Error(err))
		} else {
			log.Debug("api.response", zap.String("spent_time", time.Since(startTime).String()))
		}
	}
}
