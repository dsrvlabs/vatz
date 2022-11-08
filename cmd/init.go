package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("module", "main").Msg("init")

			template := `vatz_protocol_info:
  protocol_identifier: "Put Your Protocol here"
  port: 9090
  health_checker_schedule:
    - "0 1 * * *"
  notification_info:
    host_name: "Put your machine's host name"
    default_reminder_schedule:
      - "*/15 * * * *"
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

plugins_infos:
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
			filename, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			log.Info().Str("module", "main").Msgf("create file %s", filename)

			f, err := os.Create(filename)
			if err != nil {
				return err
			}

			_, err = f.WriteString(template)
			if err != nil {
				return err
			}

			return nil
		},
	}

	_ = cmd.PersistentFlags().StringP("output", "o", defaultFlagConfig, "New config file to create")

	return cmd
}
