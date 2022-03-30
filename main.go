package main

import (
	"context"
	"fmt"
	managerpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/manager/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"time"
	"vatz/manager/api"
	config "vatz/manager/config"
	"vatz/manager/executor"
	health "vatz/manager/healthcheck"
	notification "vatz/manager/notification"
)

const (
	serviceName = "Node Manager"
)

var (
	defaultConf     = config.CManager.GetYMLData("default.yaml", true)
	grpcClient      = config.CManager.GetGRPCClient()
	healthManager   = health.HManager
	dispatchManager = notification.DManager
	configManager   = config.CManager
	executeManager  = executor.EManager
)

//This Merge yaml file
func preLoad() error {
	retrievedConf := configManager.GetConfigFromURL()
	for k, v := range retrievedConf {
		defaultConf[k] = v
	}
	return nil
}

func runningProcess(pluginInfo interface{}, ch <-chan os.Signal) {
	verifyInterval := pluginInfo.(map[interface{}]interface{})["default_verify_interval"].(int)
	executeInterval := pluginInfo.(map[interface{}]interface{})["default_execute_interval"].(int)

	executeTicker := time.NewTicker(time.Duration(verifyInterval) * time.Second)
	verifyTicker := time.NewTicker(time.Duration(executeInterval) * time.Second)

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

}

func initiateServer(ch <-chan os.Signal) error {

	log.Println("Initialize Servers:", serviceName)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := grpc.NewServer()
	serv := api.GrpcService{}
	managerpb.RegisterManagerServer(s, &serv)
	reflection.Register(s)

	//Isn't there any better way to parse it? This isn't good to coverUp all yaml
	protocolInfo := configManager.Parse("PROTOCOL", defaultConf)
	pluginInfo := configManager.Parse("PLUGIN", defaultConf)

	addr := fmt.Sprintf(":%d", protocolInfo.(map[interface{}]interface{})["port"])

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
	}

	log.Println("Listening Port", addr)

	runningProcess(pluginInfo, ch)
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
