package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	managerpb "github.com/dsrvlabs/vatz-proto/manager/v1"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

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
				f, err := os.Create(logfile)
				if err != nil {
					return err
				}

				log.Logger = log.Output(f)
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
	cmd.PersistentFlags().StringVar(&promPort, "prometheus port", defaultPromPort, "prometheus port number.")

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

	initMetricsServer(promPort, cfg.Vatz.ProtocolIdentifier)

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

type PrometheusManager struct {
	Protocol string
	// Contains many more fields not listed in this example.
}

type PrometheusManagerCollector struct {
	PrometheusManager *PrometheusManager
}

type PrometheusValue struct {
	Up   int
	Name string
}

func initMetricsServer(port, protocol string) error {
	log.Info().Str("module", "main").Msgf("Prometheus port: %s", port)

	reg := prometheus.NewPedanticRegistry()

	NewPrometheusManager(protocol, reg)

	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Error().Str("module", "main").Msgf("Prometheus Error: %s", err)
	}

	return nil
}

func NewPrometheusManager(protocol string, reg prometheus.Registerer) *PrometheusManager {
	c := &PrometheusManager{
		Protocol: protocol,
	}
	cc := PrometheusManagerCollector{PrometheusManager: c}
	prometheus.WrapRegistererWith(prometheus.Labels{"protocol": protocol}, reg).MustRegister(cc)
	return c
}

func (cc PrometheusManagerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

func (cc PrometheusManagerCollector) Collect(ch chan<- prometheus.Metric) {
	var (
		pluginUpDesc = prometheus.NewDesc(
			"plugin_up",
			"Plugin liveness checks.",
			[]string{"plugin", "port"}, nil,
		)
	)

	upByPlugin := cc.PrometheusManager.GetPluginUp(config.GetConfig().PluginInfos.Plugins)

	for port, value := range upByPlugin {
		ch <- prometheus.MustNewConstMetric(
			pluginUpDesc,
			prometheus.GaugeValue,
			float64(value.Up),
			value.Name,
			strconv.Itoa(port),
		)
	}
}

func (c *PrometheusManager) GetPluginUp(plugins []config.Plugin) (
	pluginUp map[int]*PrometheusValue,
) {
	gClients := getClients(plugins)
	pluginUp = make(map[int]*PrometheusValue)

	for idx, plugin := range plugins {
		pluginUp[plugin.Port] = &PrometheusValue{
			Up:   1,
			Name: plugin.Name,
		}
		verify, err := gClients[idx].Verify(context.Background(), new(emptypb.Empty))
		if err != nil || verify == nil {
			pluginUp[plugin.Port].Up = 0
		}
	}

	return
}
