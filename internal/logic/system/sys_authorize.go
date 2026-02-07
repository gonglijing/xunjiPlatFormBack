package system

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/cache"
	"strings"
)

type sSysAuthorize struct {
}

func init() {
	service.RegisterSysAuthorize(sysAuthorizeNew())
}

func sysAuthorizeNew() *sSysAuthorize {
	return &sSysAuthorize{}
}

func (s *sSysAuthorize) AuthorizeQuery(ctx context.Context, itemsType string, menuIds []int) (out []*model.AuthorizeQueryTreeOut, err error) {
	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)
	if loginUserId == 0 {
		err = gerror.New("未登录或TOKEN失效,请重新登录")
		return
	}
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
	if isSuperAdmin {
		var userMenuTreeOut []*model.UserMenuTreeOut
		if strings.EqualFold(itemsType, consts.Menu) {
			userMenuTreeOut, _ = GetMenuInfo(ctx, nil)
		} else {
			userMenuTreeOut, _ = GetMenuInfo(ctx, menuIds)
		}
		if userMenuTreeOut == nil {
			err = gerror.New("无可用菜单")
			return
		}

		var authorizeQueryTreeRes []*model.AuthorizeQueryTreeOut
		if err = gconv.Scan(userMenuTreeOut, &authorizeQueryTreeRes); err != nil {
			return
		}

		if !strings.EqualFold(itemsType, consts.Menu) {
			if len(menuIds) == 0 {
				err = gerror.New("请选择菜单")
				return
			}
			//根据项目类型 菜单ID封装菜单的按钮，列表字段,API接口
			authorizeItemsTypeTreeOut, userItemsTypeTreeErr := GetAuthorizeItemsTypeTreeOut(ctx, menuIds, itemsType, authorizeQueryTreeRes)
			if userItemsTypeTreeErr != nil {
				return nil, userItemsTypeTreeErr
			}

			//获取所有的子节点
			/*childrenMenuTreeOut := GetAllAuthorizeQueryChildrenTree(authorizeItemsTypeTreeOut)
			if childrenMenuTreeOut != nil {
				//查找所有的上级ID
				childrenMenuTreeOut = GetAllAuthorizeQueryParentTree(childrenMenuTreeOut, authorizeItemsTypeTreeOut)
			}*/

			//out = GetAuthorizeMenuTree(childrenMenuTreeOut)
			out = authorizeItemsTypeTreeOut
			return
		}
		out = GetAuthorizeMenuTree(authorizeQueryTreeRes)

		return
	} else {
		var userMenuTreeOut []*model.UserMenuTreeOut
		if strings.EqualFold(itemsType, consts.Menu) {
			//根据角色ID获取角色下配置的菜单信息
			authorizeItemsInfo, authorizeItemsErr := service.SysAuthorize().GetInfoByRoleIdsAndItemsType(ctx, roleIds, consts.Menu)
			if authorizeItemsErr != nil {
				return nil, authorizeItemsErr
			}
			if authorizeItemsInfo == nil {
				err = gerror.New("菜单未配置,请联系管理员")
				return
			}
			//获取菜单ID
			var authorizeItemsIds []int
			for _, authorizeItems := range authorizeItemsInfo {
				authorizeItemsIds = append(authorizeItemsIds, authorizeItems.ItemsId)
			}
			//根据菜单ID获取菜单信息
			userMenuTreeOut, _ = GetMenuInfo(ctx, authorizeItemsIds)
		} else {
			userMenuTreeOut, _ = GetMenuInfo(ctx, menuIds)
		}
		if userMenuTreeOut == nil {
			err = gerror.New("无可用菜单")
			return
		}
		var authorizeQueryTreeOut []*model.AuthorizeQueryTreeOut
		if err = gconv.Scan(userMenuTreeOut, &authorizeQueryTreeOut); err != nil {
			return
		}

		//根据项目类型 菜单ID封装菜单的按钮，列表字段,API接口
		if !strings.EqualFold(itemsType, consts.Menu) {
			if len(menuIds) == 0 {
				err = gerror.New("请选择菜单")
				return
			}
			//根据项目类型 菜单ID封装菜单的按钮，列表字段,API接口
			authorizeItemsTypeTreeOut, userItemsTypeTreeErr := GetAuthorizeItemsTypeTreeOut(ctx, menuIds, itemsType, authorizeQueryTreeOut)
			if userItemsTypeTreeErr != nil {
				return nil, userItemsTypeTreeErr
			}

			//获取所有的子节点
			/*childrenMenuTreeOut := GetAllAuthorizeQueryChildrenTree(authorizeItemsTypeTreeOut)
			if childrenMenuTreeOut != nil {
				//查找所有的上级ID
				childrenMenuTreeOut = GetAllAuthorizeQueryParentTree(childrenMenuTreeOut, authorizeItemsTypeTreeOut)
			}*/

			//out = GetAuthorizeMenuTree(childrenMenuTreeOut)
			out = authorizeItemsTypeTreeOut
			return
		}
		out = GetAuthorizeMenuTree(authorizeQueryTreeOut)
		return
	}
}

