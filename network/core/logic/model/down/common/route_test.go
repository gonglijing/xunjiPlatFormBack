package common

import (
	"encoding/json"
	"sagooiot/internal/model"
	"sagooiot/internal/model/entity"
	"sagooiot/pkg/iotModel/sagooProtocol"
	"sagooiot/pkg/iotModel/topicModel"
	"testing"
)

func TestBuildGatewayRequestPayloadWithIdentity(t *testing.T) {
	params := map[string]interface{}{"switch": 1}
	identity := &sagooProtocol.Identity{ProductKey: "sub-pk", DeviceKey: "sub-dk"}

	payloadBytes, err := BuildGatewayRequestPayload("id-1", "1.0.0", "thing.service.turn_on", params, identity)
	if err != nil {
		t.Fatalf("build payload failed: %v", err)
	}

	var payload map[string]interface{}
	if err = json.Unmarshal(payloadBytes, &payload); err != nil {
		t.Fatalf("unmarshal payload failed: %v", err)
	}

	rawIdentity, ok := payload["identity"].(map[string]interface{})
	if !ok {
		t.Fatalf("identity missing in payload")
	}
	if rawIdentity["productKey"] != identity.ProductKey {
		t.Fatalf("identity.productKey not match")
	}
	if rawIdentity["deviceKey"] != identity.DeviceKey {
		t.Fatalf("identity.deviceKey not match")
	}
}

func TestBuildTunnelPayload(t *testing.T) {
	params := map[string]interface{}{"temperature": 26}
	identity := &sagooProtocol.Identity{ProductKey: "sub-pk", DeviceKey: "sub-dk"}

	payloadBytes, err := BuildTunnelPayload(params, identity)
	if err != nil {
		t.Fatalf("build payload failed: %v", err)
	}

	var payload map[string]interface{}
	if err = json.Unmarshal(payloadBytes, &payload); err != nil {
		t.Fatalf("unmarshal payload failed: %v", err)
	}

	if _, ok := payload["params"]; !ok {
		t.Fatalf("params missing in tunnel payload")
	}
	rawIdentity, ok := payload["identity"].(map[string]interface{})
	if !ok {
		t.Fatalf("identity missing in tunnel payload")
	}
	if rawIdentity["productKey"] != identity.ProductKey {
		t.Fatalf("identity.productKey not match")
	}
	if rawIdentity["deviceKey"] != identity.DeviceKey {
		t.Fatalf("identity.deviceKey not match")
	}
}

func TestBuildDownstreamRouteLog(t *testing.T) {
	request := topicModel.TopicDownHandlerData{DeviceDetail: &model.DeviceOutput{DevDevice: &entity.DevDevice{Key: "sub-dk", ProductKey: "sub-pk"}}}
	target := topicModel.TopicDownHandlerData{DeviceDetail: &model.DeviceOutput{DevDevice: &entity.DevDevice{Key: "gw-dk", ProductKey: "gw-pk"}}}
	identity := &sagooProtocol.Identity{ProductKey: "sub-pk", DeviceKey: "sub-dk"}

	logObj := BuildDownstreamRouteLog(request, target, identity, map[string]interface{}{"id": "1"})
	if logObj.Route["targetDeviceKey"] != "sub-dk" {
		t.Fatalf("targetDeviceKey not match")
	}
	if logObj.Route["channelDeviceKey"] != "gw-dk" {
		t.Fatalf("channelDeviceKey not match")
	}
	if logObj.Route["viaGateway"] != true {
		t.Fatalf("viaGateway should be true")
	}
}
