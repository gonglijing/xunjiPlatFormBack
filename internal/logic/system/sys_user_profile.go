package system

import (
	"context"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/utility/utils"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EditUserAvatar 修改用户头像
func (s *sSysUser) EditUserAvatar(ctx context.Context, id uint, avatar string) (err error) {
	var sysUser *entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(g.Map{
		dao.SysUser.Columns().Id: id,
	}).Scan(&sysUser)
	if sysUser == nil {
		err = gerror.New("ID错误")
		return
	}
	sysUser.Avatar = avatar

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)

	//判断是否为当前用户
	if id != uint(loginUserId) {
		err = gerror.New("无法修改其他用户头像")
		return
	}

	sysUser.UpdatedBy = uint(loginUserId)
	sysUser.UpdatedAt = gtime.Now()
	_, err = dao.SysUser.Ctx(ctx).Data(sysUser).Where(g.Map{
		dao.SysUser.Columns().Id: id,
	}).Update()
	if err != nil {
		return gerror.New("修改头像失败")
	}
	return
}

// EditUserInfo 修改用户个人资料
func (s *sSysUser) EditUserInfo(ctx context.Context, input *model.EditUserInfoInput) (err error) {
	var sysUser *entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(g.Map{
		dao.SysUser.Columns().Id: input.Id,
	}).Scan(&sysUser)
	if sysUser == nil {
		err = gerror.New("ID错误")
		return
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)

	//判断是否为当前用户
	if input.Id != uint64(loginUserId) {
		err = gerror.New("无法修改其他用户资料")
		return
	}
	sysUser.Mobile = input.Mobile
	sysUser.UserNickname = input.UserNickname
	sysUser.Birthday = gtime.NewFromStrFormat(input.Birthday, "Y-m-d")
	if input.UserPassword != "" {
		sysUser.UserPassword = utils.EncryptPassword(input.UserPassword, sysUser.UserSalt)
	}
	sysUser.UserEmail = input.UserEmail
	sysUser.Sex = input.Sex
	sysUser.Avatar = input.Avatar
	sysUser.Address = input.Address
	sysUser.Describe = input.Describe
	//获取当前登录用户ID
	sysUser.UpdatedBy = uint(loginUserId)
	sysUser.UpdatedAt = gtime.Now()
	_, err = dao.SysUser.Ctx(ctx).Data(sysUser).Where(g.Map{
		dao.SysUser.Columns().Id: input.Id,
	}).Update()
	if err != nil {
		return gerror.New("修改个人资料失败")
	}
	return
}
