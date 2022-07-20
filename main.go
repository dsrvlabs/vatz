package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	ex "github.com/dsrvlabs/vatz/manager/executor"
	notification "github.com/dsrvlabs/vatz/manager/notification"

	config "github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/manager/healthcheck"
	health "github.com/dsrvlabs/vatz/manager/healthcheck"

	managerpb "github.com/dsrvlabs/vatz-proto/manager/v1"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	serviceName = "Vatz Manager"
)

var (
	healthManager   = health.HManager
	dispatchManager = notification.GetDispatcher()

	executor ex.Executor

	defaultVerifyInterval  = 15
	defaultExecuteInterval = 30
)

func init() {
	executor = ex.NewExecutor()
}

func main() {
	var configFile string

	// TODO: How to test flag?
	flag.StringVar(&configFile, "config", "default.yaml", "-config=<FILENAME>")
	flag.Parse()

	config.InitConfig(configFile)

	ch := make(chan os.Signal, 1)
	initiateServer(ch)
}

func initiateServer(ch <-chan os.Signal) error {
	log.Println("Initialize Servers:", serviceName)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := grpc.NewServer()
	serv := api.GrpcService{}
	managerpb.RegisterManagerServer(s, &serv)
	reflection.Register(s)

	cfg := config.GetConfig()
	vatzConfig := cfg.Vatz
	addr := fmt.Sprintf(":%d", vatzConfig.Port)
	err := healthcheck.VatzHealthCheck(vatzConfig.HealthCheckerSchedule)
	if err != nil {
		log.Println(err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
	}

	log.Println("Listening Port", addr)

	startExecutor(cfg.PluginInfos, ch)

	log.Println("Node Manager Started")

	if err := s.Serve(listener); err != nil {
		log.Panic(err)
	}

	return nil
}

func startExecutor(pluginInfo config.PluginInfo, quit <-chan os.Signal) {
	//TODO:: value in map would be overridden by different plugins flag value if function name is the same
	isOkayToSend := false

	grpcClients := getClients(pluginInfo.Plugins)

	//TODO: Need updated with better way for Dynamic handlers
	for idx, singleClient := range grpcClients {
		go multiPluginExecutor(pluginInfo.Plugins[idx], singleClient, isOkayToSend, quit)
	}
}

func getClients(plugins []config.Plugin) []pluginpb.PluginClient {
	var grpcClients []pluginpb.PluginClient

	if len(plugins) > 0 {
		for _, plugin := range plugins {
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", plugin.Address, plugin.Port), grpc.WithInsecure())
			if err != nil {
				// TODO: Panic??? Error message will be enough here.
				log.Fatal(err)
			}
			grpcClients = append(grpcClients, pluginpb.NewPluginClient(conn))
		}
	} else {
		// TODO: Is this really neccessary???
		defaultConnectedTarget := "localhost:9091"
		conn, err := grpc.Dial(defaultConnectedTarget, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		//TODO: Please, Create a better client functions with static
		grpcClients = append(grpcClients, pluginpb.NewPluginClient(conn))
	}

	return grpcClients
}

func multiPluginExecutor(plugin config.Plugin,
	singleClient pluginpb.PluginClient,
	isOkayToSend bool,
	quit <-chan os.Signal) {

	verifyTicker := time.NewTicker(time.Duration(plugin.VerifyInterval) * time.Second)
	executeTicker := time.NewTicker(time.Duration(plugin.ExecuteInterval) * time.Second)

	ctx := context.Background()
	for {
		select {
		case <-verifyTicker.C:
			live, _ := healthManager.HealthCheck(singleClient, plugin)
			if live == "UP" {
				isOkayToSend = true
			} else {
				isOkayToSend = false
			}
		case <-executeTicker.C:
			if isOkayToSend == true {
				err := executor.Execute(ctx, singleClient, plugin)
				if err != nil {
					// TODO: Handle error.
				}
			}
		case <-quit:
			executeTicker.Stop()
			return
		}
	}
}
