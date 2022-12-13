package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dsrvlabs/vatz/manager/plugin"
)

/*
TODO list for plugin command.

- How to implement package repository?
- How to manage installed plugins?
- Install / Remove / Update plugin.
- How to execute?
  - Pass args.
  - Start / Stop
*/

var (
	// TODO: Should be configurable.
	pluginDir = fmt.Sprintf("%s/.vatz", os.Getenv("HOME"))

	statusCommand = &cobra.Command{
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
	}

	installCommand = &cobra.Command{
		Use:     "install",
		Short:   "Install new plugin",
		Args:    cobra.ExactArgs(2), // TODO: Can I check real git repo?
		Example: "vats plugin install github.com/dsrvlabs/<somewhere> name",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("module", "plugin").Msgf("Install new plugin %s at %s", args[0], pluginDir)

			// TODO: Handle already installed.
			// TODO: Handle invalid repo name.
			mgr := plugin.NewManager(pluginDir)
			err := mgr.Install(args[0], args[1], "latest")
			if err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}
			return nil
		},
	}

	startCommand = &cobra.Command{
		Use:     "start",
		Short:   "Start installed plugin",
		Example: "vats plugin start pluginName",
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginName := viper.GetString("plugin")
			exeArgs := viper.GetString("args")

			log.Info().Str("module", "plugin").Msgf("Start plugin %s %s", pluginName, exeArgs)

			mgr := plugin.NewManager(pluginDir)
			return mgr.Start(pluginName, exeArgs)
		},
	}
)

func createPluginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Plugin commands",
	}

	statusCommand.PersistentFlags().StringVar(&vatzRPC, "rpc", defaultRPC, "RPC address of Vatz")

	startCommand.PersistentFlags().StringP("plugin", "p", "", "Installed plugin name")
	startCommand.PersistentFlags().StringP("args", "a", "", "Arguments")

	viper.BindPFlag("plugin", startCommand.PersistentFlags().Lookup("plugin"))
	viper.BindPFlag("args", startCommand.PersistentFlags().Lookup("args"))

	cmd.AddCommand(statusCommand)
	cmd.AddCommand(installCommand)
	cmd.AddCommand(startCommand)

	return cmd
}
