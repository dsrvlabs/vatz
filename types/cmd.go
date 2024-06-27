package types

type Initializer string

const (
	TEST Initializer = "TEST"
	LIVE Initializer = "LIVE"
)

type PluginState struct {
	Status       string `json:"status"`
	PluginStatus []struct {
		Status     string `json:"status"`
		PluginName string `json:"pluginName"`
	} `json:"pluginStatus"`
}
