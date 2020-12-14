/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/7/21
   Description :
-------------------------------------------------
*/

package consts

// 默认服务名
const DefaultServiceName = "default"

// 默认组件名
const DefaultComponentName = "default"

// 默认配置
const (
	// 配置文件
	DefaultConfig_ConfigFiles = "./configs/default.toml"
	// 清理内存间隔时间(毫秒)
	DefaultConfig_App_FreeMemoryInterval int = 120000

	// 等待服务启动阶段1, 等待时间(毫秒), 如果时间到则临时认为服务启动成功并提前返回
	DefaultConfig_WaitServiceRunTime int = 1000
	// 等待服务启动阶段2, 等待服务启动阶段1时间到后继续等待服务启动, 等待时间(毫秒), 如果时间到则真正认为服务启动成功.
	DefaultConfig_ContinueWaitServiceRunTime int = 30000

	// grpc默认scheme
	DefaultConfig_GrpcClient_Registry = "local"
	// grpc默认balance
	DefaultConfig_GrpcClient_Balance = "round_robin"
	// grpc客户端默认连接超时
	DefaultConfig_GrpcClient_DialTimeout = 1000
)

// 配置文件分片名
const (
	// 日志
	ConfigShardName_Log = "log"
)

// 用户储存上下文的字段名
const (
	// log
	SaveFieldName_Logger = "_logger"
)
