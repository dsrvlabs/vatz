package main

import (
	"flag"
	"github.com/dsrvlabs/vatz/cmd"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
)

const defaultLogLevel = zerolog.InfoLevel

var logLevel zerolog.Level

func init() {
	debugPtr := flag.Bool("debug", false, "Enable debug mode")

	// Parse the command-line flags
	flag.Parse()

	// If the "debug" flag is set to true, set the log level to DebugLevel
	if *debugPtr {
		utils.SetGlobalLogLevel(zerolog.DebugLevel)
	} else {
		utils.SetGlobalLogLevel(defaultLogLevel)
	}

	log.Logger = log.Output(utils.GetConsoleWriter())
}

func main() {
	rootCmd := cmd.CreateRootCommand()

	if err := rootCmd.Execute(); err != nil {
		if strings.Contains(err.Error(), "open default.yaml") {
			msg := "Please, Check config file default.yaml path or initialize VATZ with command `vatz init` to create config file `default.yaml`."
			log.Error().Str("module", "config").Msg(msg)
		} else {
			log.Error().Msgf("VATZ CLI command Error: %s", err)
		}
		//panic(err)
	}

}
