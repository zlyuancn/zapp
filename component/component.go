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
	"sync"

	"github.com/zlyuancn/zapp/component/es7"
	"github.com/zlyuancn/zapp/component/grpc"
	"github.com/zlyuancn/zapp/component/redis"
	"github.com/zlyuancn/zapp/component/validator"
	"github.com/zlyuancn/zapp/component/xorm"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/utils"
)

type ComponentCli struct {
	app    core.IApp
	config *core.Config
	core.ILogger

	core.IValidator

	core.IGrpcComponent
	core.IXormComponent
	core.IRedisComponent
	core.IES7Component
}

func NewComponent(app core.IApp) core.IComponent {
	return &ComponentCli{
		app:     app,
		config:  app.GetConfig().Config(),
		ILogger: app.GetLogger(),

		IValidator: validator.NewValidator(),

		IGrpcComponent:  grpc.NewClient(app),
		IXormComponent:  xorm.NewXorm(app),
		IRedisComponent: redis.NewRedis(app),
		IES7Component:   es7.NewES7(app),
	}
}

func (c *ComponentCli) App() core.IApp       { return c.app }
func (c *ComponentCli) Config() *core.Config { return c.config }

func (c *ComponentCli) CtxLog(ctx context.Context) core.ILogger {
	return utils.Context.MustGetLoggerFromContext(ctx)
}

func (c *ComponentCli) Close() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.IGrpcComponent.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.IXormComponent.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.IRedisComponent.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.IES7Component.Close()
	}()

	wg.Wait()
}
