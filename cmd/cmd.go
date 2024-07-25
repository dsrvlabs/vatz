package cmd

import (
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	ex "github.com/dsrvlabs/vatz/manager/executor"
	health "github.com/dsrvlabs/vatz/manager/healthcheck"
	tp "github.com/dsrvlabs/vatz/types"
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
	healthChecker = health.GetHealthChecker()
	executor      = ex.NewExecutor()
	sigs          = utils.InitializeChannel()
	dispatchers   []dp.Dispatcher
	isDebugLevel  bool
	isTraceLevel  bool
	configFile    string
	logfile       string
	vatzRPC       string
	promPort      string
)

// GetRootCommand is Return Cobra Root command include all subcommands .
func GetRootCommand() *cobra.Command {
	rootCmd := CreateRootCommand()
	rootCmd.AddCommand(createInitCommand(tp.LIVE))
	rootCmd.AddCommand(createStartCommand())
	rootCmd.AddCommand(createStopCommand())
	rootCmd.AddCommand(createPluginCommand())
	rootCmd.AddCommand(createVersionCommand())
	return rootCmd
}

// CreateRootCommand is Create Root command which initialize root command and global flags.
func CreateRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if isDebugLevel {
			utils.SetGlobalLogLevel(zerolog.DebugLevel)
		} else if isTraceLevel {
			utils.SetGlobalLogLevel(zerolog.TraceLevel)
		}
	}}
	rootCmd.PersistentFlags().BoolVarP(&isDebugLevel, "debug", "", false, "Enable debug mode on Log")
	rootCmd.PersistentFlags().BoolVarP(&isTraceLevel, "trace", "", false, "Enable trace mode on Log")
	return rootCmd
}
