package tunnel

import (
	"fmt"
	"sagooiot/network/core/tunnel/base"
	"sagooiot/network/model"
)

// NewTunnel 创建通道
func NewTunnel(tunnel *model.Tunnel) (base.TunnelInstance, error) {
	var tnl base.TunnelInstance
	switch tunnel.Type {
	case "tcp-client":
		tnl = newTunnelClient(tunnel, "tcp")
	case "udp-client":
		tnl = newTunnelClient(tunnel, "udp")
	default:
		return nil, fmt.Errorf("Unsupport type %s ", tunnel.Type)
	}
	return tnl, nil
}
