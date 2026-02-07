package product

import (
	"context"
	"sagooiot/internal/consts"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/pkg/dcache"

	"github.com/gogf/gf/v2/util/gconv"
)

// ListForPage 设备分页列表
func (s *sDevDevice) ListForPage(ctx context.Context, in *model.ListDeviceForPageInput) (out *model.ListDeviceForPageOutput, err error) {
	out = new(model.ListDeviceForPageOutput)
	c := dao.DevDevice.Columns()
	m := dao.DevDevice.Ctx(ctx).OrderDesc(c.Id)

	if in.Status != "" {
		//获取当前在线的设备KEY列表
		onlineDevice, err := dcache.GetOnlineDeviceList()
		if in.Status == gconv.String(consts.DeviceStatueOnline) {
			if err == nil && len(onlineDevice) > 0 {
				m = m.Where(c.Key, onlineDevice)
			} else {
				return nil, nil
			}

		} else if in.Status == gconv.String(consts.DeviceStatueOffline) {
			if err == nil && len(onlineDevice) > 0 {
				m = m.WhereNotIn(c.Key, onlineDevice)
			}
		} else {
			m = m.Where(c.Status, gconv.Int(in.Status))
		}
	}
	if in.Key != "" {
		m = m.Where(c.Key, in.Key)
	}
	if in.ProductKey != "" {
		m = m.Where(c.ProductKey, in.ProductKey)
	}
	if in.TunnelId > 0 {
		m = m.Where(c.TunnelId, in.TunnelId)
	}
	if in.Name != "" {
		m = m.WhereLike(c.Name, "%"+in.Name+"%")
	}
	if len(in.DateRange) > 0 {
		m = m.WhereBetween(c.CreatedAt, in.DateRange[0], in.DateRange[1])
	}

	m = m.WhereIn(
		dao.DevDevice.Columns().ProductKey,
		dao.DevProduct.Ctx(ctx).Fields(dao.DevProduct.Columns().Key),
	)

	out.Total, _ = m.Count()
	out.CurrentPage = in.PageNum
	err = m.WithAll().Page(in.PageNum, in.PageSize).Scan(&out.Device)
	if err != nil {
		return
	}

	for i, v := range out.Device {
		if v.Product != nil {
			out.Device[i].ProductName = v.Product.Name
		}
		if out.Device[i].Status != 0 {
			out.Device[i].Status = dcache.GetDeviceStatus(ctx, out.Device[i].Key)
		}
		out.Device[i].Product.Metadata = ""
	}

	return
}

// List 已发布产品的设备列表
func (s *sDevDevice) List(ctx context.Context, productKey string, keyWord string) (list []*model.DeviceOutput, err error) {
	m := dao.DevDevice.Ctx(ctx).
		Where(dao.DevDevice.Columns().Status+" > ?", model.DeviceStatusNoEnable)
	if productKey != "" {
		m = m.Where(dao.DevDevice.Columns().ProductKey, productKey)
	}
	if keyWord != "" {
		m = m.WhereLike(dao.DevDevice.Columns().Key, "%"+keyWord+"%").WhereOrLike(dao.DevDevice.Columns().Name, "%"+keyWord+"%")
	}

	err = m.WhereIn(dao.DevDevice.Columns().ProductKey,
		dao.DevProduct.Ctx(ctx).
			Fields(dao.DevProduct.Columns().Key).
			Where(dao.DevProduct.Columns().Status, model.ProductStatusOn)).
		WithAll().
		OrderDesc(dao.DevDevice.Columns().Id).
		Scan(&list)
	if err != nil {
		return
	}

	for i, v := range list {
		if v.Product != nil {
			list[i].ProductName = v.Product.Name
		}
	}

	return
}