// GetInfoByRoleId 根据角色ID获取权限信息
func (s *sSysAuthorize) GetInfoByRoleId(ctx context.Context, roleId int) (data []*entity.SysAuthorize, err error) {
	err = dao.SysAuthorize.Ctx(ctx).Where(g.Map{
		dao.SysAuthorize.Columns().RoleId:    roleId,
		dao.SysAuthorize.Columns().IsDeleted: 0,
	}).Scan(&data)
	return
}

// GetInfoByRoleIds 根据角色ID数组获取权限信息
func (s *sSysAuthorize) GetInfoByRoleIds(ctx context.Context, roleIds []int) (data []*entity.SysAuthorize, err error) {
	//获取缓存菜单按钮信息
	for _, v := range roleIds {
		var tmpData *gvar.Var
		tmpData, err = cache.Instance().Get(ctx, consts.CacheSysAuthorize+"_"+gconv.String(v))
		if err != nil {
			return
		}
		if tmpData.Val() != nil {
			strVal, ok := tmpData.Val().(string)
			if !ok {
				continue
			}
			var sysAuthorizeInfo []*entity.SysAuthorize
			err = json.Unmarshal([]byte(strVal), &sysAuthorizeInfo)
			if err != nil {
				return
			}
			data = append(data, sysAuthorizeInfo...)
		}
	}
	if len(data) == 0 {
		err = dao.SysAuthorize.Ctx(ctx).Where(g.Map{
			dao.SysAuthorize.Columns().IsDeleted: 0,
		}).WhereIn(dao.SysAuthorize.Columns().RoleId, roleIds).Scan(&data)
	}
	return
}

// GetInfoByRoleIdsAndItemsType 根据角色ID和项目类型获取权限信息
func (s *sSysAuthorize) GetInfoByRoleIdsAndItemsType(ctx context.Context, roleIds []int, itemsType string) (data []*entity.SysAuthorize, err error) {
	err = dao.SysAuthorize.Ctx(ctx).Where(g.Map{
		dao.SysAuthorize.Columns().IsDeleted: 0,
		dao.SysAuthorize.Columns().ItemsType: itemsType,
	}).WhereIn(dao.SysAuthorize.Columns().RoleId, roleIds).Scan(&data)
	return
}

func (s *sSysAuthorize) DelByRoleId(ctx context.Context, roleId int) (err error) {
	loginUserId := service.Context().GetUserId(ctx)
	_, err = dao.SysAuthorize.Ctx(ctx).Data(g.Map{
		dao.SysAuthorize.Columns().DeletedBy: uint(loginUserId),
		dao.SysAuthorize.Columns().IsDeleted: 1,
	}).Where(dao.SysAuthorize.Columns().RoleId, roleId).Update()
	_, err = dao.SysAuthorize.Ctx(ctx).Where(dao.SysAuthorize.Columns().RoleId, roleId).Delete()
	return
}

func (s *sSysAuthorize) Add(ctx context.Context, authorize []*entity.SysAuthorize) (err error) {
	var input []*model.SysAuthorizeInput
	if err = gconv.Scan(authorize, &input); err != nil {
		return
	}
	_, err = dao.SysAuthorize.Ctx(ctx).Data(input).Insert()
	return
}
