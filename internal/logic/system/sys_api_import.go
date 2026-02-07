package system

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gtime"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/internal/model/do"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"strings"
)

// GetInfoByAddress 根据Address获取API
func (s *sSysApi) GetInfoByAddress(ctx context.Context, address string) (entity *entity.SysApi, err error) {
	err = dao.SysApi.Ctx(ctx).Where(g.Map{
		dao.SysApi.Columns().Address:   address,
		dao.SysApi.Columns().IsDeleted: 0,
		dao.SysApi.Columns().Status:    1,
	}).Scan(&entity)
	return
}

// CheckApiId 检查指定ID的数据是否存在
func CheckApiId(ctx context.Context, Id int, apiColumn *entity.SysApi) *entity.SysApi {
	_ = dao.SysApi.Ctx(ctx).Where(g.Map{
		dao.SysApi.Columns().Id:        Id,
		dao.SysApi.Columns().IsDeleted: 0,
	}).Scan(&apiColumn)
	return apiColumn
}

// CheckApiName 检查相同Api名称的数据是否存在
func CheckApiName(ctx context.Context, name string, tag int, parentId int) *entity.SysApi {
	var apiInfo *entity.SysApi
	m := dao.SysApi.Ctx(ctx)
	if tag > 0 {
		m = m.WhereNot(dao.SysApi.Columns().Id, tag)
	}
	m = m.Where(dao.SysApi.Columns().ParentId, parentId)
	_ = m.Where(g.Map{
		dao.SysApi.Columns().Name:      name,
		dao.SysApi.Columns().IsDeleted: 0,
	}).Scan(&apiInfo)
	return apiInfo
}

// CheckApiAddress 检查相同Api地址的数据是否存在
func CheckApiAddress(ctx context.Context, address string, tag int, apiTypes string) *entity.SysApi {
	var apiInfo *entity.SysApi
	m := dao.SysApi.Ctx(ctx)
	if tag > 0 {
		m = m.WhereNot(dao.SysApi.Columns().Id, tag)
	}
	m = m.Where(dao.SysApi.Columns().ApiTypes, apiTypes)
	_ = m.Where(g.Map{
		dao.SysApi.Columns().Address:   address,
		dao.SysApi.Columns().IsDeleted: 0,
	}).Scan(&apiInfo)
	return apiInfo
}

// GetInfoByNameAndTypes 根据名字和类型获取API
func (s *sSysApi) GetInfoByNameAndTypes(ctx context.Context, name string, types int) (entity *entity.SysApi, err error) {
	err = dao.SysApi.Ctx(ctx).Where(g.Map{
		dao.SysApi.Columns().Name:  name,
		dao.SysApi.Columns().Types: types,
	}).Scan(&entity)
	return
}

