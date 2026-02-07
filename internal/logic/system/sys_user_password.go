package system

import (
	"context"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/internal/model/do"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/utility/utils"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
)

// ResetPassword 重置密码
func (s *sSysUser) ResetPassword(ctx context.Context, id uint, userPassword string) (err error) {
	var sysUser *entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(g.Map{
		dao.SysUser.Columns().Id: id,
	}).Scan(&sysUser)
	if sysUser == nil {
		err = gerror.New("ID错误")
		return
	}
	if sysUser.Status == 0 {
		return gerror.New("用户已禁用,无法重置密码")
	}
	if sysUser.IsDeleted == 1 {
		return gerror.New("用户已删除,无法重置密码")
	}

	//判断是否启用了安全控制和启用了RSA
	configKeys := []string{consts.SysIsSecurityControlEnabled, consts.SysIsRsaEnabled}
	configDatas, err := service.ConfigData().GetByKeys(ctx, configKeys)
	if err != nil {
		return
	}
	var isSecurityControlEnabled = "0" //是否启用安装控制
	var isRsaEnbled = "0"              //是否启用RSA
	for _, configData := range configDatas {
		if strings.EqualFold(configData.ConfigKey, consts.SysIsSecurityControlEnabled) {
			isSecurityControlEnabled = configData.ConfigValue
		}
		if strings.EqualFold(configData.ConfigKey, consts.SysIsRsaEnabled) {
			isRsaEnbled = configData.ConfigValue
		}
	}

	if strings.EqualFold(isSecurityControlEnabled, "1") {
		if strings.EqualFold(isRsaEnbled, "1") {
			//对用户密码进行解密
			userPasswordDecrypted, err := utils.Decrypt(consts.RsaPrivateKeyFile, userPassword, consts.RsaOAEP)
			if err != nil {
				// RSA 解密失败，可能是前端发送的是明文密码，直接使用明文
				g.Log().Info(ctx, "RSA解密失败，使用明文密码:", err)
			} else {
				userPassword = userPasswordDecrypted
			}
		}

		//校验密码
		err = s.CheckPassword(ctx, userPassword)
		if err != nil {
			return
		}

		//获取用户修改密码的历史记录
		var history []*entity.SysUserPasswordHistory
		if err = dao.SysUserPasswordHistory.Ctx(ctx).Where(dao.SysUserPasswordHistory.Columns().UserId, id).OrderDesc(dao.SysUserPasswordHistory.Columns().CreatedAt).Scan(&history); err != nil {
			return
		}

		if history != nil {
			//判断与最近三次是否一样
			var num int
			if len(history) > 3 {
				num = 3
			} else {
				num = len(history)
			}
			var isExit = false
			for i := 0; i < num; i++ {
				if strings.EqualFold(history[i].AfterPassword, utils.EncryptPassword(userPassword, sysUser.UserSalt)) {
					isExit = true
					break
				}
			}
			if isExit {
				err = gerror.New("密码与前三次重复,请重新输入!")
				return
			}
		}
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)
	//判断当前登录用户是否为超级管理员角色
	var isSuperAdmin = false
	var sysUserRoles []*entity.SysUserRole
	err = dao.SysUserRole.Ctx(ctx).Where(dao.SysUserRole.Columns().UserId, loginUserId).Scan(&sysUserRoles)
	for _, sysUserRole := range sysUserRoles {
		if sysUserRole.RoleId == 1 {
			isSuperAdmin = true
			break
		}
	}
	if !isSuperAdmin {
		return gerror.New("无重置密码权限")
	}

	sysUser.UserSalt = grand.S(10)

	beforePassword := sysUser.UserPassword
	afterPassword := utils.EncryptPassword(userPassword, sysUser.UserSalt)
	sysUser.UserPassword = afterPassword
	sysUser.UpdatedBy = uint(loginUserId)

	//开启事务管理
	err = dao.SysUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		_, err = dao.SysUser.Ctx(ctx).Data(sysUser).Where(g.Map{
			dao.SysUser.Columns().Id: id,
		}).Update()

		//添加
		_, err = dao.SysUserPasswordHistory.Ctx(ctx).Data(&do.SysUserPasswordHistory{
			UserId:         sysUser.Id,
			BeforePassword: beforePassword,
			AfterPassword:  afterPassword,
			ChangeTime:     gtime.Now(),
			CreatedAt:      gtime.Now(),
			CreatedBy:      loginUserId,
		}).Insert()

		return
	})

	if err != nil {
		return gerror.New("重置密码失败")
	}
	return
}

