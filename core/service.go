/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/7/21
   Description :
-------------------------------------------------
*/

package core

import (
	"fmt"
)

// 服务
type IService interface {
	// 注入, 根据服务不同具有不同作用, 具体参考服务实现说明
	Inject(a ...interface{})
	// 开始服务
	Start() error
	// 关闭服务
	Close() error
}

// 服务建造者
type IServiceCreator interface {
	// 创建服务
	Create(app IApp) IService
}

// 服务类型
type ServiceType uint8

const (
	// api服务
	ApiService ServiceType = iota
	// grpc服务
	GrpcService
	// cron服务
	CronService
	// mysql-binlog服务
	MysqlBinlogService
)

func (t ServiceType) String() string {
	switch t {
	case ApiService:
		return "api"
	case GrpcService:
		return "grpc"
	case CronService:
		return "cron"
	case MysqlBinlogService:
		return "mysql-binlog"
	default:
		if desc, ok := serviceTypeDescribeMap[t]; ok {
			return desc
		}
	}

	return fmt.Sprintf("<undefined service type %d>", t)
}

var serviceTypeDescribeMap = make(map[ServiceType]string)

// 注册服务类型描述
func RegistryServiceTypeDescribe(t ServiceType, desc string) {
	serviceTypeDescribeMap[t] = desc
}
