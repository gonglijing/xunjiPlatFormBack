package system

import (
	"sagooiot/api/v1/common"
	"sagooiot/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/frame/g"
)

type GetJobReq struct {
	g.Meta `path:"/job_log/get" method:"get" summary:"获取任务日志详情" tags:"任务日志"`
	Id     int `json:"id" description:"任务日志ID" v:"required#任务日志ID不能为空"`
}

type SysJobLogRes struct {
	ID             int64       `json:"id"`            // 主键ID
	JobName        string      `json:"jobName"`        // 任务名称
	InvokeTarget   string      `json:"invokeTarget"`   // 调用目标字符串
	CronExpression string      `json:"cronExpression"` // cron执行表达式
	StartTime      *gtime.Time `json:"startTime"`      // 执行开始时间
	EndTime        *gtime.Time `json:"endTime"`        // 执行结束时间
	JobMessage     string      `json:"jobMessage"`     // 执行结果信息
	Status         int         `json:"status"`         // 状态（0正常 1失败）
	ExceptionInfo  string      `json:"exceptionInfo"`  // 失败原因（异常信息）
	CreatedAt      *gtime.Time `json:"createdAt"`      // 创建时间
	DeletedAt      *gtime.Time `json:"deletedAt"`      // 删除时间
}


type GetJobLogListReq struct {
	g.Meta  `path:"/job_log/list" method:"get" summary:"获取任务日志列表" tags:"任务日志"`
	JobName string `json:"jobName" description:"任务名称"`
	Status  string `json:"status" description:"状态（0正常 1暂停）"`
	*common.PaginationReq
}

type GetJobLogListRes struct {
	Data []*model.SysJobLogListRes
	common.PaginationRes
}

type DeleteJobLogReq struct {
	g.Meta `path:"/job_log/delete" method:"delete" summary:"删除任务日志" tags:"任务日志"`
	Ids    []int `json:"ids" description:"任务日志ID列表" v:"required#任务日志ID不能为空"`
}

type DeleteJobLogRes struct {
}

type JobExportReq struct {
	g.Meta  `path:"/job_log/export" method:"get" summary:"导出任务日志" tags:"任务日志"`
	JobName string `json:"jobName" description:"任务名称"`
	Status  string `json:"status" description:"状态（0正常 1暂停）"`
	*common.PaginationReq
}

type JobLogExportRes struct {
	g.Meta `mime:"text/html" example:"string"`
}
