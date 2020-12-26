/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/1
   Description :
-------------------------------------------------
*/

package cron

import (
	"go.uber.org/zap"

	"github.com/zlyuancn/zscheduler"

	"github.com/zlyuancn/zapp/consts"

	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/service"
)

func init() {
	service.RegisterCreator(core.CronService, new(cronCreator))
}

type cronCreator struct{}

func (*cronCreator) Create(app core.IApp) core.IService {
	return NewCronService(app)
}

type CronService struct {
	app       core.IApp
	scheduler zscheduler.IScheduler
}

func NewCronService(app core.IApp) *CronService {
	conf := app.GetConfig().Config().Services.Cron
	return &CronService{
		app: app,
		scheduler: zscheduler.NewScheduler(
			zscheduler.WithLogger(nil),
			zscheduler.WithGoroutinePool(conf.ThreadCount, conf.JobQueueSize),
			zscheduler.WithObserver(newObserver(app)),
		),
	}
}

func (c *CronService) Inject(a ...interface{}) {
	for _, v := range a {
		task, ok := v.(zscheduler.ITask)
		if !ok {
			c.app.Fatal("Cron服务注入类型错误, 它必须能转为 zscheduler.ITask")
		}

		if ok := c.scheduler.AddTask(task); !ok {
			c.app.Fatal("添加Cron任务失败, 可能是名称重复", zap.String("name", task.Name()))
		}
	}
}

func (c *CronService) Start() error {
	c.scheduler.Start()
	return nil
}

func (c *CronService) Close() error {
	c.scheduler.Stop()
	return nil
}

// 将log存入job, 如果meta不是nil或map[string]interface{}会panic
func SaveLoggerToJob(job zscheduler.IJob, log core.ILogger) {
	if job.Meta() == nil {
		job.SetMeta(map[string]interface{}{
			consts.SaveFieldName_Logger: log,
		})
	}
	job.Meta().(map[string]interface{})[consts.SaveFieldName_Logger] = log
}

// 从job中获取log, 如果失败会panic
func MustGetLoggerFromJob(job zscheduler.IJob) core.ILogger {
	return job.Meta().(map[string]interface{})[consts.SaveFieldName_Logger].(core.ILogger)
}
