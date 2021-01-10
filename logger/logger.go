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
	"github.com/zlyuancn/zutils"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/core"
)

var Log core.ILogger = zlog.DefaultLogger

func NewLogger(appName string, c core.IConfig, opts ...zap.Option) core.ILogger {
	conf := c.Config().Frame
	if zutils.Reflect.IsZero(conf.Log) {
		conf.Log = zlog.DefaultConfig
		conf.Log.Name = appName
	}
	if conf.Log.Name == "" {
		conf.Log.Name = appName
	}
	c.Config().Frame.Log = conf.Log

	Log = zlog.New(conf.Log, opts...)
	return Log
}
