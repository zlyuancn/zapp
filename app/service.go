/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package app

import (
	"github.com/zlyuancn/zscheduler"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
	_ "github.com/zlyuancn/zapp/service/api"
	cron_service "github.com/zlyuancn/zapp/service/cron"
	_ "github.com/zlyuancn/zapp/service/grpc"
	_ "github.com/zlyuancn/zapp/service/mysql-binlog"
)

func (app *appCli) GetService(serviceType core.ServiceType, serviceName ...string) (core.IService, bool) {
	name := consts.DefaultServiceName
	if len(serviceName) > 0 {
		name = serviceName[0]
	}

	services, ok := app.services[serviceType]
	if !ok {
		return nil, false
	}

	s, ok := services[name]
	return s, ok
}

// 注册api路由
func (app *appCli) RegistryApiRouter(a interface{}) {
	s, ok := app.GetService(core.ApiService)
	if !ok {
		if app.opt.IgnoreInjectOfDisableServer {
			return
		}
		logger.Log.Fatal("未启用api服务")
	}

	s.Inject(a)
}

func (app *appCli) RegistryCronJob(name string, expression string, enable bool, handler func(cronJob core.ICronJob) error) {
	s, ok := app.GetService(core.CronService)
	if !ok {
		if app.opt.IgnoreInjectOfDisableServer {
			return
		}
		logger.Log.Fatal("未启用cron服务")
	}

	task := zscheduler.NewTaskOfConfig(name, zscheduler.TaskConfig{
		Trigger:  zscheduler.NewCronTrigger(expression),
		Executor: zscheduler.NewExecutor(0, 0, 1),
		Handler: func(job zscheduler.IJob) error {
			return handler(cron_service.NewCronJob(job))
		},
		Enable: enable,
	})

	s.Inject(task)
}
func (app *appCli) RegistryCronJobCustom(name string, trigger zscheduler.ITrigger, executor zscheduler.IExecutor, enable bool, handler func(cronJob core.ICronJob) error) {
	s, ok := app.GetService(core.CronService)
	if !ok {
		if app.opt.IgnoreInjectOfDisableServer {
			return
		}
		logger.Log.Fatal("未启用cron服务")
	}

	task := zscheduler.NewTaskOfConfig(name, zscheduler.TaskConfig{
		Trigger:  trigger,
		Executor: executor,
		Handler: func(job zscheduler.IJob) error {
			return handler(cron_service.NewCronJob(job))
		},
		Enable: enable,
	})

	s.Inject(task)
}

// 注册grpc服务
func (app *appCli) RegistryGrpcService(a ...interface{}) {
	s, ok := app.GetService(core.GrpcService)
	if !ok {
		if app.opt.IgnoreInjectOfDisableServer {
			return
		}
		logger.Log.Fatal("未启用grpc服务")
	}

	s.Inject(a...)
}

// 注册mysql-binlog服务handler
func (app *appCli) RegistryMysqlBinlogHandler(a interface{}) {
	s, ok := app.GetService(core.MysqlBinlogService)
	if !ok {
		if app.opt.IgnoreInjectOfDisableServer {
			return
		}
		logger.Log.Fatal("未启用mysql-binlog服务")
	}

	s.Inject(a)
}
