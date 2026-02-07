package product

import (
	"context"
	"errors"
	"fmt"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/model/do"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"sagooiot/pkg/cache"
	"sagooiot/internal/consts"
	"sagooiot/pkg/iotModel"
	"sagooiot/pkg/iotModel/sagooProtocol/north"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

func (s *sDevDevice) Add(ctx context.Context, in *model.AddDeviceInput) (deviceId uint, err error) {
	id, _ := dao.DevDevice.Ctx(ctx).
		Fields(dao.DevDevice.Columns().Id).
		Where(dao.DevDevice.Columns().Key, in.Key).
		Value()
	if id.Uint() > 0 {
		err = errors.New("设备标识重复")
		return
	}

	// 设备标识不能和产品标识重复
	pout, err := service.DevProduct().Detail(ctx, in.Key)
	if err != nil {
		return
	}
	if pout != nil {
		err = errors.New("设备标识不能和产品标识重复")
		return
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)

	rs, err := dao.DevDevice.Ctx(ctx).Data(do.DevDevice{
		DeptId:        service.Context().GetUserDeptId(ctx),
		Key:           in.Key,
		Name:          in.Name,
		ProductKey:    in.ProductKey,
		Desc:          in.Desc,
		Version:       in.Version,
		Lng:           in.Lng,
		Lat:           in.Lat,
		AuthType:      in.AuthType,
		AuthUser:      in.AuthUser,
		AuthPasswd:    in.AuthPasswd,
		AccessToken:   in.AccessToken,
		CertificateId: in.CertificateId,
		Status:        0,
		CreatedBy:     uint(loginUserId),
		CreatedAt:     gtime.Now(),
	}).Insert()
	if err != nil {
		return
	}
	//北向设备添加消息
	north.WriteMessage(ctx, north.DeviceAddMessageTopic, nil, in.ProductKey, in.Key, iotModel.DeviceAddMessage{
		Timestamp: time.Now().UnixMilli(),
		Desc:      "",
	})
	newId, err := rs.LastInsertId()
	deviceId = uint(newId)

	// 设备标签
	if len(in.Tags) > 0 {
		for _, v := range in.Tags {
			intag := model.AddTagDeviceInput{
				DeviceId:  deviceId,
				DeviceKey: in.Key,
				Key:       v.Key,
				Name:      v.Name,
				Value:     v.Value,
			}
			if err = service.DevDeviceTag().Add(ctx, &intag); err != nil {
				return
			}
		}
	}

	return
}

func (s *sDevDevice) Edit(ctx context.Context, in *model.EditDeviceInput) (err error) {
	out, err := s.Detail(ctx, in.Key)
	if err != nil {
		return
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)

	var param *do.DevDevice
	err = gconv.Scan(in, &param)
	if err != nil {
		return
	}
	param.UpdatedBy = uint(loginUserId)
	param.Id = nil

	_, err = dao.DevDevice.Ctx(ctx).Data(param).Where(dao.DevDevice.Columns().Key, in.Key).Update()
	if err != nil {
		return
	}

	// 设备标签
	if len(in.Tags) > 0 {
		var l []model.AddTagDeviceInput
		for _, v := range in.Tags {
			intag := model.AddTagDeviceInput{
				DeviceId:  out.Id,
				DeviceKey: out.Key,
				Key:       v.Key,
				Name:      v.Name,
				Value:     v.Value,
			}
			l = append(l, intag)
		}
		if err = service.DevDeviceTag().Update(ctx, out.Id, l); err != nil {
			return
		}
	}
	//从缓存中删除
	_, err = cache.Instance().Remove(ctx, consts.CacheGfOrmPrefix+consts.GetDetailDeviceOutput+out.Key)

	return
}

func (s *sDevDevice) UpdateExtend(ctx context.Context, in *model.DeviceExtendInput) (err error) {
	_, err = s.Detail(ctx, in.Key)
	if err != nil {
		return
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)

	var param *do.DevDevice
	err = gconv.Scan(in, &param)
	if err != nil {
		return
	}
	param.UpdatedBy = uint(loginUserId)
	param.Id = nil

	_, err = dao.DevDevice.Ctx(ctx).Data(param).Where(dao.DevDevice.Columns().Key, in.Key).Update()

	return
}

func (s *sDevDevice) Del(ctx context.Context, keys []string) (err error) {
	var p []*entity.DevDevice
	err = dao.DevDevice.Ctx(ctx).WhereIn(dao.DevDevice.Columns().Key, keys).Scan(&p)
	if err != nil {
		return
	}
	if len(p) == 0 {
		return errors.New("设备不存在")
	}

	// 状态校验
	for _, v := range p {
		if v.Status > model.DeviceStatusNoEnable {
			return errors.New(fmt.Sprintf("设备(%s)已启用，不能删除", v.Key))
		}
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)
	for _, key := range keys {
		res, _ := s.Detail(ctx, key)
		rs, err := dao.DevDevice.Ctx(ctx).
			Data(do.DevDevice{
				DeletedBy: uint(loginUserId),
				DeletedAt: gtime.Now(),
			}).
			Where(dao.DevDevice.Columns().Status, model.DeviceStatusNoEnable).
			Where(dao.DevDevice.Columns().Key, key).
			Unscoped().
			Update()
		if err != nil {
			return err
		}

		num, _ := rs.RowsAffected()
		if num > 0 && res.MetadataTable == 1 {
			// 删除TD子表
			err = service.TSLTable().DropTable(ctx, res.Key)
			if err != nil {
				return err
			}

			// 删除网关子设备TD子表
			if res.Product.DeviceType == model.DeviceTypeGateway {
				subList, err := s.bindList(ctx, res.Key)
				if err != nil {
					return err
				}
				for _, sub := range subList {
					if sub.MetadataTable == 1 {
						// 删除TD子表
						err = service.TSLTable().DropTable(ctx, sub.Key)
						if err != nil {
							return err
						}
					}
				}
			}
		}
		//北向设备删除消息
		north.WriteMessage(ctx, north.DeviceDeleteMessageTopic, nil, res.Product.Key, res.Key, iotModel.DeviceDeleteMessage{
			Timestamp: time.Now().UnixMilli(),
			Desc:      "",
		})
	}

	return
}
