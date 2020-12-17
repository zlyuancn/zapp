/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/12/17
   Description :
-------------------------------------------------
*/

package config

type Namespace string

const (
	FrameNamespace              Namespace = "frame"
	LogNamespace                          = "log"
	ApiServiceNamespace                   = "api_service"
	GrpcServiceNamespace                  = "grpc_service"
	CronServiceNamespace                  = "cron_service"
	MysqlBinlogServiceNamespace           = "mysql_binlog_service"
	GrpcClientNamespace                   = "grpc_client"
	XormNamespace                         = "xorm"
	RedisNamespace                        = "redis"
	ES7Namespace                          = "es7"
)
