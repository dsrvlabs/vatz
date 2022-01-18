package manager

import (
	"log"
	"time"
)

const (
	defaultCheckInterval     = time.Second * 10
	defaultChannelBufferSize = 10
)

type managerWorker struct {
	CheckInterval time.Duration
	workerChannel chan bool
}

func (p *managerWorker) Start() error {
	return nil
}

func (p *managerWorker) Execute() error {

	if p.CheckInterval <= 0 {
		p.CheckInterval = defaultCheckInterval
	}

	p.workerChannel = make(chan bool, defaultChannelBufferSize)

	go p.worker(p.workerChannel)
	go p.invoker(p.workerChannel)

	return nil
}

func (p *managerWorker) Stop() error {
	return nil
}

func (p *managerWorker) worker(req <-chan bool) {
	for {
		_ = <-req
		log.Println("Do something")
		//TODO: Excute to call another plugin
	}
}

func (p *managerWorker) invoker(c chan<- bool) {
	for {
		c <- true
		time.Sleep(p.CheckInterval)
	}
}
