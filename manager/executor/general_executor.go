package executor

import (
	"context"
	"log"
	"strconv"
	"sync"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"google.golang.org/protobuf/types/known/structpb"
)

type executor struct {
	status sync.Map
}

func (s *executor) Execute(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatchers []dp.Dispatcher) error {
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
			"plugin_name":    plugin.Name,
			"plugin_port":    plugin.Port,
		}

		executeInfo, err := structpb.NewStruct(methodMap)
		if err != nil {
			log.Fatalf("failed to check command structpb: %v", err)
		}

		req := &pluginpb.ExecuteRequest{
			ExecuteInfo: executeInfo,
			Options:     options,
		}

		resp, err := s.execute(ctx, gClient, req)
		if err != nil {
			return err
		}

		firstExe, preStatus := s.updateState(plugin.Name, plugin.Port, resp)

		for _, dp := range dispatchers {
			dp.SetDispatcher(firstExe, preStatus, tp.NotifyInfo{
				Plugin:     plugin.Name,
				Method:     method.Name,
				Severity:   resp.GetSeverity(),
				State:      resp.GetState(),
				ExecuteMsg: resp.GetMessage(),
			})
		}
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

func (s *executor) updateState(pluginName string, port int, resp *pluginpb.ExecuteResponse) (bool, tp.StateFlag) {
	isFirstRun := false
	exeResp := tp.StateFlag{State: resp.GetState(), Severity: resp.GetSeverity()}
	pluginNPort := pluginName + strconv.Itoa(port)
	if _, ok := s.status.Load(pluginNPort); !ok {
		isFirstRun = true
		s.status.Store(pluginNPort, exeResp)
	} else {
		preStatus, _ := s.status.Load(pluginNPort)
		preVal := preStatus.(tp.StateFlag)
		if preVal.State != resp.State || preVal.Severity != resp.Severity {
			s.status.Store(pluginNPort, exeResp)
			exeResp = tp.StateFlag{State: preVal.State, Severity: preVal.Severity}
		} else {
			exeResp = tp.StateFlag{State: preVal.State, Severity: preVal.Severity}
		}
	}
	return isFirstRun, exeResp
}
