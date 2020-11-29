/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package app

import (
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc"

	"github.com/zlyuancn/zscheduler"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
	api_service "github.com/zlyuancn/zapp/service/api"
	cron_service "github.com/zlyuancn/zapp/service/cron"
	grpc_service "github.com/zlyuancn/zapp/service/grpc"
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

func (app *appCli) RegistryApiRouter(fn func(c core.IComponent, router iris.Party)) {
	s, ok := app.GetService(core.ApiService)
	if !ok {
		if app.opt.IgnoreInjectOfDisableServer {
			return
		}
		logger.Log.Fatal("未启用api服务")
	}

	s.Inject(api_service.RegisterApiRouterFunc(fn))
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

func (app *appCli) RegistryGrpcService(a func(c core.IComponent, server *grpc.Server)) {
	s, ok := app.GetService(core.GrpcService)
	if !ok {
		if app.opt.IgnoreInjectOfDisableServer {
			return
		}
		logger.Log.Fatal("未启用grpc服务")
	}

	s.Inject(grpc_service.RegistryGrpcServiceFunc(a))
}
