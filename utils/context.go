/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package utils

import (
	"context"

	"github.com/kataras/iris/v12"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
)

var Context = new(contextUtil)

type contextUtil struct{}

// 基于传入的标准context生成一个新的标准context并保存log
func (c *contextUtil) SaveLoggerToContext(ctx context.Context, log core.ILogger) context.Context {
	return context.WithValue(ctx, consts.SaveFieldName_Logger, log)
}

// 从标准context中获取log
func (c *contextUtil) GetLoggerFromContext(ctx context.Context) (core.ILogger, bool) {
	value := ctx.Value(consts.SaveFieldName_Logger)
	log, ok := value.(core.ILogger)
	return log, ok
}

// 从标准context中获取log, 如果失败会panic
func (c *contextUtil) MustGetLoggerFromContext(ctx context.Context) core.ILogger {
	log, ok := c.GetLoggerFromContext(ctx)
	if !ok {
		logger.Log.Panic("can't load app_context from context")
	}
	return log
}

// 将log保存在iris上下文中
func (c *contextUtil) SaveLoggerToIrisContext(ctx iris.Context, log core.ILogger) {
	ctx.Values().Set(consts.SaveFieldName_Logger, log)
}

// 从iris上下文中获取log, 如果失败会panic
func (c *contextUtil) MustGetLoggerFromIrisContext(ctx iris.Context) core.ILogger {
	return ctx.Values().Get(consts.SaveFieldName_Logger).(core.ILogger)
}
