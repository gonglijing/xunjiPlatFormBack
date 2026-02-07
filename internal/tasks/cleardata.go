package tasks

import (
	"context"
	"fmt"
	"sagooiot/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

func (t TaskJob) ClearJobLogByDays(days string) {
	ctx := context.Background()
	glog.Debugf(ctx,"执行任务：清理超过%v天的定时任务日志", days)
	startTime := gtime.Now()
	err := service.SysJobLog().ClearJobLogByDays(ctx,gconv.Int(days))
	if err != nil {
		g.Log().Error(ctx, err)
	}
	if err := t.SaveLog(ctx, startTime, fmt.Sprintf("清理超过%v天的定时任务日志", days), err); err != nil {
		g.Log().Error(ctx, err)
	}
}

// ClearOperationLogByDays 清理超过指定天数的操作日志
func (t TaskJob) ClearOperationLogByDays(days string) {
	ctx := context.Background()
	glog.Debugf(ctx, "执行任务：清理超过%v天的操作日志", days)
	startTime := gtime.Now()
	err := service.SysOperLog().ClearOperationLogByDays(ctx, gconv.Int(days))
	if err != nil {
		g.Log().Error(ctx, err)
	}
	if err := t.SaveLog(ctx, startTime, fmt.Sprintf("清理超过%v天的操作日志", days), err); err != nil {
		g.Log().Error(ctx, err)
	}
}

// ClearNoticeLogByDays 清理超过指定天数的通知服务日志
func (t TaskJob) ClearNoticeLogByDays(days string) {
	ctx := context.Background()
	glog.Debugf(ctx, "执行任务：清理超过%d天的通知服务发送日志", gconv.Int(days))
	startTime := gtime.Now()
	err := service.NoticeLog().ClearLogByDays(ctx, gconv.Int(days))
	if err != nil {
		glog.Error(ctx, err)
	}
	if err := t.SaveLog(ctx, startTime, fmt.Sprintf("清理超过%d天的通知服务发送日志", gconv.Int(days)), err); err != nil {
		g.Log().Error(ctx, err)
	}
}

// ClearAlarmLogByDays 清理超过指定天数的告警日志
func (t TaskJob) ClearAlarmLogByDays(days string) {
	ctx := context.Background()
	glog.Debugf(ctx, "执行任务：清理超过%d天的告警日志", gconv.Int(days))
	startTime := gtime.Now()
	err := service.AlarmLog().ClearLogByDays(ctx, gconv.Int(days))
	if err != nil {
		glog.Error(ctx, err)
	}
	if err := t.SaveLog(ctx, startTime, fmt.Sprintf("清理超过%d天的告警日志", gconv.Int(days)), err); err != nil {
		g.Log().Error(ctx, err)
	}
}

// ClearTDengineLogByDays 清理超过指定天数的TD日志
func (t TaskJob) ClearTDengineLogByDays(days string) {
	ctx := context.Background()
	glog.Debugf(ctx, "执行任务：清理超过%d天的TD日志", gconv.Int(days))
	startTime := gtime.Now()
	err := service.TdEngine().ClearLogByDays(ctx, gconv.Int(days))
	if err != nil {
		glog.Error(ctx, err)
	}
	if err := t.SaveLog(ctx, startTime, fmt.Sprintf("清理超过%d天的TD日志", gconv.Int(days)), err); err != nil {
		g.Log().Error(ctx, err)
	}
}
