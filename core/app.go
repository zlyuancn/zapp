/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package core

import (
	"github.com/zlyuancn/zapp/consts"
)

// app
//
// 用于将所有模块连起来
type IApp interface {
	// 运行
	//
	// 开启所有服务并挂起
	Run()
	// 退出
	//
	// 结束所有服务并退出
	Exit()

	// 获取配置
	GetConfig() IConfig
	// 获取日志组件
	GetLogger() ILogger
	// 获取组件
	GetComponent() IComponent
	// 获取上下文
	GetContext() IContext
	// 获取服务
	GetService(serviceType consts.ServiceType, serviceName ...string) (IService, bool)
}
