package executor

import (
	"log"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/manager/notification"
	"google.golang.org/protobuf/types/known/structpb"
	//"vatz/manager/config"
)

var (
	EManager         executorManager
	executorInstance Executor
	dispatchManager  = notification.GetDispatcher()
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

func (s *executorManager) Execute(gClient pluginpb.PluginClient, plugin config.Plugin, exeStatus map[interface{}]interface{}) map[interface{}]interface{} {
	//TODO: Find how to deal with multiple plugin methods.
	executeMethods := plugin.ExecutableMethods

	for _, method := range executeMethods {
		optionMap := map[string]interface{}{
			"plugin_name": plugin.Name,
		}

		options, err := structpb.NewStruct(optionMap)
		if err != nil {
			log.Fatalf("failed to check target structpb: %v", err)
		}

		//TODO: Please, add new logic to add param into Map.
		methodMap := map[string]interface{}{
			"execute_method": method.Name,
		}

		executeInfo, err := structpb.NewStruct(methodMap)

		if err != nil {
			log.Fatalf("failed to check command structpb: %v", err)
		}

		if _, ok := exeStatus[method.Name]; !ok {
			exeStatus[method.Name] = true
		}

		req := &pluginpb.ExecuteRequest{
			ExecuteInfo: executeInfo,
			Options:     options,
		}

		resp, _ := executorInstance.Execute(gClient, req)
		EManager.updateStatus(resp, method.Name, exeStatus)
		notifyInfo := dispatchManager.GetNotifyInfo(resp, plugin.Name, method.Name)

		// TODO: Temporarily convert data type.
		// This part should be removed on refactoring issue of executor #179.
		temp := make(map[interface{}]interface{})
		temp["severity"] = notifyInfo.Severity.String()
		temp["state"] = notifyInfo.State.String()
		temp["method_name"] = notifyInfo.Method
		temp["execute_message"] = notifyInfo.ExecuteMsg
		temp["plugin_name"] = notifyInfo.Plugin

		executorInstance.ExecuteNotify(temp, exeStatus)
	}

	return exeStatus
}
