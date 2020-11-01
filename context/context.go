/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/1
   Description :
-------------------------------------------------
*/

package context

import (
	"github.com/zlyuancn/zapp/core"
)

type contextCli struct {
	app    core.IApp
	config *core.Config
	core.ILogger
}

func NewContext(c core.IComponent) core.IContext {
	return &contextCli{
		app:     c.App(),
		config:  c.Config(),
		ILogger: c.App().GetLogger(),
	}
}

func (c *contextCli) App() core.IApp {
	return c.app
}

func (c *contextCli) Config() *core.Config {
	return c.config
}
