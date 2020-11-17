/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/10
   Description :
-------------------------------------------------
*/

package component

import (
	"context"

	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/utils"
)

type ComponentCli struct {
	app    core.IApp
	config *core.Config
	log    core.ILogger
}

func NewComponent(app core.IApp) core.IComponent {
	return &ComponentCli{
		app:    app,
		config: app.GetConfig().Config(),
		log:    app.GetLogger(),
	}
}

func (c *ComponentCli) App() core.IApp {
	return c.app
}

func (c *ComponentCli) Config() *core.Config {
	return c.config
}

func (c *ComponentCli) Logger() core.ILogger {
	return c.log
}

func (c *ComponentCli) CtxLog(ctx context.Context) core.ILogger {
	return utils.Context.MustGetLoggerFromContext(ctx)
}

func (c *ComponentCli) Close() {
	c.log.Error("关闭组件资源")
}
