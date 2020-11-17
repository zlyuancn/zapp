/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package app

import (
	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/service/cron"
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

func (app *appCli) RegistryCronJob(name string, expression string, job cron.Job, enable ...bool) {
	s, ok := app.GetService(core.CronService)
	if !ok {
		return
	}

	s.Inject(cron.NewTask(name, expression, job, enable...))
}
