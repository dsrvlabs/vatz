package main

import (
	"context"
	"fmt"
	managerpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	node_manager "pilot-manager/src/dsrv/node_manager/api"
	"sync"
)

const (
	serviceName = "Node Manager"
)

var (
	grpcPort = 9090
)

func preLoad() error {
	//TODOs: Check the Configs and return dict as global variable.
	return nil
}

func initiateServer(ch <-chan os.Signal) error {
	log.Println("Initialize Servers:", serviceName)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := grpc.NewServer()
	serv := node_manager.GrpcService{}

	managerpb.RegisterNodeManagerServer(s, &serv)
	reflection.Register(s)

	addr := fmt.Sprintf(":%d", grpcPort)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
	}

	log.Println("Listening Port", addr)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_ = <-ch
		cancel()
		s.GracefulStop()
		wg.Done()
	}()

	log.Println("Node Manager (Pilot) started")

	if err := s.Serve(listener); err != nil {
		log.Panic(err)
	}

	return nil
}

func main() {
	ch := make(chan os.Signal, 1)
	initiateServer(ch)
}
