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
	"sagooiot/pkg/dcache"

	"github.com/gogf/gf/v2/os/gtime"
)

// BindSubDevice 网关绑定子设备
func (s *sDevDevice) BindSubDevice(ctx context.Context, in *model.DeviceBindInput) error {
	gw, err := s.Get(ctx, in.GatewayKey)
	if err != nil {
		return err
	}
	if gw.Product.DeviceType != model.DeviceTypeGateway {
		return errors.New("非网关，不能绑定子设备")
	}
	if len(in.SubKeys) == 0 {
		return nil
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)

	for _, v := range in.SubKeys {
		sub, err := s.Get(ctx, v)
		if err != nil {
			return err
		}
		if sub.Product.DeviceType != model.DeviceTypeSub {
			return errors.New("非子设备类型，不能绑定")
		}

		rs, err := dao.DevDeviceGateway.Ctx(ctx).
			Where(dao.DevDeviceGateway.Columns().SubKey, v).
			One()
		if err != nil {
			return err
		}
		if !rs.IsEmpty() {
			return errors.New(fmt.Sprintf("%s，该子设备已被绑定", sub.Name))
		}

		_, err = dao.DevDeviceGateway.Ctx(ctx).Data(do.DevDeviceGateway{
			GatewayKey: in.GatewayKey,
			SubKey:     v,
			CreatedBy:  uint(loginUserId),
			CreatedAt:  gtime.Now(),
		}).Insert()
		if err != nil {
			return err
		}
	}
	return nil
}

// UnBindSubDevice 网关解绑子设备
func (s *sDevDevice) UnBindSubDevice(ctx context.Context, in *model.DeviceBindInput) error {
	gw, err := s.Get(ctx, in.GatewayKey)
	if err != nil {
		return err
	}
	if gw.Product.DeviceType != model.DeviceTypeGateway {
		return errors.New("非网关设备")
	}
	if len(in.SubKeys) == 0 {
		return nil
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)

	_, err = dao.DevDeviceGateway.Ctx(ctx).
		Data(do.DevDeviceGateway{
			DeletedBy: uint(loginUserId),
			DeletedAt: gtime.Now(),
		}).
		Where(dao.DevDeviceGateway.Columns().GatewayKey, in.GatewayKey).
		WhereIn(dao.DevDeviceGateway.Columns().SubKey, in.SubKeys).
		Unscoped().
		Update()
	return err
}

// bindList 已绑定列表
func (s *sDevDevice) bindList(ctx context.Context, gatewayKey string) (list []*model.DeviceOutput, err error) {
	_, err = s.Get(ctx, gatewayKey)
	if err != nil {
		return
	}

	var dgw []*entity.DevDeviceGateway
	if err = dao.DevDeviceGateway.Ctx(ctx).Where(dao.DevDeviceGateway.Columns().GatewayKey, gatewayKey).Scan(&dgw); err != nil || len(dgw) == 0 {
		return
	}

	var subKeys []string
	for _, v := range dgw {
		subKeys = append(subKeys, v.SubKey)
	}

	err = dao.DevDevice.Ctx(ctx).
		WhereIn(dao.DevDevice.Columns().Key, subKeys).
		WithAll().
		OrderDesc(dao.DevDevice.Columns().Id).
		Scan(&list)
	return
}

// BindList 已绑定列表(分页)
func (s *sDevDevice) BindList(ctx context.Context, in *model.DeviceBindListInput) (out *model.DeviceBindListOutput, err error) {
	_, err = s.Get(ctx, in.GatewayKey)
	if err != nil {
		return
	}

	var dgw []*entity.DevDeviceGateway
	if err = dao.DevDeviceGateway.Ctx(ctx).Where(dao.DevDeviceGateway.Columns().GatewayKey, in.GatewayKey).Scan(&dgw); err != nil || len(dgw) == 0 {
		return
	}

	var subKeys []string
	for _, v := range dgw {
		subKeys = append(subKeys, v.SubKey)
	}

	m := dao.DevDevice.Ctx(ctx).
		WhereIn(dao.DevDevice.Columns().Key, subKeys).
		WithAll().
		OrderDesc(dao.DevDevice.Columns().Id)

	out = &model.DeviceBindListOutput{}
	out.Total, _ = m.Count()
	out.CurrentPage = in.PageNum
	err = m.Page(in.PageNum, in.PageSize).Scan(&out.List)
	if err != nil {
		return
	}

	for i, v := range out.List {
		out.List[i].Status = dcache.GetDeviceStatus(ctx, v.Key)
	}

	return
}

