package manager

import (
	worker_presenter "pilot-manager/src/dsrv/node_manager/worker"
	"time"
)

type Manager interface {
	Start() error
	Execute() error
	Stop() error
}

func RunManager() Manager {
	return &worker_presenter.managerWorker{CheckInterval: time.Second * 5}
}
