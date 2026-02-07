// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysJobLog is the golang structure of table sys_job_log for DAO operations like Where/Data.
type SysJobLog struct {
	g.Meta         `orm:"table:sys_job_log, do:true"`
	Id             interface{} // 主键ID
	JobName        interface{} // 任务名称
	InvokeTarget   interface{} // 调用目标字符串
	CronExpression interface{} // cron执行表达式
	StartTime      *gtime.Time // 执行开始时间
	EndTime        *gtime.Time // 执行结束时间
	JobMessage     interface{} // 执行结果信息
	Status         interface{} // 状态（0正常 1失败）
	ExceptionInfo  interface{} // 失败原因（异常信息）
	CreatedAt      *gtime.Time // 创建时间
	IsDeleted      interface{} // 删除标志（0代表存在 1代表删除）
	DeletedAt      *gtime.Time // 删除时间
	DeletedBy      interface{} // 删除人id
}
