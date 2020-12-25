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
	"runtime"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"

	"github.com/zlyuancn/zapp/component"
	"github.com/zlyuancn/zapp/utils"
)

func Recover() iris.Handler {
	isProduction := !component.GlobalComponent().Config().Frame.Debug
	showDetailedErrorOfProduction := component.GlobalComponent().Config().Services.ApiService.ShowDetailedErrorOfProduction
	return func(ctx iris.Context) {
		err := utils.Recover.WarpCall(func() error {
			ctx.Next()
			return nil
		})
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
		logMessage += fmt.Sprintf("err: %s\n", err)
		logMessage += strings.Join(callers, "\n")
		log := utils.Context.MustGetLoggerFromIrisContext(ctx)
		log.Error(logMessage)
		ctx.Values().Set("error", err)

		result := map[string]interface{}{
			"err_code": 1,
			"err_msg":  strings.Split(logMessage, "\n"),
		}
		if isProduction && !showDetailedErrorOfProduction {
			result["err_msg"] = "service internal error"
		}
		_, _ = ctx.JSON(result)
		ctx.StopExecution()
	}
}

func getRequestLogs(ctx iris.Context) string {
	var status, ip, method, path string
	status = strconv.Itoa(ctx.GetStatusCode())
	path = ctx.Path()
	method = ctx.Method()
	ip = ctx.RemoteAddr()
	return fmt.Sprintf("%v %s %s %s\n", status, path, method, ip)
}
