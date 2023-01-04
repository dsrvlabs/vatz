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
	"github.com/jedib0t/go-pretty/v6/table"
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
		Example: "vatz plugin install github.com/dsrvlabs/<somewhere> name",
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
		Example: "vatz plugin start pluginName",
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginName := viper.GetString("plugin")
			exeArgs := viper.GetString("args")

			log.Info().Str("module", "plugin").Msgf("Start plugin %s %s", pluginName, exeArgs)

			logfile := viper.GetString("log")
			if logfile == "" {
				logfile = fmt.Sprintf("%s/%s.log", pluginDir, pluginName)
			}

			f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Info().Str("module", "plugin").Err(err)
				return err
			}

			log.Info().Str("module", "plugin").Msgf("Plugin log redirect to %s", logfile)

			mgr := plugin.NewManager(pluginDir)
			return mgr.Start(pluginName, exeArgs, f)
		},
	}

	listCommand = &cobra.Command{
		Use:     "list",
		Short:   "List installed plugin",
		Example: "vatz plugin list",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("module", "plugin").Msgf("List plugins")

			mgr := plugin.NewManager(pluginDir)
			plugins, err := mgr.List()
			if err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}

			w := table.NewWriter()
			w.SetOutputMirror(os.Stdout)
			w.AppendHeader(table.Row{"Name", "Install Data", "Repository", "Version"})

			for _, plugin := range plugins {
				dateStr := plugin.InstalledAt.Format("2006-01-02 15:04:05")
				w.AppendRow([]interface{}{
					plugin.Name, dateStr, plugin.Repository, plugin.Version,
				})
			}

			w.Render()

			return nil
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
	startCommand.PersistentFlags().StringP("log", "l", "", "Logfile")

	viper.BindPFlag("plugin", startCommand.PersistentFlags().Lookup("plugin"))
	viper.BindPFlag("args", startCommand.PersistentFlags().Lookup("args"))
	viper.BindPFlag("log", startCommand.PersistentFlags().Lookup("log"))

	cmd.AddCommand(statusCommand)
	cmd.AddCommand(installCommand)
	cmd.AddCommand(startCommand)
	cmd.AddCommand(listCommand)

	return cmd
}
