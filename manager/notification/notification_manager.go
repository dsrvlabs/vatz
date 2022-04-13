package notification

import (
	config "vatz/manager/config"
	message "vatz/manager/message"
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
	protocolInfo := configManager.Parse("PROTOCOL", defaultConf)
	notificationInfo := protocolInfo.(map[interface{}]interface{})["notification_info"]
	discordChannel = notificationInfo.(map[interface{}]interface{})["discord_secret"].(string)
}

type dispatcher_manager struct {
}

func (s *dispatcher_manager) SendNotification(request message.ReqMsg) error {
	if request.Severity == "CRITICAL" {
		err := notificationInstance.SendDiscord(request, discordChannel)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
