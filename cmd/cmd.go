package cmd

import (
	"github.com/spf13/cobra"

	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	ex "github.com/dsrvlabs/vatz/manager/executor"
	health "github.com/dsrvlabs/vatz/manager/healthcheck"
)

const (
	defaultFlagConfig = "default.yaml"
	defaultFlagLog    = ""
	defaultRPC        = "http://localhost:19091"
	defaultPromPort   = "18080"
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
	rootCmd.AddCommand(createInitCommand())
	rootCmd.AddCommand(createStartCommand())
	rootCmd.AddCommand(createPluginCommand())
	rootCmd.AddCommand(createVersionCommand())

	return rootCmd
}
