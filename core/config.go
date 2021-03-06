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
	"github.com/zlyuancn/zlog"
)

// 配置结构
type Config struct {
	// 框架配置
	Frame FrameConfig

	// 服务配置
	Services ServicesConfig

	// 组件配置
	Components ComponentsConfig
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

	// 获取标记列表
	Tags() []string
	// 检查是否存在某个标记, 标记是忽略大小写的
	HasTag(tag string) (ok bool)
}

// frame配置
type FrameConfig struct {
	// debug标志
	Debug bool
	// 主动清理内存间隔时间(毫秒), <= 0 表示禁用
	FreeMemoryInterval int
	// 等待服务启动阶段1, 等待时间(毫秒), 如果时间到则临时认为服务启动成功并提前返回
	WaitServiceRunTime int
	// 等待服务启动阶段2, 等待服务启动阶段1时间到后继续等待服务启动, 等待时间(毫秒), 如果时间到则真正认为服务启动成功
	ContinueWaitServiceRunTime int
	// 标记列表, 注意: tag是忽略大小写的
	Tags []string
	// log配置
	Log LogConfig
}

type LogConfig = zlog.LogConfig

// 服务配置
type ServicesConfig struct {
	Api         ApiServiceConfig
	Grpc        GrpcServiceConfig
	Cron        CronServiceConfig
	MysqlBinlog MysqlBinlogServiceConfig
}

// api服务配置
type ApiServiceConfig struct {
	Bind                          string // bind地址
	IPWithNginxForwarded          bool   // 适配nginx的Forwarded获取ip, 优先级高于nginx的Real
	IPWithNginxReal               bool   // 适配nginx的Real获取ip, 优先级高于sock连接的ip
	LogResultInDevelop            bool   // 在开发环境记录api输出结果
	ShowDetailedErrorInProduction bool   // 在生产环境显示详细的错误
}

// grpc服务配置
type GrpcServiceConfig struct {
	Bind          string // bind地址
	HeartbeatTime int    // 心跳时间(毫秒),
}

// CronService配置
type CronServiceConfig struct {
	ThreadCount  int // 线程数
	JobQueueSize int // job队列大小
}

// MysqlBinlogService配置
type MysqlBinlogServiceConfig struct {
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

// 组件配置
type ComponentsConfig struct {
	GrpcClient map[string]GrpcClientConfig
	Xorm       map[string]XormConfig
	Redis      map[string]RedisConfig
	ES7        map[string]ES7Config
	Cache      map[string]CacheConfig
}

// grpc客户端配置
type GrpcClientConfig struct {
	Address     string // 链接地址
	Registry    string // 注册器, 默认为 local
	Balance     string // 负载均衡, 默认为 round_robin
	DialTimeout int    // 连接超时(毫秒), 0表示不限, 默认为 1000
}

// xorm配置
type XormConfig struct {
	Driver          string // 驱动
	Source          string // 连接源
	MaxIdleConns    int    // 最大空闲连接数
	MaxOpenConns    int    // 最大连接池个数
	ConnMaxLifetime int    // 最大续航时间(毫秒, 0表示无限
}

// redis配置
type RedisConfig struct {
	Address      string // 地址: host1:port1,host2:port2
	Password     string // 密码
	DB           int    // db, 只有单点有效
	IsCluster    bool   // 是否为集群
	PoolSize     int    // 客户端池大小
	ReadTimeout  int64  // 超时(毫秒
	WriteTimeout int64  // 超时(毫秒
	DialTimeout  int64  // 超时(毫秒
}

// es7配置
type ES7Config struct {
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

// 缓存配置
type CacheConfig struct {
	CacheDB                  string // 缓存db; default, no, memory, redis
	Codec                    string // 编解码器; default, byte, json, jsoniter, msgpack, proto_buffer
	DirectReturnOnCacheFault bool   // 在缓存故障时直接返回缓存错误
	PanicOnLoaderExists      bool   // 注册加载器时如果加载器已存在会panic, 设为false会替换旧的加载器
	SingleFlight             string // 单跑; default, no, single
	DefaultExpire            int64  // 默认有效时间, 毫秒, <= 0 表示永久
	DefaultExpireMax         int64  // 默认最大有效时间, 毫秒, 如果 > 0 且 DefaultExpire > 0, 则默认有效时间在 [DefaultExpire, DefaultExpireMax-1] 区间随机

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
