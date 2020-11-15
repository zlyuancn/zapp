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
	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/context"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
	"github.com/zlyuancn/zapp/service"
	"github.com/zlyuancn/zapp/utils"
)

type appCli struct {
	name   string
	daemon daemon.Daemon

	config    core.IConfig
	logger    core.ILogger
	component core.IComponent

	closeChan chan struct{}
	interrupt chan os.Signal

	services map[core.ServiceType]map[string]core.IService
	opt      *option
}

// 创建一个app
//
// 根据提供的app名和选项创建一个app
// 正常启动时会初始化所有服务
func NewApp(appName string, opts ...Option) core.IApp {
	if appName == "" {
		utils.Panic("appName must not empty")
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
	app.logger = logger.NewLogger(appName, app.config)
	app.component = component.NewComponent(app)

	app.logger.Info("初始化服务")

	// app初始化前
	app.handler(BeforeInitializeHandler)

	for serviceType, names := range app.opt.Servers {
		services := make(map[string]core.IService, len(names))
		for name := range names {
			services[name] = service.NewService(serviceType, app.component)
		}
		app.services[serviceType] = services
	}

	// app初始化后
	app.handler(AfterInitializeHandler)

	app.logger.Info("服务初始化完毕")

	return app
}

func (app *appCli) run() {
	app.logger.Info("启动app")

	// app启动前
	app.handler(BeforeStartHandler)

	app.startService()

	go app.freeMemory()

	// app启动后
	app.handler(AfterStartHandler)

	app.logger.Info("app已启动")

	signal.Notify(app.interrupt, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-app.interrupt
	app.exit()
}

func (app *appCli) startService() {
	app.logger.Info("启动服务")
	for serviceType, services := range app.services {
		for name, s := range services {
			if err := s.Start(); err != nil {
				app.logger.Fatal("服务启动失败", zap.String("服务类型", serviceType.String()), zap.String("服务名", name), zap.Error(err))
			}
		}
	}
}

func (app *appCli) closeService() {
	app.logger.Warn("关闭服务")
	for serviceType, services := range app.services {
		for name, s := range services {
			if err := s.Close(); err != nil {
				app.logger.Error("服务关闭失败", zap.String("服务类型", serviceType.String()), zap.String("服务名", name), zap.Error(err))
			}
		}
	}
}

func (app *appCli) enableDaemon() {
	if len(os.Args) < 2 {
		return
	}

	d, err := daemon.New(app.name, app.name, daemon.SystemDaemon)
	utils.FailOnError(err, "守护进程模块创建失败")
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
	app.logger.Warn("app准备退出")
	close(app.closeChan)

	// app退出前
	app.handler(BeforeExitHandler)

	app.closeService()
	app.component.Close()

	// app退出后
	app.handler(AfterExitHandler)

	app.logger.Warn("app已退出")
}

// 启动服务
//
// 启动所有服务并挂起进程, 直到收到退出信号或主动结束
func (app *appCli) Run() {
	app.run()
}

// 结束服务
//
// 发送退出信号
func (app *appCli) Exit() {
	app.interrupt <- syscall.SIGTERM
}

func (app *appCli) GetLogger() core.ILogger {
	return app.logger
}

func (app *appCli) GetConfig() core.IConfig {
	return app.config
}

func (app *appCli) GetComponent() core.IComponent {
	return app.component
}

func (app *appCli) NewContext(tag ...string) core.IContext {
	return context.NewContext(app, tag...)
}

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

func (app *appCli) freeMemory() {
	t := time.NewTicker(time.Second * 120)
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
