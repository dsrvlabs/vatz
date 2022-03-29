package main

import (
	"context"
	"fmt"
	managerpb "github.com/hqueue/vatz-secret/dist/proto/vatz/manager/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
	"vatz/manager/api"
	"vatz/manager/executor"
	manager2 "vatz/manager/healthcheck"
	"vatz/manager/notification"
)

const (
	serviceName = "Node Manager"
)

var (
	defaultConf     = getConf()
	healthManager   = manager2.HManager
	dispatchManager = notification.DManager
	executeManager  = executor.EManager
)

func getConf() map[interface{}]interface{} {
	wd, _ := os.Getwd()
	confPath := fmt.Sprintf("%s/default.yaml", wd)

	yamlFile, err := ioutil.ReadFile(confPath)

	if err != nil {
		log.Fatal(err)
	}

	data := make(map[interface{}]interface{})
	err2 := yaml.Unmarshal(yamlFile, &data)

	if err2 != nil {
		log.Fatal(err2)
	}

	return data
}

func preLoad() error {
	//TODOs: Check the Configs and return dict as global variable.
	return nil
}

func initiateServer(ch <-chan os.Signal) error {

	log.Println("Initialize Servers:", serviceName)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := grpc.NewServer()
	serv := api.GrpcService{}

	managerpb.RegisterManagerServer(s, &serv)
	reflection.Register(s)
	grpcPort := defaultConf["port"]
	addr := fmt.Sprintf(":%d", grpcPort)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
	}

	log.Println("Listening Port", addr)

	executeTicker := time.NewTicker(1 * time.Second)
	verifyTicker := time.NewTicker(2 * time.Second)
	isOkayToSend := false
	go func() {
		for {
			select {
			case <-verifyTicker.C:
				live, _ := healthManager.HealthCheck()
				if live == "UP" {
					isOkayToSend = true
				}

			case <-executeTicker.C:
				if isOkayToSend == true {
					executeManager.Execute()
				}
			case <-ch:
				executeTicker.Stop()
				return
			}
		}
	}()

	log.Println("Node Manager Started")

	if err := s.Serve(listener); err != nil {
		log.Panic(err)
	}

	return nil
}

func main() {
	preLoad()
	ch := make(chan os.Signal, 1)
	initiateServer(ch)
}
