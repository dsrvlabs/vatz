package healthcheck

import (
	"context"

	"github.com/dsrvlabs/vatz/manager/config"
	msg "github.com/dsrvlabs/vatz/manager/model"
	"github.com/dsrvlabs/vatz/manager/notification"
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	dispatchManager = notification.DManager
	configManager   = config.CManager
)

type healthCheck struct {
}

func (h healthCheck) HealthCheck(gClient pluginpb.PluginClient, pluginInfo interface{}) (string, error) {
	isAlive := "UP"
	defaultPluginName := pluginInfo.(map[interface{}]interface{})["default_plugin_name"].(string)
	verify, err := gClient.Verify(context.Background(), new(emptypb.Empty))

	if err != nil || verify == nil {
		isAlive = "DOWN"
		jsonMessage := msg.ReqMsg{FuncName: "is_plugin_up", State: msg.Faliure, Msg: "is Down !!", Severity: msg.Critical, ResourceType: defaultPluginName}
		dispatchManager.SendNotification(jsonMessage)
	}

	return isAlive, nil
}

type HealthCheck interface {
	HealthCheck(gClient pluginpb.PluginClient, pluginInfo interface{}) (string, error)
}

func NewHealthChecker() HealthCheck {
	return &healthCheck{}
}
