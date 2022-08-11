package healthcheck

import (
	"context"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	tp "github.com/dsrvlabs/vatz/manager/types"
)

type HealthCheck interface {
	PluginHealthCheck(ctx context.Context, gClient pluginpb.PluginClient, dispatcher dp.Dispatcher) (tp.AliveStatus, error)
	VATZHealthCheck(schedule []string, dispatcher dp.Dispatcher) error
}

type healthChecker struct {
	healthMSG tp.ReqMsg
}
