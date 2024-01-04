package cmd

import (
	"errors"
	"fmt"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"syscall"
	"time"
)

func createStopCommand() *cobra.Command {
	log.Debug().Str("module", "cmd > stop").Msg("Stop command")
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop VATZ",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.InitConfig(configFile)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				maxAttempts = 3
				sc          = syscall.SIGINT
			)
			path, err := config.GetConfig().Vatz.AbsoluteHomePath()
			if err != nil {
				return err
			}
			filPath := fmt.Sprintf("%s/vatz.pid", path)
			/*
				As mentioned on start.go
				Stored process id into a separated file to avoid terminating the wrong vatz process
				if there are two multiple running vatz processes on the single machine.
			*/
			pidData, err := os.ReadFile(filPath)
			if err != nil {
				log.Error().Str("module", "cmd > stop").Msgf("There's no File vatz.pid at %s", path)
				return errors.New("Can't specify VATZ process id from vatz.pid, Please, Check VATZ is running.")
			}

			pid, err := strconv.Atoi(string(pidData))
			log.Debug().Str("module", "stop").Msgf("pid is %d", pid)
			if err != nil {
				log.Error().Str("module", "cmd > stop").Msgf("Please, check vatz process is running with pid %d", pid)
				return errors.New("Can't specify VATZ process id, Please, Check VATZ is running.")
			}

			if err := os.Remove(filPath); err != nil {
				log.Error().Str("module", "cmd > stop").Msgf("file deletion failed due to %s", err)
			} else {
				log.Debug().Str("module", "cmd > stop").Msgf("File successfully deleted at %s", filPath)
			}

			process, err := os.FindProcess(pid)
			if err != nil {
				log.Error().Str("module", "cmd > stop").Msgf("Failed to find process with pid %d", pid)
				log.Error().Str("module", "cmd > stop").Msg("Please, check vatz process is running")
				return err
			}

			ticker := time.NewTicker(3 * time.Second)
			defer ticker.Stop()

			for attempts := 0; attempts < maxAttempts; attempts++ {
				if attempts == 0 {
					//Show vatz stop message only once to the user
					log.Info().Msg("Sent termination signal to VATZ process, terminating ...")
				}
				if attempts == maxAttempts-1 {
					// Just kill process if terminating process at last attempt.
					sc = syscall.SIGTERM
				}

				log.Debug().Str("module", "cmd > stop").Msgf("syscall: %d", sc)
				if err := process.Signal(sc); err != nil {
					log.Error().Err(err).Msg("Failed to send termination signal")
					fmt.Printf("Failed to send SIGINT: %v\n", err)
				}

				<-ticker.C
				if err := process.Signal(syscall.Signal(0)); err != nil {
					log.Debug().Str("module", "cmd > stop").Msg("Process has exited.")
					return nil
				} else {
					log.Debug().Str("module", "cmd > stop").Msg("Process is still running. Attempting to send SIGINT again.")
				}
			}
			return nil
		},
	}
	return cmd
}
