package notification

import (
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/manager/model"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
)

var (
	notificationInstance Notification
	DManager             dispatcher_manager
)

func init() {
	notificationInstance = NewDispatcher()
}

// TODO: Rename this.
type dispatcher_manager struct {
}

func (s *dispatcher_manager) GetNotifyInfo(response *pluginpb.ExecuteResponse, pluginName string, methodName string) map[interface{}]interface{} {
	return notificationInstance.GetNotifyInfo(response, pluginName, methodName)
}

func (s *dispatcher_manager) SendNotification(request model.ReqMsg) error {
	cfg := config.GetConfig()

	err := notificationInstance.SendDiscord(request, cfg.Vatz.NotificationInfo.DiscordSecret)
	if err != nil {
		panic(err)
	}

	return nil
}
