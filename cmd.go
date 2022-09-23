package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	defaultRPC        = "http://localhost:19091"
)

var (
	configFile string
	logfile    string
	vatzRPC    string
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
    discord_secret: "Your Discord Webhook"
    pager_duty_secret: "Your Events API V2 Integration Key"
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

func createPluginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Plugin commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Get statuses of Plugin",
		RunE: func(cmd *cobra.Command, args []string) error {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/plugin_status", vatzRPC), nil)
			if err != nil {
				return err
			}

			cli := http.Client{}
			resp, err := cli.Do(req)
			if err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}

			respData, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}

			statusResp := struct {
				Status       string `json:"status"`
				PluginStatus []struct {
					Status     string `json:"status"`
					PluginName string `json:"pluginName"`
				} `json:"pluginStatus"`
			}{}

			err = json.Unmarshal(respData, &statusResp)
			if err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}

			fmt.Println("***** Plugin status *****")
			for i, plugin := range statusResp.PluginStatus {
				fmt.Printf("%d: %s [%s]\n", i+1, plugin.PluginName, plugin.Status)
			}

			return nil
		},
	})

	cmd.PersistentFlags().StringVar(&vatzRPC, "rpc", defaultRPC, "RPC address of Vatz")

	return cmd
}
