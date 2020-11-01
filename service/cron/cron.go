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

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/service"
)

func init() {
	service.RegisterCreator(consts.CronService, new(cronCreator))
}

type cronCreator struct{}

func (*cronCreator) Create(c core.IComponent) core.IService {
	return NewCronService(c)
}

type Job = zscheduler.Job

type CronService struct {
	c         core.IComponent
	scheduler zscheduler.IScheduler
}

func NewCronService(c core.IComponent) *CronService {
	return &CronService{
		c: c,
		scheduler: zscheduler.NewScheduler(
			zscheduler.WithLogger(c.App().GetLogger()),
			zscheduler.WithGoroutinePool(c.Config().Cron.ThreadCount, c.Config().Cron.JobQueueSize),
		),
	}
}

func (c *CronService) Inject(a ...interface{}) {
	for _, v := range a {
		task, ok := v.(zscheduler.ITask)
		if !ok {
			c.c.Fatal("Cron服务注入类型错误", zap.String("type", fmt.Sprintf("%T", v)))
		}

		if ok := c.scheduler.AddTask(task); !ok {
			c.c.Fatal("添加Cron任务失败, 可能是名称重复", zap.String("name", task.Name()))
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
