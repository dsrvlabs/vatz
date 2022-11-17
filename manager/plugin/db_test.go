package plugin

import (
	"database/sql"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDBWrite(t *testing.T) {
	wr, err := newWriter("vatz.db")
	assert.Nil(t, err)

	rd, err := newReader("vatz.db")
	assert.Nil(t, err)

	defer os.Remove("vatz.db")
	defer db.conn.Close()
	defer func() {
		once = sync.Once{}
	}()

	installedAt := time.Now()

	// Test insert.
	err = wr.AddPlugin(pluginEntry{
		Name:           "test",
		Repository:     "dummy",
		BinaryLocation: "home/status",
		Version:        "latest",
		InstalledAt:    installedAt,
	})
	assert.Nil(t, err)

	// Confirm insertion.
	plugin, err := rd.Get("test")

	assert.Nil(t, err)
	assert.Equal(t, "test", plugin.Name)
	assert.Equal(t, "dummy", plugin.Repository)
	assert.Equal(t, "home/status", plugin.BinaryLocation)
	assert.Equal(t, "latest", plugin.Version)
	assert.Equal(t, installedAt.UnixMilli(), plugin.InstalledAt.UnixMilli())

	// Test delete.
	err = wr.DeletePlugin("test")
	assert.Nil(t, err)

	// Confirm deleted.
	plugin, err = rd.Get("test")

	assert.Nil(t, plugin)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestDBNotExist(t *testing.T) {
	wr, err := newWriter("vatz.db")
	assert.Nil(t, err)

	rd, err := newReader("vatz.db")
	assert.Nil(t, err)

	defer os.Remove("vatz.db")
	defer db.conn.Close()
	defer func() {
		once = sync.Once{}
	}()

	// Insert
	err = wr.AddPlugin(pluginEntry{Name: "test", Repository: "dummy"})

	assert.Nil(t, err)

	plugin, err := rd.Get("not-exist")

	assert.Nil(t, plugin)
	assert.Equal(t, sql.ErrNoRows, err)
}

// TODO: Handle already exist.

func TestDBList(t *testing.T) {
	wr, err := newWriter("vatz.db")
	assert.Nil(t, err)

	rd, err := newReader("vatz.db")
	assert.Nil(t, err)

	defer os.Remove("vatz.db")
	defer db.conn.Close()
	defer func() {
		once = sync.Once{}
	}()

	installedAt := time.Now()

	// Add dummy plugins
	testPlugins := []pluginEntry{
		{
			Name:           "test",
			Repository:     "dummy",
			BinaryLocation: "home/status",
			Version:        "latest",
			InstalledAt:    installedAt,
		},
		{
			Name:           "test2",
			Repository:     "dummy",
			BinaryLocation: "home/status",
			Version:        "latest",
			InstalledAt:    installedAt,
		},
	}

	// insert.
	for _, p := range testPlugins {
		err = wr.AddPlugin(p)
		assert.Nil(t, err)
	}

	plugins, err := rd.List()

	assert.Nil(t, err)
	assert.Equal(t, len(testPlugins), len(plugins))

	for i, p := range plugins {
		assert.Equal(t, testPlugins[i].Name, p.Name)
		assert.Equal(t, testPlugins[i].Repository, p.Repository)
		assert.Equal(t, testPlugins[i].BinaryLocation, p.BinaryLocation)
		assert.Equal(t, testPlugins[i].Version, p.Version)
		assert.Equal(t, testPlugins[i].InstalledAt.Unix(), p.InstalledAt.Unix())
	}
}
