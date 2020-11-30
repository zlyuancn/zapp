/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/28
   Description :
-------------------------------------------------
*/

package cron

import (
	"go.uber.org/zap"

	"github.com/zlyuancn/zscheduler"

	"github.com/zlyuancn/zapp/core"
)

type Observer struct {
	app core.IApp
	zscheduler.Observer
}

func newObserver(app core.IApp) zscheduler.IObserver {
	return &Observer{app: app}
}

func (o *Observer) Started() { o.app.Debug("cron服务启动成功") }
func (o *Observer) Stopped() { o.app.Debug("cron服务已关闭") }
func (o *Observer) JobStart(job zscheduler.IJob) {
	log := o.app.CreateLogger(job.Task().Name())
	SaveLoggerToJob(job, log)
	log.Debug("cron.start")
}
func (o *Observer) JobErr(job zscheduler.IJob, err error) {
	log := MustGetLoggerFromJob(job)
	log.Error("cron.error! try retry", zap.Error(err))
}
func (o *Observer) JobEnd(job zscheduler.IJob, executeInfo *zscheduler.ExecuteInfo) {
	log := MustGetLoggerFromJob(job)
	if executeInfo.ExecuteSuccess {
		log.Debug("cron.success")
	} else {
		log.Error("cron.error!", zap.Error(executeInfo.ExecuteErr))
	}
}
