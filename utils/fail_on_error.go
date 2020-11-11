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
)

func Panic(a ...interface{}) {
	log.Panic(a...)
}

func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func FailOnErrorf(err error, format string, msg ...interface{}) {
	if err != nil {
		log.Panicf("%s: %s", fmt.Sprintf(format, msg...), err)
	}
}
