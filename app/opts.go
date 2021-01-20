/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package app

import (
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/config"
	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
)

type Option func(opt *option)

type option struct {
	// 配置选项
	ConfigOpts []config.Option
	// 启用守护
	EnableDaemon bool
	// 忽略未启用的服务注入
	IgnoreInjectOfDisableServer bool
	// 服务
	Services map[core.ServiceType]map[string]struct{}
	// 自定义启用服务函数
	customEnableServicesFn func(c core.IComponent) (servers map[core.ServiceType]map[string]bool)
	// handlers
	Handlers map[HandlerType][]Handler
}

func newOption() *option {
	return &option{
		EnableDaemon:                false,
		IgnoreInjectOfDisableServer: false,
		Services:                    make(map[core.ServiceType]map[string]struct{}),
		Handlers:                    make(map[HandlerType][]Handler),
	}
}

// 添加服务
func (o *option) AddService(serviceType core.ServiceType, serviceName ...string) {
	name := consts.DefaultServiceName
	if len(serviceName) > 0 {
		name = serviceName[0]
	}

	services, ok := o.Services[serviceType]
	if !ok {
		services = make(map[string]struct{})
		o.Services[serviceType] = services
	}

	if _, ok = services[name]; ok {
		logger.Log.Fatal("服务类型的服务名已存在", zap.String("serviceType", string(serviceType)), zap.String("serviceName", name))
	}

	services[name] = struct{}{}
}

// 检查自定义启用服务
func (o *option) CheckCustomEnableServices(c core.IComponent) {
	if o.customEnableServicesFn == nil {
		return
	}

	customServices := o.customEnableServicesFn(c)
	for serviceType, names := range customServices {
		for name, enable := range names {
			if !enable {
				continue
			}
			o.AddService(serviceType, name)
		}
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

// 设置config选项
func WithConfigOption(opts ...config.Option) Option {
	return func(opt *option) {
		opt.ConfigOpts = append(opt.ConfigOpts, opts...)
	}
}

// 启用api服务
func WithApiService() Option {
	return func(opt *option) {
		opt.AddService(core.ApiService)
	}
}

// 启用cron服务
func WithCronService() Option {
	return func(opt *option) {
		opt.AddService(core.CronService)
	}
}

// 启用grpc服务
func WithGrpcService() Option {
	return func(opt *option) {
		opt.AddService(core.GrpcService)
	}
}

// 启动mysql-binlog服务
func WithMysqlBinlogService() Option {
	return func(opt *option) {
		opt.AddService(core.MysqlBinlogService)
	}
}

// 启动服务
func WithAddService(serviceType core.ServiceType, serviceName ...string) Option {
	return func(opt *option) {
		opt.AddService(serviceType, serviceName...)
	}
}

// 自定义启用哪些服务
//
// 与其他启用服务选项不同, 这里已经加载了component, 用户可以方便的根据各种条件启用想要的服务.
// 示例:
//      app.WithCustomEnableService(func(c core.IComponent) (servers map[core.ServiceType]map[string]bool) {
//			return map[core.ServiceType]map[string]bool{
//				core.CronService: {"default": true},
//			}
//		})
func WithCustomEnableService(fn func(c core.IComponent) (servers map[core.ServiceType]map[string]bool)) Option {
	return func(opt *option) {
		opt.customEnableServicesFn = fn
	}
}
