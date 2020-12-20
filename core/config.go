/*
-------------------------------------------------
   Author :       zlyuancn
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
	// 框架配置
	Frame struct {
		// debug标志
		Debug bool
		// 清理内存间隔时间(毫秒)
		FreeMemoryInterval int
		// 等待服务启动阶段1, 等待时间(毫秒), 如果时间到则临时认为服务启动成功并提前返回
		WaitServiceRunTime int
		// 等待服务启动阶段2, 等待服务启动阶段1时间到后继续等待服务启动, 等待时间(毫秒), 如果时间到则真正认为服务启动成功
		ContinueWaitServiceRunTime int
	}

	// api服务
	//
	// [ApiService]
	// Bind = ":3000"
	// IPWithNginxForwarded = false
	// IPWithNginxReal = false
	// ShowDetailedErrorOfProduction = false
	ApiService struct {
		Bind                          string // bind地址
		IPWithNginxForwarded          bool   // 适配nginx的Forwarded获取ip, 优先级高于nginx的Real
		IPWithNginxReal               bool   // 适配nginx的Real获取ip, 优先级高于sock连接的ip
		ShowDetailedErrorOfProduction bool   // 生产环境显示详细的错误
	}

	// grpc服务
	//
	// [GrpcService]
	// Bind = ":3001"
	// HeartbeatTime = 20000
	GrpcService struct {
		Bind          string // bind地址
		HeartbeatTime int    // 心跳时间(毫秒),
	}

	// 定时服务
	//
	// [CronService]
	// ThreadCount = 10
	// JobQueueSize = 100
	CronService struct {
		ThreadCount  int // 线程数
		JobQueueSize int // job队列大小
	}

	// mysql-binlog服务
	//
	// [MysqlBinlogService]
	// Host = "localhost:3306"
	// UserName = "root"
	// Password = "password"
	// Charset = "utf8mb4"
	// IncludeTableRegex = []
	// ExcludeTableRegex = []
	// DiscardNoMetaRowEvent = false
	// DumpExecutionPath = ""
	// IgnoreWKBDataParseError = true
	MysqlBinlogService struct {
		Host                    string   // mysql 主机地址
		UserName                string   // 用户名, 最好是root
		Password                string   // 密码
		Charset                 *string  // 字符集, 一般为utf8mb4
		IncludeTableRegex       []string // 包含的表正则匹配, 匹配的数据为 dbName.tableName
		ExcludeTableRegex       []string // 排除的表正则匹配, 匹配的数据为 dbName.tableName
		DiscardNoMetaRowEvent   bool     // 放弃没有表元数据的row事件
		DumpExecutionPath       string   // mysqldump执行路径, 如果为空则忽略mysqldump只使用binlog, mysqldump执行路径一般为mysqldump
		IgnoreWKBDataParseError bool     // 忽略wkb数据解析错误, 一般为POINT, GEOMETRY类型
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

	// cache
	//
	// [Cache]
	// CacheDB = "memory"
	// Codec = "msgpack"
	// DirectReturnOnCacheFault = true
	// PanicOnLoaderExists = true
	Cache CacheConfig
}

// 配置
type IConfig interface {
	// 获取配置
	Config() *Config
	// 解析数据到结构中
	Parse(outPtr interface{}) error
	// 解析指定分片的数据到结构中
	ParseShard(shard string, outPtr interface{}) error
	// 获取配置viper结构
	GetViper() *viper.Viper
}

// 缓存配置
type CacheConfig struct {
	CacheDB                  string // 缓存db; default, memory, redis
	Codec                    string // 编解码器; default, byte, json, jsoniter, msgpack, proto_buffer
	DirectReturnOnCacheFault bool   // 在缓存故障时直接返回缓存错误(默认)
	PanicOnLoaderExists      bool   // 注册加载器时如果加载器已存在会panic(默认), 设为false会替换旧的加载器

	MemoryCacheDB struct {
		CleanupInterval int64 // 清除过期key时间间隔(毫秒)
	}
	RedisCacheDB struct {
		KeyPrefix    string // key前缀
		Address      string // 地址: host1:port1,host2:port2
		Password     string // 密码
		DB           int    // db, 只有单点有效
		IsCluster    bool   // 是否为集群
		PoolSize     int    // 客户端池大小
		ReadTimeout  int64  // 读取超时(毫秒
		WriteTimeout int64  // 写入超时(毫秒
		DialTimeout  int64  // 连接超时(毫秒
	}
}
