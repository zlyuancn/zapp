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
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
)

var Log core.ILogger = zlog.DefaultLogger

func NewLogger(appName string, c core.IConfig) core.ILogger {
	var conf = zlog.DefaultConfig
	conf.Name = appName

	viper := c.GetViper()
	if viper.IsSet(consts.ConfigGroupName_Log) {
		if err := viper.UnmarshalKey(consts.ConfigGroupName_Log, &conf); err != nil {
			Log.Fatal("解析log配置失败", zap.Error(err))
		}
		if conf.Name == "" { // 考虑用户主动设置name为空
			conf.Name = appName
		}
	}

	return Log
}
