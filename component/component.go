/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/10
   Description :
-------------------------------------------------
*/

package component

import (
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/service/cron"
)

type componentCli struct {
	app    core.IApp
	config *core.Config
	core.ILogger
}

func NewComponent(app core.IApp) core.IComponent {
	return &componentCli{
		app:     app,
		config:  app.GetConfig().Config(),
		ILogger: app.GetLogger(),
	}
}

func (c *componentCli) App() core.IApp {
	return c.app
}

func (c *componentCli) Config() *core.Config {
	return c.config
}

func (c *componentCli) Logger() core.ILogger {
	return c.ILogger
}

func (c *componentCli) Close() {

}

func (c *componentCli) RegistryCronJob(name string, expression string, job cron.Job, enable ...bool) {
	s, ok := c.app.GetService(core.CronService)
	if !ok {
		return
	}

	s.Inject(cron.NewTask(name, expression, job, enable...))
}