// CheckPassword 校验用户密码
func (s *sSysUser) CheckPassword(ctx context.Context, userPassword string) (err error) {
	keys := []string{consts.SysPasswordMinimumLength, consts.SysRequireComplexity, consts.SysRequireDigit, consts.SysRequireLowercaseLetter, consts.SysRequireUppercaseLetter}
	var configData []*entity.SysConfig
	configData, err = service.ConfigData().GetByKeys(ctx, keys)
	if err != nil {
		return
	}
	if configData == nil || len(configData) == 0 {
		err = gerror.New(g.I18n().T(ctx, "{#sysUserPwCheckConfig}"))
	}
	var minimumLength int
	var complexity int
	var digit int
	var lowercaseLetter int
	var uppercaseLetter int
	for _, data := range configData {
		if strings.EqualFold(data.ConfigKey, consts.SysPasswordMinimumLength) {
			minimumLength, _ = strconv.Atoi(data.ConfigValue)
		}
		if strings.EqualFold(data.ConfigKey, consts.SysRequireComplexity) {
			complexity, _ = strconv.Atoi(data.ConfigValue)
		}
		if strings.EqualFold(data.ConfigKey, consts.SysRequireDigit) {
			digit, _ = strconv.Atoi(data.ConfigValue)
		}
		if strings.EqualFold(data.ConfigKey, consts.SysRequireLowercaseLetter) {
			lowercaseLetter, _ = strconv.Atoi(data.ConfigValue)
		}
		if strings.EqualFold(data.ConfigKey, consts.SysRequireUppercaseLetter) {
			uppercaseLetter, _ = strconv.Atoi(data.ConfigValue)
		}
	}

	var flag bool
	flag, err = utils.ValidatePassword(userPassword, minimumLength, complexity, digit, lowercaseLetter, uppercaseLetter)
	if err != nil {
		return
	}
	if !flag {
		err = gerror.New(g.I18n().T(ctx, "{#sysUserPwCheckError}"))
		return
	}
	return
}

// EditPassword 修改密码
func (s *sSysUser) EditPassword(ctx context.Context, userName string, oldUserPassword string, userPassword string) (err error) {
	//判断是否启用了安全控制和启用了RSA
	configKeys := []string{consts.SysIsSecurityControlEnabled, consts.SysIsRsaEnabled}
	configDatas, err := service.ConfigData().GetByKeys(ctx, configKeys)
	if err != nil {
		return
	}
	var isSecurityControlEnabled = "0"
	var isRSAEnbled = "0"
	for _, configData := range configDatas {
		if strings.EqualFold(configData.ConfigKey, consts.SysIsSecurityControlEnabled) {
			isSecurityControlEnabled = configData.ConfigValue
		}
		if strings.EqualFold(configData.ConfigKey, consts.SysIsRsaEnabled) {
			isRSAEnbled = configData.ConfigValue
		}
	}
	if strings.EqualFold(isSecurityControlEnabled, "1") && strings.EqualFold(isRSAEnbled, "1") {
		//对旧密码进行解密
		oldUserPasswordDecrypted, err := utils.Decrypt(consts.RsaPrivateKeyFile, oldUserPassword, consts.RsaOAEP)
		if err != nil {
			// RSA 解密失败，可能是前端发送的是明文密码，直接使用明文
			g.Log().Info(ctx, "RSA解密失败，使用明文密码:", err)
		} else {
			oldUserPassword = oldUserPasswordDecrypted
		}
	}

	//获取用户信息
	userInfo, err := service.SysUser().GetUserByUsername(ctx, userName)
	if err != nil {
		return
	}
	//判断旧密码是否一致
	if !strings.EqualFold(userInfo.UserPassword, utils.EncryptPassword(oldUserPassword, userInfo.UserSalt)) {
		err = gerror.New("原密码输入错误, 请重新输入!")
		return
	}
	//重置密码
	err = s.ResetPassword(ctx, uint(userInfo.Id), userPassword)
	if err != nil {
		return
	}
	return
}
