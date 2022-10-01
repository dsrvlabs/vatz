package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

const (
	// FlagConfig is name of CLI flags for config.
	FlagConfig = "config"

	// DefaultConfigFile is default file name of config.
	DefaultConfigFile = "default.yaml"
)

var (
	configOnce = &sync.Once{}
	vatzConfig *Config
)

// Config is Vatz config structure.
type Config struct {
	Vatz struct {
		ProtocolIdentifier    string           `yaml:"protocol_identifier"`
		Port                  int              `yaml:"port"`
		NotificationInfo      NotificationInfo `yaml:"notification_info"`
		HealthCheckerSchedule []string         `yaml:"health_checker_schedule"`
	} `yaml:"vatz_protocol_info"`

	PluginInfos PluginInfo `yaml:"plugins_infos"`
}

type NotificationInfo struct {
	DiscordSecret    string `yaml:"discord_secret"`
	PagerDutySecret  string `yaml:"pager_duty_secret"`
	HostName         string `yaml:"host_name"`
	DispatchChannels []struct {
		Channel string `yaml:"channel"`
		Secret  string `yaml:"secret"`
		ChatID  string `yaml:"chat_id"`
	} `yaml:"dispatch_channels"`
}

// PluginInfo contains general plugin info.
type PluginInfo struct {
	DefaultVerifyInterval  int      `yaml:"default_verify_interval"`
	DefaultExecuteInterval int      `yaml:"default_execute_interval"`
	DefaultPluginName      string   `yaml:"default_plugin_name"`
	Plugins                []Plugin `yaml:"plugins"`
}

// Plugin contains specific plugin info.
type Plugin struct {
	Name              string `yaml:"plugin_name"`
	Address           string `yaml:"plugin_address"`
	VerifyInterval    int    `yaml:"verify_interval"`
	ExecuteInterval   int    `yaml:"execute_interval"`
	Port              int    `yaml:"plugin_port"`
	ExecutableMethods []struct {
		Name string `yaml:"method_name"`
	} `yaml:"executable_methods"`
}

type parser struct {
	rawConfig map[string]interface{}
}

func (p *parser) loadConfigFile(path string) ([]byte, error) {
	var (
		rawYAML []byte
		err     error
	)

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("invalid response status %d", resp.StatusCode)
		}

		rawYAML, err = io.ReadAll(resp.Body)
	} else {
		rawYAML, err = ioutil.ReadFile(path)
	}

	if err != nil {
		return nil, err
	}

	return rawYAML, err
}

func (p *parser) parseYAML(contents []byte) (*Config, error) {
	newConfig := Config{}
	err := yaml.Unmarshal(contents, &newConfig)
	if err != nil {
		return nil, err
	}

	p.overrideDefault(&newConfig)

	return &newConfig, nil
}

func (p *parser) overrideDefault(config *Config) {
	for i, plugin := range config.PluginInfos.Plugins {
		if plugin.VerifyInterval == 0 {
			config.PluginInfos.Plugins[i].VerifyInterval = config.PluginInfos.DefaultVerifyInterval
		}

		if plugin.ExecuteInterval == 0 {
			config.PluginInfos.Plugins[i].ExecuteInterval = config.PluginInfos.DefaultExecuteInterval
		}

		if plugin.Name == "" {
			config.PluginInfos.Plugins[i].Name = config.PluginInfos.DefaultPluginName
		}
	}
}

// InitConfig - initializes VATZ config.
func InitConfig(configFile string) *Config {
	configOnce.Do(func() {
		p := parser{}
		configData, err := p.loadConfigFile(configFile)
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				log.Error().Str("module", "config").Msgf("loadConfig Error: %s", err)
				log.Error().Str("module", "config").Msg("Please, execute `.vatz init` first or set appropriate path for default.yaml")
			}
			panic(err)
		}
		config, err := p.parseYAML(configData)
		if err != nil {
			log.Error().Str("module", "config").Msgf("parseYAML Error: %s", err)
			panic(err)
		}

		vatzConfig = config
	})

	return vatzConfig
}

// GetConfig returns current Vatz config.
func GetConfig() *Config {
	return vatzConfig
}
