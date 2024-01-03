package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dsrvlabs/vatz/manager/config"
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
	statusCommand = &cobra.Command{
		Use:   "status",
		Short: "Get statuses of Plugin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.InitConfig(configFile)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			statusRequest := fmt.Sprintf("%s/v1/plugin_status", vatzRPC)
			req, err := http.NewRequest(http.MethodGet, statusRequest, nil)
			if err != nil {
				return err
			}

			cli := http.Client{}
			resp, err := cli.Do(req)
			if err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}

			log.Debug().Str("module", "plugin").Msgf("Plugin(s) status is requested to  %s.", statusRequest)

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

			fmt.Println("***** Plugin Status *****")
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
		Example: "vatz plugin install <plugin's githubAddress> <pluginName>",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.InitConfig(configFile)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			const defaultVersion = "latest"
			pluginDir, err := config.GetConfig().Vatz.AbsoluteHomePath()
			if err != nil {
				return err
			}

			log.Debug().Str("module", "plugin").Msgf("Install new plugin %s at %s.", args[0], pluginDir)

			pluginVersion := defaultVersion
			if viper.GetString("plugin_version") != "" {
				pluginVersion = viper.GetString("plugin_version")
			}

			log.Debug().Str("module", "plugin").Msgf("Installing plugin version is %s.", pluginVersion)
			mgr := plugin.NewManager(pluginDir)
			err = mgr.Install(args[0], args[1], pluginVersion)
			if err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}

			log.Info().Str("module", "plugin").Msgf("A new plugin %s is successfully installed.", args[1])
			return nil
		},
	}

	uninstallCommand = &cobra.Command{
		Use:     "uninstall",
		Short:   "Uninstall plugin from plugin registry",
		Args:    cobra.ExactArgs(1), // TODO: Can I check real git repo?
		Example: "vatz plugin uninstall <pluginName>",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.InitConfig(configFile)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginDir, err := config.GetConfig().Vatz.AbsoluteHomePath()
			if err != nil {
				return err
			}

			log.Debug().Str("module", "plugin").Msgf("Uninstall a plugin %s from %s", args[0], pluginDir)

			// TODO: Handle invalid repo name.
			mgr := plugin.NewManager(pluginDir)
			pluginExist := false
			plugins, err := mgr.List()
			if err != nil {
				return err
			}
			for _, plugin := range plugins {
				if plugin.Name == args[0] {
					pluginExist = true
					break
				}
			}
			if !pluginExist {
				log.Error().Str("module", "plugin").Msgf("There's no plugin with the name %s from the installed plugin list. ", args[0])
				return errors.New("Please confirm plugin name again.")
			}
			if err = mgr.Uninstall(args[0]); err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}
			log.Info().Str("module", "plugin").Msgf("Plugin %s is successfully uninstalled from %s", args[0], pluginDir)
			return nil
		},
	}

	startCommand = &cobra.Command{
		Use:     "start",
		Short:   "Start installed plugin",
		Example: "vatz plugin start --plugin <pluginName> or vatz plugin start --plugin <pluginName> --args <arguments>",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.InitConfig(configFile)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginName := viper.GetString("start_plugin")
			exeArgs := viper.GetString("start_args")

			pluginDir, err := config.GetConfig().Vatz.AbsoluteHomePath()
			if err != nil {
				return err
			}

			log.Info().Str("module", "plugin").Msgf("Start plugin %s %s", pluginName, exeArgs)

			logfile := viper.GetString("start_log")
			if logfile == "" {
				logfile = fmt.Sprintf("%s/%s.log", pluginDir, pluginName)
			}

			f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Info().Str("module", "plugin").Err(err)
				return err
			}

			log.Debug().Str("module", "plugin").Msgf("Plugin log redirect to %s", logfile)

			mgr := plugin.NewManager(pluginDir)
			return mgr.Start(pluginName, exeArgs, f)
		},
	}

	stopCommand = &cobra.Command{
		Use:     "stop",
		Short:   "Stop running plugin",
		Example: "vatz plugin stop --plugin <pluginName>",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.InitConfig(configFile)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginName := viper.GetString("stop_plugin")
			pluginDir, err := config.GetConfig().Vatz.AbsoluteHomePath()
			if err != nil {
				return err
			}

			log.Info().Str("module", "plugin").Msgf("Stop plugin %s", pluginName)

			mgr := plugin.NewManager(pluginDir)
			return mgr.Stop(pluginName)
		},
	}

	enableCommand = &cobra.Command{
		Use:     "enable",
		Short:   "Enable plugin",
		Example: "vatz plugin enable --plugin <pluginName>",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.InitConfig(configFile)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginName := viper.GetString("enable_plugin")
			pluginDir, err := config.GetConfig().Vatz.AbsoluteHomePath()
			if err != nil {
				return err
			}

			log.Debug().Str("module", "plugin").Msgf("enable installed plugin %s at %s", pluginName, pluginDir)

			mgr := plugin.NewManager(pluginDir)
			err = mgr.Update(pluginName, true)
			if err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}
			return nil
		},
	}

	disableCommand = &cobra.Command{
		Use:     "disable",
		Short:   "Disable plugin",
		Example: "vatz plugin disable --plugin <pluginName>",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.InitConfig(configFile)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginName := viper.GetString("disable_plugin")
			pluginDir, err := config.GetConfig().Vatz.AbsoluteHomePath()
			if err != nil {
				return err
			}

			log.Debug().Str("module", "plugin").Msgf("disable installed plugin %s at %s", pluginName, pluginDir)

			mgr := plugin.NewManager(pluginDir)
			err = mgr.Update(pluginName, false)
			if err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}
			return nil
		},
	}

	listCommand = &cobra.Command{
		Use:     "list",
		Short:   "List installed plugin",
		Example: "vatz plugin list",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.InitConfig(configFile)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginDir, err := config.GetConfig().Vatz.AbsoluteHomePath()
			if err != nil {
				return err
			}
			log.Debug().Str("module", "plugin").Msgf("List plugins")
			mgr := plugin.NewManager(pluginDir)
			plugins, err := mgr.List()
			if err != nil {
				log.Error().Str("module", "plugin").Err(err)
				return err
			}

			w := table.NewWriter()
			w.SetOutputMirror(os.Stdout)
			w.AppendHeader(table.Row{"Name", "Is Enabled", "Install Date", "Repository", "Version"})

			for _, plugin := range plugins {
				dateStr := plugin.InstalledAt.Format("2006-01-02 15:04:05")
				w.AppendRow([]interface{}{
					plugin.Name, plugin.IsEnabled, dateStr, plugin.Repository, plugin.Version,
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

	cmd.PersistentFlags().StringVar(&configFile, "config", defaultFlagConfig, "VATZ config file.")

	statusCommand.PersistentFlags().StringVar(&vatzRPC, "rpc", defaultRPC, "RPC address of Vatz")
	startCommand.PersistentFlags().StringP("plugin", "p", "", "Installed plugin name")
	startCommand.PersistentFlags().StringP("args", "a", "", "Arguments")
	startCommand.PersistentFlags().StringP("log", "l", "", "Logfile")

	err := viper.BindPFlag("start_plugin", startCommand.PersistentFlags().Lookup("plugin"))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
	}
	err = viper.BindPFlag("start_args", startCommand.PersistentFlags().Lookup("args"))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
	}
	err = viper.BindPFlag("start_log", startCommand.PersistentFlags().Lookup("log"))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
	}

	stopCommand.PersistentFlags().StringP("plugin", "p", "", "Installed plugin name")
	err = viper.BindPFlag("stop_plugin", stopCommand.PersistentFlags().Lookup("plugin"))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
	}

	installCommand.PersistentFlags().StringP("version", "v", "", "Installed plugin version")
	err = viper.BindPFlag("plugin_version", installCommand.PersistentFlags().Lookup("version"))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
	}

	enableCommand.PersistentFlags().StringP("plugin", "p", "", "Installed plugin name")
	err = viper.BindPFlag("enable_plugin", enableCommand.PersistentFlags().Lookup("plugin"))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
	}

	disableCommand.PersistentFlags().StringP("plugin", "p", "", "Installed plugin name")
	err = viper.BindPFlag("disable_plugin", disableCommand.PersistentFlags().Lookup("plugin"))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
	}

	commands := []*cobra.Command{
		statusCommand,
		installCommand,
		uninstallCommand,
		startCommand,
		stopCommand,
		enableCommand,
		disableCommand,
		listCommand,
	}

	for _, pluginCmd := range commands {
		cmd.AddCommand(pluginCmd)
	}

	return cmd
}
