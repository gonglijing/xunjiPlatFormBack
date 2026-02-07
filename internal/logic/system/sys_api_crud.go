package system

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/model/do"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/cache"
	"strings"
)

// Add 添加Api列表
func (s *sSysApi) Add(ctx context.Context, input *model.AddApiInput) (err error) {
	if input.Types == 2 {
		if input.Address == "" || len(input.MenuIds) == 0 {
			err = gerror.New("参数错误")
			return
		}
		if input.ParentId == -1 {
			err = gerror.New("接口不能为根节点")
			return
		}
	} else {
		if input.Address != "" || len(input.MenuIds) > 0 {
			err = gerror.New("参数错误")
			return
		}
	}
	if input.ParentId != -1 {
		var parentApiInfo *entity.SysApi
		err = dao.SysApi.Ctx(ctx).Where(g.Map{
			dao.SysApi.Columns().Id: input.ParentId,
		}).Scan(&parentApiInfo)
		if parentApiInfo.IsDeleted != 0 {
			return gerror.New("上级节点已删除，无法新增")
		}
		if parentApiInfo.Status != 1 {
			return gerror.New("上级节点未启用，无法新增")
		}
		if parentApiInfo.Types != 1 {
			return gerror.New("上级节点不是分类，无法新增")
		}
	}
	err = dao.SysUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		var apiInfo *entity.SysApi
		//根据名称查看是否存在
		apiInfo = CheckApiName(ctx, input.Name, 0, input.ParentId)
		if apiInfo != nil {
			return gerror.New("同一个分类下名称不能重复")
		}
		if input.Types == 2 {
			//根据名称查看是否存在
			apiInfo = CheckApiAddress(ctx, input.Address, 0, input.ApiTypes)
			if apiInfo != nil {
				return gerror.New("同一个服务下Api地址,无法添加")
			}
		}
		//获取当前登录用户ID
		loginUserId := service.Context().GetUserId(ctx)
		apiInfo = new(entity.SysApi)
		if apiInfoErr := gconv.Scan(input, &apiInfo); apiInfoErr != nil {
			return
		}
		apiInfo.IsDeleted = 0
		apiInfo.CreatedBy = uint(loginUserId)
		result, err := dao.SysApi.Ctx(ctx).Data(do.SysApi{
			ParentId:  apiInfo.ParentId,
			Name:      apiInfo.Name,
			Types:     apiInfo.Types,
			ApiTypes:  apiInfo.ApiTypes,
			Method:    apiInfo.Method,
			Address:   apiInfo.Address,
			Remark:    apiInfo.Remark,
			Status:    apiInfo.Status,
			Sort:      apiInfo.Sort,
			IsDeleted: apiInfo.IsDeleted,
			CreatedBy: apiInfo.CreatedBy,
			CreatedAt: gtime.Now(),
		}).Insert()
		if err != nil {
			return err
		}

		if input.Types == 2 {
			var lastInsertId int64
			//获取主键ID
			lastInsertId, err = service.Sequences().GetSequences(ctx, result, dao.SysApi.Table(), dao.SysApi.Columns().Id)
			if err != nil {
				return
			}
			//绑定菜单
			err = s.AddMenuApi(ctx, "api", []int{int(lastInsertId)}, input.MenuIds)
			if err != nil {
				return
			}
		}
		//获取所有接口并添加缓存
		_, err = s.GetApiAll(ctx, "")
		return
	})
	return
}

// Detail Api列表详情
func (s *sSysApi) Detail(ctx context.Context, id int) (out *model.SysApiOut, err error) {
	var e *entity.SysApi
	_ = dao.SysApi.Ctx(ctx).Where(g.Map{
		dao.SysApi.Columns().Id: id,
	}).Scan(&e)
	if e != nil {
		if err = gconv.Scan(e, &out); err != nil {
			return nil, err
		}
		menuApiInfo, _ := service.SysMenuApi().GetInfoByApiId(ctx, out.Id)
		var menuIds []int
		for _, menuApi := range menuApiInfo {
			menuIds = append(menuIds, menuApi.MenuId)
		}
		out.MenuIds = append(out.MenuIds, menuIds...)
	}
	return
}

