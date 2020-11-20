/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package grpc

import (
	"context"
	"fmt"
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
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			UnaryServerLogInterceptor(app), // 日志
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(), // panic拦截
		)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: time.Minute, // 心跳
		}),
	)

	return &GrpcService{
		app:    app,
		server: server,
	}
}

func (g *GrpcService) Inject(a ...interface{}) {
	for _, v := range a {
		fn, ok := v.(func(c core.IComponent, server *grpc.Server))
		if !ok {
			g.app.GetLogger().Fatal("Grpc服务注入类型错误", zap.String("type", fmt.Sprintf("%T", v)))
		}

		fn(g.app.GetComponent(), g.server)
	}
}

func (g *GrpcService) Start() error {
	conf := g.app.GetConfig().Config().GrpcService

	listener, err := net.Listen("tcp", conf.Bind)
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	go func(errChan chan error) {
		if err := g.server.Serve(listener); err != nil {
			errChan <- err
		}
	}(errChan)

	select {
	case <-time.After(time.Millisecond * 500):
	case err := <-errChan:
		return err
	}

	g.app.GetLogger().Info("grpc启动成功", zap.String("bind", conf.Bind))
	return nil
}

func (g *GrpcService) Close() error {
	g.server.GracefulStop()
	return nil
}

// 日志拦截器
func UnaryServerLogInterceptor(app core.IApp) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := app.CreateLogger(info.FullMethod)
		ctx = utils.Context.SaveLoggerToContext(ctx, log)

		startTime := time.Now()
		log.Info("grpc.request", zap.Any("args", req))

		resp, err := handler(ctx, req)
		if err != nil {
			log.Error("grpc.response", zap.String("耗时", time.Since(startTime).String()), zap.Error(err))
		} else {
			log.Info("grpc.response", zap.String("耗时", time.Since(startTime).String()), zap.Any("resp", resp))
		}

		return resp, err
	}
}