package executor

import (
	"context"
	"fmt"
	"log"
	message "vatz/manager/message"
	"vatz/manager/notification"

	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	EManager        executor_manager
	dispatchManager = notification.DManager
)

func init() {
}

type executor_manager struct {
}

func (s *executor_manager) Execute(pluginInfo interface{}, gClient pluginpb.PluginClient) error {
	defaultPluginName := pluginInfo.(map[interface{}]interface{})["defult_plugin_name"].(string)
	pluginAPIs := pluginInfo.(map[interface{}]interface{})["plugins"].([]interface{})

	for _, api := range pluginAPIs {
		executeMethods := api.(map[interface{}]interface{})["executable_apis"].([]interface{})
		for _, exe := range executeMethods {
			targetMap := map[string]interface{}{
				"source": "localhost:9091",
			}

			target, err := structpb.NewStruct(targetMap)
			if err != nil {
				log.Fatalf("failed to check target structpb: %v", err)
			}

			methodName := exe.(map[interface{}]interface{})["method_name"].(string)
			commandMap := map[string]interface{}{
				"command": methodName,
			}

			commands, err := structpb.NewStruct(commandMap)
			if err != nil {
				log.Fatalf("failed to check command structpb: %v", err)
			}

			req := &pluginpb.ExecuteRequest{
				ExecuteInfo: commands,
				Options:     target,
			}

			resp, err := gClient.Execute(context.Background(), req)

			if err != nil || resp == nil {
				jsonMessage := message.ReqMsg{methodName, "FAILURE", "No response from Plugin", "CRITICAL", defaultPluginName}
				dispatchManager.SendNotification(jsonMessage)
			}

			if resp.GetSeverity().String() == "CRITICAL" {
				jsonMessage := message.ReqMsg{methodName, "FAILURE", resp.GetMessage(), "CRITICAL", defaultPluginName}
				fmt.Println(jsonMessage)
				dispatchManager.SendNotification(jsonMessage)
			}
		}
	}

	return nil
}
