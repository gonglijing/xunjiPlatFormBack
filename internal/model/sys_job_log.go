package model

import "github.com/gogf/gf/v2/os/gtime"

type GetJobLogListInput struct {
	JobName string `json:"jobName" description:"任务名称"`
	Status  string `json:"status" description:"状态（0正常 1暂停）"`
	*PaginationInput
}

type SysJobLogAddInput struct {
	JobName        string      `json:"jobName"  description:"任务名称"`
	InvokeTarget   string      `json:"invokeTarget" description:"调用目标字符串"`
	CronExpression string      `json:"cronExpression" description:"cron执行表达式"`
	StartTime      *gtime.Time `json:"startTime" description:"执行开始时间"`
	EndTime        *gtime.Time `json:"endTime" description:"执行结束时间"`
	JobMessage     string      `json:"jobMessage" description:"执行结果信息"`
	Status         int         `json:"status" description:"状态（0正常 1失败）"`
	ExceptionInfo  string      `json:"exceptionInfo" description:"失败原因（异常信息）"`
}

type SysJobLogListRes struct {
	ID             int64       `orm:"id,primary"        json:"id"`            // 主键ID
	JobName        string      `orm:"job_name"         json:"jobName"`        // 任务名称
	InvokeTarget   string      `orm:"invoke_target"    json:"invokeTarget"`   // 调用目标字符串
	CronExpression string      `orm:"cron_expression"  json:"cronExpression"` // cron执行表达式
	StartTime      *gtime.Time `orm:"start_time"       json:"startTime"`      // 执行开始时间
	EndTime        *gtime.Time `orm:"end_time"         json:"endTime"`        // 执行结束时间
	CreatedAt      *gtime.Time `orm:"created_at"       json:"createdAt"`      // 创建时间
	Status         int         `orm:"status"           json:"status"`         // 状态（0正常 1失败）
}

type SysJobLogRes struct {
	ID             int64       `orm:"id,primary"        json:"id"`            // 主键ID
	JobName        string      `orm:"job_name"         json:"jobName"`        // 任务名称
	InvokeTarget   string      `orm:"invoke_target"    json:"invokeTarget"`   // 调用目标字符串
	CronExpression string      `orm:"cron_expression"  json:"cronExpression"` // cron执行表达式
	StartTime      *gtime.Time `orm:"start_time"       json:"startTime"`      // 执行开始时间
	EndTime        *gtime.Time `orm:"end_time"         json:"endTime"`        // 执行结束时间
	JobMessage     string      `orm:"job_message"      json:"jobMessage"`     // 执行结果信息
	Status         int         `orm:"status"           json:"status"`         // 状态（0正常 1失败）
	ExceptionInfo  string      `orm:"exception_info"   json:"exceptionInfo"`  // 失败原因（异常信息）
	CreatedAt      *gtime.Time `orm:"created_at"       json:"createdAt"`      // 创建时间
	DeletedAt      *gtime.Time `orm:"deleted_at"       json:"deletedAt"`      // 删除时间}}
}

type SysJobLogOut struct {
	ID             int64       `orm:"id,primary"        json:"id"`            // 主键ID
	JobName        string      `orm:"job_name"         json:"jobName"`        // 任务名称
	InvokeTarget   string      `orm:"invoke_target"    json:"invokeTarget"`   // 调用目标字符串
	CronExpression string      `orm:"cron_expression"  json:"cronExpression"` // cron执行表达式
	StartTime      *gtime.Time `orm:"start_time"       json:"startTime"`      // 执行开始时间
	EndTime        *gtime.Time `orm:"end_time"         json:"endTime"`        // 执行结束时间
	JobMessage     string      `orm:"job_message"      json:"jobMessage"`     // 执行结果信息
	Status         int         `orm:"status"           json:"status"`         // 状态（0正常 1失败）
	ExceptionInfo  string      `orm:"exception_info"   json:"exceptionInfo"`  // 失败原因（异常信息）
	CreatedAt      *gtime.Time `orm:"created_at"       json:"createdAt"`      // 创建时间
	DeletedAt      *gtime.Time `orm:"deleted_at"       json:"deletedAt"`      // 删除时间
}

type SysJobLogExportOut struct {
	ID             int64       `orm:"id,primary"        json:"id"             description:"ID"`
	JobName        string      `orm:"job_name"         json:"jobName"        description:"任务名称"`
	InvokeTarget   string      `orm:"invoke_target"    json:"invokeTarget"   description:"功能名称"`
	CronExpression string      `orm:"cron_expression"  json:"cronExpression" description:"表达式"`
	StartTime      *gtime.Time `orm:"start_time"       json:"startTime"      description:"开始时间"`
	EndTime        *gtime.Time `orm:"end_time"         json:"endTime"        description:"结束时间"`
	JobMessage     string      `orm:"job_message"      json:"jobMessage"     description:"执行结果"`
	Status         string      `orm:"status"           json:"status"         description:"状态"`
	ExceptionInfo  string      `orm:"exception_info"   json:"exceptionInfo"  description:"失败原因"`
}
