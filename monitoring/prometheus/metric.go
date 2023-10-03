package prometheus

import (
	"context"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *prometheusManager) getPluginUp(plugins []config.Plugin, hostName string, grpcClients []pluginpb.PluginClient) (pluginUp map[int]*prometheusValue) {
	pluginUp = make(map[int]*prometheusValue)
	for idx, plugin := range plugins {
		pluginUp[plugin.Port] = &prometheusValue{
			Up:         1,
			PluginName: plugin.Name,
			HostName:   hostName,
		}
		verify, err := grpcClients[idx].Verify(context.Background(), new(emptypb.Empty))
		if err != nil || verify == nil {
			pluginUp[plugin.Port].Up = 0
		}
	}

	return
}
