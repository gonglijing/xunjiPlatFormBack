package product

import (
	"context"
	"errors"
	"sagooiot/api/v1/product"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/pkg/dcache"
	"sagooiot/pkg/iotModel"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// UpdateDeviceStatusInfo 更新设备状态信息，设备上线、离线、注册
func (s *sDevDevice) UpdateDeviceStatusInfo(ctx context.Context, deviceKey string, status int, timestamp time.Time) (err error) {
	device, err := s.Get(ctx, deviceKey)
	if err != nil {
		return
	}
	if device == nil {
		err = errors.New("设备不存在")
		return
	}
	var data = g.Map{}
	switch status {
	case consts.DeviceStatueOnline:
		if device.RegistryTime == nil {
			data[dao.DevDevice.Columns().RegistryTime] = timestamp
		} else {

			data[dao.DevDevice.Columns().LastOnlineTime] = timestamp
		}
		data[dao.DevDevice.Columns().Status] = consts.DeviceStatueOnline

	case consts.DeviceStatueOffline:
		data[dao.DevDevice.Columns().LastOnlineTime] = timestamp
		data[dao.DevDevice.Columns().Status] = consts.DeviceStatueOffline
	}
	_, err = dao.DevDevice.Ctx(ctx).
		Data(data).
		Where(dao.DevDevice.Columns().Key, deviceKey).
		Update()
	return
}

// BatchUpdateDeviceStatusInfo 批量更新设备状态信息，设备上线、离线、注册
func (s *sDevDevice) BatchUpdateDeviceStatusInfo(ctx context.Context, deviceStatusLogList []iotModel.DeviceStatusLog) (err error) {

	onLineData := g.Map{}
	onlineDeviceKeyList := make([]string, 0)

	offLineData := g.Map{}
	offLineDeviceKeyList := make([]string, 0)

	registryData := g.Map{}
	registryDeviceKeyList := make([]string, 0)

	for _, statusLog := range deviceStatusLogList {
		switch statusLog.Status {
		case consts.DeviceStatueOnline:
			device, err := dcache.GetDeviceDetailInfo(statusLog.DeviceKey)
			if err != nil {
				continue
			}

			if device.RegistryTime == nil {
				registryData[dao.DevDevice.Columns().RegistryTime] = statusLog.Timestamp
				registryDeviceKeyList = append(registryDeviceKeyList, statusLog.DeviceKey)
			}

			onLineData[dao.DevDevice.Columns().LastOnlineTime] = statusLog.Timestamp
			onLineData[dao.DevDevice.Columns().Status] = consts.DeviceStatueOnline
			onlineDeviceKeyList = append(onlineDeviceKeyList, statusLog.DeviceKey)
		case consts.DeviceStatueOffline:
			offLineData[dao.DevDevice.Columns().LastOnlineTime] = statusLog.Timestamp
			offLineData[dao.DevDevice.Columns().Status] = consts.DeviceStatueOffline
			offLineDeviceKeyList = append(offLineDeviceKeyList, statusLog.DeviceKey)
		}
	}
	if len(registryDeviceKeyList) > 0 {
		_, err = dao.DevDevice.Ctx(ctx).
			Data(registryData).
			WhereIn(dao.DevDevice.Columns().Key, registryDeviceKeyList).
			Update()
	}
	if len(onlineDeviceKeyList) > 0 {
		_, err = dao.DevDevice.Ctx(ctx).
			Data(onLineData).
			WhereIn(dao.DevDevice.Columns().Key, onlineDeviceKeyList).
			Update()
	}
	if len(offLineDeviceKeyList) > 0 {
		_, err = dao.DevDevice.Ctx(ctx).
			Data(offLineData).
			WhereIn(dao.DevDevice.Columns().Key, offLineDeviceKeyList).
			Update()
	}

	return
}

func (s *sDevDevice) SetDevicesStatus(ctx context.Context, req *product.SetDeviceStatusReq) (res product.SetDeviceStatusRes, err error) {
	if req.Status == 1 {
		for _, v := range req.Keys {
			err = s.Deploy(ctx, v)
			if err != nil {
				return
			}
		}
		return
	}
	if req.Status == 0 {
		for _, v := range req.Keys {
			err = s.Undeploy(ctx, v)
			if err != nil {
				return
			}
		}
	}
	return
}
