package common

import (
	"context"
	"encoding/json"
	"fmt"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/service"
	"sagooiot/pkg/dcache"
	"sagooiot/pkg/iotModel/sagooProtocol"
	"sagooiot/pkg/iotModel/topicModel"
)

type DownstreamRouteLog struct {
	Payload interface{}            `json:"payload"`
	Route   map[string]interface{} `json:"route"`
}

// ResolveDownstreamTarget 解析下发目标设备。
// 普通设备返回原设备，子设备返回其绑定网关和子设备标识。
func ResolveDownstreamTarget(ctx context.Context, request topicModel.TopicDownHandlerData) (target topicModel.TopicDownHandlerData, identity *sagooProtocol.Identity, err error) {
	target = request
	if request.DeviceDetail == nil {
		err = fmt.Errorf("device detail is nil")
		return
	}
	if request.DeviceDetail.Product == nil {
		err = fmt.Errorf("device product detail is nil")
		return
	}

	if request.DeviceDetail.Product.DeviceType != model.DeviceTypeSub {
		return
	}

	gatewayKey, err := getBindGatewayKey(ctx, request.DeviceDetail.Key)
	if err != nil {
		return target, nil, err
	}

	gatewayDevice, getErr := dcache.GetDeviceDetailInfo(gatewayKey)
	if getErr != nil {
		return target, nil, getErr
	}
	if gatewayDevice == nil {
		gatewayDevice, err = service.DevDevice().Get(ctx, gatewayKey)
		if err != nil {
			return target, nil, err
		}
		_ = dcache.SetDeviceDetailInfo(gatewayKey, gatewayDevice)
	}
	if gatewayDevice == nil {
		return target, nil, fmt.Errorf("gateway[%s] not found", gatewayKey)
	}
	if gatewayDevice.Product == nil || gatewayDevice.Product.DeviceType != model.DeviceTypeGateway {
		return target, nil, fmt.Errorf("device[%s] is not a gateway", gatewayKey)
	}
	if dcache.GetDeviceStatus(ctx, gatewayDevice.Key) != model.DeviceStatusOn {
		return target, nil, fmt.Errorf("gateway[%s] is offline", gatewayDevice.Key)
	}

	productKey := request.DeviceDetail.ProductKey
	if request.DeviceDetail.Product != nil && request.DeviceDetail.Product.Key != "" {
		productKey = request.DeviceDetail.Product.Key
	}

	target.DeviceDetail = gatewayDevice
	identity = &sagooProtocol.Identity{
		ProductKey: productKey,
		DeviceKey:  request.DeviceDetail.Key,
	}
	return
}

// BuildGatewayRequestPayload 构造平台下发报文。
// 对子设备场景在根节点增加 identity，便于网关识别目标子设备。
func BuildGatewayRequestPayload(id, version, method string, params map[string]interface{}, identity *sagooProtocol.Identity) ([]byte, error) {
	payload := map[string]interface{}{
		"id":      id,
		"version": version,
		"params":  params,
		"method":  method,
	}
	if identity != nil {
		payload["identity"] = identity
	}
	return json.Marshal(payload)
}

// BuildTunnelPayload 构造隧道下发参数。
// 对子设备场景增加 identity 和 params 包装。
func BuildTunnelPayload(params map[string]interface{}, identity *sagooProtocol.Identity) ([]byte, error) {
	if identity == nil {
		return json.Marshal(params)
	}
	payload := map[string]interface{}{
		"params":   params,
		"identity": identity,
	}
	return json.Marshal(payload)
}

// BuildDownstreamRouteLog 构建下发日志内容，记录真实目标与实际通道。
func BuildDownstreamRouteLog(request, target topicModel.TopicDownHandlerData, identity *sagooProtocol.Identity, payload interface{}) DownstreamRouteLog {
	viaGateway := identity != nil
	return DownstreamRouteLog{
		Payload: payload,
		Route: map[string]interface{}{
			"targetProductKey":  getProductKey(request.DeviceDetail),
			"targetDeviceKey":   getDeviceKey(request.DeviceDetail),
			"channelProductKey": getProductKey(target.DeviceDetail),
			"channelDeviceKey":  getDeviceKey(target.DeviceDetail),
			"viaGateway":        viaGateway,
		},
	}
}

func getBindGatewayKey(ctx context.Context, subKey string) (string, error) {
	result, err := dao.DevDeviceGateway.Ctx(ctx).
		Fields(dao.DevDeviceGateway.Columns().GatewayKey).
		Where(dao.DevDeviceGateway.Columns().SubKey, subKey).
		One()
	if err != nil {
		return "", err
	}
	if result.IsEmpty() {
		return "", fmt.Errorf("sub device[%s] not bind gateway", subKey)
	}
	gatewayKey := result[dao.DevDeviceGateway.Columns().GatewayKey].String()
	if gatewayKey == "" {
		return "", fmt.Errorf("sub device[%s] bind gateway key is empty", subKey)
	}
	return gatewayKey, nil
}

func getProductKey(device *model.DeviceOutput) string {
	if device == nil {
		return ""
	}
	if device.Product != nil && device.Product.Key != "" {
		return device.Product.Key
	}
	return device.ProductKey
}

func getDeviceKey(device *model.DeviceOutput) string {
	if device == nil {
		return ""
	}
	return device.Key
}
