package system

import (
	"context"
	"errors"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/cache"
	"sagooiot/pkg/utility/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// GetAdminUserByUsernamePassword 根据用户名和密码获取用户信息
func (s *sSysUser) GetAdminUserByUsernamePassword(ctx context.Context, userName string, password string) (user *entity.SysUser, err error) {
	//判断账号是否存在
	user, err = s.GetUserByUsername(ctx, userName)
	if err != nil {
		return nil, err
	}
	//验证密码是否正确
	if utils.EncryptPassword(password, user.UserSalt) != user.UserPassword {
		//判断获取是否启动安全控制和允许再次登录时间
		configKeys := []string{consts.SysIsSecurityControlEnabled, consts.SysAgainLoginDate}
		var configDatas []*entity.SysConfig
		configDatas, err = service.ConfigData().GetByKeys(ctx, configKeys)
		if err != nil {
			return
		}
		isSecurityControlEnabled := "0" //是否启动安全控制
		againLoginDate := 1             //允许再次登录时间
		for _, configData := range configDatas {
			if strings.EqualFold(configData.ConfigKey, consts.SysIsSecurityControlEnabled) {
				isSecurityControlEnabled = configData.ConfigValue
			}
			if strings.EqualFold(configData.ConfigKey, consts.SysAgainLoginDate) {
				againLoginDate, err = strconv.Atoi(configData.ConfigValue)
				if err != nil {
					err = gerror.New("允许再次登录时间配置错误")
					return nil, err
				}
			}
		}
		if strings.EqualFold(isSecurityControlEnabled, "1") {
			tmpData, _ := cache.Instance().Get(ctx, consts.CacheSysErrorPrefix+"_"+gconv.String(userName))
			var num = 1
			if tmpData.Val() != nil {
				if strVal, ok := tmpData.Val().(string); ok {
					tempValue, _ := strconv.Atoi(strVal)
					num = tempValue + 1
				}
			}
			//存入缓存
			err := cache.Instance().Set(ctx, consts.CacheSysErrorPrefix+"_"+gconv.String(userName), num, time.Duration(againLoginDate*60)*time.Second)
			if err != nil {
				return nil, err
			}
		}
		err = gerror.New("密码错误")
		return
	}
	return
}

// UpdateLoginInfo 更新用户登录信息
func (s *sSysUser) UpdateLoginInfo(ctx context.Context, id uint64, ip string) (err error) {
	_, err = dao.SysUser.Ctx(ctx).WherePri(id).Update(g.Map{
		dao.SysUser.Columns().LastLoginIp:   ip,
		dao.SysUser.Columns().LastLoginTime: gtime.Now(),
	})
	if err != nil {
		return errors.New("更新用户登录信息失败")
	}

	return
}
