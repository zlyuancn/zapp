/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package app

import (
	"github.com/zlyuancn/zscheduler"
	"google.golang.org/grpc"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	_ "github.com/zlyuancn/zapp/service/cron"
	_ "github.com/zlyuancn/zapp/service/grpc"
	"github.com/zlyuancn/zapp/utils"
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

func (app *appCli) RegistryCronJob(name string, expression string, enable bool, handler func(log core.ILogger) error) {
	s, ok := app.GetService(core.CronService)
	if !ok {
		if app.opt.IgnoreInjectOfDisableServer {
			return
		}
		utils.Fatal("未启用cron服务")
	}

	task := zscheduler.NewTaskOfConfig(name, zscheduler.TaskConfig{
		Trigger:  zscheduler.NewCronTrigger(expression),
		Executor: zscheduler.NewExecutor(0, 0, 1),
		Job: func() error {
			return handler(app.CreateLogger(name))
		},
		Enable: enable,
	})

	s.Inject(task)
}

func (app *appCli) RegistryGrpcService(a func(c core.IComponent, server *grpc.Server)) {
	s, ok := app.GetService(core.GrpcService)
	if !ok {
		utils.Fatal("未启用grpc服务")
	}

	s.Inject(a)
}
