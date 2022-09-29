package main

import (
	"os"
	"time"

	config "github.com/dsrvlabs/vatz/manager/config"
	"github.com/spf13/cobra"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultFlagConfig = "default.yaml"
	defaultFlagLog    = ""
)

var (
	configFile string
	logfile    string
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
  notification_info:
    host_name: "Your machine name"
    dispatch_channels:
      - channel: "discord"
        secret: "Your channel secret"
      - channel: "telegram"
        secret: "Your channel secret"
		chat_id: "482109801"
  health_checker_schedule:
    - "0 1 * * *"
plugins_infos:
  default_verify_interval: 15
  default_execute_interval: 30
  default_plugin_name: "vatz-plugin"
  plugins:
    - plugin_name: "sample1"
      plugin_address: "localhost"
      plugin_port: 9001
      executable_methods:
        - method_name: "sampleMethod1"
    - plugin_name: "sample2"
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

func createStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start VATZ",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if logfile == defaultFlagLog {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
			} else {
				f, err := os.Create(logfile)
				if err != nil {
					return err
				}

				log.Logger = log.Output(f)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("module", "main").Msg("start")
			log.Info().Str("module", "main").Msgf("load config %s", configFile)
			log.Info().Str("module", "main").Msgf("logfile %s", logfile)

			config.InitConfig(configFile)

			ch := make(chan os.Signal, 1)
			return initiateServer(ch)
		},
	}

	cmd.PersistentFlags().StringVar(&configFile, "config", defaultFlagConfig, "VATZ config file.")
	cmd.PersistentFlags().StringVar(&logfile, "log", defaultFlagLog, "log file export to.")

	return cmd
}
