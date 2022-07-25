package healthcheck

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/manager/notification"
	msg "github.com/dsrvlabs/vatz/manager/notification"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	dispatchManager = notification.GetDispatcher()
)

type healthCheck struct {
	isSending map[string]bool
}

func (h *healthCheck) PluginHealthCheck(gClient pluginpb.PluginClient, plugin config.Plugin) (bool, error) {
	// TODO: Magic value always wrong.
	isAlive := true
	verify, err := gClient.Verify(context.Background(), new(emptypb.Empty))

	if err != nil || verify == nil {
		isAlive = false
		if !h.isSending[plugin.Name] {
			jsonMessage := msg.ReqMsg{
				FuncName:     "is_plugin_up",
				State:        pluginpb.STATE_FAILURE,
				Msg:          "is Down !!",
				Severity:     pluginpb.SEVERITY_CRITICAL,
				ResourceType: plugin.Name,
			}
			dispatchManager.SendNotification(jsonMessage)
			h.isSending[plugin.Name] = true
		}
	} else {
		h.isSending[plugin.Name] = false
	}

	return isAlive, nil
}

func (v *healthCheck) VatzHealthCheck(HealthCheckerSchedule []string) error {
	c := cron.New(cron.WithLocation(time.UTC))
	jsonMessage := msg.ReqMsg{
		FuncName:     "vatz_healthcheck",
		State:        pluginpb.STATE_SUCCESS,
		Msg:          "Vatz is alive!",
		Severity:     pluginpb.SEVERITY_CRITICAL,
		ResourceType: "Vatz",
	}
	for i := 0; i < len(HealthCheckerSchedule); i++ {
		log.Info().Str("module", "VatzHealthCheck").Msgf("%d, %s", i, HealthCheckerSchedule[i])
		c.AddFunc(HealthCheckerSchedule[i], func() { dispatchManager.SendNotification(jsonMessage) })
	}
	c.Start()
	return nil
}

type HealthCheck interface {
	PluginHealthCheck(gClient pluginpb.PluginClient, plugin config.Plugin) (bool, error)
	VatzHealthCheck(schedule []string) error
}

func NewHealthChecker() HealthCheck {
	return &healthCheck{
		isSending: map[string]bool{},
	}
}
