package tasks

import (
	"context"
	"sagooiot/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceLogClear 设备日志清理
func (t TaskJob) DeviceLogClear() {
	ctx := context.Background()
	startTime := gtime.Now()
	err := service.TdLogTable().Clear(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	if err := t.SaveLog(ctx, startTime,"设备日志清理成功",err); err != nil {
		g.Log().Error(ctx, err)
	}
	glog.Debug(ctx, "执行任务：清理设备日志")
}
