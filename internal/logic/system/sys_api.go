package system

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/cache"
)

type sSysApi struct {
}

func sysApiNew() *sSysApi {
	return &sSysApi{}
}

func init() {
	service.RegisterSysApi(sysApiNew())
}

// GetInfoByIds 根据接口APIID数组获取接口信息
func (s *sSysApi) GetInfoByIds(ctx context.Context, ids []int) (data []*entity.SysApi, err error) {
	//获取缓存信息
	var tmpData *gvar.Var
	tmpData, err = cache.Instance().Get(ctx, consts.CacheSysApi)
	if err != nil {
		return
	}

	if tmpData == nil {
		err = dao.SysApi.Ctx(ctx).Where(g.Map{
			dao.SysApi.Columns().IsDeleted: 0,
			dao.SysApi.Columns().Status:    1,
		}).Scan(&data)

	}

	var tmpSysApiInfo []*entity.SysApi

	var apiInfo []*entity.SysApi
	//根据菜单ID数组获取菜单列表信息
	if tmpData.Val() != nil {
		strVal, ok := tmpData.Val().(string)
		if !ok {
			err = gerror.New("缓存数据类型错误")
			return
		}
		if err = json.Unmarshal([]byte(strVal), &tmpSysApiInfo); err != nil {
			return
		}

		for _, id := range ids {
			for _, menuTmp := range tmpSysApiInfo {
				if id == int(menuTmp.Id) {
					apiInfo = append(apiInfo, menuTmp)
					continue
				}
			}
		}
	}
	if len(apiInfo) > 0 {
		data = apiInfo
		return
	}
	err = dao.SysApi.Ctx(ctx).Where(g.Map{
		dao.SysApi.Columns().IsDeleted: 0,
		dao.SysApi.Columns().Status:    1,
	}).WhereIn(dao.SysApi.Columns().Id, ids).Scan(&data)
	return
}

// GetApiByMenuId 根据ApiID获取接口信息
func (s *sSysApi) GetApiByMenuId(ctx context.Context, apiId int) (data []*entity.SysApi, err error) {
	var apiApi []*entity.SysMenuApi
	err = dao.SysMenuApi.Ctx(ctx).Where(g.Map{
		dao.SysMenuApi.Columns().MenuId:    apiId,
		dao.SysMenuApi.Columns().IsDeleted: 0,
	}).Scan(&apiApi)
	//获取接口ID数组
	if apiApi != nil {
		var ids []int
		for _, api := range apiApi {
			ids = append(ids, api.ApiId)
		}
		err = dao.SysApi.Ctx(ctx).Where(g.Map{
			dao.SysApi.Columns().IsDeleted: 0,
			dao.SysApi.Columns().Status:    1,
		}).WhereIn(dao.SysApi.Columns().Id, ids).Scan(&data)
	}
	return
}

// GetInfoById 根据ID获取API
func (s *sSysApi) GetInfoById(ctx context.Context, id int) (entity *entity.SysApi, err error) {
	err = dao.SysApi.Ctx(ctx).Where(g.Map{
		dao.SysApi.Columns().Id: id,
	}).Scan(&entity)
	return
}

// GetApiAll 获取所有接口
func (s *sSysApi) GetApiAll(ctx context.Context, method string) (data []*entity.SysApi, err error) {
	m := dao.SysApi.Ctx(ctx)
	if method != "" {
		m = m.Where(dao.SysApi.Columns().Method, method)
	}
	err = m.Where(g.Map{
		dao.SysApi.Columns().IsDeleted: 0,
		dao.SysApi.Columns().Status:    1,
		dao.SysApi.Columns().Types:     2,
	}).Scan(&data)
	if method == "" {
		if data != nil && len(data) > 0 {
			err = cache.Instance().Set(ctx, consts.CacheSysApi, data, 0)
		} else {
			_, err = cache.Instance().Remove(ctx, consts.CacheSysApi)
		}
	}
	return
}

// GetApiTree 获取Api数结构数据
func (s *sSysApi) GetApiTree(ctx context.Context, name string, address string, status int, types int) (out []*model.SysApiTreeOut, err error) {
	var es []*model.SysApiTreeOut
	m := dao.SysApi.Ctx(ctx)
	if name != "" {
		m = m.WhereLike(dao.SysApi.Columns().Name, "%"+name+"%")
	}
	if address != "" {
		m = m.WhereLike(dao.SysApi.Columns().Address, "%"+address+"%")
	}
	if status != -1 {
		m = m.Where(dao.SysApi.Columns().Status, status)
	}
	if types != -1 {
		m = m.Where(dao.SysApi.Columns().Types, types)
	}
	m = m.Where(dao.SysApi.Columns().IsDeleted, 0)

	err = m.OrderAsc(dao.SysApi.Columns().Sort).Scan(&es)
	for _, e := range es {
		menuApiInfo, _ := service.SysMenuApi().GetInfoByApiId(ctx, int(e.Id))
		var menuIds []int
		for _, menuApi := range menuApiInfo {
			menuIds = append(menuIds, menuApi.MenuId)
		}
		e.MenuIds = append(e.MenuIds, menuIds...)
	}

	if len(es) > 0 {
		out, err = GetApiTree(es)
		if err != nil {
			return
		}
	}
	return
}
