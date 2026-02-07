package system

import (
	"context"
	"sagooiot/api/v1/system"
	"sagooiot/internal/model"
	"sagooiot/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var SysJobLog = cSysJobLog{}

type cSysJobLog struct{}

// Get 任务日志详情
func (a *cSysJobLog) Get(ctx context.Context, req *system.GetJobReq) (res *system.SysJobLogRes, err error) {
	var out *model.SysJobLogOut
	if out, err = service.SysJobLog().GetJobLog(ctx, req.Id); err != nil {
		g.Log().Error(ctx, "获取任务日志详情失败", err)
		return
	}
	res = new(system.SysJobLogRes)
	if err = gconv.Scan(out, &res); err != nil {
		g.Log().Error(ctx, "转换任务日志详情失败", err)
		return
	}
	g.Log().Debug(ctx, "获取任务日志详情", req)
	return
}

// List 任务日志列表
func (a *cSysJobLog) List(ctx context.Context, req *system.GetJobLogListReq) (res *system.GetJobLogListRes, err error) {
	var input *model.GetJobLogListInput
	if err = gconv.Scan(req, &input); err != nil {
		return
	}

	total, out, err := service.SysJobLog().JobLogList(ctx, input)

	if err != nil {
		return
	}
	res = new(system.GetJobLogListRes)
	res.Total = total
	res.CurrentPage = req.PageNum

	if out != nil {
		if err = gconv.Scan(out, &res.Data); err != nil {
			return
		}
	}

	return
}

func (a *cSysJobLog) Delete(ctx context.Context, req *system.DeleteJobLogReq) (res *system.DeleteJobLogRes, err error) {
	err = service.SysJobLog().DelJobLogByIds(ctx, req.Ids)
	return
}

func (a *cSysJobLog) Export(ctx context.Context, req *system.JobExportReq) (res *system.JobLogExportRes, err error) {
	var input *model.GetJobLogListInput
	if err = gconv.Scan(req, &input); err != nil {
		return
	}
	err = service.SysJobLog().Export(ctx, input)
	if err != nil {
		return
	}

	return
}
