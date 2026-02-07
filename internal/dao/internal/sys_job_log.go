// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysJobLogDao is the data access object for the table sys_job_log.
type SysJobLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysJobLogColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysJobLogColumns defines and stores column names for the table sys_job_log.
type SysJobLogColumns struct {
	Id             string // 主键ID
	JobName        string // 任务名称
	InvokeTarget   string // 调用目标字符串
	CronExpression string // cron执行表达式
	StartTime      string // 执行开始时间
	EndTime        string // 执行结束时间
	JobMessage     string // 执行结果信息
	Status         string // 状态（0正常 1失败）
	ExceptionInfo  string // 失败原因（异常信息）
	CreatedAt      string // 创建时间
	IsDeleted      string // 删除标志（0代表存在 1代表删除）
	DeletedAt      string // 删除时间
	DeletedBy      string // 删除人id
}

// sysJobLogColumns holds the columns for the table sys_job_log.
var sysJobLogColumns = SysJobLogColumns{
	Id:             "id",
	JobName:        "job_name",
	InvokeTarget:   "invoke_target",
	CronExpression: "cron_expression",
	StartTime:      "start_time",
	EndTime:        "end_time",
	JobMessage:     "job_message",
	Status:         "status",
	ExceptionInfo:  "exception_info",
	CreatedAt:      "created_at",
	IsDeleted:      "is_deleted",
	DeletedAt:      "deleted_at",
	DeletedBy:      "deleted_by",
}

// NewSysJobLogDao creates and returns a new DAO object for table data access.
func NewSysJobLogDao(handlers ...gdb.ModelHandler) *SysJobLogDao {
	return &SysJobLogDao{
		group:    "default",
		table:    "sys_job_log",
		columns:  sysJobLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysJobLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysJobLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysJobLogDao) Columns() SysJobLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysJobLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysJobLogDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SysJobLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
