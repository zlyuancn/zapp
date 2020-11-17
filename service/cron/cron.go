/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/1
   Description :
-------------------------------------------------
*/

package cron

import (
	"fmt"

	"github.com/zlyuancn/zscheduler"
	"go.uber.org/zap"

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

type Job = zscheduler.Job

type CronService struct {
	app       core.IApp
	scheduler zscheduler.IScheduler
}

func NewCronService(app core.IApp) *CronService {
	return &CronService{
		app: app,
		scheduler: zscheduler.NewScheduler(
			zscheduler.WithLogger(app.GetLogger()),
			zscheduler.WithGoroutinePool(app.GetConfig().Config().Cron.ThreadCount, app.GetConfig().Config().Cron.JobQueueSize),
		),
	}
}

func (c *CronService) Inject(a ...interface{}) {
	for _, v := range a {
		task, ok := v.(zscheduler.ITask)
		if !ok {
			c.app.GetLogger().Fatal("Cron服务注入类型错误", zap.String("type", fmt.Sprintf("%T", v)))
		}

		if ok := c.scheduler.AddTask(task); !ok {
			c.app.GetLogger().Fatal("添加Cron任务失败, 可能是名称重复", zap.String("name", task.Name()))
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

func NewTask(name string, expression string, job Job, enable ...bool) zscheduler.ITask {
	return zscheduler.NewTaskOfConfig(name, zscheduler.TaskConfig{
		Trigger:  zscheduler.NewCronTrigger(expression),
		Executor: zscheduler.NewExecutor(0, 0, 1),
		Job:      job,
		Enable:   len(enable) == 0 || enable[0],
	})
}
