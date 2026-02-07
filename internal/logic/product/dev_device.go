package product

import (
	"context"
	"encoding/json"
	"errors"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/dcache"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gconv"
)

type sDevDevice struct{}

func init() {
	service.RegisterDevDevice(deviceNew())
}

func deviceNew() *sDevDevice {
	return &sDevDevice{}
}

// Get 获取设备详情
func (s *sDevDevice) Get(ctx context.Context, key string) (out *model.DeviceOutput, err error) {
	err = dao.DevDevice.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 0,
		Name:     consts.GetDetailDeviceOutput + key,
		Force:    false,
	}).WithAll().Where(dao.DevDevice.Columns().Key, key).Scan(&out)
	if err != nil {
		return
	}
	if out == nil {
		err = errors.New("设备不存在")
		return
	}
	if out.Status != 0 {
		out.Status = dcache.GetDeviceStatus(ctx, out.Key) //查询设备状态
	}
	if out.Product != nil {
		out.ProductName = out.Product.Name
		if out.Product.Metadata != "" {
			err = json.Unmarshal([]byte(out.Product.Metadata), &out.TSL)
		}
	}
	return
}

// GetAll 获取所有设备
func (s *sDevDevice) GetAll(ctx context.Context) (out []*entity.DevDevice, err error) {
	m := dao.DevDevice.Ctx(ctx)
	err = m.Scan(&out)
	return
}

func (s *sDevDevice) Detail(ctx context.Context, key string) (out *model.DeviceOutput, err error) {
	err = dao.DevDevice.Ctx(ctx).WithAll().Where(dao.DevDevice.Columns().Key, key).Scan(&out)
	if err != nil {
		return
	}
	if out == nil {
		err = errors.New("设备不存在")
		return
	}
	if out.Status != 0 {
		out.Status = dcache.GetDeviceStatus(ctx, out.Key) //查询设备状态
	}

	//如果未设置，获取系统设置的默认超时时间
	if out.OnlineTimeout == 0 {
		defaultTimeout, err := service.ConfigData().GetConfigByKey(ctx, consts.DeviceDefaultTimeoutTime)
		if err != nil || defaultTimeout == nil {
			defaultTimeout = &entity.SysConfig{
				ConfigValue: "30",
			}
		}
		out.OnlineTimeout = gconv.Int(defaultTimeout.ConfigValue)
	}

	if out.Product != nil {
		out.ProductName = out.Product.Name
		if out.Product.Metadata != "" {
			err = json.Unmarshal([]byte(out.Product.Metadata), &out.TSL)
		}
	}
	return
}
