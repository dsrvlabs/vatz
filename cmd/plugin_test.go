package cmd

import (
	"os"
	"testing"

	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPluginInstall(t *testing.T) {
	defer os.Remove("cosmos-status")
	defer os.Remove("./vatz-test.db")

	_, err := config.InitConfig("../default.yaml")
	assert.Nil(t, err)

	cfg := config.GetConfig()
	cfg.Vatz.HomePath = os.Getenv("PWD")
	// pluginDir = os.Getenv("PWD")

	root := cobra.Command{}
	root.AddCommand(createPluginCommand())
	root.SetArgs([]string{
		"plugin",
		"install",
		"github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_active_status",
		"cosmos-status"})

	err = root.Execute()
	assert.Nil(t, err)
}

// TODO: Test Start.
