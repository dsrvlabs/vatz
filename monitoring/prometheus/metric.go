package prometheus

import (
	"context"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"sync"
)

func InitMetricsServer(addr, port, protocol string) error {
	log.Info().Str("module", "main").Msgf("start metric server: %s:%s", addr, port)

	reg := prometheus.NewPedanticRegistry()

	var prometheusOnce sync.Once

	prometheusOnce.Do(func() {
		newPrometheusManager(protocol, reg)
	})

	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	err := http.ListenAndServe(addr+":"+port, nil) //nolint:gosec

	if err != nil {
		log.Error().Str("module", "main").Msgf("Prometheus Error: %s", err)
	}

	return nil
}

func (c *prometheusManager) getPluginUp(plugins []config.Plugin, hostName string) (
	pluginUp map[int]*prometheusValue,
) {
	gClients := utils.GetClients(plugins)
	pluginUp = make(map[int]*prometheusValue)

	for idx, plugin := range plugins {
		pluginUp[plugin.Port] = &prometheusValue{
			Up:         1,
			PluginName: plugin.Name,
			HostName:   hostName,
		}
		verify, err := gClients[idx].Verify(context.Background(), new(emptypb.Empty))
		if err != nil || verify == nil {
			pluginUp[plugin.Port].Up = 0
		}
	}

	return
}
