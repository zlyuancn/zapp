/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package utils

import (
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/logger"
)

func Panic(a ...interface{}) {
	logger.Log.Panic(a...)
}

func PanicOnError(err error, a ...interface{}) {
	if err != nil {
		logger.Log.Panic(append(append([]interface{}{}, a...), zap.Error(err)))
	}
}

func Fatal(a ...interface{}) {
	logger.Log.Fatal(a...)
}
func FatalOnError(err error, a ...interface{}) {
	if err != nil {
		logger.Log.Fatal(append(append([]interface{}{}, a...), zap.Error(err)))
	}
}
