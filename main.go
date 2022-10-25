package main

import (
	"os"
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
		panic(err)
	}
}
