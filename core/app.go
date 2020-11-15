/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package core

// app
//
// 用于将所有模块连起来
type IApp interface {
	// 运行
	//
	// 开启所有服务并挂起
	Run()
	// 退出
	//
	// 结束所有服务并退出
	Exit()

	// 获取配置
	GetConfig() IConfig
	// 获取日志组件
	GetLogger() ILogger
	// 获取组件
	GetComponent() IComponent
	// 创建上下文,
	NewContext(tag ...string) IContext
	// 获取服务
	GetService(serviceType ServiceType, serviceName ...string) (IService, bool)
}