func (s *sSysApi) AddMenuApi(ctx context.Context, addPageSource string, apiIds []int, menuIds []int) (err error) {
	loginUserId := service.Context().GetUserId(ctx)

	if addPageSource != "" && strings.EqualFold(addPageSource, "api") {
		//解除旧绑定关系
		_, err = dao.SysMenuApi.Ctx(ctx).Data(g.Map{
			dao.SysMenuApi.Columns().IsDeleted: 1,
			dao.SysMenuApi.Columns().DeletedBy: loginUserId,
			dao.SysMenuApi.Columns().DeletedAt: gtime.Now(),
		}).WhereIn(dao.SysMenuApi.Columns().ApiId, apiIds).Update()
	}
	if addPageSource != "" && strings.EqualFold(addPageSource, "menu") {
		//解除旧绑定关系
		_, err = dao.SysMenuApi.Ctx(ctx).Data(g.Map{
			dao.SysMenuApi.Columns().IsDeleted: 1,
			dao.SysMenuApi.Columns().DeletedBy: loginUserId,
			dao.SysMenuApi.Columns().DeletedAt: gtime.Now(),
		}).WhereIn(dao.SysMenuApi.Columns().MenuId, menuIds).Update()
	}

	if menuIds != nil && len(menuIds) > 0 && apiIds != nil && len(apiIds) > 0 {
		//添加菜单
		var sysMenuApis []*entity.SysMenuApi
		for _, menuId := range menuIds {
			var menuInfo *entity.SysMenu
			err = dao.SysMenu.Ctx(ctx).Where(dao.SysMenu.Columns().Id, menuId).Scan(&menuInfo)
			if menuInfo == nil {
				err = gerror.New("菜单ID错误")
				return
			}
			if menuInfo != nil && menuInfo.IsDeleted == 1 {
				err = gerror.New(menuInfo.Name + "已删除,无法绑定")
				return
			}
			if menuInfo != nil && menuInfo.Status == 0 {
				err = gerror.New(menuInfo.Name + "已禁用,无法绑定")
				return
			}

			for _, id := range apiIds {
				apiInfo, _ := s.Detail(ctx, id)
				if apiInfo == nil {
					return gerror.New("API接口不存在")
				}
				if apiInfo.Types != 2 {
					return gerror.New("参数错误")
				}
				var sysMenuApi = new(entity.SysMenuApi)
				sysMenuApi.MenuId = menuId
				sysMenuApi.ApiId = id
				sysMenuApi.IsDeleted = 0
				sysMenuApi.CreatedBy = uint(loginUserId)
				sysMenuApis = append(sysMenuApis, sysMenuApi)
			}
		}
		if sysMenuApis != nil {
			//添加
			var sysMenuApisInput []*model.SysMenuApiInput
			if err = gconv.Scan(sysMenuApis, &sysMenuApisInput); err != nil {
				return
			}

			_, addErr := dao.SysMenuApi.Ctx(ctx).Data(sysMenuApisInput).Insert()
			if addErr != nil {
				err = gerror.New("添加失败")
				return
			}
			//查询菜单ID绑定的所有接口ID
			var menuApiInfos []*entity.SysMenuApi
			err = dao.SysMenuApi.Ctx(ctx).Where(g.Map{
				dao.SysMenuApi.Columns().IsDeleted: 0,
			}).WhereIn(dao.SysMenuApi.Columns().MenuId, menuIds).Scan(&menuApiInfos)
			//添加缓存
			for _, menuId := range menuIds {
				var menuApi []*entity.SysMenuApi
				for _, menuApiInfo := range menuApiInfos {
					if menuId == menuApiInfo.MenuId {
						menuApi = append(menuApi, menuApiInfo)
					}
				}
				if menuApi != nil && len(menuApi) > 0 {
					err := cache.Instance().Set(ctx, consts.CacheSysMenuApi+"_"+gconv.String(menuId), menuApi, 0)
					if err != nil {
						return err
					}
				}

			}
			//获取所有信息
			var menuApiInfoAll []*entity.SysMenuApi
			err = dao.SysMenuApi.Ctx(ctx).Where(g.Map{
				dao.SysMenuApi.Columns().IsDeleted: 0,
			}).WhereIn(dao.SysMenuApi.Columns().MenuId, menuIds).Scan(&menuApiInfoAll)
			if menuApiInfoAll != nil && len(menuApiInfoAll) > 0 {
				err := cache.Instance().Set(ctx, consts.CacheSysMenuApi, menuApiInfoAll, 0)
				if err != nil {
					return err
				}
			}
		}
	}

	return
}

