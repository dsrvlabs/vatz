package healthcheck

import (
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
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
