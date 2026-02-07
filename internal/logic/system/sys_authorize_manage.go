package system

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"sagooiot/internal/consts"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/cache"
	"strings"
)

func (s *sSysAuthorize) AddAuthorize(ctx context.Context, roleId int, menuIds []string, buttonIds []string, columnIds []string, apiIds []string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		//删除原有权限
		err = service.SysAuthorize().DelByRoleId(ctx, roleId)

		var authorizeInfo []*entity.SysAuthorize

		//封装菜单权限
		menuEntries, buildErr := buildAuthorizeEntries(menuIds, consts.Menu, roleId)
		if buildErr != nil {
			err = buildErr
			return
		}
		authorizeInfo = append(authorizeInfo, menuEntries...)

		//封装按钮权限
		buttonEntries, buildErr := buildAuthorizeEntries(buttonIds, consts.Button, roleId)
		if buildErr != nil {
			err = buildErr
			return
		}
		authorizeInfo = append(authorizeInfo, buttonEntries...)

		//封装列表权限
		columnEntries, buildErr := buildAuthorizeEntries(columnIds, consts.Column, roleId)
		if buildErr != nil {
			err = buildErr
			return
		}
		authorizeInfo = append(authorizeInfo, columnEntries...)

		//封装接口权限
		apiEntries, buildErr := buildAuthorizeEntries(apiIds, consts.Api, roleId)
		if buildErr != nil {
			err = buildErr
			return
		}
		authorizeInfo = append(authorizeInfo, apiEntries...)

		err = s.Add(ctx, authorizeInfo)
		if err != nil {
			err = gerror.New("添加权限失败")
			return
		}
		//添加缓存信息
		err := cache.Instance().Set(ctx, consts.CacheSysAuthorize+"_"+gconv.String(roleId), authorizeInfo, 0)
		if err != nil {
			return
		}
	})
	return
}
func (s *sSysAuthorize) IsAllowAuthorize(ctx context.Context, roleId int) (isAllow bool, err error) {
	//判断角色ID是否为1
	if roleId == 1 {
		err = gerror.New("无法更改超级管理员权限")
		return
	}

	isAllow = false
	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)
	//查看当前登录用户所有的角色
	//查询用户角色信息
	userRoleInfo, err := service.SysUserRole().GetInfoByUserId(ctx, loginUserId)
	if err != nil {
		return
	}
	//判断是否为超级管理员
	if userRoleInfo == nil {
		err = gerror.New("无权限,请联系管理员")
		return
	}
	var isSuperAdmin = false
	var roleIds []int
	for _, userRole := range userRoleInfo {
		if userRole.RoleId == 1 {
			isSuperAdmin = true
		}
		roleIds = append(roleIds, userRole.RoleId)
	}
	//判断当前用户是否为超级管理员
	if isSuperAdmin {
		//如果是超级管理员则可以对所有角色授权
		isAllow = true
	} else {
		//根据角色ID获取菜单配置
		authorizeInfo, authorizeErr := s.GetInfoByRoleId(ctx, roleId)
		if authorizeErr != nil {
			return
		}
		//如果需要授权的角色ID无任何权限，则当前用户可以对其授权
		if authorizeInfo == nil {
			isAllow = true
		} else {
			//菜单Ids
			var menuIds []int
			//按钮Ids
			var menuButtonIds []int
			//列表Ids
			var menuColumnIds []int
			//API Ids
			var menuApiIds []int
			for _, authorize := range authorizeInfo {
				if strings.EqualFold(authorize.ItemsType, consts.Menu) {
					menuIds = append(menuIds, authorize.ItemsId)
				} else if strings.EqualFold(authorize.ItemsType, consts.Button) {
					menuButtonIds = append(menuButtonIds, authorize.ItemsId)
				} else if strings.EqualFold(authorize.ItemsType, consts.Column) {
					menuColumnIds = append(menuColumnIds, authorize.ItemsId)
				} else if strings.EqualFold(authorize.ItemsType, consts.Api) {
					menuApiIds = append(menuApiIds, authorize.ItemsId)
				}
			}
			//获取当前用户所有的权限
			nowUserAuthorizeInfo, nowUserAuthorizeErr := s.GetInfoByRoleIds(ctx, roleIds)
			if nowUserAuthorizeErr != nil {
				return
			}

			//获取系统列表开关参数
			var sysColumnSwitchConfig *entity.SysConfig
			sysColumnSwitchConfig, err = service.ConfigData().GetConfigByKey(ctx, "sys.column.switch")
			sysColumnSwitch := 0
			if sysColumnSwitchConfig != nil {
				sysColumnSwitch = gconv.Int(sysColumnSwitchConfig.ConfigValue)
			}
			//获取系统按钮开关参数
			var sysButtonSwitchConfig *entity.SysConfig
			sysButtonSwitchConfig, err = service.ConfigData().GetConfigByKey(ctx, "sys.button.switch")
			sysButtonSwitch := 0
			if sysButtonSwitchConfig != nil {
				sysButtonSwitch = gconv.Int(sysButtonSwitchConfig.ConfigValue)
			}

			//获取系统API开关参数
			var sysApiSwitchConfig *entity.SysConfig
			sysApiSwitchConfig, err = service.ConfigData().GetConfigByKey(ctx, "sys.api.switch")
			sysApiSwitch := 0
			if sysApiSwitchConfig != nil {
				sysApiSwitch = gconv.Int(sysApiSwitchConfig.ConfigValue)
			}

			//菜单Ids
			var nowUserMenuIds []int
			//按钮Ids
			var nowUserMenuButtonIds []int
			//列表Ids
			var nowUserMenuColumnIds []int
			//API Ids
			var nowUserMenuApiIds []int
			for _, authorize := range nowUserAuthorizeInfo {
				if strings.EqualFold(authorize.ItemsType, consts.Menu) {
					nowUserMenuIds = append(nowUserMenuIds, authorize.ItemsId)
				} else if strings.EqualFold(authorize.ItemsType, consts.Button) && sysButtonSwitch == 1 {
					nowUserMenuButtonIds = append(nowUserMenuButtonIds, authorize.ItemsId)
				} else if strings.EqualFold(authorize.ItemsType, consts.Column) && sysColumnSwitch == 1 {
					nowUserMenuColumnIds = append(nowUserMenuColumnIds, authorize.ItemsId)
				} else if strings.EqualFold(authorize.ItemsType, consts.Api) && sysApiSwitch == 1 {
					nowUserMenuApiIds = append(nowUserMenuApiIds, authorize.ItemsId)
				}
			}

			//判断按钮、列表、API开关状态，如果关闭则获取菜单对应的所有信息
			//判断按钮开关
			if sysButtonSwitch == 0 {
				//获取所有按钮
				var menuButtons []*entity.SysMenuButton
				menuButtons, err = service.SysMenuButton().GetInfoByMenuIds(ctx, menuIds)
				if err != nil {
					return
				}
				if len(menuButtons) > 0 {
					for _, menuButton := range menuButtons {
						nowUserMenuButtonIds = append(nowUserMenuButtonIds, int(menuButton.Id))
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
						nowUserMenuColumnIds = append(nowUserMenuColumnIds, int(menuColumn.Id))
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
						nowUserMenuApiIds = append(nowUserMenuApiIds, int(menuApi.Id))
					}
				}
			}

			//判断当前登录用户是否大于授权角色的权限
			//获取当前登录用户的菜单信息
			nowUserMenuInfo, _ := service.SysMenu().GetInfoByMenuIds(ctx, nowUserMenuIds)
			//获取授权角色的菜单信息
			menuInfo, _ := service.SysMenu().GetInfoByMenuIds(ctx, menuIds)
			if len(menuInfo) == 0 {
				isAllow = true
			} else {
				var menuInfoIsAllow = true
				//判断当前登录用户所拥有的菜单是包含授权角色的菜单
				for _, menu := range menuInfo {
					var nowMenuIsAllow = false
					for _, nowUserMenu := range nowUserMenuInfo {
						if menu.Id == nowUserMenu.Id {
							nowMenuIsAllow = true
							break
						}
					}
					if !nowMenuIsAllow {
						menuInfoIsAllow = false
						break
					}
				}
				//判断是否都包含
				if menuInfoIsAllow {
					//判断按钮是否都包含
					//获取当前登录用户按钮单信息
					nowUserMenuButtonInfo, _ := service.SysMenuButton().GetInfoByMenuIds(ctx, nowUserMenuButtonIds)
					//获取授权角色的按钮信息
					menuButtonInfo, _ := service.SysMenuButton().GetInfoByMenuIds(ctx, menuButtonIds)
					var menuButtonInfoIsAllow = true
					if len(menuButtonInfo) != 0 {
						//判断当前登录用户所拥有的按钮是包含授权角色的按钮
						for _, menuButton := range menuButtonInfo {
							var nowMenuButtonIsAllow = false
							for _, nowUserMenuButton := range nowUserMenuButtonInfo {
								if menuButton.Id == nowUserMenuButton.Id {
									nowMenuButtonIsAllow = true
									break
								}
							}
							if !nowMenuButtonIsAllow {
								menuButtonInfoIsAllow = false
								break
							}
						}
						if menuButtonInfoIsAllow {
							//判断列表是否都包含
							//获取当前登录用户的列表信息
							nowUserMenuColumnInfo, _ := service.SysMenuColumn().GetInfoByMenuIds(ctx, nowUserMenuColumnIds)
							//获取授权角色的列表信息
							menuColumnInfo, _ := service.SysMenuColumn().GetInfoByMenuIds(ctx, menuColumnIds)
							var menuColumnInfoIsAllow = true
							if len(menuColumnInfo) == 0 {
								//判断当前登录用户所拥有的列表是包含授权角色的列表
								for _, menuColumn := range menuColumnInfo {
									var nowMenuColumnIsAllow = false
									for _, nowUserMenuColumn := range nowUserMenuColumnInfo {
										if menuColumn.Id == nowUserMenuColumn.Id {
											nowMenuColumnIsAllow = true
											break
										}
									}
									if !nowMenuColumnIsAllow {
										menuColumnInfoIsAllow = false
										break
									}
								}

								if menuColumnInfoIsAllow {
									isAllow = true
								}
							} else {
								isAllow = true
							}
						}

					} else {
						isAllow = true
					}

				}
			}
		}
	}
	return
}
