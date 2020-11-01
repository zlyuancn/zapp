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

	"github.com/zlyuancn/zapp/consts"
)

type Option func(opt *option)

type option struct {
	// 启用守护
	EnableDaemon bool
	// 服务
	Servers map[consts.ServiceType]map[string]struct{}
	// handlers
	Handlers map[HandlerType][]Handler
}

// 添加服务
func (o *option) AddService(serviceType consts.ServiceType, serviceName ...string) {
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
		panic(fmt.Errorf("服务类型[%s]的服务名[%s]已存在", serviceType.String(), name))
	}

	services[name] = struct{}{}
}

func newOption() *option {
	return &option{
		EnableDaemon: false,
		Servers:      make(map[consts.ServiceType]map[string]struct{}),
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

// 添加Cron服务
func WithCron(serviceName ...string) Option {
	return func(opt *option) {
		opt.AddService(consts.CronService, serviceName...)
	}
}
