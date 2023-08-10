package main

import (
	"github.com/dsrvlabs/vatz/cmd"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

var rootCmd *cobra.Command

func init() {
	//Set to Log level to Info which reduce log that doesn't be recorded and save log volumes.
	utils.InitializeLogger()
	rootCmd = cmd.GetRootCommand()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		if strings.Contains(err.Error(), "open default.yaml") {
			msg := "Please, Check config file default.yaml path or initialize VATZ with command `vatz init` to create config file `default.yaml`."
			log.Error().Str("module", "config").Msg(msg)
		} else {
			log.Error().Msgf("VATZ CLI command Error: %s", err)
		}
	}
}
