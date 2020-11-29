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
	"runtime"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"

	"github.com/zlyuancn/zapp/utils"
)

func Recover() iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			if ctx.IsStopped() { // handled by other middleware.
				return
			}

			var callers []string
			for i := 1; ; i++ {
				_, file, line, got := runtime.Caller(i)
				if !got {
					break
				}

				callers = append(callers, fmt.Sprintf("%s:%d", file, line))
			}

			logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
			logMessage += fmt.Sprint(getRequestLogs(ctx))
			logMessage += fmt.Sprintf("%s\n", err)
			logMessage += fmt.Sprintf("%s\n", strings.Join(callers, "\n"))
			log := utils.Context.MustGetLoggerFromIrisContext(ctx)
			log.Warn(logMessage)
			ctx.Values().Set("error", err)
			ctx.Values().Set("panic", err)

			_, _ = ctx.JSON(map[string]interface{}{
				"err_code": 1,
				"err_msg":  "Internal server error",
			})
			ctx.StopExecution()
		}()
		ctx.Next()
	}
}

func getRequestLogs(ctx iris.Context) string {
	var status, ip, method, path string
	status = strconv.Itoa(ctx.GetStatusCode())
	path = ctx.Path()
	method = ctx.Method()
	ip = ctx.RemoteAddr()
	return fmt.Sprintf("%v %s %s %s", status, path, method, ip)
}
