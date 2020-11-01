/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package core

import (
	"github.com/zlyuancn/zscheduler"
)

// 组件, 如db, rpc等
type IComponent interface {
	IContext
	// 关闭
	Close()

	// 注册cron任务
	RegistryCronJob(name string, expression string, handler zscheduler.Job, enable ...bool)
}
