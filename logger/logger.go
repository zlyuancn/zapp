/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/7/10
   Description :
-------------------------------------------------
*/

package logger

import (
	"github.com/zlyuancn/zlog"

	"github.com/zlyuancn/zapp/core"
)

var Log core.ILogger = zlog.DefaultLogger

func NewLogger(appName string, c core.IConfig) core.ILogger {
	Log = zlog.New(c.Config().Log)
	return Log
}
