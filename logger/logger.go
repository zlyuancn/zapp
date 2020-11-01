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

func NewLogger(c core.IConfig) core.ILogger {
	var conf zlog.LogConfig
	utils.FailOnErrorf(c.ParseShard(consts.LogConfigShardName, &conf), "解析log配置失败")
	return zlog.New(conf)
}
