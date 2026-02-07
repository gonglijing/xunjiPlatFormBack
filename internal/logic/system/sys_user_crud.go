package system

import (
	"context"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/model/do"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"

	"sagooiot/pkg/utility/utils"
)

// Add 添加
func (s *sSysUser) Add(ctx context.Context, input *model.AddUserInput) (err error) {
	//判断账户是否存在
	num, _ := dao.SysUser.Ctx(ctx).Where(g.Map{
		dao.SysUser.Columns().UserName:  input.UserName,
		dao.SysUser.Columns().IsDeleted: 0,
	}).Count()
	if num > 0 {
		return gerror.New("账户已存在")
	}
	num, _ = dao.SysUser.Ctx(ctx).Where(g.Map{
		dao.SysUser.Columns().Mobile:    input.Mobile,
		dao.SysUser.Columns().IsDeleted: 0,
	}).Count()
	if num > 0 {
		return gerror.New("手机号已存在")
	}

	//判断是否有权限选择当前部门
	_, err = service.SysDept().Detail(ctx, int64(input.DeptId))
	if err != nil {
		return
	}

	timestamp := gtime.Now().TimestampMilli()
	code := "U_" + gconv.String(timestamp) + "_" + grand.S(6)

	//开启事务管理
	err = dao.SysUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		//添加用户
		var sysUser *entity.SysUser
		if err = gconv.Scan(input, &sysUser); err != nil {
			return err
		}
		//创建用户密码
		sysUser.UserSalt = grand.S(10)
		sysUser.UserPassword = utils.EncryptPassword(input.UserPassword, sysUser.UserSalt)
		//初始化为未删除
		sysUser.IsDeleted = 0
		//获取当前登录用户ID
		loginUserId := service.Context().GetUserId(ctx)
		//添加用户信息
		result, err := dao.SysUser.Ctx(ctx).Data(do.SysUser{
			Code:         code,
			UserName:     sysUser.UserName,
			UserTypes:    sysUser.UserTypes,
			Mobile:       sysUser.Mobile,
			UserNickname: sysUser.UserNickname,
			Birthday:     sysUser.Birthday,
			UserPassword: sysUser.UserPassword,
			UserSalt:     sysUser.UserSalt,
			UserEmail:    sysUser.UserEmail,
			Sex:          sysUser.Sex,
			Avatar:       sysUser.Avatar,
			DeptId:       sysUser.DeptId,
			Remark:       sysUser.Remark,
			IsAdmin:      sysUser.IsAdmin,
			Address:      sysUser.Address,
			Describe:     sysUser.Describe,
			Status:       sysUser.Status,
			IsDeleted:    sysUser.IsDeleted,
			CreatedBy:    uint(loginUserId),
			CreatedAt:    gtime.Now(),
		}).Insert()
		if err != nil {
			return gerror.New("添加用户失败")
		}

		//获取主键ID
		lastInsertId, err := service.Sequences().GetSequences(ctx, result, dao.SysUser.Table(), dao.SysUser.Columns().Id)
		if err != nil {
			return
		}

		//绑定岗位
		err = service.SysUserPost().BindUserAndPost(ctx, int(lastInsertId), input.PostIds)

		//绑定角色
		err = service.SysUserRole().BindUserAndRole(ctx, int(lastInsertId), input.RoleIds)

		return err
	})
	return
}

func (s *sSysUser) Edit(ctx context.Context, input *model.EditUserInput) (err error) {
	//根据ID获取用户信息
	var sysUser *entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, input.Id).Scan(&sysUser)
	if sysUser == nil {
		return gerror.New("Id错误")
	}
	if sysUser.IsDeleted == 1 {
		return gerror.New("该用户已删除")
	}

	//判断账户是否存在
	var sysUserByUserName *entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(g.Map{
		dao.SysUser.Columns().UserName:  input.UserName,
		dao.SysUser.Columns().IsDeleted: 0,
	}).Scan(&sysUserByUserName)
	if sysUserByUserName != nil && sysUserByUserName.Id != input.Id {
		return gerror.New("账户已存在")
	}
	//判断手机号是否存在
	var sysUserByMobile *entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(g.Map{
		dao.SysUser.Columns().Mobile:    input.Mobile,
		dao.SysUser.Columns().IsDeleted: 0,
	}).Scan(&sysUserByMobile)
	if sysUserByMobile != nil && sysUserByMobile.Id != input.Id {
		return gerror.New("手机号已存在")
	}

	//开启事务管理
	err = dao.SysUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		//编辑用户
		sysUser.UserNickname = input.UserNickname
		sysUser.DeptId = input.DeptId
		sysUser.Mobile = input.Mobile
		sysUser.Status = input.Status
		sysUser.Address = input.Address
		sysUser.Avatar = input.Avatar
		sysUser.Birthday = gtime.NewFromStrFormat(input.Birthday, "Y-m-d")
		sysUser.Describe = input.Describe
		sysUser.UserTypes = input.UserTypes
		sysUser.UserEmail = input.UserEmail
		sysUser.Sex = input.Sex
		sysUser.IsAdmin = input.IsAdmin
		//获取当前登录用户ID
		loginUserId := service.Context().GetUserId(ctx)
		sysUser.UpdatedBy = uint(loginUserId)
		//编辑用户信息
		_, err = dao.SysUser.Ctx(ctx).Data(sysUser).Where(dao.SysUser.Columns().Id, input.Id).Update()
		if err != nil {
			return gerror.New("编辑用户失败")
		}

		//绑定岗位
		err = service.SysUserPost().BindUserAndPost(ctx, int(input.Id), input.PostIds)

		//绑定角色
		err = service.SysUserRole().BindUserAndRole(ctx, int(input.Id), input.RoleIds)
		return err
	})
	return
}

