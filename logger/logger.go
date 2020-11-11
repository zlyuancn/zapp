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

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/utils"
)

func NewLogger(appName string, c core.IConfig) core.ILogger {
	var conf = zlog.DefaultConfig
	conf.Name = appName

	viper := c.GetViper()
	if viper.IsSet(consts.LogConfigShardName) {
		utils.FailOnErrorf(viper.UnmarshalKey(consts.LogConfigShardName, &conf), "解析log配置失败")
	}

	log := zlog.New(conf)
	utils.SetLog(log)
	return log
}
