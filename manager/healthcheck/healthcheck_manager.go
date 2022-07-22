package healthcheck

import (
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
)

var (
	healthCheckInstance HealthCheck
	HManager            healthManager
)

func init() {
	healthCheckInstance = NewHealthChecker()
}

type healthManager struct {
}


func (s *healthManager) PluginHealthCheck(gClient pluginpb.PluginClient, plugin config.Plugin) (bool, error) {
	return healthCheckInstance.PluginHealthCheck(gClient, plugin)
}

func (s *healthManager) VatzHealthCheck(schedule []string) error {
	return healthCheckInstance.VatzHealthCheck(schedule)
}
