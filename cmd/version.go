package cmd

import (
	"fmt"

	"github.com/dsrvlabs/vatz/utils"
	"github.com/spf13/cobra"
)

func createVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "VATZ Version",
		RunE: func(cmd *cobra.Command, args []string) error {
			verStr := utils.GetVersion()
			fmt.Println(verStr)
			return nil
		},
	}

	return cmd
}
