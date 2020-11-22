/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package core

import (
	"github.com/spf13/viper"
)

// 配置结构
type Config struct {
	Debug              bool
	FreeMemoryInterval int // 清理内存间隔时间(毫秒)

	// grpc服务
	//
	// [GrpcService]
	// Bind = ":3001"
	GrpcService struct {
		Bind string
	}

	// 定时器
	CronService struct {
		ThreadCount  int // 线程数
			JobQueueSize int // job队列大小
	}

	// grpc客户端
	//
	// [GrpcClient.test]
	// Address = "localhost:3001"
	GrpcClient map[string]struct {
		Address  string // 链接地址
		Registry string // 注册器, 默认为 local
		Balance  string // 负载均衡, 默认为 round_robin
	}
}

// 配置
type IConfig interface {
	// 获取配置
	Config() *Config
	// 加载的配置文件
	ConfigFiles() []string
	// 解析数据到结构中
	Parse(outPtr interface{}) error
	// 解析指定分片的数据到结构中
	ParseShard(shard string, outPtr interface{}) error
	// 获取配置viper结构
	GetViper() *viper.Viper
}
