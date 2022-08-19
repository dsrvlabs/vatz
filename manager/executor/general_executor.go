package executor

import (
	"context"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"sync"
)

type executor struct {
	status sync.Map
}

func (s *executor) Execute(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatcher dp.Dispatcher) error {
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

		if _, ok := s.status.Load(method.Name); !ok {
			s.status.Store(method.Name, true)
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
			s.status.Store(method.Name, false)
		}

		notifyInfo := tp.NotifyInfo{
			Plugin:     plugin.Name,
			Method:     method.Name,
			Severity:   resp.GetSeverity(),
			State:      resp.GetState(),
			ExecuteMsg: resp.GetMessage(),
		}

		s.executeNotify(notifyInfo, dispatcher)
	}

	return nil
}

func (s *executor) execute(ctx context.Context, gClient pluginpb.PluginClient, in *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error) {
	resp, err := gClient.Execute(ctx, in)
	if err != nil || resp == nil {
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

//executeNotify function has to be moved to dispatcher.
func (s *executor) executeNotify(notifyInfo tp.NotifyInfo, dispatcher dp.Dispatcher) error {
	// if response's state is not `SUCCESS` and then we consider all execute call has failed.
	methodName := notifyInfo.Method

	if notifyInfo.State != pluginpb.STATE_SUCCESS {
		s.status.Store(methodName, false)
		if notifyInfo.Severity == pluginpb.SEVERITY_ERROR {
			jsonMessage := tp.ReqMsg{
				FuncName:     notifyInfo.Method,
				State:        pluginpb.STATE_FAILURE,
				Msg:          "No response from Plugin",
				Severity:     pluginpb.SEVERITY_CRITICAL, // TODO: Error or Critical?
				ResourceType: notifyInfo.Plugin,
			}

			dispatcher.SendNotification(jsonMessage)
		} else if notifyInfo.Severity == pluginpb.SEVERITY_CRITICAL {
			jsonMessage := tp.ReqMsg{
				FuncName:     notifyInfo.Method,
				State:        pluginpb.STATE_FAILURE,
				Msg:          notifyInfo.ExecuteMsg,
				Severity:     pluginpb.SEVERITY_CRITICAL,
				ResourceType: notifyInfo.Plugin,
			}
			dispatcher.SendNotification(jsonMessage)
		}
	} else {
		if status, ok := s.status.Load(methodName); ok && status == false {
			jsonMessage := tp.ReqMsg{
				FuncName:     notifyInfo.Method,
				State:        pluginpb.STATE_SUCCESS,
				Msg:          notifyInfo.ExecuteMsg,
				Severity:     pluginpb.SEVERITY_INFO,
				ResourceType: notifyInfo.Plugin,
			}
			dispatcher.SendNotification(jsonMessage)
			s.status.Store(methodName, true)
		}
	}
	return nil
}
