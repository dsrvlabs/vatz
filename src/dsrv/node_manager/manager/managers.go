package manager

import (
	"context"
	managerpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/v1"
	worker_presenter "pilot-manager/src/dsrv/node_manager/worker"
	"time"
)

type Manager interface {
	Init() error
	Execute() error
	Verify(ctx context.Context, in *managerpb.VerifyRequest) (*managerpb.VerifyInfo, error)
	End() error
	UpdateConfig(ctx context.Context, in *managerpb.UpdateRequest) (*managerpb.UpdateResponse, error)
}

func RunManager() Manager {
	return &worker_presenter.ManagerWorker{CheckInterval: time.Second * 5}
}
