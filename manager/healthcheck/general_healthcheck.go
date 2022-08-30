package healthcheck

import (
	"context"
	"sync"
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	tp "github.com/dsrvlabs/vatz/manager/types"
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

func (h *healthChecker) PluginHealthCheck(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatchers []dp.Dispatcher) (tp.AliveStatus, error) {
	isAlive := tp.AliveStatusUp
	verify, err := gClient.Verify(ctx, new(emptypb.Empty))
	if err != nil || verify == nil {
		isAlive = tp.AliveStatusDown
		failErrorMessage := tp.ReqMsg{
			FuncName:     "isPluginUp",
			State:        pluginpb.STATE_FAILURE,
			Msg:          "Plugin is DOWN!!",
			Severity:     pluginpb.SEVERITY_CRITICAL,
			ResourceType: plugin.Name,
		}

		for _, dispatcher := range dispatchers {
			dispatcher.SendNotification(failErrorMessage)
		}
	}

	h.pluginStatus.Store(plugin.Name, &tp.PluginStatus{
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
