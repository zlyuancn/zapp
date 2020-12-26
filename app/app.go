/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/takama/daemon"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/component"
	"github.com/zlyuancn/zapp/config"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
	"github.com/zlyuancn/zapp/service"
)

type appCli struct {
	name string

	opt *option

	closeChan     chan struct{}
	interrupt     chan os.Signal
	baseCtx       context.Context
	baseCtxCancel context.CancelFunc

	config core.IConfig

	loggerId uint64
	core.ILogger

	component core.IComponent
	services  map[core.ServiceType]map[string]core.IService
}

// 创建一个app
//
// 根据提供的app名和选项创建一个app
// 正常启动时会初始化所有服务
func NewApp(appName string, opts ...Option) core.IApp {
	if appName == "" {
		logger.Log.Fatal("appName must not empty")
	}

	app := &appCli{
		name:      appName,
		closeChan: make(chan struct{}),
		interrupt: make(chan os.Signal, 1),
		services:  make(map[core.ServiceType]map[string]core.IService),
		opt:       newOption(),
	}
	app.baseCtx, app.baseCtxCancel = context.WithCancel(context.Background())

	// 初始化选项
	for _, o := range opts {
		o(app.opt)
	}

	// 处理选项
	if app.opt.EnableDaemon {
		app.enableDaemon()
	}

	app.config = config.NewConfig(appName, app.opt.configOpts...)
	app.ILogger = logger.NewLogger(appName, app.config)

	app.Debug("app初始化")
	app.handler(BeforeInitializeHandler)

	// 初始化组件
	app.component = component.NewComponent(app)

	// 初始化服务
	for serviceType, names := range app.opt.Servers {
		services := make(map[string]core.IService, len(names))
		for name := range names {
			services[name] = service.NewService(serviceType, app)
		}
		app.services[serviceType] = services
	}

	app.handler(AfterInitializeHandler)
	app.Debug("app初始化完毕")

	return app
}

func (app *appCli) run() {
	app.Debug("启动app")
	app.handler(BeforeStartHandler)

	// 启动服务
	app.startService()

	go app.freeMemory()

	app.handler(AfterStartHandler)
	app.Info("app已启动")

	signal.Notify(app.interrupt, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-app.interrupt
	app.exit()
}

func (app *appCli) startService() {
	app.Debug("启动服务")
	for serviceType, services := range app.services {
		for name, s := range services {
			if err := s.Start(); err != nil {
				app.Fatal("服务启动失败", zap.String("serviceType", serviceType.String()), zap.String("serviceName", name), zap.Error(err))
			}
		}
	}
}

func (app *appCli) closeService() {
	app.Debug("关闭服务")
	for serviceType, services := range app.services {
		for name, s := range services {
			if err := s.Close(); err != nil {
				app.Error("服务关闭失败", zap.String("serviceType", serviceType.String()), zap.String("serviceName", name), zap.Error(err))
			}
		}
	}
}

func (app *appCli) enableDaemon() {
	if len(os.Args) < 2 {
		return
	}

	switch os.Args[1] {
	case "install":
	case "remove":
	case "start":
	case "stop":
	case "status":
	default:
		return
	}

	d, err := daemon.New(app.name, app.name, daemon.SystemDaemon)
	if err != nil {
		logger.Log.Fatal("守护进程模块创建失败", zap.Error(err))
	}

	var out string
	switch os.Args[1] {
	case "install":
		out, err = d.Install(os.Args[2:]...)
	case "remove":
		out, err = d.Remove()
	case "start":
		out, err = d.Start()
	case "stop":
		out, err = d.Stop()
	case "status":
		out, err = d.Status()
	}

	if err != nil {
		fmt.Println(out, err)
		os.Exit(1)
	}

	fmt.Println(out)
	os.Exit(0)
}

func (app *appCli) exit() {
	app.Debug("app准备退出")
	close(app.closeChan)

	// app退出前
	app.handler(BeforeExitHandler)

	// 关闭基础上下文
	app.baseCtxCancel()
	// 关闭服务
	app.closeService()
	// 释放组件资源
	app.closeComponentResource()

	// app退出后
	app.handler(AfterExitHandler)
	app.Warn("app已退出")
}

func (app *appCli) Name() string {
	return app.name
}

// 启动服务
//
// 启动所有服务并挂起进程, 直到收到退出信号或主动结束
func (app *appCli) Run() {
	app.run()
}

// 关闭app
func (app *appCli) Close() {
	app.interrupt <- syscall.SIGTERM
}

func (app *appCli) GetConfig() core.IConfig {
	return app.config
}

func (app *appCli) BaseContext() context.Context {
	return app.baseCtx
}

func (app *appCli) freeMemory() {
	interval := app.config.Config().Frame.FreeMemoryInterval
	if interval <= 0 {
		return
	}

	t := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for {
		select {
		case <-app.closeChan:
			t.Stop()
			return
		case <-t.C:
			debug.FreeOSMemory()
		}
	}
}

func (app *appCli) handler(t HandlerType) {
	for _, h := range app.opt.Handlers[t] {
		h(app, t)
	}
}
