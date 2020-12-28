/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package core

import (
	"context"

	"github.com/zlyuancn/zscheduler"
)

// app
//
// 用于将所有模块连起来
type IApp interface {
	// app名
	Name() string
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
	// 基础上下文, 这个用于监听服务结束, app会在关闭服务之前调用cancel()
	BaseContext() context.Context

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
	//
	// a 必须是 func(c core.IComponent, router iris.Party)
	// 示例:
	//    a.RegistryApiRouter(func(c core.IComponent, router iris.Party) {
	//        router.Get("/", api.Wrap(func(ctx *api.Context) interface{} { return "hello" }))
	//    })
	RegistryApiRouter(a interface{})
	// 注册cron任务
	RegistryCronJob(name string, expression string, enable bool, handler func(job ICronJob) error)
	// 注册cron任务自定义
	RegistryCronJobCustom(name string, trigger zscheduler.ITrigger, executor zscheduler.IExecutor, enable bool, handler func(job ICronJob) error)
	// 注册grpc服务
	//
	// a 必须是 func(c core.IComponent, server *grpc.Server)
	// 用户必须在函数 a 中对 grpc 服务实体注册
	// 示例:
	//     app.RegistryGrpcService(func(c core.IComponent, service *grpc.Server) {
	//         pb.RegisterXXXServer(service, &srvObj)
	//     })
	RegistryGrpcService(a ...interface{})
	// 注册mysql-binlog服务handler
	//
	// a 必须是 func(c core.IComponent) mysql_binlog.IEventHandler
	// 示例:
	//     type Handler struct {
	//         c core.IComponent
	//         mysql_binlog.BaseEventHandler
	//     }
	//     app.RegistryMysqlBinlogHandler(func(c core.IComponent) mysql_binlog.IEventHandler {
	//         return &Handler{c: c}
	//     })
	RegistryMysqlBinlogHandler(a interface{})
}
