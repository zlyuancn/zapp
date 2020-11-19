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

	"github.com/zlyuancn/zapp/component/grpc"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/utils"
)

type ComponentCli struct {
	app    core.IApp
	config *core.Config
	log    core.ILogger

	core.IGrpcClient
}

func NewComponent(app core.IApp) core.IComponent {
	return &ComponentCli{
		app:    app,
		config: app.GetConfig().Config(),
		log:    app.GetLogger(),

		IGrpcClient: grpc.NewClient(app),
	}
}

func (c *ComponentCli) App() core.IApp       { return c.app }
func (c *ComponentCli) Config() *core.Config { return c.config }
func (c *ComponentCli) Logger() core.ILogger { return c.log }

func (c *ComponentCli) CtxLog(ctx context.Context) core.ILogger {
	return utils.Context.MustGetLoggerFromContext(ctx)
}

func (c *ComponentCli) Close() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.IGrpcClient.Close()
	}()

	wg.Wait()
}
