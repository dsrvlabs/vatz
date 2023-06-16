package main

import (
	"os"
	"strings"
	"time"

	"github.com/dsrvlabs/vatz/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
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
