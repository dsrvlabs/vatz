package plugin

import (
	"database/sql"
	"os"
	"sync"
	"testing"

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

	// Test insert.
	err = wr.AddPlugin(pluginEntry{Name: "test", Repository: "dummy"})
	assert.Nil(t, err)

	// Confirm insertion.
	plugin, err := rd.Get("test")

	assert.Equal(t, "test", plugin.Name)
	assert.Equal(t, "dummy", plugin.Repository)

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
