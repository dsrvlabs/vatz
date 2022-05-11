package config

import (
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	model "github.com/dsrvlabs/vatz/manager/model"
)

var (
	configInstance Config
	CManager       configManager
)

func init() {
	configInstance = NewConfig()
}

type configManager struct {
}

func (c *configManager) Parse(parseKey model.Type, data map[interface{}]interface{}) interface{} {
	return configInstance.parse(parseKey, data)
}

func (c *configManager) GetYMLData(str string, isDefault bool) map[interface{}]interface{} {
	return configInstance.getYMLData(str, isDefault)
}

func (c *configManager) GetConfigFromURL() map[interface{}]interface{} {
	return configInstance.getConfigFromURL()
}

func (c *configManager) GetGRPCClients(pluginInfo interface{}) []pluginpb.PluginClient {
	return configInstance.getClients(pluginInfo)
}

func (c *configManager) GetPingIntervals(pluginInfo interface{}, IntervalKey string) []int {
	return configInstance.getPingIntervals(pluginInfo, IntervalKey)
}
