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

const (
	serviceName = "Vatz Manager"
)

var (
	grpcClients            []pluginpb.PluginClient
	defaultConf            = make(map[interface{}]interface{})
	healthManager          = health.HManager
	dispatchManager        = notification.DManager
	configManager          = config.CManager
	executeManager         = executor.EManager
	defaultVerifyInterval  = 15
	defaultExecuteInterval = 30
)

func preLoad() error {
	// Get a Default Info from default Yaml
	defaultConf = config.CManager.GetYMLData("default.yaml", true)
	retrievedConf := configManager.GetConfigFromURL()
	pluginInfo := configManager.Parse(model.Plugin, defaultConf)

	defaultVerifyInterval = pluginInfo.(map[interface{}]interface{})["default_verify_interval"].(int)
	defaultExecuteInterval = pluginInfo.(map[interface{}]interface{})["default_execute_interval"].(int)

	if defaultVerifyInterval > defaultExecuteInterval || defaultVerifyInterval == defaultExecuteInterval {
		defaultVerifyInterval = defaultExecuteInterval - 1
	}

	// Get a Default Info from default Yaml
	if !reflect.DeepEqual(retrievedConf, make(map[interface{}]interface{})) {
		for k, v := range retrievedConf {
			defaultConf[k] = v
		}
	}

	grpcClients = configManager.GetGRPCClients(pluginInfo)
	fmt.Println("grpcClients", grpcClients)
	return nil
}

func runningProcess(pluginInfo interface{}, quit <-chan os.Signal) {

	verifyIntervals := configManager.GetPingIntervals(pluginInfo, "verify_interval")
	executeIntervals := configManager.GetPingIntervals(pluginInfo, "execute_interval")

	//TODO:: value in map would be overridden by different plugins flag value if function name is the same
	autoUpdateNotification := make(map[interface{}]interface{})
	isOkayToSend := false

	//TODO: Need updated with better way for Dynamic handlers
	for idx, singleClient := range grpcClients {
		go multiPluginExecutor(pluginInfo, singleClient, verifyIntervals[idx], executeIntervals[idx], isOkayToSend, autoUpdateNotification, quit)
	}
}

func multiPluginExecutor(pluginInfo interface{},
	singleClient pluginpb.PluginClient,
	verifyInterval int,
	executeInterval int,
	isOkayToSend bool,
	autoUpdateNotification map[interface{}]interface{},
	quit <-chan os.Signal) {

	verifyTicker := time.NewTicker(time.Duration(verifyInterval) * time.Second)
	executeTicker := time.NewTicker(time.Duration(executeInterval) * time.Second)

	for {
		select {
		case <-verifyTicker.C:
			live, _ := healthManager.HealthCheck(singleClient, pluginInfo)
			if live == "UP" {
				isOkayToSend = true
			} else {
				isOkayToSend = false
			}

		case <-executeTicker.C:
			if isOkayToSend == true {
				executedAfter := executeManager.Execute(singleClient, pluginInfo, autoUpdateNotification)
				autoUpdateNotification = executedAfter
			}

		case <-quit:
			executeTicker.Stop()
			return
		}
	}
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
