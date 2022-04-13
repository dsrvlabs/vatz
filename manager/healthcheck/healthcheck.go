package healthcheck

import (
	"context"
	"vatz/manager/config"
	message "vatz/manager/message"
	"vatz/manager/notification"

	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	dispatchManager = notification.DManager
	configManager   = config.CManager
)

type healthCheck struct {
}

func (h healthCheck) HealthCheck(pluginInfo interface{}, gClient pluginpb.PluginClient) (string, error) {
	isAlive := "UP"
	defaultPluginName := pluginInfo.(map[interface{}]interface{})["defult_plugin_name"].(string)
	verify, err := gClient.Verify(context.Background(), new(emptypb.Empty))
	if err != nil || verify == nil {
		isAlive = "DOWN"
		jsonMessage := message.ReqMsg{"is_plugin_up", "FAILURE", "is Down !!", "CRITICAL", defaultPluginName}
		dispatchManager.SendNotification(jsonMessage)
	}
	return isAlive, nil
}

type HealthCheck interface {
	HealthCheck(pluginInfo interface{}, gClient pluginpb.PluginClient) (string, error)
}

func NewHealthChecker() HealthCheck {
	return &healthCheck{}
}
