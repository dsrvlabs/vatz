package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPluginInstall(t *testing.T) {
	root := cobra.Command{}
	root.AddCommand(createPluginCommand())
	root.SetArgs([]string{"plugin", "install", "github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/active_status@latest"})

	err := root.Execute()

	fmt.Println(err)

	assert.Nil(t, err)
}
