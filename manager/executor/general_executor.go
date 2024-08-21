package executor

import (
	"context"
	tp "github.com/dsrvlabs/vatz/types"
	"os"
	"sync"

	"github.com/rs/zerolog"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/structpb"
)

type executor struct {
	status sync.Map
}

func (s *executor) Execute(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatchers []dp.Dispatcher) error {
	executeMethods := plugin.ExecutableMethods
	log.Debug().Str("module", "executor").Msg("Execute is called.")
	for _, method := range executeMethods {
		optionMap := map[string]interface{}{
			"plugin_name": plugin.Name,
		}

		options, err := structpb.NewStruct(optionMap)
		if err != nil {
			log.Error().Str("module", "executor").Msgf("failed to check target structpb: %v", err)
			os.Exit(1)
		}

		//TODO: Please, add new logic to add param into Map.

		methodMap := map[string]interface{}{
			"execute_method": method.Name,
		}

		executeInfo, err := structpb.NewStruct(methodMap)
		if err != nil {
			log.Error().Str("module", "executor").Msgf("failed to check target structpb: %v", err)
			os.Exit(1)
		}

		req := &pluginpb.ExecuteRequest{
			ExecuteInfo: executeInfo,
			Options:     options,
		}

		if zerolog.GlobalLevel() == zerolog.DebugLevel {
			log.Debug().Str("module", "executor").Msgf("request (Plugin Name: %s, Method Name: %s)", plugin.Name, method.Name)
		} else {
			log.Info().Str("module", "executor").Msgf("Executor send request to %s", plugin.Name)
		}

		resp, err := s.execute(ctx, gClient, req)
		if err != nil {
			return err
		}

		pUnique := utils.MakeUniqueValue(plugin.Name, plugin.Address, plugin.Port)
		firstExe, preStatus := s.updateState(pUnique, resp)

		for _, dpSingle := range dispatchers {
			err = dpSingle.SetDispatcher(firstExe, method.Flag, preStatus, tp.NotifyInfo{
				Plugin:     plugin.Name,
				Method:     method.Name,
				Address:    plugin.Address,
				Port:       plugin.Port,
				Severity:   resp.GetSeverity(),
				State:      resp.GetState(),
				ExecuteMsg: resp.GetMessage(),
			})
			if err != nil {
				log.Error().Str("module", "dispatcher").Msgf("failed to set dispatcher: %v", err)
			}
		}
	}

	return nil
}

func (s *executor) execute(ctx context.Context, gClient pluginpb.PluginClient, in *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error) {
	log.Debug().Str("module", "executor").Msgf("func execute")
	resp, err := gClient.Execute(ctx, in)
	if err != nil || resp == nil {
		return &pluginpb.ExecuteResponse{
			State:        pluginpb.STATE_FAILURE,
			Message:      "API Execution Failed",
			Severity:     pluginpb.SEVERITY_ERROR,
			ResourceType: "ResourceType Unknown",
		}, nil
	}

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		log.Debug().Str("module", "executor").Msgf("response (res message:%s, res State: %s) ", resp.Message, resp.State)
	} else {
		log.Info().Str("module", "executor").Msgf("response: %s", resp.State)
	}

	return resp, err
}

func (s *executor) updateState(unique string, resp *pluginpb.ExecuteResponse) (bool, tp.StateFlag) {
	log.Debug().Str("module", "executor").Msgf("func updateState")
	isFirstRun := false
	exeResp := tp.StateFlag{State: resp.GetState(), Severity: resp.GetSeverity()}
	if _, ok := s.status.Load(unique); !ok {
		isFirstRun = true
		s.status.Store(unique, exeResp)
	} else {
		preStatus, _ := s.status.Load(unique)
		preVal := preStatus.(tp.StateFlag)
		if preVal.State != resp.State || preVal.Severity != resp.Severity {
			s.status.Store(unique, exeResp)
			exeResp = tp.StateFlag{State: preVal.State, Severity: preVal.Severity}
		} else {
			exeResp = tp.StateFlag{State: preVal.State, Severity: preVal.Severity}
		}
	}
	return isFirstRun, exeResp
}
