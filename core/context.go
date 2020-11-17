/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/1
   Description :
-------------------------------------------------
*/

package core

// 上下文
type IContext interface {
	// 获取app
	App() IApp
	// 获取配置
	Config() *Config
	ILogger

	// 保存数据
	Set(k string, v interface{})
	// 加载数据
	Get(k string) (interface{}, bool)
	// 设置元数据
	SetMetadata(data interface{})
	// 获取元数据
	GetMetadata() interface{}
}
