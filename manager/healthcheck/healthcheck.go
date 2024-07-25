package healthcheck

import (
	"context"
	"github.com/dsrvlabs/vatz/types"
	"sync"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
)

// HealthCheck provides interfaces to check health.
type HealthCheck interface {
	PluginHealthCheck(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatcher []dp.Dispatcher) (types.AliveStatus, error)
	VATZHealthCheck(schedule []string, dispatcher []dp.Dispatcher) error
	PluginStatus(ctx context.Context) []types.PluginStatus
}

// GetHealthChecker creates instance of HealthCheck
func GetHealthChecker() HealthCheck {
	healthCheckerOnce.Do(func() {
		option := map[string]interface{}{"pUnique": "VATZHealthChecker"}
		healthCheckerSingle = healthChecker{
			healthMSG: types.ReqMsg{
				FuncName:     "VATZHealthCheck",
				State:        pluginpb.STATE_SUCCESS,
				Msg:          "VATZ is Alive!!",
				Severity:     pluginpb.SEVERITY_INFO,
				ResourceType: "VATZ",
				Options:      option,
			},
			pluginStatus: sync.Map{},
		}
	})

	return &healthCheckerSingle
}
