package main

import (
	"context"
	"fmt"
	managerpb "github.com/dsrvlabs/vatz-proto/manager/v1"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/api"
	"github.com/dsrvlabs/vatz/manager/config"
	ex "github.com/dsrvlabs/vatz/manager/executor"
	health "github.com/dsrvlabs/vatz/manager/healthcheck"
	"github.com/dsrvlabs/vatz/manager/notification"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"time"
)

const (
	serviceName = "VATZ Manager"
)

var (
	healthChecker = health.NewHealthChecker()
	dispatcher    = notification.GetDispatcher()
	executor      = ex.NewExecutor()

	defaultVerifyInterval  = 15
	defaultExecuteInterval = 30
)

func init() {
	executor = ex.NewExecutor()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}

func main() {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(createInitCommand())
	rootCmd.AddCommand(createStartCommand())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func initiateServer(ch <-chan os.Signal) error {
	log.Info().Str("module", "main").Msgf("Initialize Servers: %s", serviceName)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := grpc.NewServer()
	serv := api.GrpcService{}
	managerpb.RegisterManagerServer(s, &serv)
	reflection.Register(s)

	cfg := config.GetConfig()
	vatzConfig := cfg.Vatz
	addr := fmt.Sprintf(":%d", vatzConfig.Port)
	err := healthChecker.VATZHealthCheck(vatzConfig.HealthCheckerSchedule, dispatcher)
	if err != nil {
		log.Error().Str("module", "main").Msgf("VATZHealthCheck Error: %s", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error().Str("module", "main").Msgf("listener Error: %s", err)
	}
	log.Info().Str("module", "main").Msgf("Listening Port: %s", addr)
	startExecutor(cfg.PluginInfos, ch)
	log.Info().Str("module", "main").Msg("Node Manager Started")

	InitHealthServer(s)
	if err := s.Serve(listener); err != nil {
		log.Panic().Str("module", "main").Msgf("Serve Error: %s", err)
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
				log.Fatal().Str("module", "main").Msgf("gRPC Dial Error(%s): %s", plugin.Name, err)
			}
			grpcClients = append(grpcClients, pluginpb.NewPluginClient(conn))
		}
	} else {
		// TODO: Is this really neccessary???
		defaultConnectedTarget := "localhost:9091"
		conn, err := grpc.Dial(defaultConnectedTarget, grpc.WithInsecure())
		if err != nil {
			log.Fatal().Str("module", "main").Msgf("gRPC Dial Error: %s", err)
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
			live, _ := healthChecker.PluginHealthCheck(ctx, singleClient, plugin, dispatcher)
			if live == tp.AliveStatusUp {
				isOkayToSend = true
			} else {
				isOkayToSend = false
			}
		case <-executeTicker.C:
			if isOkayToSend == true {
				err := executor.Execute(ctx, singleClient, plugin)
				if err != nil {
					log.Error().Str("module", "main").Msgf("Executor Error: %s", err)
				}
			}
		case <-quit:
			executeTicker.Stop()
			return
		}
	}
}

func InitHealthServer(s *grpc.Server) {
	gRPCHealthServer := grpchealth.NewServer()
	gRPCHealthServer.SetServingStatus("vatz-health-status", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, gRPCHealthServer)
}
