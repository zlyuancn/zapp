/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package app

import (
	"github.com/zlyuancn/zapp/core"
)

// 关闭组件内加载的资源
func (app *appCli) closeComponentResource() {
	app.Debug("释放组件加载的资源")
	app.component.Close()
}

func (app *appCli) GetComponent() core.IComponent {
	return app.component
}
