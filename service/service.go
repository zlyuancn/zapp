/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/21
   Description :
-------------------------------------------------
*/

package service

import (
	"fmt"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
)

var creators = make(map[consts.ServiceType]core.IServiceCreator)

// 注册建造者
func RegisterCreator(serviceType consts.ServiceType, creator core.IServiceCreator) {
	if _, ok := creators[serviceType]; ok {
		panic(fmt.Errorf("重复注册建造者: %s" + serviceType.String()))
	}
	creators[serviceType] = creator
}

// 创建服务,
func NewService(serviceType consts.ServiceType, c core.IComponent) core.IService {
	if creator, ok := creators[serviceType]; ok {
		return creator.Create(c)
	}
	panic(fmt.Errorf("使用了未注册的建造者: %s", serviceType.String()))
}
