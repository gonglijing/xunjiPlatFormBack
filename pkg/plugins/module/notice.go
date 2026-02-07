package module

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	gplugin "github.com/hashicorp/go-plugin"
	"net/rpc"
	"sagooiot/pkg/plugins/model"
)

// Notice 通知服务插件接口
type Notice interface {
	Info() model.PluginInfo
	Send(data []byte) model.JsonRes
}

// NoticeRPC 基于RPC实现
type NoticeRPC struct {
	Client *rpc.Client
}

func (n *NoticeRPC) Info() model.PluginInfo {
	var resp model.PluginInfo
	err := n.Client.Call("Plugin.Info", new(interface{}), &resp)
	if err != nil {
		g.Log().Error(context.Background(), err)
	}
	return resp
}
func (n *NoticeRPC) Send(data []byte) model.JsonRes {
	var resp model.JsonRes
	err := n.Client.Call("Plugin.Send", data, &resp)
	if err != nil {
		resp.Code = 1
		resp.Message = fmt.Sprintf("notice.go Send %s", err.Error())
		g.Log().Error(context.Background(), err)
	}

	return resp
}

// NoticeRPCServer  GreeterRPC的RPC服务器，符合 net/rpc的要求
type NoticeRPCServer struct {
	Impl Notice
}

func (s *NoticeRPCServer) Info(args interface{}, resp *model.PluginInfo) error {
	*resp = s.Impl.Info()
	return nil
}
func (s *NoticeRPCServer) Send(data []byte, resp *model.JsonRes) error {
	*resp = s.Impl.Send(data)
	return nil
}

// NoticePlugin 插件的虚拟实现。用于PluginMap的插件接口。在运行时，来自插件实现的实际实现会覆盖
type NoticePlugin struct{}

// Server 此方法由插件进程延迟的调用
func (NoticePlugin) Server(*gplugin.MuxBroker) (interface{}, error) {
	checkParentAlive()
	return &NoticeRPCServer{}, nil
}

// Client 此方法由宿主进程调用
func (NoticePlugin) Client(b *gplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &NoticeRPC{Client: c}, nil
}
