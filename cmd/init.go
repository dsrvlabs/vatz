package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/manager/plugin"
	tp "github.com/dsrvlabs/vatz/manager/types"
)

func createInitCommand(initializer tp.Initializer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("module", "main").Msg("init")

			template := `vatz_protocol_info:
  home_path: "%s"
  protocol_identifier: "Put Your Protocol here"
  port: 9090
  health_checker_schedule:
    - "0 1 * * *"
  notification_info:
    host_name: "Put your machine's host name"
    default_reminder_schedule:
      - "*/30 * * * *"
    dispatch_channels:
      - channel: "discord"
        secret: "Put your Discord Webhook"
      - channel: "pagerduty"
        secret: "Put your PagerDuty's Integration Key (Events API v2)"
      - channel: "telegram"
        secret: "Put Your Bot's Token"
        chat_id: "Put Your Chat's chat_id"
        reminder_schedule:
          - "*/5 * * * *"
  rpc_info:
    enabled: true
    address: "127.0.0.1"
    grpc_port: 19090
    http_port: 19091
  monitoring_info:
    prometheus:
      enabled: true
      address: "127.0.0.1"
      port: 18080
`
			samplePluginOptionTemplate := `plugins_infos:
  default_verify_interval: 15
  default_execute_interval: 30
  default_plugin_name: "vatz-plugin"
  plugins:
    - plugin_name: "samplePlugin1"
      plugin_address: "localhost"
      plugin_port: 9001
      executable_methods:
        - method_name: "sampleMethod1"
    - plugin_name: "samplePlugin2"
      plugin_address: "localhost"
      verify_interval: 7
      execute_interval: 9
      plugin_port: 10002
      executable_methods:
        - method_name: "sampleMethod2"
`
			defaultPluginOptionTemplate := `plugins_infos:
  default_verify_interval: 15
  default_execute_interval: 30
  default_plugin_name: "vatz-plugin"
  plugins:
    - plugin_name: "vatz_cpu_monitor"
      plugin_address: "localhost"
      plugin_port: 9001
      executable_methods:
        - method_name: "cpu_monitor"
    - plugin_name: "vatz_mem_monitor"
      plugin_address: "localhost"
      plugin_port: 9002
      executable_methods:
        - method_name: "mem_monitor"
    - plugin_name: "vatz_disk_monitor"
      plugin_address: "localhost"
      plugin_port: 9003
      executable_methods:
        - method_name: "disk_monitor"
    - plugin_name: "vatz_net_monitor"
      plugin_address: "localhost"
      plugin_port: 9004
      executable_methods:
        - method_name: "net_monitor"
    - plugin_name: "vatz_block_sync"
      plugin_address: "localhost"
      plugin_port: 10001
      executable_methods:
        - method_name: "node_block_sync"
    - plugin_name: "vatz_node_is_alived"
      plugin_address: "localhost"
      plugin_port: 10002
      executable_methods:
        - method_name: "node_is_alived"
    - plugin_name: "vatz_peer_count"
      plugin_address: "localhost"
      plugin_port: 10003
      executable_methods:
        - method_name: "node_peer_count"
    - plugin_name: "vatz_active_status"
      plugin_address: "localhost"
      plugin_port: 10004
      executable_methods:
        - method_name: "node_active_status"
    - plugin_name: "vatz_gov_alarm"
      plugin_address: "localhost"
      plugin_port: 10005
      executable_methods:
        - method_name: "node_governance_alarm"`
      
			filename, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			homePath, err := cmd.Flags().GetString("home")
			if err != nil {
				return err
			}

			template = fmt.Sprintf(template, homePath)
			log.Info().Str("module", "main").Msgf("home path %s", homePath)
			log.Info().Str("module", "main").Msgf("create file %s", filename)

			f, err := os.Create(filename)
			if err != nil {
				return err
			}

			configOption, err := cmd.Flags().GetBool("all")
			if err != nil {
				return err
			}
			if configOption {
				template = template + defaultPluginOptionTemplate
			} else {
				template = template + samplePluginOptionTemplate
			}
			_, err = f.WriteString(template)
			if err != nil {
				return err
			}

			config.InitConfig(filename)

			pluginDir, err := config.GetConfig().Vatz.AbsoluteHomePath()
			if err != nil {
				return err
			}

			log.Info().Str("module", "main").Msgf("Plugin dir %s", pluginDir)
			mgr := plugin.NewManager(pluginDir)
			return mgr.Init(initializer)
		},
	}

	_ = cmd.PersistentFlags().StringP("output", "o", defaultFlagConfig, "New config file to create")
	_ = cmd.PersistentFlags().StringP("home", "p", defaultHomePath, "Home directory of VATZ")
  _ = cmd.PersistentFlags().BoolP("all", "a", false, "Create config yaml with all default setting of official plugins.")

	return cmd
}
