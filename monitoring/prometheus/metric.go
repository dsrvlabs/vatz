package prometheus

import (
	"context"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *prometheusManager) getPluginUp(plugins []config.Plugin, hostName string) (pluginUp map[int]*prometheusValue) {
	pluginUp = make(map[int]*prometheusValue)
	gClients := utils.GetClients(plugins)
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
