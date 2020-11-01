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
}
