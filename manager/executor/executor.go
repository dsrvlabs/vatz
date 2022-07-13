package executor

import (
	"context"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	message "github.com/dsrvlabs/vatz/manager/notification"
)

type executor struct {
}

func (v executor) Execute(gClient pluginpb.PluginClient, in *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error) {
	resp, err := gClient.Execute(context.Background(), in)
	if err != nil || resp == nil {
		return &pluginpb.ExecuteResponse{
			State:        pluginpb.STATE_FAILURE,
			Message:      "API Execution Failed",
			AlertType:    []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD, pluginpb.ALERT_TYPE_PAGER_DUTY},
			Severity:     pluginpb.SEVERITY_ERROR,
			ResourceType: "ResourceType Unknown",
		}, nil
	}
	return resp, nil
}

func (v executor) ExecuteNotify(notifyInfo map[interface{}]interface{}, exeStatus map[interface{}]interface{}) error {
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
			exeStatus[notifyInfo["method_name"]] = true
		}
	}
	return nil
}

func (v executor) GenerateExecutorRequest() error {
	return nil
}

type Executor interface {
	GenerateExecutorRequest() error
	Execute(gClient pluginpb.PluginClient, in *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error)
	ExecuteNotify(notifyInfo map[interface{}]interface{}, exeStatus map[interface{}]interface{}) error
}

func NewExecutor() Executor {
	return &executor{}
}
