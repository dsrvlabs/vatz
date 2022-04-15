package config

import (
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
	model "vatz/manager/model"
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

func (c *configManager) GetGRPCClient(pluginInfo interface{}) pluginpb.PluginClient {
	return configInstance.getClient(pluginInfo)
}
