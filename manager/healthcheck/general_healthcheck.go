package healthcheck

import (
	"context"
	"fmt"
	"github.com/dsrvlabs/vatz/types"
	"sync"
	"time"

	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	healthCheckerOnce   = sync.Once{}
	healthCheckerSingle = healthChecker{}
)

type healthChecker struct {
	healthMSG    types.ReqMsg
	pluginStatus sync.Map
}

func (h *healthChecker) PluginHealthCheck(ctx context.Context, gClient pb.PluginClient, plugin config.Plugin, dispatchers []dp.Dispatcher) (types.AliveStatus, error) {
	isAlive := types.AliveStatusUp
	sendMSG := false
	pUnique := utils.MakeUniqueValue(plugin.Name, plugin.Address, plugin.Port)
	verify, err := gClient.Verify(ctx, new(emptypb.Empty))

	option := map[string]interface{}{"pUnique": pUnique}

	deliverMSG := types.ReqMsg{
		FuncName:     "isPluginUp",
		State:        pb.STATE_FAILURE,
		Msg:          "Plugin is DOWN!!",
		Severity:     pb.SEVERITY_CRITICAL,
		ResourceType: plugin.Name,
		Options:      option,
	}

	if _, ok := h.pluginStatus.Load(pUnique); !ok {
		if err != nil || verify == nil {
			isAlive = types.AliveStatusDown
			sendMSG = true
		}
	} else {
		plStat, _ := h.pluginStatus.Load(pUnique)
		pStruct := plStat.(*types.PluginStatus)
		if err != nil || verify == nil {
			isAlive = types.AliveStatusDown
			if pStruct.IsAlive == types.AliveStatusUp {
				sendMSG = true
			}
		} else {
			if pStruct.IsAlive == types.AliveStatusDown {
				sendMSG = true
				deliverMSG.UpdateSeverity(pb.SEVERITY_INFO)
				deliverMSG.UpdateState(pb.STATE_SUCCESS)
				deliverMSG.UpdateMSG("Plugin is Alive.")
			}
		}
	}

	if sendMSG {
		errorCount := 0
		for _, dispatcher := range dispatchers {
			sendNotificationError := dispatcher.SendNotification(deliverMSG)
			if sendNotificationError != nil {
				log.Error().Str("module", "healthcheck").Msgf("failed to send notification: %v", err)
				errorCount = errorCount + 1
			}
		}

		if len(dispatchers) == errorCount {
			log.Error().Str("module", "healthcheck").Msg("All Dispatchers failed to send a notifications, Please, Check your dispatcher configs.")
			return isAlive, fmt.Errorf("Failed to send all configured notifications. ")
		}
	}

	h.pluginStatus.Store(pUnique, &types.PluginStatus{
		Plugin:    plugin,
		IsAlive:   isAlive,
		LastCheck: time.Now(),
	})

	return isAlive, nil
}

// VATZHealthCheck send a notification at a specific time that the vatz is alive.
func (h *healthChecker) VATZHealthCheck(healthCheckerSchedule []string, dispatchers []dp.Dispatcher) error {
	c := cron.New(cron.WithLocation(time.UTC))
	for i := 0; i < len(healthCheckerSchedule); i++ {
		_, err := c.AddFunc(healthCheckerSchedule[i], func() {
			for _, dispatcher := range dispatchers {
				err := dispatcher.SendNotification(h.healthMSG)
				if err != nil {
					log.Error().Str("module", "dispatcher").Msgf("failed to send notification: %v", err)
				}
			}
		})
		if err != nil {
			log.Error().Str("module", "healthcheck").Msgf("failed to add function to cron: %v", err)
		}
	}
	c.Start()
	return nil
}

func (h *healthChecker) PluginStatus(ctx context.Context) []types.PluginStatus {
	status := make([]types.PluginStatus, 0)

	h.pluginStatus.Range(func(k, value any) bool {
		curStatus := value.(*types.PluginStatus)
		status = append(status, *curStatus)
		return true
	})

	return status
}
