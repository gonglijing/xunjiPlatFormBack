package router

import (
	"context"
	analysisController "sagooiot/internal/controller/analysis"
	"sagooiot/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Analysis 分析统计相关的接口
func Analysis(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/analysis", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().Auth)
		group.Bind(
			analysisController.Device,     // 设备相关统计
			analysisController.Alarm,      // 设备相关相关统计
			analysisController.Product,    // 产品相关统计
			analysisController.DeviceData, // 设备数据相关统计

		)
	})
}
