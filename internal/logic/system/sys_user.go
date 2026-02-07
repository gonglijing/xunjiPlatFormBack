package system

import (
	"context"
	"sagooiot/internal/dao"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

type sSysUser struct {
}

func init() {
	service.RegisterSysUser(sysUserNew())
}

func sysUserNew() *sSysUser {
	return &sSysUser{}
}

// GetUserByUsername 通过用户名获取用户信息
func (s *sSysUser) GetUserByUsername(ctx context.Context, userName string) (data *entity.SysUser, err error) {
	var user *entity.SysUser
	err = dao.SysUser.Ctx(ctx).Fields(user).Where(dao.SysUser.Columns().UserName, userName).Scan(&user)
	if user == nil {
		return nil, gerror.New("账号不存在,请重新输入!")
	}
	if user.IsDeleted == 1 {
		return nil, gerror.New("该账号已删除")
	}
	if user.Status == 0 {
		return nil, gerror.New("该账号已禁用")
	}
	if user.Status == 2 {
		return nil, gerror.New("该账号未验证")
	}
	return user, nil
}
