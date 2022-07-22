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

func (s *healthManager) HealthCheck(gClient pluginpb.PluginClient, plugin config.Plugin) (string, error) {
	return healthCheckInstance.HealthCheck(gClient, plugin)
}

func (s *healthManager) VatzHealthCheck(schedule []string) error {
	return healthCheckInstance.VatzHealthCheck(schedule)
}
