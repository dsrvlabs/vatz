package config

import (
	"net/http"
	"os"
	"sync"
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
	ExpectHostName               string
	DispatchChannels             []testChannelsExpect
	ExpectHealthCheckerSchedule  []string
	ExpectDefaultVerifyInterval  int
	ExpectDefaultExecuteInterval int
	ExpectDefaultPluginName      string
	ExpectRPCInfo                RPCInfo

	Plugins []testPluginExpects
}

type testChannelsExpect struct {
	ExpectChannel string
	ExpectSecret  string
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
		MockContents:          configDefaultContents,
		ExpectProtocolID:      "vatz",
		ExpectVatzPort:        9090,
		ExpectDiscordSecret:   "XXXXX",
		ExpectPagerDutySecret: "YYYYY",
		ExpectHostName:        "xxx-xxxx-xxxx",
		DispatchChannels: []testChannelsExpect{
			{
				ExpectChannel: "discord",
				ExpectSecret:  "https://xxxxx.xxxxxx",
			},
			{
				ExpectChannel: "telegram",
				ExpectSecret:  "https://yyyyy.yyyyyy",
			},
			{
				ExpectChannel: "pagerduty",
				ExpectSecret:  "https://zzzzz.zzzzzz",
			},
		},
		ExpectHealthCheckerSchedule:  []string{"* 1 * * *"},
		ExpectDefaultVerifyInterval:  15,
		ExpectDefaultExecuteInterval: 30,
		ExpectDefaultPluginName:      "vatz-plugin",
		ExpectRPCInfo: RPCInfo{
			Enabled:  true,
			Address:  "127.0.0.1",
			GRPCPort: 19090,
			HTTPPort: 19091,
		},
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
	assert.Equal(t, test.ExpectHostName, cfg.Vatz.NotificationInfo.HostName)
	for i, dispatchChannel := range test.DispatchChannels {
		assert.Equal(t, dispatchChannel.ExpectChannel, cfg.Vatz.NotificationInfo.DispatchChannels[i].Channel)
		assert.Equal(t, dispatchChannel.ExpectSecret, cfg.Vatz.NotificationInfo.DispatchChannels[i].Secret)
	}
	assert.Equal(t, test.ExpectHealthCheckerSchedule, cfg.Vatz.HealthCheckerSchedule)

	assert.Equal(t, test.ExpectDefaultVerifyInterval, cfg.PluginInfos.DefaultVerifyInterval)
	assert.Equal(t, test.ExpectDefaultExecuteInterval, cfg.PluginInfos.DefaultExecuteInterval)

	assert.True(t, test.ExpectRPCInfo.Enabled)
	assert.Equal(t, "127.0.0.1", test.ExpectRPCInfo.Address)
	assert.Equal(t, 19090, test.ExpectRPCInfo.GRPCPort)
	assert.Equal(t, 19091, test.ExpectRPCInfo.HTTPPort)

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

func TestOverrideDefaultValues(t *testing.T) {
	test := testConfigExpects{
		MockContents: configNoIntervalContents,

		ExpectProtocolID:      "vatz",
		ExpectVatzPort:        9090,
		ExpectDiscordSecret:   "hello",
		ExpectPagerDutySecret: "world",
		ExpectHostName:        "dummy0",
		DispatchChannels: []testChannelsExpect{
			{
				ExpectChannel: "discord",
				ExpectSecret:  "dummy1",
			},
			{
				ExpectChannel: "telegram",
				ExpectSecret:  "dummy2",
			},
			{
				ExpectChannel: "pagerduty",
				ExpectSecret:  "dummy3",
			},
		},
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
	assert.Equal(t, test.ExpectHostName, cfg.Vatz.NotificationInfo.HostName)
	for i, dispatchChannel := range test.DispatchChannels {
		assert.Equal(t, dispatchChannel.ExpectChannel, cfg.Vatz.NotificationInfo.DispatchChannels[i].Channel)
		assert.Equal(t, dispatchChannel.ExpectSecret, cfg.Vatz.NotificationInfo.DispatchChannels[i].Secret)
	}

	// On testcase, there is no RPCPort information and that means default value shoule be set.
	assert.Equal(t, DefaultGRPCPort, cfg.Vatz.RPCInfo.GRPCPort)
	assert.Equal(t, DefaultHTTPPort, cfg.Vatz.RPCInfo.HTTPPort)

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

func TestNotExistConfigFile(t *testing.T) {
	defer func() {
		recover()
	}()

	_ = InitConfig("not_existing_file.yaml")

	// DO NOT REACH HERE
	t.Error("no panic occures")
}

func TestInvalidYAMLFormat(t *testing.T) {
	defer func() {
		recover()
	}()

	configOnce = &sync.Once{} // Overrice Once

	f, err := createDefaultConfigFile(configInvalidYAMLContents)
	defer os.Remove(f.Name())

	assert.Nil(t, err)

	_ = InitConfig(DefaultConfigFile)

	// DO NOT REACH HERE
	t.Errorf("no panic")
}

func createDefaultConfigFile(contents string) (*os.File, error) {
	f, err := os.Create(DefaultConfigFile)
	if err != nil {
		return nil, err
	}

	_, err = f.WriteString(contents)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// This test should be the last one, because "GetConfig" function can execute once.
func TestGetConfig(t *testing.T) {
	test := testConfigExpects{
		MockContents: configDefaultContents,

		ExpectProtocolID:      "vatz",
		ExpectVatzPort:        9090,
		ExpectDiscordSecret:   "XXXXX",
		ExpectPagerDutySecret: "YYYYY",
		ExpectHostName:        "xxx-xxxx-xxxx",
		DispatchChannels: []testChannelsExpect{
			{
				ExpectChannel: "discord",
				ExpectSecret:  "https://xxxxx.xxxxxx",
			},
			{
				ExpectChannel: "telegram",
				ExpectSecret:  "https://yyyyy.yyyyyy",
			},
			{
				ExpectChannel: "pagerduty",
				ExpectSecret:  "https://zzzzz.zzzzzz",
			},
		},
		ExpectHealthCheckerSchedule:  []string{"* 1 * * *"},
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

	// Prepare files for testing.
	configOnce = &sync.Once{} // Overrice Once
	f, err := createDefaultConfigFile(test.MockContents)
	defer os.Remove(f.Name())

	assert.Nil(t, err)

	// Init Config.
	cfg := InitConfig(DefaultConfigFile)

	// Asserts.
	assert.Equal(t, test.ExpectProtocolID, cfg.Vatz.ProtocolIdentifier)
	assert.Equal(t, test.ExpectVatzPort, cfg.Vatz.Port)
	assert.Equal(t, test.ExpectHostName, cfg.Vatz.NotificationInfo.HostName)
	for i, dispatchChannel := range test.DispatchChannels {
		assert.Equal(t, dispatchChannel.ExpectChannel, cfg.Vatz.NotificationInfo.DispatchChannels[i].Channel)
		assert.Equal(t, dispatchChannel.ExpectSecret, cfg.Vatz.NotificationInfo.DispatchChannels[i].Secret)
	}
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
