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
	Execute(ctx context.Context, gClient pluginpb.PluginClient, plugin config.Plugin, dispatchers []dp.Dispatcher) error
}

// NewExecutor create new executor instance.
func NewExecutor() Executor {
	return &executor{
		status: sync.Map{},
	}
}
