package notification

import (
	model "github.com/dsrvlabs/vatz/manager/model"

	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"

	config "github.com/dsrvlabs/vatz/manager/config"
)

var (
	notificationInstance Notification
	DManager             dispatcher_manager
	defaultConf          = config.CManager.GetYMLData("default.yaml", true)
	configManager        = config.CManager
	discordChannel       string
)

func init() {
	notificationInstance = NewDispatcher()
	//message.ConfigType
	protocolInfo := configManager.Parse(model.Protocol, defaultConf)
	notificationInfo := protocolInfo.(map[interface{}]interface{})["notification_info"]
	discordChannel = notificationInfo.(map[interface{}]interface{})["discord_secret"].(string)
}

type dispatcher_manager struct {
}

func (s *dispatcher_manager) GetNotifyInfo(response *pluginpb.ExecuteResponse, pluginName string, methodName string) map[interface{}]interface{} {
	return notificationInstance.GetNotifyInfo(response, pluginName, methodName)
}

func (s *dispatcher_manager) SendNotification(request model.ReqMsg) error {
	if request.Severity == model.Critical {
		err := notificationInstance.SendDiscord(request, discordChannel)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
