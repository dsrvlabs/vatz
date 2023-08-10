package utils

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var (
	cmdConsoleWriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
)

func InitializeLogger() {
	SetGlobalLogLevel(zerolog.InfoLevel)
	log.Logger = log.Output(cmdConsoleWriter)
}

func SetGlobalLogLevel(level zerolog.Level) {
	/*
			zerolog allows for logging at the following levels (from highest to lowest):
					panic (zerolog.PanicLevel, 5)
					fatal (zerolog.FatalLevel, 4)
					error (zerolog.ErrorLevel, 3)
					warn (zerolog.WarnLevel, 2)
					info (zerolog.InfoLevel, 1)
					debug (zerolog.DebugLevel, 0)
					trace (zerolog.TraceLevel, -1)
			VATZ's global loglevel is Info, which hide debug and trace and ignores all other cases except
		    Info, Debug, and Trace.
	*/
	switch {
	case level == zerolog.DebugLevel:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case level == zerolog.TraceLevel:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func SetLog(logfile string, defaultFlagLog string) error {
	if logfile == defaultFlagLog {
		log.Logger = log.Output(cmdConsoleWriter)
	} else {
		f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: f, TimeFormat: time.RFC3339})
	}
	return nil
}
