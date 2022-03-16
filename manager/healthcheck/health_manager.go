package healthcheck

var (
	healthInstance Health
	HManager       healthManager
)

func init() {
	healthInstance = NewHealthChecker()
}

type healthManager struct {
}

func (s *healthManager) HealthCheck() (string, error) {
	Aliveness, nil := healthInstance.HealthCheck()
	return Aliveness, nil
}
