package tasks

import (
	"context"
	"sagooiot/internal/model"
	"sagooiot/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

type TaskJob struct {
	ID             string        //任务ID
	TaskType       string        //任务类型
	MethodName     string        //方法名
	Params         []interface{} //参数
	Explain        string        //任务描述
	CronExpression string        //cron表达式
}

func (t TaskJob) SaveLog(ctx context.Context, StartTime *gtime.Time, res string, err error) error {
	var status int
	var exceptionInfo string
	if err != nil {
		status = 1 // 失败
		exceptionInfo = err.Error()
		res = "任务执行失败"
	} else {
		status = 0 // 成功
		exceptionInfo = ""
	}
	return service.SysJobLog().AddJobLog(ctx, &model.SysJobLogAddInput{
		JobName:        t.Explain,
		InvokeTarget:   t.MethodName,
		CronExpression: t.CronExpression,
		StartTime:      StartTime,
		EndTime:        gtime.Now(),
		JobMessage:     res,
		Status:         status,
		ExceptionInfo:  exceptionInfo,
	})
}

func (t TaskJob) GetFuncNameList() (res map[string]string) {
	res = map[string]string{
		"ClearJobLogByDays":       "清理超过指定天数的定时任务日志",
		"ClearOperationLogByDays": "清理超过指定天数的操作日志",
		"ClearNoticeLogByDays":    "清理超过指定天数的通知服务日志",
		"ClearAlarmLogByDays":     "清理超过指定天数的告警日志",
		"ClearTDengineLogByDays":  "清理超过指定天数的TD日志",
		"GetAccessURL":            "访问URL",
		"DataSourceSync":          "数据源同步",
		"DataTemplate":            "数据模型聚合数据",
		"DeviceLogClear":          "设备日志清理",
	}
	return
}
