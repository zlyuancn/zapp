/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package app

import (
	"github.com/zlyuancn/zapp/component"
	"github.com/zlyuancn/zapp/core"
)

// 关闭组件内加载的资源
func (app *appCli) closeComponentResource() {
	app.Info("释放组件加载的资源")
	c, ok := app.component.(*component.ComponentCli)
	if ok {
		c.Close()
	}
}

func (app *appCli) GetComponent() core.IComponent {
	return app.component
}
