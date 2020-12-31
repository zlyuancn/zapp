/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/17
   Description :
-------------------------------------------------
*/

package app

import (
	"sync/atomic"

	"github.com/zlyuancn/zlog"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/core"
)

// 获取下个日志id
//
// 将数值转为32进制, 因为求余2的次幂可以用位运算所以采用 数字+22位英文字母
func (app *appCli) nextLoggerId() string {
	id := atomic.AddUint32(&app.loggerId, 1)
	bs := []byte{48, 48, 48, 48, 48, 48, 48}

	i := len(bs) - 1
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

func (app *appCli) GetLogger() core.ILogger {
	return app.ILogger
}

func (app *appCli) CreateLogger(tag ...string) core.ILogger {
	log, _ := zlog.WrapZapFieldsWithLoger(app.ILogger, zap.String("logId", app.nextLoggerId()), zap.Strings("logTag", tag))
	return log
}
