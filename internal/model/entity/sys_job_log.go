// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysJobLog is the golang structure for table sys_job_log.
type SysJobLog struct {
	Id             int64       `json:"id"             orm:"id"              ` // 主键ID
	JobName        string      `json:"jobName"        orm:"job_name"        ` // 任务名称
	InvokeTarget   string      `json:"invokeTarget"   orm:"invoke_target"   ` // 调用目标字符串
	CronExpression string      `json:"cronExpression" orm:"cron_expression" ` // cron执行表达式
	StartTime      *gtime.Time `json:"startTime"      orm:"start_time"      ` // 执行开始时间
	EndTime        *gtime.Time `json:"endTime"        orm:"end_time"        ` // 执行结束时间
	JobMessage     string      `json:"jobMessage"     orm:"job_message"     ` // 执行结果信息
	Status         int         `json:"status"         orm:"status"          ` // 状态（0正常 1失败）
	ExceptionInfo  string      `json:"exceptionInfo"  orm:"exception_info"  ` // 失败原因（异常信息）
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` // 创建时间
	IsDeleted      int         `json:"isDeleted"      orm:"is_deleted"      ` // 删除标志（0代表存在 1代表删除）
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"      ` // 删除时间
	DeletedBy      int64       `json:"deletedBy"      orm:"deleted_by"      ` // 删除人id
}
