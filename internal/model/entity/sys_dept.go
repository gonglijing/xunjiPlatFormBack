// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDept is the golang structure for table sys_dept.
type SysDept struct {
	DeptId         int64       `json:"deptId"         orm:"dept_id"         description:"部门id"`
	OrganizationId int         `json:"organizationId" orm:"organization_id" description:"组织ID"`
	ParentId       int64       `json:"parentId"       orm:"parent_id"       description:"父部门id"`
	Ancestors      string      `json:"ancestors"      orm:"ancestors"       description:"祖级列表"`
	DeptName       string      `json:"deptName"       orm:"dept_name"       description:"部门名称"`
	OrderNum       int         `json:"orderNum"       orm:"order_num"       description:"显示顺序"`
	Leader         string      `json:"leader"         orm:"leader"          description:"负责人"`
	Phone          string      `json:"phone"          orm:"phone"           description:"联系电话"`
	Email          string      `json:"email"          orm:"email"           description:"邮箱"`
	Status         uint        `json:"status"         orm:"status"          description:"部门状态（0停用 1正常）"`
	IsDeleted      int         `json:"isDeleted"      orm:"is_deleted"      description:"是否删除 0未删除 1已删除"`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`
	CreatedBy      uint        `json:"createdBy"      orm:"created_by"      description:"创建人"`
	UpdatedBy      int         `json:"updatedBy"      orm:"updated_by"      description:"修改人"`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:"修改时间"`
	DeletedBy      int         `json:"deletedBy"      orm:"deleted_by"      description:"删除人"`
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"      description:"删除时间"`
	DeptCode       string      `json:"deptCode"       orm:"dept_code"       description:"部门编码"`
}
