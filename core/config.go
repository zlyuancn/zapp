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
	ApiBind            string // api的bind地址
	FreeMemoryInterval int    // 清理内存间隔时间(毫秒)
	Cron               struct {
		ThreadCount  int // 线程数
		JobQueueSize int // job队列大小
	} // 定时器
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
