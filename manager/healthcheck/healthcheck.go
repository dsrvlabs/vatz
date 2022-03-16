package healthcheck

import (
	"fmt"
)

type healthCheck struct {
}

func (h healthCheck) HealthCheck() (string, error) {
	fmt.Println("this is sending Notification")
	//TODO: Set Client to Plugin
	return "UP", nil
}

type HealthCheck interface {
	HealthCheck() (string, error)
}

func NewHealthChecker() HealthCheck {
	return &healthCheck{}
}
