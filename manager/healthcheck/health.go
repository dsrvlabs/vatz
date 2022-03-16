package healthcheck

import (
	"fmt"
)

type health struct {
}

func (h health) HealthCheck() (string, error) {
	fmt.Println("this is sending Notification")
	//TODO: Set Client to Plugin
	return "UP", nil
}

type Health interface {
	HealthCheck() (string, error)
}

func NewHealthChecker() Health {
	return &health{}
}
