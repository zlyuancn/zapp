/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/29
   Description :
-------------------------------------------------
*/

package cron

import (
	"github.com/zlyuancn/zscheduler"

	"github.com/zlyuancn/zapp/core"
)

type Job struct {
	core.ILogger
	zscheduler.IJob
}

func NewCronJob(job zscheduler.IJob) core.ICronJob {
	return &Job{
		ILogger: MustGetLoggerFromJob(job),
		IJob:    job,
	}
}
