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
	"github.com/zlyuancn/zutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/zlyuancn/zapp/core"
)

const logIdKey = "logId"

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
	log, _ := zlog.WrapZapFieldsWithLoger(app.ILogger, zap.String(logIdKey, app.nextLoggerId()), zap.Strings("logTag", tag))
	return log
}

func (app *appCli) withColoursMessageOfLoggerId() zap.Option {
	isTerminal := app.config.Config().Frame.Log.IsTerminal
	return zlog.WithHook(func(ent *zapcore.Entry, fields []zapcore.Field) (cancel bool) {
		if !isTerminal || ent.Message == "" {
			return
		}

		for _, field := range fields {
			if field.Key == logIdKey {
				ent.Message = app.makeColorMessageOfLoggerId(field.String, ent.Message)
				break
			}
		}
		return
	})
}

func (app *appCli) makeColorMessageOfLoggerId(logId string, message string) string {
	var id uint32
	for _, c := range logId {
		id <<= 5
		if c >= 'a' {
			id += uint32(c) - 87
		} else {
			id += uint32(c) - 48
		}
	}

	color := zutils.ColorType(id%7) + zutils.ColorRed
	return zutils.Color.MakeColorText(color, message)
}
