package healthcheck

import (
	"context"
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	notif "github.com/dsrvlabs/vatz/manager/notification"
	"github.com/robfig/cron/v3"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// Alive status
const (
	AliveStatusUp   AliveStatus = "UP"
	AliveStatusDown AliveStatus = "DOWN"
)

// AliveStatus is type to describe liveness.
type AliveStatus string

var (
	dispatchManager = notif.GetDispatcher()
)

// HealthCheck is...
type HealthCheck interface {
	PluginHealthCheck(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin) (AliveStatus, error)
	VatzHealthCheck(schedule []string) error
}

type healthCheck struct {
}

func (h *healthCheck) PluginHealthCheck(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin) (AliveStatus, error) {
	// TODO: plugin parameter is used only for add notification message to get plugin's name.
	// But it could be better to get plugin's name from verify message.

	isAlive := AliveStatusUp
	verify, err := gClient.Verify(ctx, new(emptypb.Empty))

	// TODO: Do I have to send notification from here?
	if err != nil || verify == nil {
		isAlive = AliveStatusDown
		jsonMessage := notif.ReqMsg{
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

func (h *healthCheck) VatzHealthCheck(healthCheckerSchedule []string) error {
	c := cron.New(cron.WithLocation(time.UTC))
	jsonMessage := notif.ReqMsg{
		FuncName:     "vatzHealthCheck",
		State:        pluginpb.STATE_SUCCESS,
		Msg:          "VATZ is alive!.",
		Severity:     pluginpb.SEVERITY_INFO,
		ResourceType: "VATZ",
	}
	for i := 0; i < len(healthCheckerSchedule); i++ {
		c.AddFunc(healthCheckerSchedule[i], func() { dispatchManager.SendNotification(jsonMessage) })
	}
	c.Start()
	return nil
}

// NewHealthChecker creates new healthcheck instance
func NewHealthChecker() HealthCheck {
	return &healthCheck{}
}
