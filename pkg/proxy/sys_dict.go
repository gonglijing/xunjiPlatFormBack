package proxy

import (
	"context"
	"sagooiot/internal/model"
	"sagooiot/internal/service"
	proxyModel "sagooiot/pkg/proxy/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// GetDictWithDataByType 通过字典键类型获取选项
func GetDictWithDataByType(ctx context.Context, input *proxyModel.GetDictInput) (output *proxyModel.GetDictOut, err error) {
	var in *model.GetDictInput
	if err = gconv.Scan(input, &in); err != nil {
		return
	}
	if g.IsEmpty(in.DictType) {
		return
	}

	dict, err := service.DictData().GetDictWithDataByType(ctx, in)
	if err != nil {
		return
	}
	if !g.IsEmpty(dict) && !g.IsEmpty(dict.Data) && !g.IsEmpty(dict.Values) {
		err = gconv.Scan(dict, &output)
	}
	return
}
