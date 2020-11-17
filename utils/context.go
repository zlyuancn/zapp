/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package utils

import (
	"context"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
)

var Context = new(contextUtil)

type contextUtil struct{}

// 基于传入的标准context生成一个新的标准context并保存app上下文
func (c *contextUtil) SaveAppContextToContext(ctx context.Context, appCtx core.IContext) context.Context {
	return context.WithValue(ctx, consts.SaveFieldName_AppContext, appCtx)
}

// 从标准context中加载app上下文
func (c *contextUtil) LoadAppContextFromContext(ctx context.Context) (core.IContext, bool) {
	value := ctx.Value(consts.SaveFieldName_AppContext)
	appCtx, ok := value.(core.IContext)
	return appCtx, ok
}

// 从标准context中加载app上下文, 如果失败会panic
func (c *contextUtil) MustLoadAppContextFromContext(ctx context.Context) core.IContext {
	out, ok := c.LoadAppContextFromContext(ctx)
	if !ok {
		Panic("can't load app_context from context")
	}
	return out
}
