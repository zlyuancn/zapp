/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package app

import (
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
)

type Option func(opt *option)

type option struct {
	// 启用守护
	EnableDaemon bool
	// 忽略未启用的服务注入
	IgnoreInjectOfDisableServer bool
	// 服务
	Servers map[core.ServiceType]map[string]struct{}
	// handlers
	Handlers map[HandlerType][]Handler
}

// 添加服务
func (o *option) AddService(serviceType core.ServiceType, serviceName ...string) {
	name := consts.DefaultServiceName
	if len(serviceName) > 0 {
		name = serviceName[0]
	}

	services, ok := o.Servers[serviceType]
	if !ok {
		services = make(map[string]struct{})
		o.Servers[serviceType] = services
	}

	if _, ok = services[name]; ok {
		logger.Log.Fatal("服务类型的服务名已存在", zap.String("serviceType", serviceType.String()), zap.String("serviceName", name))
	}

	services[name] = struct{}{}
}

func newOption() *option {
	return &option{
		EnableDaemon: false,
		Servers:      make(map[core.ServiceType]map[string]struct{}),
		Handlers:     make(map[HandlerType][]Handler),
	}
}

// 启用守护进程模块
func WithEnableDaemon() Option {
	return func(opt *option) {
		opt.EnableDaemon = true
	}
}

// 添加handler
func WithHandler(t HandlerType, hs ...Handler) Option {
	return func(opt *option) {
		handlers, ok := opt.Handlers[t]
		if !ok {
			handlers = make([]Handler, 0)
		}
		handlers = append(handlers, hs...)
		opt.Handlers[t] = handlers
	}
}

// 忽略未启用的服务注入
func WithIgnoreInjectOfDisableServer(ignore ...bool) Option {
	return func(opt *option) {
		opt.IgnoreInjectOfDisableServer = len(ignore) == 0 || ignore[0]
	}
}

// 启用Cron服务
func WithCron() Option {
	return func(opt *option) {
		opt.AddService(core.CronService)
	}
}

// 启用grpc服务
func WithGrpc() Option {
	return func(opt *option) {
		opt.AddService(core.GrpcService)
	}
}
