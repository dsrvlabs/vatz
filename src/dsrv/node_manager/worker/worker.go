package worker

import (
	"context"
	managerpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/v1"
	"log"
	"time"
)

const (
	defaultCheckInterval     = time.Second * 10
	defaultChannelBufferSize = 10
)

type ManagerWorker struct {
	CheckInterval time.Duration
	workerChannel chan bool
}

func (p *ManagerWorker) Init() error {
	return nil
}

func (p *ManagerWorker) Execute() error {

	if p.CheckInterval <= 0 {
		p.CheckInterval = defaultCheckInterval
	}

	p.workerChannel = make(chan bool, defaultChannelBufferSize)

	go p.worker(p.workerChannel)
	go p.invoker(p.workerChannel)

	return nil
}

func (p *ManagerWorker) Verify(ctx context.Context, in *managerpb.VerifyRequest) (*managerpb.VerifyInfo, error) {
	return nil, nil
}

func (p *ManagerWorker) UpdateConfig(ctx context.Context, in *managerpb.UpdateRequest) (*managerpb.UpdateResponse, error) {
	return nil, nil
}

func (p *ManagerWorker) End() error {
	return nil
}

func (p *ManagerWorker) worker(req <-chan bool) {
	for {
		_ = <-req
		log.Println("Do something")
		//TODO: Execute to call another plugin
	}
}

func (p *ManagerWorker) invoker(c chan<- bool) {
	for {
		c <- true
		time.Sleep(p.CheckInterval)
	}
}
