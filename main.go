package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"time"

	notification "github.com/dsrvlabs/vatz/manager/notification"

	config "github.com/dsrvlabs/vatz/manager/config"
	executor "github.com/dsrvlabs/vatz/manager/executor"
	health "github.com/dsrvlabs/vatz/manager/healthcheck"

	managerpb "github.com/dsrvlabs/vatz-proto/manager/v1"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/api"
	"github.com/dsrvlabs/vatz/manager/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Grpc struct {
	client pluginpb.PluginClient
}

const (
	serviceName = "Vatz Manager"
)

var (
	defaultConf     = make(map[interface{}]interface{})
	grpcClient      = Grpc{}
	healthManager   = health.HManager
	dispatchManager = notification.DManager
	configManager   = config.CManager
	executeManager  = executor.EManager
)

func preLoad() error {
	// Get a Default Info from default Yaml
	defaultConf = config.CManager.GetYMLData("default.yaml", true)
	retrievedConf := configManager.GetConfigFromURL()
	pluginInfo := configManager.Parse(model.Plugin, defaultConf)
	// Get a Default Info from default Yaml
	if !reflect.DeepEqual(retrievedConf, make(map[interface{}]interface{})) {
		for k, v := range retrievedConf {
			defaultConf[k] = v
		}
	}
	grpcClient = Grpc{configManager.GetGRPCClient(pluginInfo)}
	return nil
}

func runningProcess(pluginInfo interface{}, quit <-chan os.Signal) {
	verifyInterval := pluginInfo.(map[interface{}]interface{})["default_verify_interval"].(int)
	executeInterval := pluginInfo.(map[interface{}]interface{})["default_execute_interval"].(int)

	if verifyInterval > executeInterval || verifyInterval == executeInterval {
		verifyInterval = executeInterval - 1
	}

	verifyTicker := time.NewTicker(time.Duration(verifyInterval) * time.Second)
	executeTicker := time.NewTicker(time.Duration(executeInterval) * time.Second)

	autoUpdateNotification := make(map[interface{}]interface{})
	isOkayToSend := false

	go func() {
		for {
			select {
			case <-verifyTicker.C:
				live, _ := healthManager.HealthCheck(grpcClient.client, pluginInfo)
				if live == "UP" {
					isOkayToSend = true
				} else {
					isOkayToSend = false
				}

			//TODO: Dynamic handler for execute APIs with different time ticker.
			case <-executeTicker.C:
				if isOkayToSend == true {
					executedAfter := executeManager.Execute(grpcClient.client, pluginInfo, autoUpdateNotification)
					autoUpdateNotification = executedAfter
				}

			case <-quit:
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

	protocolInfo := configManager.Parse(model.Protocol, defaultConf)
	pluginInfo := configManager.Parse(model.Plugin, defaultConf)

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
