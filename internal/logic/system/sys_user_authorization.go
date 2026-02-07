package system

import (
	"context"
	"encoding/json"
	"sagooiot/internal/consts"
	"sagooiot/internal/model"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/cache"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

func (s *sSysUser) CurrentUser(ctx context.Context) (userInfoOut *model.UserInfoOut, menuTreeOut []*model.UserMenuTreeOut, err error) {
	//获取当前登录用户信息
	loginUserId := service.Context().GetUserId(ctx)
	if loginUserId == 0 {
		err = gerror.New("无登录用户信息,请先登录!")
		return
	}
	tmpUserAuthorize, err := cache.Instance().Get(ctx, consts.CacheUserAuthorize+"_"+gconv.String(loginUserId))
	tmpUserInfo, err := cache.Instance().Get(ctx, consts.CacheUserInfo+"_"+gconv.String(loginUserId))
	if err != nil {
		return
	}
	if tmpUserAuthorize == nil || tmpUserInfo == nil {
		err = gerror.New("无登录用户信息,请先登录!")
	}

	if tmpUserAuthorize.Val() != nil && tmpUserInfo.Val() != nil {
		if err = json.Unmarshal([]byte(tmpUserAuthorize.Val().(string)), &menuTreeOut); err != nil {
			return
		}
		err = json.Unmarshal([]byte(tmpUserInfo.Val().(string)), &userInfoOut)
		return
	}

	//获取当前登录用户信息
	userInfo, err := service.SysUser().GetUserById(ctx, uint(loginUserId))
	if err = gconv.Scan(userInfo, &userInfoOut); err != nil {
		return
	}
	if userInfo != nil {
		err := cache.Instance().Set(ctx, consts.CacheUserInfo+"_"+gconv.String(loginUserId), userInfo, time.Hour)
		if err != nil {
			return nil, nil, err
		}
	}

	//根据当前登录用户ID查询用户角色信息
	userRoleInfo, err := service.SysUserRole().GetInfoByUserId(ctx, loginUserId)
	if userRoleInfo == nil {
		err = gerror.New("用户无权限,禁止访问")
		return
	}
	var isSuperAdmin = false
	var roleIds []int
	//获取角色ID
	for _, role := range userRoleInfo {
		if role.RoleId == 1 {
			isSuperAdmin = true
		}
		roleIds = append(roleIds, role.RoleId)
	}
	if isSuperAdmin {
		//获取所有的菜单
		//根据菜单ID数组获取菜单列表信息
		userMenuTreeOut, menuTreeError := GetMenuInfo(ctx, nil)
		if menuTreeError != nil {
			err = menuTreeError
			return
		}
		var menuIds []int
		//获取菜单绑定的按钮信息
		for _, userMenu := range userMenuTreeOut {
			menuIds = append(menuIds, int(userMenu.Id))
		}
		//获取所有的按钮
		userMenuTreeOut, userItemsTypeTreeErr := GetUserItemsTypeTreeOut(ctx, menuIds, consts.Button, userMenuTreeOut)
		if userItemsTypeTreeErr != nil {
			err = userItemsTypeTreeErr
			return
		}
		//获取所有的列表
		userMenuTreeOut, userItemsTypeTreeErr = GetUserItemsTypeTreeOut(ctx, menuIds, consts.Column, userMenuTreeOut)
		if userItemsTypeTreeErr != nil {
			err = userItemsTypeTreeErr
			return
		}
		//获取所有的接口
		userMenuTreeOut, userItemsTypeTreeErr = GetUserItemsTypeTreeOut(ctx, menuIds, consts.Api, userMenuTreeOut)
		if userItemsTypeTreeErr != nil {
			err = userItemsTypeTreeErr
			return
		}

		//对菜单进行树状重组
		menuTreeOut = GetUserMenuTree(userMenuTreeOut)

		if menuTreeOut != nil {
			err := cache.Instance().Set(ctx, consts.CacheUserAuthorize+"_"+gconv.String(loginUserId), menuTreeOut, 0)
			if err != nil {
				return nil, nil, err
			}
		}
		return
	} else {
		//获取缓存配置信息
		var authorizeInfo []*entity.SysAuthorize
		authorizeInfo, err = service.SysAuthorize().GetInfoByRoleIds(ctx, roleIds)
		if err != nil {
			return
		}
		if authorizeInfo == nil {
			err = gerror.New("无权限配置,请联系管理员")
			return
		}

		isSecurityControlEnabled := 0 //是否启用安全控制
		sysColumnSwitch := 0          //列表开关
		sysButtonSwitch := 0          //按钮开关
		sysApiSwitch := 0             //api开关
		//判断是否启用了安全控制、列表开关、按钮开关、api开关
		configKeys := []string{consts.SysIsSecurityControlEnabled, consts.SysColumnSwitch, consts.SysButtonSwitch, consts.SysApiSwitch}
		var configDatas []*entity.SysConfig
		configDatas, err = service.ConfigData().GetByKeys(ctx, configKeys)
		if err != nil {
			return
		}
		for _, configData := range configDatas {
			if strings.EqualFold(configData.ConfigKey, consts.SysIsSecurityControlEnabled) {
				isSecurityControlEnabled = gconv.Int(configData.ConfigValue)
			}
			if strings.EqualFold(configData.ConfigKey, consts.SysColumnSwitch) {
				sysColumnSwitch = gconv.Int(configData.ConfigValue)
			}
			if strings.EqualFold(configData.ConfigKey, consts.SysButtonSwitch) {
				sysButtonSwitch = gconv.Int(configData.ConfigValue)
			}
			if strings.EqualFold(configData.ConfigKey, consts.SysApiSwitch) {
				sysApiSwitch = gconv.Int(configData.ConfigValue)
			}
		}
		//菜单Ids
		var menuIds []int
		var menuButtonIds []int
		//列表Ids
		var menuColumnIds []int
		//API Ids
		var menuApiIds []int

		for _, authorize := range authorizeInfo {
			if strings.EqualFold(authorize.ItemsType, consts.Menu) {
				menuIds = append(menuIds, authorize.ItemsId)
			} else if strings.EqualFold(authorize.ItemsType, consts.Button) && sysButtonSwitch == 1 {
				menuButtonIds = append(menuButtonIds, authorize.ItemsId)
			} else if strings.EqualFold(authorize.ItemsType, consts.Column) && sysColumnSwitch == 1 {
				menuColumnIds = append(menuColumnIds, authorize.ItemsId)
			} else if strings.EqualFold(authorize.ItemsType, consts.Api) && sysApiSwitch == 1 {
				menuApiIds = append(menuApiIds, authorize.ItemsId)
			}
		}
		//判断按钮、列表、API开关状态，如果关闭则获取菜单对应的所有信息
		//判断按钮开关
		if isSecurityControlEnabled == 0 {
			if sysButtonSwitch == 0 {
				//获取所有按钮
				var menuButtons []*entity.SysMenuButton
				menuButtons, err = service.SysMenuButton().GetInfoByMenuIds(ctx, menuIds)
				if err != nil {
					return
				}
				if len(menuButtons) > 0 {
					for _, menuButton := range menuButtons {
						menuButtonIds = append(menuButtonIds, int(menuButton.Id))
					}
				}
			}

			//判断列表开关
			if sysColumnSwitch == 0 {
				//获取所有列表
				var menuColumns []*entity.SysMenuColumn
				menuColumns, err = service.SysMenuColumn().GetInfoByMenuIds(ctx, menuIds)
				if err != nil {
					return
				}
				if len(menuColumns) > 0 {
					for _, menuColumn := range menuColumns {
						menuColumnIds = append(menuColumnIds, int(menuColumn.Id))
					}
				}
			}
			//判断API开关
			if sysApiSwitch == 0 {
				//获取所有API
				var menuApis []*entity.SysMenuApi
				menuApis, err = service.SysMenuApi().GetInfoByMenuIds(ctx, menuIds)
				if err != nil {
					return
				}
				if len(menuApis) > 0 {
					for _, menuApi := range menuApis {
						menuApiIds = append(menuApiIds, int(menuApi.Id))
					}
				}
			}
		}

		//获取所有菜单信息
		var menuInfo []*entity.SysMenu
		menuInfo, err = service.SysMenu().GetInfoByMenuIds(ctx, menuIds)
		if err != nil {
			return
		}
		if menuInfo == nil {
			err = gerror.New("未配置菜单, 请联系系统管理员")
			return
		}
		//根据按钮ID数组获取按钮信息
		menuButtonInfo, menuButtonErr := service.SysMenuButton().GetInfoByButtonIds(ctx, menuButtonIds)
		if menuButtonErr != nil {
			return
		}
		//根据列表字段ID数据获取字段信息
		menuColumnInfo, menuColumnErr := service.SysMenuColumn().GetInfoByColumnIds(ctx, menuColumnIds)
		if menuColumnErr != nil {
			return
		}
		//根据菜单接口ID数组获取接口信息
		menuApiInfo, menuApiErr := service.SysMenuApi().GetInfoByIds(ctx, menuApiIds)
		if menuApiErr != nil {
			return
		}

		var userMenuTreeOut []*model.UserMenuTreeOut

		for _, menu := range menuInfo {
			var userMenuTree *model.UserMenuTreeOut
			if err = gconv.Scan(menu, &userMenuTree); err != nil {
				return
			}

			//获取菜单按钮权限
			if menuButtonInfo != nil {
				menuButtonTreeData, _ := GetUserMenuButton(int(menu.Id), menuButtonInfo)
				userMenuTree.Button = append(userMenuTree.Button, menuButtonTreeData...)
			}

			//获取列表权限
			if menuColumnInfo != nil {
				menuColumnTreeData, _ := GetUserMenuColumn(int(menu.Id), menuColumnInfo)
				userMenuTree.Column = append(userMenuTree.Column, menuColumnTreeData...)
			}

			//获取接口权限
			if menuApiIds != nil {
				//获取相关接口ID
				var apiIds []int
				for _, menuApi := range menuApiInfo {
					if menuApi.MenuId == int(menu.Id) {
						apiIds = append(apiIds, menuApi.ApiId)
					}
				}
				//获取相关接口信息
				apiInfo, _ := service.SysApi().GetInfoByIds(ctx, apiIds)
				if apiInfo != nil {
					var apiInfoOut []*model.UserApiOut
					if err = gconv.Scan(apiInfo, &apiInfoOut); err != nil {
						return
					}
					userMenuTree.Api = append(userMenuTree.Api, apiInfoOut...)
				}
			}
			userMenuTreeOut = append(userMenuTreeOut, userMenuTree)
		}

		//对菜单进行树状重组
		menuTreeOut = GetUserMenuTree(userMenuTreeOut)
		if menuTreeOut != nil {
			err := cache.Instance().Set(ctx, consts.CacheUserAuthorize+"_"+gconv.String(loginUserId), menuTreeOut, 0)
			if err != nil {
				return nil, nil, err
			}
		}
		return
	}
}
