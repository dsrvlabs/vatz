package executor

import (
	"context"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	notif "github.com/dsrvlabs/vatz/manager/notification"
	"sync"
)

// Executor provides interfaces to execute plugin features.
type Executor interface {
	Execute(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatcher notif.Notification) error
}

type executor struct {
	status sync.Map
}
