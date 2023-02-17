package prometheus

import (
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
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

	upByPlugin := cc.prometheusManager.getPluginUp(config.GetConfig().PluginInfos.Plugins, config.GetConfig().Vatz.NotificationInfo.HostName)

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
