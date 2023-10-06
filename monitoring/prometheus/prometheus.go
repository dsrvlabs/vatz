package prometheus

import (
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"sync"
)

type prometheusManager struct {
	Protocol string
	// Contains many more fields not listed in this example.
}

type prometheusManagerCollector struct {
	prometheusManager *prometheusManager
}

type prometheusValue struct {
	Up         int
	PluginName string
	HostName   string
}

func newPrometheusManager(protocol string, reg prometheus.Registerer) *prometheusManager {
	c := &prometheusManager{
		Protocol: protocol,
	}
	cc := prometheusManagerCollector{prometheusManager: c}
	prometheus.WrapRegistererWith(prometheus.Labels{"protocol": protocol}, reg).MustRegister(cc)
	return c
}

func (cc prometheusManagerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

func (cc prometheusManagerCollector) Collect(ch chan<- prometheus.Metric) {
	var (
		pluginUpDesc = prometheus.NewDesc(
			"plugin_up",
			"Plugin liveness checks.",
			[]string{"plugin", "port", "host_name"}, nil,
		)
	)
	gClientInfos := utils.GetClients(config.GetConfig().PluginInfos.Plugins)
	upByPlugin := cc.prometheusManager.getPluginUp(config.GetConfig().Vatz.NotificationInfo.HostName, gClientInfos)

	for port, value := range upByPlugin {
		ch <- prometheus.MustNewConstMetric(
			pluginUpDesc,
			prometheus.GaugeValue,
			float64(value.Up),
			value.PluginName,
			strconv.Itoa(port),
			value.HostName,
		)
	}
}

func InitPrometheusServer(addr, port, protocol string) error {
	log.Info().Str("module", "main").Msgf("start metric server: %s:%s", addr, port)

	reg := prometheus.NewPedanticRegistry()

	var prometheusOnce sync.Once
	prometheusOnce.Do(func() {
		newPrometheusManager(protocol, reg)
	})

	reg.MustRegister(

		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	err := http.ListenAndServe(addr+":"+port, nil) //nolint:gosec

	if err != nil {
		log.Error().Str("module", "main").Msgf("Prometheus Error: %s", err)
	}

	return nil
}
