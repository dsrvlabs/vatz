package cmd

import (
	"context"
	"errors"
	"fmt"
	managerpb "github.com/dsrvlabs/vatz-proto/manager/v1"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/monitoring/prometheus"
	"github.com/dsrvlabs/vatz/rpc"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/dsrvlabs/vatz/manager/api"
	config "github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	tp "github.com/dsrvlabs/vatz/manager/types"
)

func createStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start VATZ",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if logfile == defaultFlagLog {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
			} else {
				f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					return err
				}

				log.Logger = log.Output(zerolog.ConsoleWriter{Out: f, TimeFormat: time.RFC3339})
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("module", "main").Msg("start")
			log.Info().Str("module", "main").Msgf("load config %s", configFile)
			log.Info().Str("module", "main").Msgf("logfile %s", logfile)

			_, err := config.InitConfig(configFile)
			if err != nil {
				log.Error().Str("module", "config").Msgf("loadConfig Error: %s", err)
				if errors.Is(err, os.ErrNotExist) {
					msg := "Please, initialize VATZ with command `./vatz init` to create config file `default.yaml` first or set appropriate path for config file default.yaml."
					log.Error().Str("module", "config").Msg(msg)
				}
				return nil
			}

			ch := make(chan os.Signal, 1)
			return initiateServer(ch)
		},
	}

	cmd.PersistentFlags().StringVar(&configFile, "config", defaultFlagConfig, "VATZ config file.")
	cmd.PersistentFlags().StringVar(&logfile, "log", defaultFlagLog, "log file export to.")
	cmd.PersistentFlags().StringVar(&promPort, "prometheus", defaultPromPort, "prometheus port number.")

	return cmd
}

func initiateServer(ch <-chan os.Signal) error {
	log.Info().Str("module", "main").Msgf("Initialize Servers: %s", "VATZ Manager")
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.GetConfig()
	dispatchers = dp.GetDispatchers(cfg.Vatz.NotificationInfo)

	s := grpc.NewServer()
	serv := api.GrpcService{}
	managerpb.RegisterManagerServer(s, &serv)
	reflection.Register(s)

	vatzConfig := cfg.Vatz
	addr := fmt.Sprintf(":%d", vatzConfig.Port)
	err := healthChecker.VATZHealthCheck(vatzConfig.HealthCheckerSchedule, dispatchers)
	if err != nil {
		log.Error().Str("module", "main").Msgf("VATZHealthCheck Error: %s", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error().Str("module", "main").Msgf("VATZ Listener Error: %s", err)
	}

	log.Info().Str("module", "main").Msgf("VATZ Listening Port: %s", addr)
	startExecutor(cfg.PluginInfos, ch)

	rpcServ := rpc.NewRPCService()
	go func() {
		rpcServ.Start(cfg.Vatz.RPCInfo.Address, cfg.Vatz.RPCInfo.GRPCPort, cfg.Vatz.RPCInfo.HTTPPort)
	}()

	if cfg.Vatz.MonitoringInfo.Prometheus.Enabled {
		if defaultPromPort == promPort {
			prometheus.InitMetricsServer(cfg.Vatz.MonitoringInfo.Prometheus.Address, strconv.Itoa(cfg.Vatz.MonitoringInfo.Prometheus.Port), cfg.Vatz.ProtocolIdentifier)
		} else {
			prometheus.InitMetricsServer(cfg.Vatz.MonitoringInfo.Prometheus.Address, promPort, cfg.Vatz.ProtocolIdentifier)
		}
	}

	log.Info().Str("module", "main").Msg("VATZ Manager Started")
	initHealthServer(s)
	if err := s.Serve(listener); err != nil {
		log.Panic().Str("module", "main").Msgf("Serve Error: %s", err)
	}

	return nil
}

func startExecutor(pluginInfo config.PluginInfo, quit <-chan os.Signal) {
	//TODO:: value in map would be overridden by different plugins flag value if function name is the same
	isOkayToSend := false
	grpcClients := utils.GetClients(pluginInfo.Plugins)
	//TODO: Need updated with better way for Dynamic handlers
	for idx, singleClient := range grpcClients {
		go multiPluginExecutor(pluginInfo.Plugins[idx], singleClient, isOkayToSend, quit)
	}
}

func multiPluginExecutor(plugin config.Plugin,
	singleClient pluginpb.PluginClient,
	okToSend bool,
	quit <-chan os.Signal) {

	verifyTicker := time.NewTicker(time.Duration(plugin.VerifyInterval) * time.Second)
	executeTicker := time.NewTicker(time.Duration(plugin.ExecuteInterval) * time.Second)

	ctx := context.Background()
	for {
		select {
		case <-verifyTicker.C:
			live, _ := healthChecker.PluginHealthCheck(ctx, singleClient, plugin, dispatchers)
			if live == tp.AliveStatusUp {
				okToSend = true
			} else {
				okToSend = false
			}
		case <-executeTicker.C:
			if okToSend == true {
				err := executor.Execute(ctx, singleClient, plugin, dispatchers)
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

func initHealthServer(s *grpc.Server) {
	gRPCHealthServer := grpchealth.NewServer()
	gRPCHealthServer.SetServingStatus("vatz-health-status", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, gRPCHealthServer)
}
