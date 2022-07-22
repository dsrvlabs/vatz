package healthcheck

import (
	"context"
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/manager/notification"
	msg "github.com/dsrvlabs/vatz/manager/notification"
	"github.com/robfig/cron/v3"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	dispatchManager = notification.GetDispatcher()
)

type healthCheck struct {
}

func (h *healthCheck) PluginHealthCheck(gClient pluginpb.PluginClient, plugin config.Plugin) (string, error) {
	isAlive := "UP"
	verify, err := gClient.Verify(context.Background(), new(emptypb.Empty))

	if err != nil || verify == nil {
		isAlive = "DOWN"
		jsonMessage := msg.ReqMsg{
			FuncName:     "isPluginUp",
			State:        pluginpb.STATE_FAILURE,
			Msg:          "Plugin is DOWN!!",
			Severity:     pluginpb.SEVERITY_CRITICAL,
			ResourceType: plugin.Name,
		}

		dispatchManager.SendNotification(jsonMessage)
	}

	return isAlive, nil
}

func (v *healthCheck) VatzHealthCheck(HealthCheckerSchedule []string) error {
	c := cron.New(cron.WithLocation(time.UTC))
	jsonMessage := msg.ReqMsg{
		FuncName:     "vatzHealthCheck",
		State:        pluginpb.STATE_SUCCESS,
		Msg:          "VATZ is alive!.",
		Severity:     pluginpb.SEVERITY_INFO,
		ResourceType: "VATZ",
	}
	for i := 0; i < len(HealthCheckerSchedule); i++ {
		c.AddFunc(HealthCheckerSchedule[i], func() { dispatchManager.SendNotification(jsonMessage) })
	}
	c.Start()
	return nil
}

type HealthCheck interface {
	PluginHealthCheck(gClient pluginpb.PluginClient, plugin config.Plugin) (string, error)
	VatzHealthCheck(schedule []string) error
}

func NewHealthChecker() HealthCheck {
	return &healthCheck{}
}
