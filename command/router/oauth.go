package router

import (
	"context"
	oauthController "sagooiot/internal/controller/oauth"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Analysis 分析统计相关的接口
func OAuth(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/oauth", func(group *ghttp.RouterGroup) {
		group.Bind(
			oauthController.OProvider, // 第三方授权配置提供
			oauthController.OUser,     // 第三方授权用户登录
		)
	})
}