// ListForSub 子设备
func (s *sDevDevice) ListForSub(ctx context.Context, in *model.ListForSubInput) (out *model.ListDeviceForPageOutput, err error) {
	m := dao.DevDevice.Ctx(ctx)
	if in.ProductKey != "" {
		m = m.Where(dao.DevDevice.Columns().ProductKey, in.ProductKey)
	}
	if in.GatewayKey != "" {
		rs, err := dao.DevDeviceGateway.Ctx(ctx).
			Fields(dao.DevDeviceGateway.Columns().SubKey).
			Where(dao.DevDeviceGateway.Columns().GatewayKey, in.GatewayKey).
			Array()
		if err == nil && len(rs) > 0 {
			m = m.WhereNotIn(dao.DevDevice.Columns().Key, rs)
		}
	}

	m = m.WithAll().OrderDesc(dao.DevDevice.Columns().Id)
	out = &model.ListDeviceForPageOutput{}
	out.Total, _ = m.Count()
	out.CurrentPage = in.PageNum
	err = m.Page(in.PageNum, in.PageSize).Scan(&out.Device)
	return
}

// CheckBind 检查网关、子设备绑定关系
func (s *sDevDevice) CheckBind(ctx context.Context, in *model.CheckBindInput) (bool, error) {
	count, err := dao.DevDeviceGateway.Ctx(ctx).
		Where(dao.DevDeviceGateway.Columns().GatewayKey, in.GatewayKey).
		Where(dao.DevDeviceGateway.Columns().SubKey, in.SubKey).
		Count()
	if err != nil || count == 0 {
		return false, err
	}
	return true, nil
}

// DelSub 子设备删除
func (s *sDevDevice) DelSub(ctx context.Context, key string) (err error) {
	subDev, err := s.Detail(ctx, key)
	if err != nil {
		return
	}
	if subDev.Product.DeviceType != model.DeviceTypeSub {
		return errors.New("该设备不是子设备类型")
	}
	if subDev.Status > model.DeviceStatusNoEnable {
		return errors.New("设备已启用，不能删除")
	}

	//获取当前登录用户ID
	loginUserId := service.Context().GetUserId(ctx)
	now := gtime.Now()

	rs, err := dao.DevDevice.Ctx(ctx).
		Data(do.DevDevice{
			DeletedBy: uint(loginUserId),
			DeletedAt: now,
		}).
		Where(dao.DevDevice.Columns().Key, key).
		Unscoped().
		Update()
	if err != nil {
		return err
	}

	// 删除绑定关系
	_, err = dao.DevDeviceGateway.Ctx(ctx).
		Data(do.DevDeviceGateway{
			DeletedBy: uint(loginUserId),
			DeletedAt: now,
		}).
		Where(dao.DevDeviceGateway.Columns().SubKey, subDev.Key).
		Unscoped().
		Update()
	if err != nil {
		return
	}

	num, _ := rs.RowsAffected()
	if num > 0 && subDev.MetadataTable == 1 {
		// 删除TD子表
		err = service.TSLTable().DropTable(ctx, subDev.Key)
		if err != nil {
			return err
		}
	}

	return
}

// AuthInfo 获取认证信息
func (s *sDevDevice) AuthInfo(ctx context.Context, in *model.AuthInfoInput) (*model.AuthInfoOutput, error) {
	if in.DeviceKey == "" && in.ProductKey == "" {
		return nil, errors.New("缺少必要参数")
	}

	out := &model.AuthInfoOutput{}

	if in.DeviceKey != "" {
		device, err := s.Get(ctx, in.DeviceKey)
		if err != nil {
			return nil, err
		}
		out.AuthType = device.AuthType
		out.AuthUser = device.AuthUser
		out.AuthPasswd = device.AuthPasswd
		out.AccessToken = device.AccessToken
		out.CertificateId = device.CertificateId

		if out.AuthUser == "" && out.AuthPasswd == "" && out.AccessToken == "" && out.CertificateId == 0 {
			in.ProductKey = device.Product.Key
		} else {
			in.ProductKey = ""
		}
	}

	if in.ProductKey != "" {
		productData, err := service.DevProduct().Detail(ctx, in.ProductKey)
		if err != nil {
			return nil, err
		}
		out.AuthType = productData.AuthType
		out.AuthUser = productData.AuthUser
		out.AuthPasswd = productData.AuthPasswd
		out.AccessToken = productData.AccessToken
		out.CertificateId = productData.CertificateId
	}

	if out.CertificateId > 0 {
		if err := dao.SysCertificate.Ctx(ctx).Where(dao.SysCertificate.Columns().Id, out.CertificateId).Scan(&out.Certificate); err != nil {
			return nil, err
		}
	}

	return out, nil
}
