package plugin

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginManager(t *testing.T) {
	defer os.Remove("./active_status")
	defer os.Remove("./cosmos-active")
	defer os.Remove("./" + pluginDBName)

	// TODO: could be better using mocks.
	repo := "github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_active_status"

	binName := "cosmos-active"

	mgr := NewManager(os.Getenv("PWD"))
	err := mgr.Install(repo, binName, "latest")
	assert.Nil(t, err)

	_, err = os.Open("./active_status")
	assert.True(t, errors.Is(err, os.ErrNotExist))

	_, err = os.Open(binName)
	assert.Nil(t, err)

	// Test Execute
	err = mgr.Start(binName, "-valoperAddr=dummy")
	assert.Nil(t, err)

	// Test DB.
	rd, err := newReader("./" + pluginDBName)
	assert.Nil(t, err)

	e, err := rd.Get(binName)
	assert.Nil(t, err)

	assert.Equal(t, binName, e.Name)
	assert.Equal(t, repo, e.Repository)
}
