package executor

import (
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"vatz/manager/notification"
)

var (
	EManager         executorManager
	executorInstance Executor
	dispatchManager  = notification.DManager
)

func init() {
	executorInstance = NewExecutor()
}

type executorManager struct {
}

func (s *executorManager) updateStatus(resp *pluginpb.ExecuteResponse, methodName string, exeStatus map[interface{}]interface{}) error {
	if resp.GetState().String() != "SUCCESS" {
		exeStatus[methodName] = false
	}
	return nil
}

func (s *executorManager) Execute(gClient pluginpb.PluginClient, pluginInfo interface{}, exeStatus map[interface{}]interface{}) map[interface{}]interface{} {
	defaultPluginName := pluginInfo.(map[interface{}]interface{})["default_plugin_name"].(string)
	pluginAPIs := pluginInfo.(map[interface{}]interface{})["plugins"].([]interface{})

	//TODO: Find how to deal with multiple plugin methods.
	executeMethods := pluginAPIs[0].(map[interface{}]interface{})["executable_methods"].([]interface{})

	for _, method := range executeMethods {

		optionMap := map[string]interface{}{
			"plugin_name": defaultPluginName,
		}

		options, err := structpb.NewStruct(optionMap)
		if err != nil {
			log.Fatalf("failed to check target structpb: %v", err)
		}

		methodName := method.(map[interface{}]interface{})["method_name"].(string)

		//TODO: Please, add new logic to add param into Map.
		methodMap := map[string]interface{}{
			"execute_method": methodName,
		}

		executeInfo, err := structpb.NewStruct(methodMap)

		if err != nil {
			log.Fatalf("failed to check command structpb: %v", err)
		}

		if _, ok := exeStatus[methodName]; !ok {
			exeStatus[methodName] = true
		}

		req := &pluginpb.ExecuteRequest{
			ExecuteInfo: executeInfo,
			Options:     options,
		}

		resp, err := executorInstance.Execute(gClient, req)

		EManager.updateStatus(resp, methodName, exeStatus)

		notifyInfo := make(map[interface{}]interface{})
		notifyInfo["severity"] = resp.GetSeverity().String()
		notifyInfo["state"] = resp.GetState().String()
		notifyInfo["method_name"] = methodName
		notifyInfo["execute_message"] = resp.GetMessage()
		notifyInfo["plugin_name"] = defaultPluginName

		executorInstance.ExecuteNotify(notifyInfo, exeStatus)
	}
	return exeStatus
}
