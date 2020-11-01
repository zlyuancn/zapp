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

func FailOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%s: %s", msg, err))
	}
}

func FailOnErrorf(err error, format string, msg ...interface{}) {
	if err != nil {
		panic(fmt.Errorf("%s: %s", fmt.Sprintf(format, msg...), err))
	}
}
