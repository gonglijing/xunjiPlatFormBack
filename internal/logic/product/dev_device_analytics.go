package product

import (
	"context"
	"fmt"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/model/entity"
	"sagooiot/pkg/tsd"
	"sagooiot/pkg/tsd/comm"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// initMonthlyMap creates a map with 12 months initialized to 0.
func initMonthlyMap() map[int]int {
	data := make(map[int]int, 12)
	for i := 1; i <= 12; i++ {
		data[i] = 0
	}
	return data
}

// initDailyMap creates a map with the last 30 days initialized to 0.
func initDailyMap() map[string]int {
	start := time.Now().AddDate(0, -1, 1)
	data := make(map[string]int)
	for i := start; i.Before(time.Now()); i = i.Add(24 * time.Hour) {
		data[i.Format("2006-01-02")] = 0
	}
	return data
}

// getAllDeviceKeys returns quoted device keys for SQL IN clauses.
func getAllDeviceKeys(devices []*entity.DevDevice) []string {
	var keys []string
	for _, device := range devices {
		keys = append(keys, "'"+device.Key+"'")
	}
	return keys
}

// TotalByProductKey 统计产品下的设备数量
func (s *sDevDevice) TotalByProductKey(ctx context.Context, productKeys []string) (totals map[string]int, err error) {
	m := dao.DevDevice.Ctx(ctx)
	r, err := m.Fields(dao.DevDevice.Columns().ProductKey+", count(*) as total").
		WhereIn(dao.DevDevice.Columns().ProductKey, productKeys).
		Group(dao.DevDevice.Columns().ProductKey).
		All()
	if err != nil || r.Len() == 0 {
		return
	}

	totals = make(map[string]int, r.Len())
	for _, v := range r {
		t := gconv.Int(v["total"])
		productKey := gconv.String(v[dao.DevDevice.Columns().ProductKey])
		totals[productKey] = t
	}

	for _, key := range productKeys {
		if _, ok := totals[key]; !ok {
			totals[key] = 0
		}
	}

	return
}

// getTotalForMonthsData 获取统计设备月度消息数量
func (s *sDevDevice) getTotalForMonthsData(ctx context.Context) (data map[int]int, err error) {

	data = initMonthlyMap()
	devices, err := s.GetAll(ctx)
	if err != nil {
		return
	}
	if devices != nil {
		deviceKeys := getAllDeviceKeys(devices)

		tsdDb := tsd.DB()
		defer tsdDb.Close()

		// 设备消息总量 TDengine
		var list gdb.Result
		sql := "select substr(to_iso8601(ts), 1, 7) as ym, count(*) as num from device_log where device in (?) and  ts >= '?' group by substr(to_iso8601(ts), 1, 7)"
		list, err = tsdDb.GetTableDataAll(ctx, sql, strings.Join(deviceKeys, ","), gtime.Now().Format("Y-01-01 00:00:00"))
		if err != nil {
			return
		}

		for _, v := range list {
			m := gstr.SubStr(v["ym"].String(), 5)
			data[gconv.Int(m)] = v["num"].Int()
		}
	}
	return
}

// getTotalForDayData 统计设备最近一月消息数量
func (s *sDevDevice) getTotalForDayData(ctx context.Context) (data map[string]int, err error) {

	data = initDailyMap()

	devices, err := s.GetAll(ctx)
	if err != nil {
		return
	}
	if devices != nil {
		deviceKeys := getAllDeviceKeys(devices)

		tsdDb := tsd.DB()
		defer tsdDb.Close()

		// 设备消息总量 TDengine
		var list gdb.Result
		sql := "select substr(to_iso8601(ts), 1, 10) as ymd, count(*) as num from device_log where device in (?) partition by substr(to_iso8601(ts), 1, 10) interval(4w)"
		list, err = tsdDb.GetTableDataAll(ctx, sql, strings.Join(deviceKeys, ","))
		if err != nil {
			return
		}

		for _, v := range list {
			data[v["ymd"].String()] = v["num"].Int()
		}
	}

	return
}

// getAlarmTotalForMonthsData 统计设备月度告警数量
func (s *sDevDevice) getAlarmTotalForMonthsData(ctx context.Context) (data map[int]int, err error) {

	data = initMonthlyMap()

	devices, err := s.GetAll(ctx)
	if err != nil {
		return
	}
	var deviceKeys []string
	for _, device := range devices {
		deviceKeys = append(deviceKeys, device.Key)
	}

	at := dao.AlarmLog.Columns().CreatedAt
	list, err := dao.AlarmLog.Ctx(ctx).
		Fields("date_format("+at+", '%Y%m') as ym, count(*) as num").
		Where(at+">=?", gtime.Now().Format("Y-01-01 00:00:00")).
		WhereIn(dao.AlarmLog.Columns().DeviceKey, deviceKeys).
		Group("date_format(" + at + ", '%Y%m')").
		All()
	if err != nil {
		return
	}

	for _, v := range list {
		m := gstr.SubStr(v["ym"].String(), 4)
		data[gconv.Int(m)] = v["num"].Int()
	}

	return
}

// getAlarmTotalForDayData 获取统计设备最近一个月告警数量
func (s *sDevDevice) getAlarmTotalForDayData(ctx context.Context) (data map[string]int, err error) {

	start := time.Now().AddDate(0, -1, 1)
	data = initDailyMap()

	devices, err := s.GetAll(ctx)
	if err != nil {
		return
	}
	var deviceKeys []string
	for _, device := range devices {
		deviceKeys = append(deviceKeys, "'"+device.Key+"'")
	}

	at := dao.AlarmLog.Columns().CreatedAt
	list, err := dao.AlarmLog.Ctx(ctx).
		Fields("date_format("+at+", '%Y-%m-%d') as ymd, count(*) as num").
		Where(at+">=?", start.Format("20060102")).
		WhereIn(dao.AlarmLog.Columns().DeviceKey, deviceKeys).
		Group("date_format(" + at + ", '%Y-%m-%d')").
		All()
	if err != nil {
		return
	}

	for _, v := range list {
		data[v["ymd"].String()] = v["num"].Int()
	}

	return
}

// GetDeviceDataList 获取设备属性聚合数据列表
func (s *sDevDevice) GetDeviceDataList(ctx context.Context, in *model.DeviceDataListInput) (out *model.DeviceDataListOutput, err error) {
	device, err := s.Get(ctx, in.DeviceKey)
	if err != nil {
		return
	}
	propertys := device.TSL.Properties
	if len(propertys) == 0 {
		return
	}

	fields := []string{"_wstart", "_wend"}
	for _, v := range propertys {
		key := comm.TsdColumnName(v.Key)
		fields = append(fields, fmt.Sprintf("mode(%s) as %s", key, key))
	}

	timeUnit := "m"
	switch in.TimeUnit {
	case 1:
		timeUnit = "s"
	case 2:
		timeUnit = "m"
	case 3:
		timeUnit = "h"
	case 4:
		timeUnit = "d"
	}
	interval := fmt.Sprintf("interval(%d%s)", in.Interval, timeUnit)

	tsdDb := tsd.DB()
	defer tsdDb.Close()

	deviceTable := comm.DeviceTableName(device.Key)
	sql := "select count(*) as num from (select count(*) from ? ?)"
	rs, err := tsdDb.GetTableDataOne(ctx, sql, deviceTable, interval)
	if err != nil {
		return
	}
	out = new(model.DeviceDataListOutput)
	out.Total = rs["num"].Int()
	out.CurrentPage = in.PageNum

	sql = "select ? from ? ? order by _wstart desc limit ?, ?"
	out.List, err = tsdDb.GetTableDataAll(
		ctx,
		sql,
		strings.Join(fields, ","),
		deviceTable,
		interval,
		(in.PageNum-1)*in.PageSize,
		in.PageSize,
	)
	return
}

// GetAllForProduct 获取指定产品所有设备
func (s *sDevDevice) GetAllForProduct(ctx context.Context, productKey string) (list []*entity.DevDevice, err error) {
	err = dao.DevDevice.Ctx(ctx).Where(dao.DevDevice.Columns().ProductKey, productKey).Scan(&list)
	return
}
