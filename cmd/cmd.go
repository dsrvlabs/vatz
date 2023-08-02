package cmd

import (
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	ex "github.com/dsrvlabs/vatz/manager/executor"
	health "github.com/dsrvlabs/vatz/manager/healthcheck"
	tp "github.com/dsrvlabs/vatz/manager/types"
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

	configFile string
	logfile    string
	vatzRPC    string
	promPort   string
)

// CreateRootCommand creates root command of Cobra.
func CreateRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{}
	rootCmd.PersistentFlags().BoolP("debug", "d", true, "Set Log level to Debug.")

	rootCmd.AddCommand(createInitCommand(tp.LIVE))
	rootCmd.AddCommand(createStartCommand())
	rootCmd.AddCommand(createPluginCommand())
	rootCmd.AddCommand(createVersionCommand())

	return rootCmd
}