// GetUserById 根据ID获取用户信息
func (s *sSysUser) GetUserById(ctx context.Context, id uint) (out *model.UserInfoOut, err error) {
	var e *entity.SysUser
	m := dao.SysUser.Ctx(ctx)

	err = m.Where(g.Map{
		dao.SysUser.Columns().Id:        id,
		dao.SysUser.Columns().IsDeleted: 0,
	}).Scan(&e)
	if err != nil {
		return
	}
	if e != nil {
		if err = gconv.Scan(e, &out); err != nil {
			return
		}

		//获取用户角色ID
		userRoleInfo, userRoleErr := service.SysUserRole().GetInfoByUserId(ctx, int(e.Id))
		if userRoleErr != nil {
			return nil, userRoleErr
		}
		if userRoleInfo != nil {
			var roleIds []int
			for _, userRole := range userRoleInfo {
				roleIds = append(roleIds, userRole.RoleId)
			}
			out.RoleIds = roleIds
		}
		//获取用户岗位ID
		userPostInfo, userPostErr := service.SysUserPost().GetInfoByUserId(ctx, int(e.Id))
		if userPostErr != nil {
			return nil, userPostErr
		}
		if userPostInfo != nil {
			var postIds []int
			for _, userPost := range userPostInfo {
				postIds = append(postIds, userPost.PostId)
			}
			out.PostIds = postIds
		}
	}

	return
}

// DelInfoById 根据ID删除信息
func (s *sSysUser) DelInfoById(ctx context.Context, id uint) (err error) {
	var sysUser *entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(g.Map{
		dao.SysUser.Columns().Id: id,
	}).Scan(&sysUser)
	if sysUser == nil {
		err = gerror.New("ID错误")
		return
	}
	if sysUser.IsDeleted == 1 {
		return gerror.New("用户已删除,无须重复删除")
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)
	_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
		dao.SysUser.Columns().IsDeleted: 1,
		dao.SysUser.Columns().DeletedBy: loginUserId,
		dao.SysUser.Columns().DeletedAt: gtime.Now(),
	}).Where(g.Map{
		dao.SysUser.Columns().Id: id,
	}).Update()
	if err != nil {
		return gerror.New("删除用户失败")
	}
	return
}

// EditUserStatus 修改用户状态
func (s *sSysUser) EditUserStatus(ctx context.Context, id uint, status uint) (err error) {
	var sysUser *entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(g.Map{
		dao.SysUser.Columns().Id: id,
	}).Scan(&sysUser)
	if sysUser == nil {
		err = gerror.New("ID错误")
		return
	}
	if sysUser.Status == status {
		return gerror.New("无须重复修改状态")
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)
	sysUser.Status = status
	sysUser.UpdatedBy = uint(loginUserId)
	_, err = dao.SysUser.Ctx(ctx).Data(sysUser).Where(g.Map{
		dao.SysUser.Columns().Id: id,
	}).Update()
	if err != nil {
		return gerror.New("修改状态失败")
	}
	return
}

// 获取搜索的部门ID数组
func (s *sSysUser) getSearchDeptIds(ctx context.Context, deptId int64) (deptIds []int64, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		deptAll, err := sysDeptNew().GetFromCache(ctx)
		if err != nil {
			return
		}
		deptWithChildren := sysDeptNew().FindSonByParentId(deptAll, gconv.Int64(deptId))
		deptIds = make([]int64, len(deptWithChildren))
		for k, v := range deptWithChildren {
			deptIds[k] = v.DeptId
		}
		deptIds = append(deptIds, deptId)
	})
	return
}
