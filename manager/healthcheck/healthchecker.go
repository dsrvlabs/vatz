package healthcheck

import (
	"context"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	notif "github.com/dsrvlabs/vatz/manager/notification"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/robfig/cron/v3"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (h *healthChecker) PluginHealthCheck(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatcher notif.Notification) (tp.AliveStatus, error) {
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

func (h *healthChecker) VATZHealthCheck(healthCheckerSchedule []string, dispatcher notif.Notification) error {
	c := cron.New(cron.WithLocation(time.UTC))
	for i := 0; i < len(healthCheckerSchedule); i++ {
		c.AddFunc(healthCheckerSchedule[i], func() { dispatcher.SendNotification(h.healthMSG) })
	}
	c.Start()
	return nil
}

type healthChecker struct {
	healthMSG tp.ReqMsg
}
