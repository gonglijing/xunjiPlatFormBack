// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysDeptDao is the data access object for the table sys_dept.
type SysDeptDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysDeptColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysDeptColumns defines and stores column names for the table sys_dept.
type SysDeptColumns struct {
	DeptId         string // 部门id
	OrganizationId string // 组织ID
	ParentId       string // 父部门id
	Ancestors      string // 祖级列表
	DeptName       string // 部门名称
	OrderNum       string // 显示顺序
	Leader         string // 负责人
	Phone          string // 联系电话
	Email          string // 邮箱
	Status         string // 部门状态（0停用 1正常）
	IsDeleted      string // 是否删除 0未删除 1已删除
	CreatedAt      string // 创建时间
	CreatedBy      string // 创建人
	UpdatedBy      string // 修改人
	UpdatedAt      string // 修改时间
	DeletedBy      string // 删除人
	DeletedAt      string // 删除时间
	DeptCode       string // 部门编码
}

// sysDeptColumns holds the columns for the table sys_dept.
var sysDeptColumns = SysDeptColumns{
	DeptId:         "dept_id",
	OrganizationId: "organization_id",
	ParentId:       "parent_id",
	Ancestors:      "ancestors",
	DeptName:       "dept_name",
	OrderNum:       "order_num",
	Leader:         "leader",
	Phone:          "phone",
	Email:          "email",
	Status:         "status",
	IsDeleted:      "is_deleted",
	CreatedAt:      "created_at",
	CreatedBy:      "created_by",
	UpdatedBy:      "updated_by",
	UpdatedAt:      "updated_at",
	DeletedBy:      "deleted_by",
	DeletedAt:      "deleted_at",
	DeptCode:       "dept_code",
}

// NewSysDeptDao creates and returns a new DAO object for table data access.
func NewSysDeptDao(handlers ...gdb.ModelHandler) *SysDeptDao {
	return &SysDeptDao{
		group:    "default",
		table:    "sys_dept",
		columns:  sysDeptColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysDeptDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysDeptDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysDeptDao) Columns() SysDeptColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysDeptDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysDeptDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysDeptDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
