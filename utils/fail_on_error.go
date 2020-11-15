/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package utils

import (
	"fmt"

	"github.com/zlyuancn/zapp/logger"
)

func Panic(a ...interface{}) {
	logger.Log.Panic(a...)
}

func Panicf(format string, v ...interface{}) {
	logger.Log.Panicf(format, v...)
}

func FailOnError(err error, msg string) {
	if err != nil {
		logger.Log.Panicf("%s: %s", msg, err)
	}
}

func FailOnErrorf(err error, format string, msg ...interface{}) {
	if err != nil {
		logger.Log.Panicf("%s: %s", fmt.Sprintf(format, msg...), err)
	}
}
