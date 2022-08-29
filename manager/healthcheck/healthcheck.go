package healthcheck

import (
	"context"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	tp "github.com/dsrvlabs/vatz/manager/types"
)

// HealthCheck provides interfaces to check health.
type HealthCheck interface {
	PluginHealthCheck(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatcher []dp.Dispatcher) (tp.AliveStatus, error)
	VATZHealthCheck(schedule []string, dispatcher []dp.Dispatcher) error

	PluginStatus(ctx context.Context) []tp.PluginStatus
}

func NewHealthChecker() HealthCheck {
	return &healthChecker{
		healthMSG: tp.ReqMsg{
			FuncName:     "VATZHealthCheck",
			State:        pluginpb.STATE_SUCCESS,
			Msg:          "VATZ is Alive!!",
			Severity:     pluginpb.SEVERITY_INFO,
			ResourceType: "VATZ",
		},
		pluginStatus: map[string]tp.PluginStatus{},
	}
}
