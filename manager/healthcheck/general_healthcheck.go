package healthcheck

import (
	"context"
	"sync"
	"time"

	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/robfig/cron/v3"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	healthCheckerOnce   = sync.Once{}
	healthCheckerSingle = healthChecker{}
)

type healthChecker struct {
	healthMSG    tp.ReqMsg
	pluginStatus sync.Map
}

func (h *healthChecker) PluginHealthCheck(ctx context.Context, gClient pb.PluginClient, plugin config.Plugin, dispatchers []dp.Dispatcher) (tp.AliveStatus, error) {
	isAlive := tp.AliveStatusUp
	sendMSG := false
	pUnique := utils.MakeUniqueValue(plugin.Name, plugin.Address, plugin.Port)
	verify, err := gClient.Verify(ctx, new(emptypb.Empty))

	deliverMSG := tp.ReqMsg{
		FuncName:     "isPluginUp",
		State:        pb.STATE_FAILURE,
		Msg:          "Plugin is DOWN!!",
		Severity:     pb.SEVERITY_CRITICAL,
		ResourceType: plugin.Name,
	}

	if _, ok := h.pluginStatus.Load(pUnique); !ok {
		if err != nil || verify == nil {
			isAlive = tp.AliveStatusDown
			sendMSG = true
		}
	} else {
		plStat, _ := h.pluginStatus.Load(pUnique)
		pStruct := plStat.(*tp.PluginStatus)
		if err != nil || verify == nil {
			isAlive = tp.AliveStatusDown
			if pStruct.IsAlive == tp.AliveStatusUp {
				sendMSG = true
			}
		} else {
			if pStruct.IsAlive == tp.AliveStatusDown {
				sendMSG = true
				deliverMSG.UpdateSeverity(pb.SEVERITY_INFO)
				deliverMSG.UpdateState(pb.STATE_SUCCESS)
				deliverMSG.UpdateMSG("Plugin is Alive.")
			}
		}
	}

	if sendMSG {
		for _, dispatcher := range dispatchers {
			dispatcher.SendNotification(deliverMSG)
		}
	}

	h.pluginStatus.Store(pUnique, &tp.PluginStatus{
		Plugin:    plugin,
		IsAlive:   isAlive,
		LastCheck: time.Now(),
	})

	return isAlive, nil
}

func (h *healthChecker) VATZHealthCheck(healthCheckerSchedule []string, dispatchers []dp.Dispatcher) error {
	c := cron.New(cron.WithLocation(time.UTC))
	for i := 0; i < len(healthCheckerSchedule); i++ {
		c.AddFunc(healthCheckerSchedule[i], func() {
			for _, dispatcher := range dispatchers {
				dispatcher.SendNotification(h.healthMSG)
			}
		})
	}
	c.Start()
	return nil
}

func (h *healthChecker) PluginStatus(ctx context.Context) []tp.PluginStatus {
	status := make([]tp.PluginStatus, 0)

	h.pluginStatus.Range(func(k, value any) bool {
		curStatus := value.(*tp.PluginStatus)
		status = append(status, *curStatus)
		return true
	})

	return status
}
