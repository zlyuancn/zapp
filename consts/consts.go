/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/21
   Description :
-------------------------------------------------
*/

package consts

import (
	"fmt"
)

const (
	// 日志配置文件分片名
	LogConfigShardName = "log"
)

// 服务类型
type ServiceType uint8

const (
	// http服务
	HttpService ServiceType = iota
	// grpc服务
	GrpcService
	// cron服务
	CronService
)

func (t ServiceType) String() string {
	switch t {
	case HttpService:
		return "http"
	case GrpcService:
		return "grpc"
	case CronService:
		return "cron"
	}

	return fmt.Sprintf("ServiceType<%d>", t)
}

// 默认服务名
const DefaultServiceName = "default"
