/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package core

import (
	"context"

	"google.golang.org/grpc"
)

// 组件, 如db, rpc, cache, mq等
type IComponent interface {
	// 获取app
	App() IApp
	// 获取配置
	Config() *Config

	// 从标准context获取日志
	CtxLog(ctx context.Context) ILogger

	IGrpcClient
}

type IGrpcClient interface {
	// 获取grpc客户端
	GetGrpcClient(name string, creator func(conn *grpc.ClientConn) interface{}) interface{}
	Close()
}
