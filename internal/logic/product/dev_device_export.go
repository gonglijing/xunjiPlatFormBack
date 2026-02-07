package product

import (
	"context"
	"encoding/json"
	"errors"
	"sagooiot/api/v1/product"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/cache"
	"sagooiot/pkg/dcache"
	"sagooiot/pkg/response"
	"sagooiot/pkg/utility"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
)

// GetDeviceOnlineTimeOut 获取设备在线超时时长
func (s *sDevDevice) GetDeviceOnlineTimeOut(ctx context.Context, deviceKey string) (timeOut int) {
	// 获取设备在线超时时长
	timeOut = consts.DeviceOnlineTimeOut
	tv, err := dao.DevDevice.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: time.Second * 5,
		Name:     "DeviceOnlineTimeOut" + deviceKey,
		Force:    false,
	}).Where(dao.DevDevice.Columns().Key, deviceKey).Fields(dao.DevDevice.Columns().OnlineTimeout).Value()
	if err == nil && tv.Int() > 0 {
		timeOut = tv.Int()
	}
	return
}

// ExportDevices 导出设备
func (s *sDevDevice) ExportDevices(ctx context.Context, req *product.ExportDevicesReq) (res product.ExportDevicesRes, err error) {
	var data []model.DeviceOutput
	m := dao.DevDevice.Ctx(ctx).WithAll().Where(dao.DevDevice.Columns().ProductKey, req.ProductKey)
	if err = m.Scan(&data); err != nil {
		return
	}

	var outList []*model.DeviceExport
	for _, v := range data {
		var status string
		switch v.DevDevice.Status {
		case model.DeviceStatusNoEnable:
			status = "未启用"
		case model.DeviceStatusOff:
			status = "离线"
		case model.DeviceStatusOn:
			status = "在线"
		}
		var reqData = new(model.DeviceExport)
		reqData.Desc = v.DevDevice.Desc
		reqData.DeviceKey = v.DevDevice.Key
		reqData.DeviceName = v.DevDevice.Name
		reqData.ProductName = v.Product.Name
		reqData.Version = v.DevDevice.Version
		reqData.DeviceType = v.Product.DeviceType
		reqData.Status = status
		outList = append(outList, reqData)
	}
	if len(outList) == 0 {
		var reqData = new(model.DeviceExport)
		outList = append(outList, reqData)
	}

	//处理数据并导出
	var outData []interface{}
	for _, d := range outList {
		outData = append(outData, d)
	}
	dataRes := utility.ToExcel(outData)
	var request = g.RequestFromCtx(ctx)
	response.ToXls(request, dataRes, "设备列表")

	return
}

// ImportDevices 导入设备
func (s *sDevDevice) ImportDevices(ctx context.Context, req *product.ImportDevicesReq) (res product.ImportDevicesRes, err error) {
	file, err := req.File.Open()
	if err != nil {
		return
	}
	xlsx, err := excelize.OpenReader(file)

	if err != nil {
		return
	}

	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		return
	}
	device := new(entity.DevDevice)
	if len(rows) < 2 {
		err = errors.New("请添加设备")
		return
	}
	productData := new(entity.DevProduct)
	err = dao.DevProduct.Ctx(ctx).Where(dao.DevProduct.Columns().Key, req.ProductKey).Scan(&productData)
	if err != nil {
		err = errors.New("产品不存在")
		return
	}
	for rIndex, row := range rows {
		if len(row) < 7 {
			continue
		}
		if rIndex == 0 {
			continue
		}
		bl := gregex.IsMatch("^[A-Za-z_]+[A-Za-z0-9_]*|[0-9]+$", []byte(row[1]))
		if !bl {
			res.Fail++
			res.DevicesKey = append(res.DevicesKey, row[1])
			continue
		}
		device.DeptId = service.Context().GetUserDeptId(ctx)
		device.ProductKey = productData.Key
		device.Name = row[0]
		device.Key = row[1]
		device.Desc = row[2]
		device.Version = row[3]
		device.Lng = row[4]
		device.Lat = row[5]
		device.OnlineTimeout = gconv.Int(row[6])

		_, err = dao.DevDevice.Ctx(ctx).Insert(device)
		if err != nil {
			res.Fail++
			res.DevicesKey = append(res.DevicesKey, device.Key)
		} else {
			res.Success++
		}
	}
	return
}

// CacheDeviceDetailList 缓存所有设备详情数据
func (s *sDevDevice) CacheDeviceDetailList(ctx context.Context) (err error) {
	productList, err := service.DevProduct().List(context.Background())
	if err != nil {
		return
	}
	for _, p := range productList {
		deviceList, err := service.DevDevice().List(context.Background(), p.Key, "")
		if err != nil {
			g.Log().Error(ctx, err.Error())
		}

		//缓存产品详细信息
		var detailProduct = new(model.DetailProductOutput)
		err = gconv.Scan(p, detailProduct)
		if err == nil {
			err = dcache.SetProductDetailInfo(p.Key, detailProduct)
		} else {
			g.Log().Error(ctx, err.Error())
		}

		for _, d := range deviceList {
			if d.Product.Metadata != "" {
				err = json.Unmarshal([]byte(d.Product.Metadata), &d.TSL)
				d.Product.Metadata = ""
				if err != nil {
					continue
				}
			}
			//缓存设备详细信息
			err := cache.Instance().Set(context.Background(), consts.DeviceDetailInfoPrefix+d.Key, d, 0)
			if err != nil {
				g.Log().Error(ctx, err.Error())
			}
		}
	}
	return

}
