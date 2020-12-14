/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/7/21
   Description :
-------------------------------------------------
*/

package service

import (
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
)

var creators = make(map[core.ServiceType]core.IServiceCreator)

// 注册建造者
func RegisterCreator(serviceType core.ServiceType, creator core.IServiceCreator) {
	if _, ok := creators[serviceType]; ok {
		logger.Log.Fatal("重复注册建造者", zap.String("serviceType", serviceType.String()))
	}
	creators[serviceType] = creator
}

// 创建服务,
func NewService(serviceType core.ServiceType, app core.IApp) core.IService {
	if creator, ok := creators[serviceType]; ok {
		return creator.Create(app)
	}
	logger.Log.Fatal("使用了未注册的建造者", zap.String("serviceType", serviceType.String()))
	return nil
}

// 注册服务类型描述
func RegistryServiceTypeDescribe(t core.ServiceType, desc string) {
	core.RegistryServiceTypeDescribe(t, desc)
}
