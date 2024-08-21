package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	// DefaultGRPCPort is default port number of grpc service.
	DefaultGRPCPort = 19090

	// DefaultHTTPPort is default port number of http service.
	DefaultHTTPPort = 19091

	// DefaultHomePath  default home directory of VATZ.
	DefaultHomePath = "~/.vatz"
)

var (
	configOnce = &sync.Once{}
	vatzConfig *Config
)

// Config is Vatz config structure.
type Config struct {
	Vatz        VatzProtocolInfo `yaml:"vatz_protocol_info"`
	PluginInfos PluginInfo       `yaml:"plugins_infos"`
}

// VatzProtocolInfo is VATZ information.
type VatzProtocolInfo struct {
	ProtocolIdentifier    string           `yaml:"protocol_identifier"`
	Port                  int              `yaml:"port"`
	NotificationInfo      NotificationInfo `yaml:"notification_info"`
	HealthCheckerSchedule []string         `yaml:"health_checker_schedule"`
	RPCInfo               RPCInfo          `yaml:"rpc_info"`
	MonitoringInfo        MonitoringInfo   `yaml:"monitoring_info"`
	HomePath              string           `yaml:"home_path"`
}

// AbsoluteHomePath is the default home path
func (i VatzProtocolInfo) AbsoluteHomePath() (string, error) {
	if strings.HasPrefix(i.HomePath, "~") {
		homePath := os.Getenv("HOME")
		absPath := fmt.Sprintf("%s/%s", homePath, strings.Trim(i.HomePath, "~"))

		// Prevent double slash
		absPath, err := filepath.Abs(absPath)
		return absPath, err
	}

	absPath, err := filepath.Abs(i.HomePath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// NotificationInfo is notification structure.
type NotificationInfo struct {
	HostName                string   `yaml:"host_name"`
	DefaultReminderSchedule []string `yaml:"default_reminder_schedule"`
	DispatchChannels        []struct {
		Channel          string   `yaml:"channel"`
		Secret           string   `yaml:"secret"`
		ChatID           string   `yaml:"chat_id"`
		Flag             string   `yaml:"flag,omitempty"`
		ReminderSchedule []string `yaml:"reminder_schedule"`
	} `yaml:"dispatch_channels"`
}

// RPCInfo is structure for RPC service configuration.
type RPCInfo struct {
	Enabled  bool   `yaml:"enabled"`
	Address  string `yaml:"address"`
	GRPCPort int    `yaml:"grpc_port"`
	HTTPPort int    `yaml:"http_port"`
}

// MonitoringInfo is structure for RPC service configuration.
type MonitoringInfo struct {
	GCP struct {
		GCPCloudLogging GCPCloudLoggingInfo `yaml:"gcp_cloud_logging_info"`
	} `yaml:"gcp"`
	Prometheus struct {
		Enabled bool   `yaml:"enabled"`
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	} `yaml:"prometheus"`
}

type GCPCloudLoggingInfo struct {
	Enabled           bool                       `yaml:"enabled"`
	GCPCredentialInfo CloudLoggingCredentialInfo `yaml:"cloud_logging_credential_info"`
}

type CloudLoggingCredentialInfo struct {
	ProjectID       string   `yaml:"project_id"`
	CredentialsType string   `yaml:"credentials_type"`
	Credentials     string   `yaml:"credentials"`
	CheckerSchedule []string `yaml:"checker_schedule"`
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
		Flag string `yaml:"flag,omitempty"`
	} `yaml:"executable_methods"`
}

type parser struct {
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
		if err != nil {
			return nil, err
		}
	} else {
		rawYAML, err = os.ReadFile(path)
		if err != nil {
			return nil, err
		}
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

	p.overwrite(&newConfig)

	return &newConfig, nil
}

func (p *parser) overwrite(config *Config) {
	if config.Vatz.RPCInfo.GRPCPort == 0 {
		config.Vatz.RPCInfo.GRPCPort = DefaultGRPCPort
	}

	if config.Vatz.RPCInfo.HTTPPort == 0 {
		config.Vatz.RPCInfo.HTTPPort = DefaultHTTPPort
	}

	if config.Vatz.HomePath == "" {
		config.Vatz.HomePath = DefaultHomePath
	}

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

func (p *parser) duplicatedPlugin(config *Config) {
	b := make(map[string][]int)
	for _, p := range config.PluginInfos.Plugins {
		b[p.Name] = append(b[p.Name], p.Port)
	}

	for pName, port := range b {
		if len(port) > 1 {
			log.Warn().Str("module", "config").
				Msgf(fmt.Sprintf("The plugin(%s) with the same name are currently up and running on %v ports.", pName, port))
		}
	}
}

// InitConfig - initializes VATZ config.
func InitConfig(configFile string) (*Config, error) {
	if vatzConfig != nil {
		log.Info().Str("module", "config").Msgf("Config already loaded")
		return vatzConfig, nil
	}

	var configError error

	wg := sync.WaitGroup{}
	wg.Add(1)

	configOnce.Do(func() {
		// TODO: How do I add default values?
		log.Debug().Str("module", "config").Msgf("Load Config %s", configFile)

		defer wg.Done()
		var configData []byte

		p := parser{}
		configData, configError = p.loadConfigFile(configFile)
		if configError != nil {
			return
		}

		vatzConfig, configError = p.parseYAML(configData)
		if configError != nil {
			log.Error().Str("module", "config").Msgf("parseYAML Error: %s", configError)
			return
		}

		p.duplicatedPlugin(vatzConfig)
	})

	wg.Wait()

	return vatzConfig, configError
}

// GetConfig returns current Vatz config.
func GetConfig() *Config {
	return vatzConfig
}
