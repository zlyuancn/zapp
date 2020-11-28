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

	"github.com/go-redis/redis"
	"xorm.io/xorm"
)

// 组件, 如db, rpc, cache, mq等
type IComponent interface {
	// 获取app
	App() IApp
	// 获取配置
	Config() *Config

	// 从标准context获取日志
	CtxLog(ctx context.Context) ILogger
	// 关闭所有组件
	Close()

	IGrpcComponent
	IXormComponent
	IRedisComponent
}

type IGrpcComponent interface {
	// 注册grpc客户端建造者, 这个方法应该在app.Run之前调用
	RegistryGrpcClientCreator(name string, creator interface{})
	// 获取grpc客户端, 如果未注册grpc客户端建造者会panic
	GetGrpcClient(name string) interface{}
	// 关闭客户端
	Close()
}

type IXormComponent interface {
	// 获取
	GetXorm(name ...string) *xorm.Engine
	// 释放
	Close()
}

type IRedisComponent interface {
	// 获取redis客户端
	GetRedis(name ...string) redis.UniversalClient
	// 关闭
	Close()
}
