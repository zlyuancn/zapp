/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/11
   Description :
-------------------------------------------------
*/

package utils

import (
	"github.com/zlyuancn/zlog"

	"github.com/zlyuancn/zapp/core"
)

var log = zlog.DefaultLogger

func SetLog(l core.ILogger) {
	log = l
}
