/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package grpc

import (
	"context"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/service"
	"github.com/zlyuancn/zapp/utils"
)

type RegistryGrpcServiceFunc = func(c core.IComponent, server *grpc.Server)

func init() {
	service.RegisterCreator(core.GrpcService, new(grpcCreator))
}

type grpcCreator struct{}

func (*grpcCreator) Create(app core.IApp) core.IService {
	return NewGrpcService(app)
}

type GrpcService struct {
	app    core.IApp
	server *grpc.Server
}

func NewGrpcService(app core.IApp) core.IService {
	conf := app.GetConfig().Config().Services.Grpc
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			UnaryServerLogInterceptor(app), // 日志
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(), // panic拦截
		)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: time.Duration(conf.HeartbeatTime) * time.Millisecond, // 心跳
		}),
	)

	return &GrpcService{
		app:    app,
		server: server,
	}
}

func (g *GrpcService) Inject(a ...interface{}) {
	for _, v := range a {
		fn, ok := v.(RegistryGrpcServiceFunc)
		if !ok {
			g.app.Fatal("Grpc服务注入类型错误, 它必须能转为 grpc.RegistryGrpcServiceFunc")
		}

		fn(g.app.GetComponent(), g.server)
	}
}

func (g *GrpcService) Start() error {
	conf := g.app.GetConfig().Config().Services.Grpc

	listener, err := net.Listen("tcp", conf.Bind)
	if err != nil {
		return err
	}

	err = service.WaitRun(&service.WaitRunOption{
		ServiceName:      "grpc",
		IgnoreErrs:       nil,
		ExitOnErrOfWait2: true,
		RunServiceFn: func() error {
			return g.server.Serve(listener)
		},
	})
	if err != nil {
		return err
	}

	g.app.Debug("grpc服务启动成功", zap.String("bind", conf.Bind))
	return nil
}

func (g *GrpcService) Close() error {
	g.server.GracefulStop()
	g.app.Debug("grpc服务已关闭")
	return nil
}

// 日志拦截器
func UnaryServerLogInterceptor(app core.IApp) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := app.CreateLogger(info.FullMethod)
		ctx = utils.Context.SaveLoggerToContext(ctx, log)

		startTime := time.Now()
		log.Debug("grpc.request", zap.Any("req", req))

		resp, err := handler(ctx, req)
		if err != nil {
			log.Error("grpc.response", zap.String("latency", time.Since(startTime).String()), zap.Error(err))
		} else {
			log.Debug("grpc.response", zap.String("latency", time.Since(startTime).String()), zap.Any("resp", resp))
		}

		return resp, err
	}
}
