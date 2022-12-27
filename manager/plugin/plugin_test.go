package plugin

import (
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPluginManager(t *testing.T) {
	defer os.Remove("./active_status")
	defer os.Remove("./cosmos-active")
	defer os.Remove("./" + pluginDBName)
	defer func() {
		once = sync.Once{}
		db = nil
	}()

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

func TestPluginList(t *testing.T) {
	defer os.Remove("./" + pluginDBName)
	defer func() {
		once = sync.Once{}
		db = nil
	}()

	wr, err := newWriter("./" + pluginDBName)
	assert.Nil(t, err)

	// Add dummy plugins
	testPlugins := []pluginEntry{
		{
			Name:           "test",
			Repository:     "dummy",
			BinaryLocation: "home/status",
			Version:        "latest",
			InstalledAt:    time.Now(),
		},
		{
			Name:           "test2",
			Repository:     "dummy",
			BinaryLocation: "home/status",
			Version:        "latest",
			InstalledAt:    time.Now(),
		},
	}

	// Insert.
	for _, p := range testPlugins {
		err = wr.AddPlugin(p)
		assert.Nil(t, err)
	}

	pluginManager := NewManager(os.Getenv("PWD"))
	plugins, err := pluginManager.List()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(plugins))

	for i, p := range plugins {
		assert.Equal(t, testPlugins[i].Name, p.Name)
		assert.Equal(t, testPlugins[i].Repository, p.Repository)
		assert.Equal(t, testPlugins[i].BinaryLocation, p.Location)
		assert.Equal(t, testPlugins[i].Version, p.Version)
		assert.Equal(t, testPlugins[i].InstalledAt.Unix(), p.InstalledAt.Unix())
	}
}
