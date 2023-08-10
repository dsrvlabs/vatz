package cmd

import (
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	ex "github.com/dsrvlabs/vatz/manager/executor"
	health "github.com/dsrvlabs/vatz/manager/healthcheck"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	defaultFlagConfig = "default.yaml"
	defaultFlagLog    = ""
	defaultRPC        = "http://localhost:19091"
	defaultPromPort   = "18080"
	defaultHomePath   = "~/.vatz"
)

var (
	healthChecker          = health.GetHealthChecker()
	executor               = ex.NewExecutor()
	dispatchers            []dp.Dispatcher
	defaultVerifyInterval  = 15
	defaultExecuteInterval = 30
	IsDebugLevel           bool
	IsTraceLevel           bool
	configFile             string
	logfile                string
	vatzRPC                string
	promPort               string
)

/*	GetRootCommand: Return Cobra Root command include all subcommands .*/
func GetRootCommand() *cobra.Command {
	rootCmd := CreateRootCommand()
	rootCmd.AddCommand(createInitCommand(tp.LIVE))
	rootCmd.AddCommand(createStartCommand())
	rootCmd.AddCommand(createPluginCommand())
	rootCmd.AddCommand(createVersionCommand())
	return rootCmd
}

/*	CreateRootCommand: Create Root command which initialize root command and global flags. */
func CreateRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if IsDebugLevel {
			utils.SetGlobalLogLevel(zerolog.DebugLevel)
		} else if IsTraceLevel {
			utils.SetGlobalLogLevel(zerolog.TraceLevel)
		}
	}}
	rootCmd.PersistentFlags().BoolVarP(&IsDebugLevel, "debug", "", false, "Enable debug mode on Log.")
	rootCmd.PersistentFlags().BoolVarP(&IsTraceLevel, "trace", "", false, "Enable Trace mode on Log.")
	return rootCmd
}
