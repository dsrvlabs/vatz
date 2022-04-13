package config

import (
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
)

var (
	configInstance Config
	CManager       config_manager
)

func init() {
	configInstance = NewConfig()
}

type config_manager struct {
}

func (c *config_manager) Parse(retrievalInfo string, data map[interface{}]interface{}) interface{} {
	return configInstance.parse(retrievalInfo, data)
}

func (c *config_manager) GetYMLData(str string, isDefault bool) map[interface{}]interface{} {
	return configInstance.getYMLData(str, isDefault)
}

func (c *config_manager) GetConfigFromURL() map[interface{}]interface{} {
	return configInstance.getConfigFromURL()
}

func (c *config_manager) GetGRPCClient() pluginpb.PluginClient {
	return configInstance.getClient()
}
