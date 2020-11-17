/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package core

import (
	"context"

	"github.com/zlyuancn/zscheduler"
)

// 组件, 如db, rpc等
type IComponent interface {
	// 获取app
	App() IApp
	// 获取配置
	Config() *Config
	// 日志
	Logger(ctx ...context.Context) ILogger
	// 关闭
	Close()

	// 注册cron任务
	RegistryCronJob(name string, expression string, handler zscheduler.Job, enable ...bool)
}
