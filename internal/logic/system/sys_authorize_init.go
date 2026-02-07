package system

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/cache"
)

// InitAuthorize 初始化系统权限
func (s *sSysAuthorize) InitAuthorize(ctx context.Context) (err error) {
	//获取所有菜单信息
	menuInfos, err := service.SysMenu().GetAll(ctx)
	if err != nil {
		return
	}
	if menuInfos != nil && len(menuInfos) > 0 {
		//根据菜单信息初始化对应按钮信息
		var menuIds []uint
		for _, menuInfo := range menuInfos {
			menuIds = append(menuIds, menuInfo.Id)
		}
		//获取所有的按钮列表
		var menuButtonInfos []*entity.SysMenuButton
		err = dao.SysMenuButton.Ctx(ctx).Where(g.Map{
			dao.SysMenuButton.Columns().IsDeleted: 0,
			dao.SysMenuButton.Columns().Status:    1,
		}).Scan(&menuButtonInfos)
		if err != nil {
			g.Log().Debug(ctx, "获取菜单按钮信息失败", err.Error())
			return
		}
		if menuButtonInfos != nil && len(menuButtonInfos) > 0 {
			err = cache.Instance().Set(ctx, consts.CacheSysMenuButton, menuButtonInfos, 0)
			if err != nil {
				return
			}
		}
		//获取所有的菜单列表
		var menuColumnInfos []*entity.SysMenuColumn
		err = dao.SysMenuColumn.Ctx(ctx).Where(g.Map{
			dao.SysMenuColumn.Columns().IsDeleted: 0,
			dao.SysMenuColumn.Columns().Status:    1,
		}).Scan(&menuColumnInfos)
		if err != nil {
			return
		}
		if menuColumnInfos != nil && len(menuColumnInfos) > 0 {
			err = cache.Instance().Set(ctx, consts.CacheSysMenuColumn, menuColumnInfos, 0)
			if err != nil {
				return
			}
		}
		//获取所有的菜单接口
		var menuApiInfos []*entity.SysMenuApi
		err = dao.SysMenuApi.Ctx(ctx).Where(g.Map{
			dao.SysMenuApi.Columns().IsDeleted: 0,
		}).Scan(&menuApiInfos)
		if err != nil {
			return
		}
		if menuApiInfos != nil && len(menuApiInfos) > 0 {
			err = cache.Instance().Set(ctx, consts.CacheSysMenuApi, menuApiInfos, 0)
			if err != nil {
				return
			}
		}
		//添加缓存信息
		for _, menuId := range menuIds {
			var tmpMenuButton []*entity.SysMenuButton
			for _, menuButtonInfo := range menuButtonInfos {
				if int(menuId) == menuButtonInfo.MenuId {
					tmpMenuButton = append(tmpMenuButton, menuButtonInfo)
				}
			}
			//添加按钮缓存
			if tmpMenuButton != nil && len(tmpMenuButton) > 0 {
				err = cache.Instance().Set(ctx, consts.CacheSysMenuButton+"_"+gconv.String(menuId), tmpMenuButton, 0)
				if err != nil {
					return
				}
			}

			var tmpMenuColumn []*entity.SysMenuColumn
			for _, menuColumnInfo := range menuColumnInfos {
				if int(menuId) == menuColumnInfo.MenuId {
					tmpMenuColumn = append(tmpMenuColumn, menuColumnInfo)
				}
			}
			//添加列表缓存
			if tmpMenuColumn != nil && len(tmpMenuColumn) > 0 {
				err = cache.Instance().Set(ctx, consts.CacheSysMenuColumn+"_"+gconv.String(menuId), tmpMenuColumn, 0)
				if err != nil {
					return
				}
			}

			var tmpMenuApi []*entity.SysMenuApi
			for _, menuApiInfo := range menuApiInfos {
				if int(menuId) == menuApiInfo.MenuId {
					tmpMenuApi = append(tmpMenuApi, menuApiInfo)
				}
			}
			//添加菜单与接口绑定关系缓存
			if tmpMenuApi != nil && len(tmpMenuApi) > 0 {
				err = cache.Instance().Set(ctx, consts.CacheSysMenuApi+"_"+gconv.String(menuId), tmpMenuApi, 0)
				if err != nil {
					return
				}
			}
		}
		//获取所有的接口信息
		var sysApiInfos []*entity.SysApi
		err = dao.SysApi.Ctx(ctx).Where(g.Map{
			dao.SysApi.Columns().IsDeleted: 0,
			dao.SysApi.Columns().Status:    1,
		}).Scan(&sysApiInfos)
		if err != nil {
			return
		}
		if sysApiInfos != nil && len(sysApiInfos) > 0 {
			err = cache.Instance().Set(ctx, consts.CacheSysApi, sysApiInfos, 0)
			if err != nil {
				return
			}
		}

		//获取所有的角色ID
		var roleInfos []*entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(g.Map{
			dao.SysRole.Columns().IsDeleted: 0,
			dao.SysRole.Columns().Status:    1,
		}).Scan(&roleInfos)
		if roleInfos != nil && len(roleInfos) > 0 {
			//获取所有的权限配置
			var authorizeInfos []*entity.SysAuthorize
			err = dao.SysAuthorize.Ctx(ctx).Where(g.Map{
				dao.SysAuthorize.Columns().IsDeleted: 0,
			}).Scan(&authorizeInfos)

			for _, roleInfo := range roleInfos {
				var tmpAuthorizeInfos []*entity.SysAuthorize
				for _, authorizeInfo := range authorizeInfos {
					if int(roleInfo.Id) == authorizeInfo.RoleId {
						tmpAuthorizeInfos = append(tmpAuthorizeInfos, authorizeInfo)
					}
				}
				if tmpAuthorizeInfos != nil && len(tmpAuthorizeInfos) > 0 {
					err = cache.Instance().Set(ctx, consts.CacheSysAuthorize+"_"+gconv.String(roleInfo.Id), tmpAuthorizeInfos, 0)
					if err != nil {
						return
					}
				}
			}
		}
	}
	return
}
