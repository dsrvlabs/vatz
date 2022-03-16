package healthcheck

var (
	healthCheckInstance HealthCheck
	HManager            healthManager
)

func init() {
	healthCheckInstance = NewHealthChecker()
}

type healthManager struct {
}

func (s *healthManager) HealthCheck() (string, error) {
	Aliveness, nil := healthCheckInstance.HealthCheck()
	return Aliveness, nil
}
