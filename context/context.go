/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/1
   Description :
-------------------------------------------------
*/

package context

import (
	"sync/atomic"

	"github.com/zlyuancn/zlog"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/core"
)

var contextId uint64

// 获取下个上下文id
//
// 将数值转为32进制, 因为求余2的次幂可以用位运算所以采用 数字+22位英文字母
func nextContextId() string {
	id := atomic.AddUint64(&contextId, 1)
	bs := []byte{48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48} // 补0位, 32^13等于2^65
	i := 12
	for {
		v := byte(31 & id)
		if v < 10 {
			bs[i] = v + 48
		} else {
			bs[i] = v + 87
		}
		if id < 32 {
			return string(bs)
		}
		i--
		id = id >> 5
	}
}

type contextCli struct {
	app    core.IApp
	config *core.Config
	core.ILogger
}

func NewContext(app core.IApp, tag ...string) core.IContext {
	log, _ := zlog.WrapZapFields(app.GetLogger(), zap.String("ctx_id", nextContextId()), zap.Strings("ctx_tag", tag))

	return &contextCli{
		app:     app,
		config:  app.GetConfig().Config(),
		ILogger: log,
	}
}

func (c *contextCli) App() core.IApp {
	return c.app
}

func (c *contextCli) Config() *core.Config {
	return c.config
}
