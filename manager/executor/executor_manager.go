package executor

import (
	"context"
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	message "vatz/manager/model"
	"vatz/manager/notification"
)

var (
	EManager         executorManager
	executorInstance Executor
	dispatchManager  = notification.DManager
)

func init() {
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

		// There's invalid memory trouble
		// panic: runtime error: invalid memory address or nil pointer dereference
		// [signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x1445d4c]
		//resp, err := executorInstance.Execute(gClient, req)

		resp, exeErr := gClient.Execute(context.Background(), req)

		if exeErr != nil || resp == nil {
			resp = &pluginpb.ExecuteResponse{
				State:     pluginpb.State_FAILURE,
				Message:   "API Execution Failed",
				AlertType: pluginpb.ALERT_TYPE_DISCORD,
				Severity:  pluginpb.Severity_ERROR,
			}
		}

		EManager.updateStatus(resp, methodName, exeStatus)
		notifyInfo := make(map[interface{}]interface{})

		notifyInfo["severity"] = resp.GetSeverity().String()
		notifyInfo["state"] = resp.GetState().String()
		notifyInfo["method_name"] = methodName
		notifyInfo["execute_message"] = resp.GetMessage()
		notifyInfo["plugin_name"] = defaultPluginName

		EManager.ExecuteNotify(notifyInfo, exeStatus)
		//executorInstance.ExecuteNotify(notifyInfo, exeStatus)

	}

	return exeStatus
}

func (s *executorManager) ExecuteNotify(notifyInfo map[interface{}]interface{}, exeStatus map[interface{}]interface{}) error {

	// if response's state is not `SUCCESS` and then we consider all execute call has failed.
	if notifyInfo["state"] != string(message.Success) {
		exeStatus[notifyInfo["method_name"]] = false

		if notifyInfo["severity"] == string(message.Error) {
			jsonMessage := message.ReqMsg{FuncName: notifyInfo["method_name"].(string), State: message.Faliure, Msg: "No response from Plugin", Severity: message.Critical, ResourceType: notifyInfo["plugin_name"].(string)}
			dispatchManager.SendNotification(jsonMessage)
		}

		if notifyInfo["severity"] == string(message.Critical) {
			jsonMessage := message.ReqMsg{FuncName: notifyInfo["method_name"].(string), State: message.Faliure, Msg: notifyInfo["execute_message"].(string), Severity: message.Critical, ResourceType: notifyInfo["plugin_name"].(string)}
			dispatchManager.SendNotification(jsonMessage)
		}

	} else {
		if exeStatus[notifyInfo["method_name"]] == false {
			jsonMessage := message.ReqMsg{FuncName: notifyInfo["method_name"].(string), State: message.Success, Msg: notifyInfo["execute_message"].(string), Severity: message.Info, ResourceType: notifyInfo["plugin_name"].(string)}
			dispatchManager.SendNotification(jsonMessage)
		}
	}

	return nil
}
