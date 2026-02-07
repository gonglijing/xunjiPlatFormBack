package product

import (
	"context"
	"errors"
	"fmt"
	"sagooiot/internal/consts"
	"sagooiot/internal/model"
	"sagooiot/pkg/dcache"
	"sagooiot/pkg/tsd"
	"sagooiot/pkg/tsd/comm"
	"strings"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RunStatus 运行状态
func (s *sDevDevice) RunStatus(ctx context.Context, deviceKey string) (out *model.DeviceRunStatusOutput, err error) {

	device, err := dcache.GetDeviceDetailInfo(deviceKey)
	if err != nil {
		return nil, errors.New("设备不存在")
	}
	out = new(model.DeviceRunStatusOutput)
	out.Status = dcache.GetDeviceStatus(ctx, deviceKey)

	//获取数据
	deviceValueList := dcache.GetDeviceDetailData(context.Background(), deviceKey, consts.MsgTypePropertyReport, consts.MsgTypeGatewayBatch)
	if len(deviceValueList) == 0 {
		return
	}

	for _, d := range deviceValueList[0] {
		out.LastOnlineTime = gtime.New(d.CreateTime)
	}

	var properties []model.DevicePropertiy
	for _, v := range device.TSL.Properties {
		unit := ""
		if v.ValueType.Unit != nil {
			unit = *v.ValueType.Unit
		}

		var valueList = make([]*g.Var, 0)
		for _, d := range deviceValueList {
			valueList = append(valueList, g.NewVar(d[v.Key].Value))
		}

		pro := model.DevicePropertiy{
			Key:   v.Key,
			Name:  v.Name,
			Type:  v.ValueType.Type,
			Unit:  unit,
			Value: valueList[0],
			List:  valueList,
		}
		properties = append(properties, pro)
	}
	out.Properties = properties
	return
}

// GetLatestProperty 获取设备最新的属性值
func (s *sDevDevice) GetLatestProperty(ctx context.Context, key string) (list []model.DeviceLatestProperty, err error) {
	p, err := s.Get(ctx, key)
	if err != nil {
		return
	}
	if p.Status == model.DeviceStatusNoEnable {
		return
	}

	deviceTable := comm.DeviceTableName(p.Key)

	tsdDb := tsd.DB()
	defer tsdDb.Close()

	for _, v := range p.TSL.Properties {
		ckey := comm.TsdColumnName(v.Key)

		// 获取属性最近有效值
		sql := "select ? from ? where ? is not null order by ts desc limit 1"
		rs, err := tsdDb.GetTableDataOne(ctx, sql, ckey, deviceTable, ckey)
		if err != nil {
			return nil, err
		}
		value := rs[strings.ToLower(v.Key)]
		if value.IsEmpty() {
			continue
		}

		unit := ""
		if v.ValueType.Unit != nil {
			unit = *v.ValueType.Unit
		}

		pro := model.DeviceLatestProperty{
			Key:   v.Key,
			Name:  v.Name,
			Type:  v.ValueType.Type,
			Unit:  unit,
			Value: value,
		}
		list = append(list, pro)
	}
	return
}

// GetProperty 获取指定属性值
func (s *sDevDevice) GetProperty(ctx context.Context, in *model.DeviceGetPropertyInput) (out *model.DevicePropertiy, err error) {
	p, err := s.Detail(ctx, in.DeviceKey)
	if err != nil {
		return
	}
	if p.Status == model.DeviceStatusNoEnable {
		err = errors.New("设备未启用")
		return
	}

	tsdDb := tsd.DB()
	defer tsdDb.Close()

	sKey := in.PropertyKey
	in.PropertyKey = strings.ToLower(in.PropertyKey)
	col := comm.TsdColumnName(in.PropertyKey)

	deviceTable := comm.DeviceTableName(p.Key)

	// 属性上报时间
	ctime := in.PropertyKey + "_time"
	ctime = comm.TsdColumnName(ctime)

	// 属性值获取
	sql := "select ? from ? where ? is not null order by ? desc limit 1"
	rs, err := tsdDb.GetTableDataOne(ctx, sql, col, deviceTable, col, ctime)
	if err != nil {
		return
	}

	var name string
	var valueType string
	for _, v := range p.TSL.Properties {
		if strings.ToLower(v.Key) == in.PropertyKey {
			name = v.Name
			valueType = v.ValueType.Type
			break
		}
	}

	out = new(model.DevicePropertiy)
	out.Key = sKey
	out.Name = name
	out.Type = valueType
	out.Value = rs[in.PropertyKey]

	// 获取当天属性值列表
	sql = "select ? from ? where ? >= '?' order by ? desc"
	ls, _ := tsdDb.GetTableDataAll(ctx, sql, col, deviceTable, ctime, gtime.Now().Format("Y-m-d"), ctime)
	out.List = ls.Array(in.PropertyKey)

	return
}

// GetPropertyList 设备属性详情列表
func (s *sDevDevice) GetPropertyList(ctx context.Context, in *model.DeviceGetPropertyListInput) (out *model.DeviceGetPropertyListOutput, err error) {
	resultList, total, currentPage := dcache.GetDeviceDetailDataByPage(ctx, in.DeviceKey, in.PageNum, in.PageSize, consts.MsgTypePropertyReport, consts.MsgTypeGatewayBatch)
	out = new(model.DeviceGetPropertyListOutput)
	out.Total = total
	out.CurrentPage = currentPage

	var pro []*model.DevicePropertiyOut
	for _, d := range resultList {
		pro = append(pro, &model.DevicePropertiyOut{
			Ts:    gtime.New(d[in.PropertyKey].CreateTime),
			Value: gvar.New(d[in.PropertyKey].Value),
		})
	}
	out.List = pro

	return
}

// GetData 获取设备指定日期属性数据
func (s *sDevDevice) GetData(ctx context.Context, in *model.DeviceGetDataInput) (list []model.DevicePropertiyOut, err error) {
	p, err := s.Detail(ctx, in.DeviceKey)
	if err != nil {
		return
	}
	if p.Status == model.DeviceStatusNoEnable {
		err = errors.New("设备未启用")
		return
	}

	deviceTable := comm.DeviceTableName(p.Key)

	in.PropertyKey = strings.ToLower(in.PropertyKey)
	col := comm.TsdColumnName(in.PropertyKey)

	// 属性上报时间
	ctime := in.PropertyKey + "_time"
	ctime = comm.TsdColumnName(ctime)

	where := fmt.Sprintf("%s >= %q", ctime, gtime.Now().Format("Y-m-d"))
	if len(in.DateRange) > 1 {
		where = fmt.Sprintf("%s >= %q and %s <= %q", ctime, in.DateRange[0], ctime, in.DateRange[1])
	}

	desc := "asc"
	if in.IsDesc == 1 {
		desc = "desc"
	}

	tsdDb := tsd.DB()
	defer tsdDb.Close()

	// TDengine
	sql := "select ?, ? from ? where ? order by ? ?"
	ls, _ := tsdDb.GetTableDataAll(
		ctx,
		sql,
		fmt.Sprintf("distinct %s as ts", ctime),
		col,
		deviceTable,
		where,
		ctime,
		desc,
	)
	for _, v := range ls.List() {
		list = append(list, model.DevicePropertiyOut{
			Ts:    gvar.New(v["ts"]).GTime(),
			Value: gvar.New(v[in.PropertyKey]),
		})
	}
	return
}
