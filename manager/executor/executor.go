package executor

import (
	"context"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	"sync"
)

// Executor provides interfaces to execute plugin features.
type Executor interface {
	Execute(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatcher dp.Dispatcher) error
}

type executor struct {
	status sync.Map
}
