package healthcheck

import (
	"context"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/manager/notification"
	msg "github.com/dsrvlabs/vatz/manager/notification"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	dispatchManager = notification.GetDispatcher()
	isSending       = false
)

type healthCheck struct {
}

func (h healthCheck) HealthCheck(gClient pluginpb.PluginClient, plugin config.Plugin) (string, error) {
	// TODO: Magic value always wrong.
	isAlive := "UP"
	verify, err := gClient.Verify(context.Background(), new(emptypb.Empty))

	if err != nil || verify == nil {
		if !isSending {
			isAlive = "DOWN"
			jsonMessage := msg.ReqMsg{
				FuncName:     "is_plugin_up",
				State:        pluginpb.STATE_FAILURE,
				Msg:          "is Down !!",
				Severity:     pluginpb.SEVERITY_CRITICAL,
				ResourceType: plugin.Name,
			}
			dispatchManager.SendNotification(jsonMessage)
			isSending = true
		}
	} else {
		isSending = false
	}

	return isAlive, nil
}

type HealthCheck interface {
	HealthCheck(gClient pluginpb.PluginClient, plugin config.Plugin) (string, error)
}

func NewHealthChecker() HealthCheck {
	return &healthCheck{}
}
