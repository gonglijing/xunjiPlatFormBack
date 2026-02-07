package product

import (
	"context"
	"errors"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/service"
	"sagooiot/pkg/dcache"
	"sagooiot/pkg/tsd/comm"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Deploy 设备启用
func (s *sDevDevice) Deploy(ctx context.Context, key string) (err error) {
	//获取设备信息
	device, err := s.Detail(ctx, key)
	if err != nil {
		return
	}
	if device.Status > model.DeviceStatusNoEnable {
		return
	}

	pd, err := service.DevProduct().Detail(ctx, device.ProductKey)
	if err != nil {
		return err
	}
	if pd == nil {
		return errors.New("产品不存在")
	}
	if pd.Status != model.ProductStatusOn {
		return errors.New("产品未发布，请先发布产品")
	}

	err = dao.DevDevice.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 是否创建TD表结构
		var isCreate bool

		if device.MetadataTable == 0 {
			isCreate = true
		}
		// 检查TD表是否存在，不存在，则补建
		if device.MetadataTable == 1 {
			table := comm.DeviceTableName(device.Key)
			b, _ := service.TSLTable().CheckTable(ctx, table)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					isCreate = true
				} else {
					return err
				}
			}
			if !b {
				isCreate = true
			}
		}

		// 创建TD表结构
		if isCreate {
			err = service.TSLTable().CreateTable(ctx, pd.Key, device.Key)
			if err != nil {
				return err
			}
		}

		// 更新状态
		_, err = dao.DevDevice.Ctx(ctx).
			Data(g.Map{
				dao.DevDevice.Columns().Status:        model.DeviceStatusOff,
				dao.DevDevice.Columns().MetadataTable: 1,
			}).
			Where(dao.DevDevice.Columns().Key, key).
			Update()

		//设备启用后，更新缓存数据
		device.Status = model.DeviceStatusOff
		err = dcache.SetDeviceDetailInfo(device.Key, device)
		if err != nil {
			g.Log().Debug(ctx, "Deploy 设备数据存入缓存失败", err.Error())
		}

		// 网关启用子设备
		if pd.DeviceType == model.DeviceTypeGateway {
			subList, err := s.bindList(ctx, device.Key)
			if err != nil {
				return err
			}

			for _, sub := range subList {
				if err := s.Deploy(ctx, sub.Key); err != nil {
					return err
				}
				//设备启用后，更新缓存数据
				sub.Status = model.DeviceStatusOff
				err = dcache.SetDeviceDetailInfo(sub.Key, sub)
				if err != nil {
					g.Log().Debug(ctx, "Deploy 子设备数据存入缓存失败", err.Error())
				}
			}
		}

		return err
	})

	return
}

// Undeploy 设备禁用
func (s *sDevDevice) Undeploy(ctx context.Context, key string) (err error) {
	//获取设备信息
	device, err := s.Detail(ctx, key)
	if err != nil {
		return
	}
	if device.Status == model.DeviceStatusNoEnable {
		return
	}

	err = dao.DevDevice.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = dao.DevDevice.Ctx(ctx).
			Data(g.Map{dao.DevDevice.Columns().Status: model.DeviceStatusNoEnable}).
			Where(dao.DevDevice.Columns().Key, key).
			Update()
		if err != nil {
			return err
		}

		//设备停用后，更新缓存数据
		device.Status = model.DeviceStatusNoEnable
		err = dcache.SetDeviceDetailInfo(device.Key, device)
		if err != nil {
			g.Log().Debug(ctx, "Deploy 设备数据存入缓存失败", err.Error())
		}

		// 网关禁用子设备
		if device.Product.DeviceType == model.DeviceTypeGateway {
			subList, err := s.bindList(ctx, device.Key)
			if err != nil {
				return err
			}

			for _, sub := range subList {
				if err := s.Undeploy(ctx, sub.Key); err != nil {
					return err
				}
				//设备停用后，更新缓存数据
				sub.Status = model.DeviceStatusNoEnable
				err = dcache.SetDeviceDetailInfo(sub.Key, sub)
				if err != nil {
					g.Log().Debug(ctx, "Deploy 子设备数据存入缓存失败", err.Error())
				}
			}
		}

		return err
	})

	return
}
