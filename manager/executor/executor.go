package executor

import (
	"context"
	"log"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/manager/notification"
	message "github.com/dsrvlabs/vatz/manager/notification"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	dispatchManager = notification.GetDispatcher()
)

func init() {
}

// Executor provides interfaces to execute plugin features.
type Executor interface {
	Execute(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin) error
}

type executor struct {
	status map[string]bool
}

func (s *executor) Execute(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin) error {
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

		if _, ok := s.status[method.Name]; !ok {
			s.status[method.Name] = true
		}

		req := &pluginpb.ExecuteRequest{
			ExecuteInfo: executeInfo,
			Options:     options,
		}

		resp, err := s.execute(ctx, gClient, req)
		if err != nil {
			return err
		}

		if resp.GetState() != pluginpb.STATE_SUCCESS {
			s.status[method.Name] = false
		}

		notifyInfo := dispatchManager.GetNotifyInfo(resp, plugin.Name, method.Name)
		s.executeNotify(notifyInfo)
	}

	return nil
}

func (s *executor) execute(ctx context.Context, gClient pluginpb.PluginClient, in *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error) {
	resp, err := gClient.Execute(ctx, in)
	if err != nil || resp == nil {
		// TODO: why below codes chnage response?
		return &pluginpb.ExecuteResponse{
			State:        pluginpb.STATE_FAILURE,
			Message:      "API Execution Failed",
			AlertType:    []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD, pluginpb.ALERT_TYPE_PAGER_DUTY},
			Severity:     pluginpb.SEVERITY_ERROR,
			ResourceType: "ResourceType Unknown",
		}, nil
	}
	return resp, err
}

func (s *executor) executeNotify(notifyInfo message.NotifyInfo) error {
	// if response's state is not `SUCCESS` and then we consider all execute call has failed.
	methodName := notifyInfo.Method

	if notifyInfo.State != pluginpb.STATE_SUCCESS {
		s.status[methodName] = false
		if notifyInfo.Severity == pluginpb.SEVERITY_ERROR {
			jsonMessage := message.ReqMsg{
				FuncName:     notifyInfo.Method,
				State:        pluginpb.STATE_FAILURE,
				Msg:          "No response from Plugin",
				Severity:     pluginpb.SEVERITY_CRITICAL, // TODO: Error or Critical?
				ResourceType: notifyInfo.Plugin,
			}

			dispatchManager.SendNotification(jsonMessage)
		} else if notifyInfo.Severity == pluginpb.SEVERITY_CRITICAL {
			jsonMessage := message.ReqMsg{
				FuncName:     notifyInfo.Method,
				State:        pluginpb.STATE_FAILURE,
				Msg:          notifyInfo.ExecuteMsg,
				Severity:     pluginpb.SEVERITY_CRITICAL,
				ResourceType: notifyInfo.Plugin,
			}
			dispatchManager.SendNotification(jsonMessage)
		}
	} else {
		if s.status[methodName] == false {
			jsonMessage := message.ReqMsg{
				FuncName:     notifyInfo.Method,
				State:        pluginpb.STATE_SUCCESS,
				Msg:          notifyInfo.ExecuteMsg,
				Severity:     pluginpb.SEVERITY_INFO,
				ResourceType: notifyInfo.Plugin,
			}

			dispatchManager.SendNotification(jsonMessage)
			s.status[methodName] = true
		}
	}
	return nil
}

// NewExecutor create new executor instance.
func NewExecutor() Executor {
	return &executor{
		status: map[string]bool{},
	}
}
