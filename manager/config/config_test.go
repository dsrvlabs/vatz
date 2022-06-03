package config

import (
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func createConfig(content string) *Config {
	p := parser{}
	config, err := p.parseYAML([]byte(content))
	if err != nil {
		panic(err)
	}

	vatzConfig = config
	return config
}

type testConfigExpects struct {
	MockContents                 string
	ExpectProtocolID             string
	ExpectVatzPort               int
	ExpectDiscordSecret          string
	ExpectPagerDutySecret        string
	ExpectDefaultVerifyInterval  int
	ExpectDefaultExecuteInterval int
	ExpectDefaultPluginName      string

	Plugins []testPluginExpects
}

type testPluginExpects struct {
	ExpectPluginName      string
	ExpectPluginAddr      string
	ExpectVerifyInterval  int
	ExpectExecuteInterval int
	ExpectPort            int
	ExpectMethods         []string
}

func TestDefaultConfig(t *testing.T) {
	test := testConfigExpects{
		MockContents: configDefaultContents,

		ExpectProtocolID:             "vatz",
		ExpectVatzPort:               9090,
		ExpectDiscordSecret:          "XXXXX",
		ExpectPagerDutySecret:        "YYYYY",
		ExpectDefaultVerifyInterval:  15,
		ExpectDefaultExecuteInterval: 30,
		ExpectDefaultPluginName:      "vatz-plugin",
		Plugins: []testPluginExpects{
			{
				ExpectPluginName:      "vatz-plugin-node-checker",
				ExpectPluginAddr:      "localhost",
				ExpectVerifyInterval:  7,
				ExpectExecuteInterval: 9,
				ExpectPort:            9091,
				ExpectMethods: []string{
					"isUp", "getBlockHeight", "getNumberOfPeers",
				},
			},
			{
				ExpectPluginName:      "vatz-plugin-machine-checker",
				ExpectPluginAddr:      "localhost",
				ExpectVerifyInterval:  8,
				ExpectExecuteInterval: 10,
				ExpectPort:            9092,
				ExpectMethods: []string{
					"getMemory", "getDiscSize", "getCPUInfo",
				},
			},
		},
	}

	cfg := createConfig(test.MockContents)

	// Asserts.
	assert.Equal(t, test.ExpectProtocolID, cfg.Vatz.ProtocolIdentifier)
	assert.Equal(t, test.ExpectVatzPort, cfg.Vatz.Port)
	assert.Equal(t, test.ExpectDiscordSecret, cfg.Vatz.NotificationInfo.DiscordSecret)
	assert.Equal(t, test.ExpectPagerDutySecret, cfg.Vatz.NotificationInfo.PagerDutySecret)

	assert.Equal(t, test.ExpectDefaultVerifyInterval, cfg.PluginInfos.DefaultVerifyInterval)
	assert.Equal(t, test.ExpectDefaultExecuteInterval, cfg.PluginInfos.DefaultExecuteInterval)

	assert.Equal(t, len(test.Plugins), len(cfg.PluginInfos.Plugins))

	for i, plugin := range test.Plugins {
		assert.Equal(t, plugin.ExpectPluginName, cfg.PluginInfos.Plugins[i].Name)
		assert.Equal(t, plugin.ExpectPluginAddr, cfg.PluginInfos.Plugins[i].Address)
		assert.Equal(t, plugin.ExpectVerifyInterval, cfg.PluginInfos.Plugins[i].VerifyInterval)
		assert.Equal(t, plugin.ExpectExecuteInterval, cfg.PluginInfos.Plugins[i].ExecuteInterval)
		assert.Equal(t, plugin.ExpectPort, cfg.PluginInfos.Plugins[i].Port)

		assert.Equal(t, len(plugin.ExpectMethods), len(cfg.PluginInfos.Plugins[i].ExecutableMethods))

		for _, method := range cfg.PluginInfos.Plugins[i].ExecutableMethods {
			assert.Contains(t, plugin.ExpectMethods, method.Name)
		}
	}
}

func TestOverrideDefaultConfig(t *testing.T) {
	test := testConfigExpects{
		MockContents: configNoIntervalContents,

		ExpectProtocolID:             "vatz",
		ExpectVatzPort:               9090,
		ExpectDiscordSecret:          "hello",
		ExpectPagerDutySecret:        "world",
		ExpectDefaultVerifyInterval:  15,
		ExpectDefaultExecuteInterval: 30,
		ExpectDefaultPluginName:      "vatz-plugin",
		Plugins: []testPluginExpects{
			{
				ExpectPluginName:      "vatz-plugin", // Same as default.
				ExpectPluginAddr:      "localhost",
				ExpectVerifyInterval:  15, // Same as default.
				ExpectExecuteInterval: 30, // Same as default.
				ExpectPort:            9091,
				ExpectMethods: []string{
					"isUp", "getBlockHeight", "getNumberOfPeers",
				},
			},
		},
	}

	cfg := createConfig(test.MockContents)

	// Asserts.
	assert.Equal(t, test.ExpectProtocolID, cfg.Vatz.ProtocolIdentifier)
	assert.Equal(t, test.ExpectVatzPort, cfg.Vatz.Port)
	assert.Equal(t, test.ExpectDiscordSecret, cfg.Vatz.NotificationInfo.DiscordSecret)
	assert.Equal(t, test.ExpectPagerDutySecret, cfg.Vatz.NotificationInfo.PagerDutySecret)

	assert.Equal(t, test.ExpectDefaultVerifyInterval, cfg.PluginInfos.DefaultVerifyInterval)
	assert.Equal(t, test.ExpectDefaultExecuteInterval, cfg.PluginInfos.DefaultExecuteInterval)

	assert.Equal(t, len(test.Plugins), len(cfg.PluginInfos.Plugins))

	for i, plugin := range test.Plugins {
		assert.Equal(t, plugin.ExpectPluginName, cfg.PluginInfos.Plugins[i].Name)
		assert.Equal(t, plugin.ExpectPluginAddr, cfg.PluginInfos.Plugins[i].Address)
		assert.Equal(t, plugin.ExpectVerifyInterval, cfg.PluginInfos.Plugins[i].VerifyInterval)
		assert.Equal(t, plugin.ExpectExecuteInterval, cfg.PluginInfos.Plugins[i].ExecuteInterval)
		assert.Equal(t, plugin.ExpectPort, cfg.PluginInfos.Plugins[i].Port)

		assert.Equal(t, len(plugin.ExpectMethods), len(cfg.PluginInfos.Plugins[i].ExecutableMethods))

		for _, method := range cfg.PluginInfos.Plugins[i].ExecutableMethods {
			assert.Contains(t, plugin.ExpectMethods, method.Name)
		}
	}
}

func TestURLConfig(t *testing.T) {
	tests := []struct {
		Desc              string
		MockConfigContent string
		MockStatus        int
		ExpectSuccess     bool
	}{
		{
			Desc:              "Get HTTP success",
			MockConfigContent: configDefaultContents,
			MockStatus:        http.StatusOK,
			ExpectSuccess:     true,
		},
		{
			Desc:              "HTTP Response failed",
			MockConfigContent: configDefaultContents,
			MockStatus:        http.StatusInternalServerError,
			ExpectSuccess:     false,
		},
	}

	for _, test := range tests {
		// Prepare http mock.
		httpmock.Activate()

		dummyURL := "https://github.com/config.yaml"
		httpmock.RegisterResponder(
			http.MethodGet,
			dummyURL,
			httpmock.NewStringResponder(test.MockStatus, test.MockConfigContent),
		)

		// Load config.
		p := parser{}
		configData, err := p.loadConfigFile(dummyURL)

		if test.ExpectSuccess {
			_ = createConfig(string(configData))
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		httpmock.DeactivateAndReset()
	}
}

// This test should be the last one, because "GetConfig" function can execute once.
func TestGetConfig(t *testing.T) {
	test := testConfigExpects{
		MockContents: configDefaultContents,

		ExpectProtocolID:             "vatz",
		ExpectVatzPort:               9090,
		ExpectDiscordSecret:          "XXXXX",
		ExpectPagerDutySecret:        "YYYYY",
		ExpectDefaultVerifyInterval:  15,
		ExpectDefaultExecuteInterval: 30,
		ExpectDefaultPluginName:      "vatz-plugin",
		Plugins: []testPluginExpects{
			{
				ExpectPluginName:      "vatz-plugin-node-checker",
				ExpectPluginAddr:      "localhost",
				ExpectVerifyInterval:  7,
				ExpectExecuteInterval: 9,
				ExpectPort:            9091,
				ExpectMethods: []string{
					"isUp", "getBlockHeight", "getNumberOfPeers",
				},
			},
			{
				ExpectPluginName:      "vatz-plugin-machine-checker",
				ExpectPluginAddr:      "localhost",
				ExpectVerifyInterval:  8,
				ExpectExecuteInterval: 10,
				ExpectPort:            9092,
				ExpectMethods: []string{
					"getMemory", "getDiscSize", "getCPUInfo",
				},
			},
		},
	}

	// Call GetConfig
	f, err := os.Create(defaultConfigFile)
	if err != nil {
		assert.Fail(t, "Config file doesn't created")
	}
	defer os.Remove(defaultConfigFile)

	_, err = f.WriteString(test.MockContents)
	if err != nil {
		assert.Fail(t, "Config file doesn't created")
	}

	cfg := GetConfig()

	// Asserts.
	assert.Equal(t, test.ExpectProtocolID, cfg.Vatz.ProtocolIdentifier)
	assert.Equal(t, test.ExpectVatzPort, cfg.Vatz.Port)
	assert.Equal(t, test.ExpectDiscordSecret, cfg.Vatz.NotificationInfo.DiscordSecret)
	assert.Equal(t, test.ExpectPagerDutySecret, cfg.Vatz.NotificationInfo.PagerDutySecret)

	assert.Equal(t, test.ExpectDefaultVerifyInterval, cfg.PluginInfos.DefaultVerifyInterval)
	assert.Equal(t, test.ExpectDefaultExecuteInterval, cfg.PluginInfos.DefaultExecuteInterval)

	assert.Equal(t, len(test.Plugins), len(cfg.PluginInfos.Plugins))

	for i, plugin := range test.Plugins {
		assert.Equal(t, plugin.ExpectPluginName, cfg.PluginInfos.Plugins[i].Name)
		assert.Equal(t, plugin.ExpectPluginAddr, cfg.PluginInfos.Plugins[i].Address)
		assert.Equal(t, plugin.ExpectVerifyInterval, cfg.PluginInfos.Plugins[i].VerifyInterval)
		assert.Equal(t, plugin.ExpectExecuteInterval, cfg.PluginInfos.Plugins[i].ExecuteInterval)
		assert.Equal(t, plugin.ExpectPort, cfg.PluginInfos.Plugins[i].Port)

		assert.Equal(t, len(plugin.ExpectMethods), len(cfg.PluginInfos.Plugins[i].ExecutableMethods))

		for _, method := range cfg.PluginInfos.Plugins[i].ExecutableMethods {
			assert.Contains(t, plugin.ExpectMethods, method.Name)
		}
	}
}
