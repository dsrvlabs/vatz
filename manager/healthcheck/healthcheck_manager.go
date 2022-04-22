package healthcheck

import (
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
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

func (s *healthManager) HealthCheck(gClient pluginpb.PluginClient, pluginInfo interface{}) (string, error) {
	Aliveness, nil := healthCheckInstance.HealthCheck(gClient, pluginInfo)
	return Aliveness, nil
}
