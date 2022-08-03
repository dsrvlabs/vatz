package healthcheck

import (
	"context"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	notif "github.com/dsrvlabs/vatz/manager/notification"
	tp "github.com/dsrvlabs/vatz/manager/types"
)

type HealthCheck interface {
	PluginHealthCheck(ctx context.Context, gClient pluginpb.PluginClient, dispatcher notif.Notification) (tp.AliveStatus, error)
	VATZHealthCheck(schedule []string, dispatcher notif.Notification) error
}

func NewHealthChecker() *healthChecker {
	return &healthChecker{
		healthMSG: tp.ReqMsg{
			FuncName:     "VATZHealthCheck",
			State:        pluginpb.STATE_SUCCESS,
			Msg:          "VATZ is alive!.",
			Severity:     pluginpb.SEVERITY_INFO,
			ResourceType: "VATZ",
		},
	}
}
