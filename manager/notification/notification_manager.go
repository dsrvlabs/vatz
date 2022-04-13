package notification

import (
	"encoding/json"
	config "vatz/manager/config"
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

type ReqMsg struct {
	FuncName     string `json:"func_name"`
	State        string `json:"state"`
	Msg          string `json:"msg"`
	Severity     string `json:"severity"`
	ResourceType string `json:"resource_type"`
}

func (s *dispatcher_manager) SendNotification(request string) error {

	rMsg := ReqMsg{}
	err := json.Unmarshal([]byte(request), &rMsg)
	if err != nil {
		panic(err)
	}
	if rMsg.Severity == "CRITICAL" {
		err := notificationInstance.SendDiscord(rMsg, discordChannel)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