// Edit 修改Api列表
func (s *sSysApi) Edit(ctx context.Context, input *model.EditApiInput) (err error) {
	if input.Types == 2 {
		if input.Address == "" || len(input.MenuIds) == 0 {
			err = gerror.New("参数错误")
			return
		}
		if input.ParentId == -1 {
			err = gerror.New("接口不能为根节点")
			return
		}
	} else {
		if input.Address != "" || len(input.MenuIds) > 0 {
			err = gerror.New("参数错误")
			return
		}
	}
	if input.ParentId != -1 {
		var parentApiInfo *entity.SysApi
		err = dao.SysApi.Ctx(ctx).Where(g.Map{
			dao.SysApi.Columns().Id: input.ParentId,
		}).Scan(&parentApiInfo)
		if parentApiInfo.IsDeleted != 0 {
			return gerror.New("上级节点已删除，无法新增")
		}
		if parentApiInfo.Status != 1 {
			return gerror.New("上级节点已启用，无法新增")
		}
		if parentApiInfo.Types != 1 {
			return gerror.New("上级节点不是分类，无法新增")
		}
	}
	err = dao.SysUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		var apiInfo, apiInfo2 *entity.SysApi
		//根据ID查看Api列表是否存在
		apiInfo = CheckApiId(ctx, input.Id, apiInfo)
		if apiInfo == nil {
			return gerror.New("Api列表不存在")
		}
		apiInfo2 = CheckApiName(ctx, input.Name, input.Id, input.ParentId)
		if apiInfo2 != nil {
			return gerror.New("同一个分类下名称不能重复")
		}
		if input.Types == 2 {
			apiInfo2 = CheckApiAddress(ctx, input.Address, input.Id, input.ApiTypes)
			if apiInfo2 != nil {
				return gerror.New("Api地址已存在,无法修改")
			}
		}
		//获取当前登录用户ID
		loginUserId := service.Context().GetUserId(ctx)
		if apiInfoErr := gconv.Scan(input, &apiInfo); apiInfoErr != nil {
			return
		}
		apiInfo.UpdatedBy = uint(loginUserId)
		_, err = dao.SysApi.Ctx(ctx).Data(apiInfo).
			Where(dao.SysApi.Columns().Id, input.Id).Update()
		if err != nil {
			return gerror.New("修改失败")
		}
		//绑定菜单
		err = s.AddMenuApi(ctx, "api", []int{input.Id}, input.MenuIds)

		//获取所有接口并添加缓存
		_, err = s.GetApiAll(ctx, "")
		return
	})

	return
}

// Del 根据ID删除Api列表信息
func (s *sSysApi) Del(ctx context.Context, Id int) (err error) {
	var apiColumn *entity.SysApi
	_ = dao.SysApi.Ctx(ctx).Where(g.Map{
		dao.SysApi.Columns().Id: Id,
	}).Scan(&apiColumn)
	if apiColumn == nil {
		return gerror.New("ID错误")
	}
	//判断是否为分类
	if apiColumn.Types == 1 {
		//判断是否存在子节点
		num, _ := dao.SysApi.Ctx(ctx).Where(g.Map{
			dao.SysApi.Columns().ParentId:  Id,
			dao.SysApi.Columns().IsDeleted: 0,
		}).Count()
		if num > 0 {
			return gerror.New("存在子节点，无法删除")
		}
	}
	loginUserId := service.Context().GetUserId(ctx)
	//获取当前时间
	time, err := gtime.StrToTimeFormat(gtime.Datetime(), "2006-01-02 15:04:05")
	if err != nil {
		return
	}
	//开启事务管理
	err = dao.SysApi.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		//更新Api列表信息
		_, err = dao.SysApi.Ctx(ctx).
			Data(g.Map{
				dao.SysApi.Columns().DeletedBy: uint(loginUserId),
				dao.SysApi.Columns().IsDeleted: 1,
				dao.SysApi.Columns().DeletedAt: time,
			}).Where(dao.SysApi.Columns().Id, Id).
			Update()
		//删除于菜单关系绑定
		_, err = dao.SysMenuApi.Ctx(ctx).Data(g.Map{
			dao.SysMenuApi.Columns().IsDeleted: 1,
			dao.SysMenuApi.Columns().DeletedBy: loginUserId,
			dao.SysMenuApi.Columns().DeletedAt: time,
		}).Where(dao.SysMenuApi.Columns().ApiId, Id).Update()

		//获取所有接口并添加缓存
		_, err = s.GetApiAll(ctx, "")

		return
	})
	return
}

// EditStatus 修改状态
func (s *sSysApi) EditStatus(ctx context.Context, id int, status int) (err error) {
	var apiInfo *entity.SysApi
	_ = dao.SysApi.Ctx(ctx).Where(g.Map{
		dao.SysApi.Columns().Id: id,
	}).Scan(&apiInfo)
	if apiInfo == nil {
		return gerror.New("ID错误")
	}
	if apiInfo != nil && apiInfo.IsDeleted == 1 {
		return gerror.New("列表字段已删除,无法修改")
	}
	if apiInfo != nil && apiInfo.Status == status {
		return gerror.New("API已禁用或启用,无须重复修改")
	}
	loginUserId := service.Context().GetUserId(ctx)
	apiInfo.Status = status
	apiInfo.UpdatedBy = uint(loginUserId)

	_, err = dao.SysApi.Ctx(ctx).Data(apiInfo).Where(g.Map{
		dao.SysApi.Columns().Id: id,
	}).Update()

	//获取所有接口并添加缓存
	_, err = s.GetApiAll(ctx, "")

	return
}
