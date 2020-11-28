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
	// debug标志
	Debug bool

	// 清理内存间隔时间(毫秒)
	FreeMemoryInterval int

	// grpc服务
	//
	// [GrpcService]
	// Bind = ":3001"
	// HeartbeatTime = 20000
	GrpcService struct {
		Bind          string // bind地址
		HeartbeatTime int    // 心跳时间(毫秒),
	}

	// 定时器
	// [CronService]
	// ThreadCount = 10
	// JobQueueSize = 100
	CronService struct {
		ThreadCount  int // 线程数
		JobQueueSize int // job队列大小
	}

	// grpc客户端
	//
	// [GrpcClient.default]
	// Address = "localhost:3001"
	// Registry = "local"
	// Balance = "round_robin"
	// DialTimeout = 1000
	GrpcClient map[string]struct {
		Address     string // 链接地址
		Registry    string // 注册器, 默认为 local
		Balance     string // 负载均衡, 默认为 round_robin
		DialTimeout int    // 连接超时(毫秒), 0表示不限, 默认为 1000
	}

	// xorm
	//
	// [Xorm.test]
	// Driver = "sqlite3"
	// Source = "test.db"
	// MaxIdleConns = 1
	// MaxOpenConns = 1
	// ConnMaxLifetime = 0
	Xorm map[string]struct {
		Driver          string // 驱动
		Source          string // 连接源
		MaxIdleConns    int    // 最大空闲连接数
		MaxOpenConns    int    // 最大连接池个数
		ConnMaxLifetime int    // 最大续航时间(毫秒, 0表示无限
	}

	// redis
	//
	// [Redis.default]
	// Address = "127.0.0.1:6379"
	// Password = ""
	// DB = 0
	// IsCluster = false
	// PoolSize = 20
	// ReadTimeout = 5000
	// WriteTimeout = 5000
	// DialTimeout = 5000
	Redis map[string]struct {
		Address      string // 地址: host1:port1,host2:port2
		Password     string // 密码
		DB           int    // db, 只有单点有效
		IsCluster    bool   // 是否为集群
		PoolSize     int    // 客户端池大小
		ReadTimeout  int64  // 超时(毫秒
		WriteTimeout int64  // 超时(毫秒
		DialTimeout  int64  // 超时(毫秒
	}

	// es7
	//
	// [ES7.default]
	// Address = "http://127.0.0.1:9200"
	// UserName = ""
	// Password = ""
	// DialTimeout = 5000
	// Sniff = false
	// Healthcheck = true
	// Retry = 0
	// RetryInterval = 0
	// GZip = false
	ES7 map[string]struct {
		Address       string // 地址: http://localhost1:9200,http://localhost2:9200
		UserName      string // 用户名
		Password      string // 密码
		DialTimeout   int64  // 连接超时(毫秒
		Sniff         bool   // 开启嗅探器
		Healthcheck   *bool  // 心跳检查(默认true
		Retry         int    // 重试次数
		RetryInterval int    // 重试间隔(毫秒)
		GZip          bool   // 启用gzip压缩
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
