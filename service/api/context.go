/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/30
   Description :
-------------------------------------------------
*/

package api

import (
	"reflect"

	"github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/component"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/utils"
)

type Context struct {
	core.ILogger
	*iris_context.Context
}

func makeContext(ctx iris.Context) *Context {
	return &Context{
		ILogger: utils.Context.MustGetLoggerFromIrisContext(ctx),
		Context: ctx,
	}
}

//  bind api数据, 它会将api数据反序列化到a中, 如果a是结构体会验证a
func (c *Context) Bind(a interface{}) error {
	if err := c.ReadBody(a); err != nil {
		return ParamError.WithError(err)
	}

	c.Debug("api.request.arg", zap.Any("arg", a))

	val := reflect.ValueOf(a)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	err := component.GlobalComponent().Valid(a)
	if err != nil {
		return ParamError.WithError(err)
	}
	return nil
}
