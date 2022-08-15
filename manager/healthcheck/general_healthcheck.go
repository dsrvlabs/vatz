package healthcheck

import (
	"context"
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/robfig/cron/v3"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (h *healthChecker) PluginHealthCheck(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatcher dp.Dispatcher) (tp.AliveStatus, error) {
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
		dispatcher.SendNotification(failErrorMessage)
	}
	return isAlive, nil
}

func (h *healthChecker) VATZHealthCheck(healthCheckerSchedule []string, dispatcher dp.Dispatcher) error {
	c := cron.New(cron.WithLocation(time.UTC))
	for i := 0; i < len(healthCheckerSchedule); i++ {
		c.AddFunc(healthCheckerSchedule[i], func() { dispatcher.SendNotification(h.healthMSG) })
	}
	c.Start()
	return nil
}

// NewHealthChecker creates instance of HealchChecker
func NewHealthChecker() HealthCheck {
	return &healthChecker{
		healthMSG: tp.ReqMsg{
			FuncName:     "VATZHealthCheck",
			State:        pluginpb.STATE_SUCCESS,
			Msg:          "VATZ is Alive!!",
			Severity:     pluginpb.SEVERITY_INFO,
			ResourceType: "VATZ",
		},
	}
}
