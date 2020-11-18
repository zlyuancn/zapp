/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package app

import (
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
	"github.com/zlyuancn/zapp/utils"
)

type appCli struct {
	name   string
	daemon daemon.Daemon

	opt *option

	closeChan chan struct{}
	interrupt chan os.Signal

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
		utils.Fatal("appName must not empty")
	}
	app := &appCli{
		name:      appName,
		closeChan: make(chan struct{}),
		interrupt: make(chan os.Signal, 1),
		services:  make(map[core.ServiceType]map[string]core.IService),
		opt:       newOption(),
	}

	// 初始化选项
	for _, o := range opts {
		o(app.opt)
	}

	// 处理选项
	if app.opt.EnableDaemon {
		app.enableDaemon()
	}

	app.config = config.NewConfig()
	app.ILogger = logger.NewLogger(appName, app.config)
	app.component = component.NewComponent(app)

	app.Info("初始化服务")

	// app初始化前
	app.handler(BeforeInitializeHandler)

	for serviceType, names := range app.opt.Servers {
		services := make(map[string]core.IService, len(names))
		for name := range names {
			services[name] = service.NewService(serviceType, app)
		}
		app.services[serviceType] = services
	}

	// app初始化后
	app.handler(AfterInitializeHandler)

	app.Info("服务初始化完毕")

	return app
}

func (app *appCli) run() {
	app.Info("启动app")

	// app启动前
	app.handler(BeforeStartHandler)

	app.startService()

	go app.freeMemory()

	// app启动后
	app.handler(AfterStartHandler)

	app.Info("app已启动")

	signal.Notify(app.interrupt, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-app.interrupt
	app.exit()
}

func (app *appCli) startService() {
	app.Info("启动服务")
	for serviceType, services := range app.services {
		for name, s := range services {
			if err := s.Start(); err != nil {
				app.Fatal("服务启动失败", zap.String("serviceType", serviceType.String()), zap.String("serviceName", name), zap.Error(err))
			}
		}
	}
}

func (app *appCli) closeService() {
	app.Info("关闭服务")
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

	d, err := daemon.New(app.name, app.name, daemon.SystemDaemon)
	utils.FatalOnError(err, "守护进程模块创建失败")
	app.daemon = d

	var out string
	switch os.Args[1] {
	case "install":
		out, err = app.daemon.Install(os.Args[2:]...)
	case "remove":
		out, err = app.daemon.Remove()
	case "start":
		out, err = app.daemon.Start()
	case "stop":
		out, err = app.daemon.Stop()
	case "status":
		out, err = app.daemon.Status()
	default:
		return
	}

	if err != nil {
		fmt.Println(out, err)
		os.Exit(1)
	}

	fmt.Println(out)
	os.Exit(0)
}

func (app *appCli) exit() {
	app.Warn("app准备退出")
	close(app.closeChan)

	// app退出前
	app.handler(BeforeExitHandler)

	app.closeService()
	app.closeComponentResource()

	// app退出后
	app.handler(AfterExitHandler)

	app.Warn("app已退出")
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

func (app *appCli) freeMemory() {
	t := time.NewTicker(time.Millisecond * time.Duration(app.config.Config().FreeMemoryInterval))
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
