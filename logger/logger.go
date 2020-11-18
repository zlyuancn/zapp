/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/10
   Description :
-------------------------------------------------
*/

package logger

import (
	"github.com/zlyuancn/zlog"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
)

var Log core.ILogger = zlog.DefaultLogger

func NewLogger(appName string, c core.IConfig) core.ILogger {
	var conf = zlog.DefaultConfig
	conf.Name = appName

	viper := c.GetViper()
	if viper.IsSet(consts.ConfigShardName_Log) {
		if err := viper.UnmarshalKey(consts.ConfigShardName_Log, &conf); err != nil {
			Log.Fatal("解析log配置失败", zap.Error(err))
		}
	}

	Log = zlog.New(conf)
	return Log
}