// ImportApiFile 导入API文件
func (s *sSysApi) ImportApiFile(ctx context.Context) (err error) {
	//获取服务端口号
	port := g.Cfg().MustGet(ctx, "server.address").String()
	if port == "" {
		err = gerror.New("服务地址不能为空")
		return
	}

	//获取API路径
	openApiPath := g.Cfg().MustGet(ctx, "server.openapiPath").String()
	if openApiPath == "" {
		err = gerror.New("openApi路径不能为空")
		return
	}

	url := "http://127.0.0.1" + port + openApiPath
	resp, err := g.Client().Get(ctx, url)
	if err != nil {
		return
	}
	respContent := resp.ReadAllString()
	status := resp.Status
	defer func(resp *gclient.Response) {
		if err = resp.Close(); err != nil {
			g.Log().Error(ctx, err)
		}
	}(resp)
	if !strings.Contains(status, "200") {
		err = gerror.New("请求失败,请联系管理员")
		return
	}
	var apiJsonContent map[string]interface{}
	if respContent != "" {
		err = json.Unmarshal([]byte(respContent), &apiJsonContent)
		if err != nil {
			return
		}
	}
	//封装数据、整理入库
	//获取所有的接口
	var address []string
	var tags []string
	var apiPathsAll = make(map[string][]map[string]string)
	paths := apiJsonContent["paths"].(map[string]interface{})
	for key, value := range paths {
		//接口信息
		apiInfo := make(map[string]string)
		//获取对应的分组
		valueMap := value.(map[string]interface{})
		var valueKey string
		for valueMapKey := range valueMap {
			if !strings.EqualFold(valueMapKey, "summary") {
				valueKey = valueMapKey
				break
			}
		}
		apiInfo["address"] = key
		address = append(address, key)
		//方法
		apiInfo["method"] = valueKey
		//名字
		apiInfo["name"] = valueMap[valueKey].(map[string]interface{})["summary"].(string)

		tags := valueMap[valueKey].(map[string]interface{})["tags"].([]interface{})
		fmt.Printf("key:%s---tags:%s\n", key, tags)
		for _, tag := range tags {
			if _, ok := apiPathsAll[tag.(string)]; !ok {
				apiPathsAll[tag.(string)] = []map[string]string{apiInfo}
			} else {
				apiPathsAll[tag.(string)] = append(apiPathsAll[tag.(string)], apiInfo)
			}
		}
	}
	g.Log().Debugf(ctx, "全部API信息:%s", apiPathsAll)

	//获取数据字典类型
	dictData, err := service.DictData().GetDictDataByType(ctx, consts.ApiTypes)
	if err != nil {
		return err
	}
	apiTpyes := "IOT"
	if dictData.Data != nil && dictData.Values != nil {
		for _, value := range dictData.Values {
			if strings.EqualFold(value.DictLabel, "sagoo_iot") {
				apiTpyes = value.DictValue
			}
		}
	}
	err = dao.SysApi.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		//开始入库
		for key, value := range apiPathsAll {
			tags = append(tags, key)
			//获取分组ID
			var categoryId int
			//判断分组是否存在,不存在则创建分组
			categoryInfo, _ := s.GetInfoByNameAndTypes(ctx, key, 1)
			if categoryInfo == nil {
				//创建分组
				var result sql.Result
				result, err = dao.SysApi.Ctx(ctx).Data(do.SysApi{
					ParentId:  -1,
					Name:      key,
					Types:     1,
					Status:    1,
					IsDeleted: 0,
				}).Insert()
				if err != nil {
					return err
				}
				//获取主键ID
				var lastInsertId int64
				lastInsertId, err = service.Sequences().GetSequences(ctx, result, dao.SysApi.Table(), dao.SysApi.Columns().Id)
				if err != nil {
					return
				}
				categoryId = int(lastInsertId)
			} else {
				categoryId = int(categoryInfo.Id)
			}

			for _, apiValue := range value {
				if strings.Contains(apiValue["address"], "openapi") {
					continue
				}
				//判断接口是否存在
				var sysApi *entity.SysApi
				sysApi, err = s.GetInfoByAddress(ctx, apiValue["address"])
				if err != nil {
					return
				}
				if sysApi != nil {
					_, err = dao.SysApi.Ctx(ctx).Data(do.SysApi{
						ParentId:  categoryId,
						Name:      apiValue["name"],
						Types:     2,
						ApiTypes:  apiTpyes,
						Method:    apiValue["method"],
						Address:   apiValue["address"],
						Status:    1,
						IsDeleted: 0,
					}).Where(dao.SysApi.Columns().Address, apiValue["address"]).Update()
					if err != nil {
						return err
					}
				} else {
					_, err = dao.SysApi.Ctx(ctx).Data(do.SysApi{
						ParentId:  categoryId,
						Name:      apiValue["name"],
						Types:     2,
						ApiTypes:  apiTpyes,
						Method:    apiValue["method"],
						Address:   apiValue["address"],
						Status:    1,
						IsDeleted: 0,
					}).Insert()
					if err != nil {
						return
					}
				}
			}
		}

		//根据地址获取不存在的APIID
		var sysApiInfos []*entity.SysApi
		err = dao.SysApi.Ctx(ctx).Where(dao.SysApi.Columns().ApiTypes, apiTpyes).WhereNotIn(dao.SysApi.Columns().Address, address).Scan(&sysApiInfos)
		if err != nil {
			return
		}
		var apiIds []uint
		for _, info := range sysApiInfos {
			apiIds = append(apiIds, info.Id)
		}
		//删除API
		_, err = dao.SysApi.Ctx(ctx).Data(g.Map{
			dao.SysApi.Columns().IsDeleted: 1,
			dao.SysApi.Columns().Status:    0,
			dao.SysApi.Columns().DeletedAt: gtime.Now(),
			dao.SysApi.Columns().DeletedBy: service.Context().GetUserId(ctx),
		}).WhereIn(dao.SysApi.Columns().Id, apiIds).Update()

		if err != nil {
			return
		}
		//删除于菜单关系绑定
		_, err = dao.SysMenuApi.Ctx(ctx).Data(g.Map{
			dao.SysMenuApi.Columns().IsDeleted: 1,
			dao.SysMenuApi.Columns().DeletedAt: gtime.Now(),
			dao.SysMenuApi.Columns().DeletedBy: service.Context().GetUserId(ctx),
		}).WhereIn(dao.SysMenuApi.Columns().Id, apiIds).Update()
		if err != nil {
			return
		}

		//删除不存在的分类
		_, err = dao.SysApi.Ctx(ctx).Data(g.Map{
			dao.SysApi.Columns().IsDeleted: 1,
			dao.SysApi.Columns().Status:    0,
			dao.SysApi.Columns().DeletedAt: gtime.Now(),
			dao.SysApi.Columns().DeletedBy: service.Context().GetUserId(ctx),
		}).WhereNotIn(dao.SysApi.Columns().Name, tags).Where(dao.SysApi.Columns().Types, 1).Update()

		return
	})
	return
}
