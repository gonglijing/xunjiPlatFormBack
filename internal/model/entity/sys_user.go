// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysUser is the golang structure for table sys_user.
type SysUser struct {
	Id            uint64      `json:"id"            orm:"id"              description:""`
	UserName      string      `json:"userName"      orm:"user_name"       description:"用户名"`
	UserTypes     string      `json:"userTypes"     orm:"user_types"      description:"系统 system 企业 company"`
	Mobile        string      `json:"mobile"        orm:"mobile"          description:"中国手机不带国家代码，国际手机号格式为：国家代码-手机号"`
	UserNickname  string      `json:"userNickname"  orm:"user_nickname"   description:"用户昵称"`
	Birthday      *gtime.Time `json:"birthday"      orm:"birthday"        description:"生日"`
	UserPassword  string      `json:"userPassword"  orm:"user_password"   description:"登录密码;cmf_password加密"`
	UserSalt      string      `json:"userSalt"      orm:"user_salt"       description:"加密盐"`
	UserEmail     string      `json:"userEmail"     orm:"user_email"      description:"用户登录邮箱"`
	Sex           int         `json:"sex"           orm:"sex"             description:"性别;0:保密,1:男,2:女"`
	Avatar        string      `json:"avatar"        orm:"avatar"          description:"用户头像"`
	DeptId        uint64      `json:"deptId"        orm:"dept_id"         description:"部门id"`
	Remark        string      `json:"remark"        orm:"remark"          description:"备注"`
	IsAdmin       int         `json:"isAdmin"       orm:"is_admin"        description:"是否后台管理员 1 是  0   否"`
	Address       string      `json:"address"       orm:"address"         description:"联系地址"`
	Describe      string      `json:"describe"      orm:"describe"        description:"描述信息"`
	LastLoginIp   string      `json:"lastLoginIp"   orm:"last_login_ip"   description:"最后登录ip"`
	LastLoginTime *gtime.Time `json:"lastLoginTime" orm:"last_login_time" description:"最后登录时间"`
	Status        uint        `json:"status"        orm:"status"          description:"用户状态;0:禁用,1:正常,2:未验证"`
	IsDeleted     int         `json:"isDeleted"     orm:"is_deleted"      description:"是否删除 0未删除 1已删除"`
	CreatedBy     uint        `json:"createdBy"     orm:"created_by"      description:"创建者"`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:"创建日期"`
	UpdatedBy     uint        `json:"updatedBy"     orm:"updated_by"      description:"更新者"`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      description:"修改日期"`
	DeletedBy     int         `json:"deletedBy"     orm:"deleted_by"      description:"删除人"`
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"      description:"删除时间"`
	Code          string      `json:"code"          orm:"code"            description:"编码"`
}
