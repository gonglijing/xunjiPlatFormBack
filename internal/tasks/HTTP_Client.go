package tasks

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gtime"
)

// GetAccessURL 执行访问URL
func (t TaskJob) GetAccessURL(accessURL string) {
	ctx := context.Background()
	g.Log().Debug(ctx, "访问URL：", accessURL)
	startTime := gtime.Now()
	res, err := g.Client().Get(ctx, accessURL)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	if err := t.SaveLog(ctx, startTime,fmt.Sprintf("访问URL：%s成功",accessURL),err); err != nil {
		g.Log().Error(ctx, err)
	}
	defer func(res *gclient.Response) {
		if err := res.Close(); err != nil {
			g.Log().Error(ctx, err)
		}
	}(res)
}
