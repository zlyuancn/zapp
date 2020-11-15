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
		bs[i] = byte(31&id) + 48 // 从字符0开始
		if bs[i] > 57 {          // 超过数字用字母表示
			bs[i] += 39
		}
		if id < 32 {
			return string(bs)
		}
		i--
		id >>= 5
	}
}

type contextCli struct {
	app    core.IApp
	config *core.Config
	core.ILogger
}

func NewContext(app core.IApp, tag ...string) core.IContext {
	log, _ := zlog.WrapZapFields(app.GetLogger(), zap.String("ctxId", nextContextId()), zap.Strings("ctxTag", tag))

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
