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
	GrpcService        struct { // grpc服务
		Bind string
	}
	Cron struct {
		ThreadCount  int // 线程数
		JobQueueSize int // job队列大小
	} // 定时器
	GrpcClient map[string]struct { // grpc客户端
		Address  string // 链接地址
		Registry string // 注册器, 默认为 local
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
