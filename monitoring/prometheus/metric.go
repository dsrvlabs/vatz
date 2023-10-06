package prometheus

import (
	"context"
	"github.com/dsrvlabs/vatz/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *prometheusManager) getPluginUp(hostName string, grpcClientAndPluginInfo []utils.GClientWithPlugin) (pluginUp map[int]*prometheusValue) {
	pluginUp = make(map[int]*prometheusValue)

	for _, info := range grpcClientAndPluginInfo {
		pluginInfo := info.PluginInfo
		pluginUp[pluginInfo.Port] = &prometheusValue{
			Up:         1,
			PluginName: pluginInfo.Name,
			HostName:   hostName,
		}
		grpcClient := info.GRPCClient
		verify, err := grpcClient.Verify(context.Background(), new(emptypb.Empty))
		if err != nil || verify == nil {
			pluginUp[pluginInfo.Port].Up = 0
		}
	}

	return
}
