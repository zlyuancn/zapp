/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/29
   Description :
-------------------------------------------------
*/

package api

import (
	"context"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/service"
	"github.com/zlyuancn/zapp/service/api/middleware"
)

// api注入函数定义
type RegisterApiRouterFunc = func(c core.IComponent, router iris.Party)

func init() {
	service.RegisterCreator(core.ApiService, new(apiCreator))
}

type apiCreator struct{}

func (*apiCreator) Create(app core.IApp) core.IService {
	return NewHttpService(app)
}

type ApiService struct {
	app core.IApp
	*iris.Application
}

func NewHttpService(app core.IApp) core.IService {
	irisApp := iris.New()
	irisApp.Logger().SetLevel("disable")
	irisApp.Use(
		middleware.LoggerMiddleware(app),
		cors.AllowAll(),
		middleware.Recover(),
	)
	irisApp.AllowMethods(iris.MethodOptions)

	return &ApiService{app: app, Application: irisApp}
}

func (a *ApiService) Inject(sc ...interface{}) {
	if len(sc) != 1 {
		a.app.Fatal("api服务注入数量必须为1个")
	}

	fn, ok := sc[0].(RegisterApiRouterFunc)
	if !ok {
		a.app.Fatal("api服务注入类型错误, 它必须能转为 api.RegisterApiRouterFunc")
	}

	fn(a.app.GetComponent(), a.Party("/"))
}

func (a *ApiService) Start() error {
	conf := a.app.GetConfig().Config().Services.ApiService

	err := service.WaitRun(&service.WaitRunOption{
		ServiceName:       "api",
		IgnoreErrs:        []error{iris.ErrServerClosed},
		FatalOnErrOfWait2: true,
		RunServiceFn: func() error {
			opts := []iris.Configurator{
				iris.WithoutBodyConsumptionOnUnmarshal, // 重复消费
				iris.WithoutPathCorrection,             // 不自动补全斜杠
				iris.WithOptimizations,                 // 启用性能优化
				iris.WithoutStartupLog,                 // 不要打印iris启动信息
				iris.WithPathEscape,                    // 解析path转义
			}
			if conf.IPWithNginxForwarded {
				opts = append(opts, iris.WithRemoteAddrHeader("X-Forwarded-For"))
			}
			if conf.IPWithNginxReal {
				opts = append(opts, iris.WithRemoteAddrHeader("X-Real-IP"))
			}
			return a.Run(iris.Addr(conf.Bind), opts...)
		},
	})
	if err != nil {
		return err
	}

	a.app.Debug("app服务启动成功", zap.String("bind", conf.Bind))
	return nil
}

func (a *ApiService) Close() error {
	err := a.Shutdown(context.Background())
	a.app.Debug("api服务已关闭")
	return err
}
