package types

import (
	"time"

	"github.com/dsrvlabs/vatz/manager/config"
)

// AliveStatus is aliveness of plugin.
type AliveStatus string

// AliveStatus is type that describes aliveness flags.
const (
	AliveStatusUp   AliveStatus = "UP"
	AliveStatusDown AliveStatus = "DOWN"
)

// PluginStatus describes detail status of plugin.
type PluginStatus struct {
	Plugin    config.Plugin `json:"plugin"`
	IsAlive   AliveStatus   `json:"is_alive"`
	LastCheck time.Time     `json:"last_check"`
}
