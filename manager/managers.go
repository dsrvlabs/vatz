package manager

import (
	"time"
)

type Manager interface {
	Start() error
	Execute() error
	Stop() error
}

func RunManager() Manager {
	return &managerWorker{CheckInterval: time.Second * 5}
}
