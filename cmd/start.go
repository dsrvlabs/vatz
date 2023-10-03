package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	managerPb "github.com/dsrvlabs/vatz-proto/manager/v1"
	pluginPb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/api"
	config "github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	pl "github.com/dsrvlabs/vatz/manager/plugin"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/dsrvlabs/vatz/monitoring/prometheus"
	"github.com/dsrvlabs/vatz/rpc"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func createStartCommand() *cobra.Command {
	log.Debug().Str("module", "cmd > start").Msg("start command")
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start VATZ",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().Str("module", "cmd start").Msgf("Set logfile %s", logfile)
			return utils.SetLog(logfile, defaultFlagLog)
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			log.Debug().Str("module", "main").Msgf("load config %s", configFile)
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
	log.Info().Str("module", "main").Msg("Initialize Server")
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.GetConfig()
	dispatchers = dp.GetDispatchers(cfg.Vatz.NotificationInfo)

	s := grpc.NewServer()
	serv := api.GrpcService{}
	managerPb.RegisterManagerServer(s, &serv)
	reflection.Register(s)

	vatzConfig := cfg.Vatz
	addr := fmt.Sprintf(":%d", vatzConfig.Port)
	err := healthChecker.VATZHealthCheck(vatzConfig.HealthCheckerSchedule, dispatchers)
	if err != nil {
		log.Error().Str("module", "cmd > start").Msgf("VATZHealthCheck Error: %s", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error().Str("module", "cmd > start").Msgf("VATZ Listener Error: %s", err)
	}

	log.Info().Str("module", "main").Msgf("Start VATZ Server on Listening Port: %s", addr)
	grpcClients := utils.GetClients(cfg.PluginInfos.Plugins)
	startExecutor(grpcClients, cfg.PluginInfos, ch)

	rpcServ := rpc.NewRPCService()
	go func() {
		rpcServ.Start(cfg.Vatz.RPCInfo.Address, cfg.Vatz.RPCInfo.GRPCPort, cfg.Vatz.RPCInfo.HTTPPort)
	}()
	monitoringInfo := cfg.Vatz.MonitoringInfo
	if monitoringInfo.Prometheus.Enabled {
		if defaultPromPort == promPort {
			prometheus.InitPrometheusServer(
				monitoringInfo.Prometheus.Address,
				strconv.Itoa(monitoringInfo.Prometheus.Port),
				vatzConfig.ProtocolIdentifier,
				grpcClients)
		} else {
			prometheus.InitPrometheusServer(monitoringInfo.Prometheus.Address, promPort, vatzConfig.ProtocolIdentifier, grpcClients)
		}
	}

	log.Info().Str("module", "main").Msg("VATZ Manager Started")
	initHealthServer(s)
	if err := s.Serve(listener); err != nil {
		log.Panic().Str("module", "main").Msgf("Serve Error: %s", err)
	}

	return nil
}

func startExecutor(grpcClients []pluginpb.PluginClient, pluginInfo config.PluginInfo, quit <-chan os.Signal) {
	// TODO:: value in map would be overridden by different plugins flag value if function name is the same
	isOkayToSend := false
	if len(grpcClients) == 0 {
		log.Error().Str("module", "cmd:Start").Msg("No Plugins are set, Check your Configs.")
		os.Exit(1)
	}
	for idx, singleClient := range grpcClients {
		go multiPluginExecutor(pluginInfo.Plugins[idx], singleClient, isOkayToSend, quit)
	}
}

func multiPluginExecutor(plugin config.Plugin, singleClient pluginPb.PluginClient, okToSend bool, quit <-chan os.Signal) {
	verifyTicker := time.NewTicker(time.Duration(plugin.VerifyInterval) * time.Second)
	executeTicker := time.NewTicker(time.Duration(plugin.ExecuteInterval) * time.Second)

	ctx := context.Background()

	pluginDir, err := config.GetConfig().Vatz.AbsoluteHomePath()
	if err != nil {
		return
	}

	mgr := pl.NewManager(pluginDir)
	for {
		pluginState, pluginStateErr := mgr.Get(plugin.Name)
		select {
		case <-verifyTicker.C:
			if pluginState.IsEnabled {
				live, _ := healthChecker.PluginHealthCheck(ctx, singleClient, plugin, dispatchers)
				if live == tp.AliveStatusUp {
					okToSend = true
				} else {
					okToSend = false
				}
			}
		case <-executeTicker.C:
			if pluginState.IsEnabled {
				if okToSend == true {
					if pluginStateErr != nil {
						log.Error().Str("module", "cmd > start").Msgf("Executor Error: %s", pluginStateErr)
					}
					err := executor.Execute(ctx, singleClient, plugin, dispatchers)
					if err != nil {
						log.Error().Str("module", "cmd > start").Msgf("Executor Error: %s", err)
					}
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
