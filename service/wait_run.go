/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/12/14
   Description :
-------------------------------------------------
*/

package service

import (
	"time"

	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/component"
	"github.com/zlyuancn/zapp/logger"
)

// 等待运行选项
type WaitRunOption struct {
	// 服务名
	ServiceName string
	// 如果错误是这些值则忽略
	IgnoreErrs []error
	// 如果等待阶段2返回错误是否在打印错误后退出
	ExitOnErrOfWait2 bool
	// 启动服务函数
	RunServiceFn func() error
}

func WaitRun(opt *WaitRunOption) error {
	if opt.ServiceName == "" {
		logger.Log.Fatal("serverName must not empty")
	}

	errChan := make(chan error, 1)
	go func(errChan chan error) {
		errChan <- opt.RunServiceFn()
	}(errChan)

	wait := time.NewTimer(time.Duration(component.GlobalComponent().Config().Frame.WaitServiceRunTime) * time.Millisecond) // 等待启动提前返回
	select {
	case <-wait.C:
	case <-component.GlobalComponent().App().BaseContext().Done():
		wait.Stop()
		return nil
	case err := <-errChan:
		wait.Stop()
		for _, e := range opt.IgnoreErrs {
			if err == e {
				return nil
			}
		}
		return err
	}

	// 开始等待服务启动阶段2
	go func(errChan chan error) {
		wait = time.NewTimer(time.Duration(component.GlobalComponent().Config().Frame.ContinueWaitServiceRunTime) * time.Millisecond)
		select {
		case <-wait.C:
		case <-component.GlobalComponent().App().BaseContext().Done():
			wait.Stop()
		case err := <-errChan:
			wait.Stop()
			if err == nil {
				return
			}
			for _, e := range opt.IgnoreErrs {
				if err == e {
					return
				}
			}

			if opt.ExitOnErrOfWait2 {
				logger.Log.Fatal("服务启动失败", zap.String("serviceName", opt.ServiceName), zap.Error(err))
			} else {
				logger.Log.Error("服务启动失败", zap.String("serviceName", opt.ServiceName), zap.Error(err))
			}
		}
	}(errChan)

	return nil
}
