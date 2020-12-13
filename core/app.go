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
	Close()
	// 获取配置
	GetConfig() IConfig

	// 日志组件
	ILogger
	// 获取日志组件
	GetLogger() ILogger
	// 创建日志组件副本
	CreateLogger(tag ...string) ILogger

	// 获取组件
	GetComponent() IComponent

	// 获取服务
	GetService(serviceType ServiceType, serviceName ...string) (IService, bool)
	// 注册api路由
	RegistryApiRouter(a interface{})
	// 注册cron任务
	RegistryCronJob(name string, expression string, enable bool, handler func(job ICronJob) error)
	// 注册cron任务自定义
	RegistryCronJobCustom(name string, trigger zscheduler.ITrigger, executor zscheduler.IExecutor, enable bool, handler func(job ICronJob) error)
	// 注册grpc服务
	RegistryGrpcService(a ...interface{})
	// 注册mysql-binlog服务handler
	RegistryMysqlBinlogHandler(a interface{})
}
