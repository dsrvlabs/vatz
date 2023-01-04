package plugin

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	// TODO: Should be configurable?
	pluginDBName = "vatz.db"
)

// VatzPlugin describes plugin information
type VatzPlugin struct {
	ID          string    `json:"id"` // TODO: Handle this after #323
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Repository  string    `json:"repository"`
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installed_at"`
}

// VatzPluginManager provides management functions for plugin.
type VatzPluginManager interface {
	Install(repo, name, version string) error
	List() ([]VatzPlugin, error)

	Update() error

	Start(name, args string, logfile *os.File) error
}

type vatzPluginManager struct {
	home string
}

func (m *vatzPluginManager) Install(repo, name, version string) error {
	var fixedRepo string
	strTokens := strings.Split(repo, "://")
	if len(strTokens) >= 2 {
		fixedRepo = strTokens[1]
	} else {
		fixedRepo = strTokens[0]
	}

	log.Info().Str("module", "plugin").Msgf("Install new plugin %s", fixedRepo)

	var stdout, stderr bytes.Buffer

	os.Setenv("GOBIN", m.home)

	exeCmd := exec.Command("go", "install", fixedRepo+"@"+version)
	exeCmd.Stdout = &stdout
	exeCmd.Stderr = &stderr

	err := exeCmd.Run()
	if err != nil {
		log.Error().Str("module", "plugin").Msg(string(stderr.Bytes()))
		return err
	}

	dirTokens := strings.Split(repo, "/")
	binName := dirTokens[len(dirTokens)-1]

	origPath := fmt.Sprintf("%s/%s", m.home, binName)
	newPath := fmt.Sprintf("%s/%s", m.home, name)

	// Binary name should be changed.
	err = os.Rename(origPath, newPath)
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
		return err
	}

	dbWr, err := newWriter(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
		return err
	}

	err = dbWr.AddPlugin(pluginEntry{
		Name:           name,
		Repository:     repo,
		BinaryLocation: newPath,
		Version:        version,
		InstalledAt:    time.Now(),
	})

	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
		return err
	}

	return nil
}

func (m *vatzPluginManager) List() ([]VatzPlugin, error) {
	log.Info().Str("module", "plugin").Msgf("List")

	dbRd, err := newReader(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
		return nil, err
	}

	dbPlugins, err := dbRd.List()
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
		return nil, err
	}

	plugins := make([]VatzPlugin, len(dbPlugins))

	for i, p := range dbPlugins {
		plugins[i].ID = "" // TODO
		plugins[i].Name = p.Name
		plugins[i].Repository = p.Repository
		plugins[i].Location = p.BinaryLocation
		plugins[i].Version = p.Version
		plugins[i].InstalledAt = p.InstalledAt
	}

	return plugins, nil
}

func (m *vatzPluginManager) Update() error {
	return nil
}

func (m *vatzPluginManager) Start(name, args string, logfile *os.File) error {
	log.Info().Str("module", "plugin").Msgf("Start plugin %s", name)

	dbRd, err := newReader(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		return err
	}

	e, err := dbRd.Get(name)
	if err != nil {
		return err
	}

	cmd := exec.Command(e.BinaryLocation, args)
	cmd.Stdout = logfile

	return cmd.Start()
}

// NewManager creates new plugin manager.
func NewManager(vatzHome string) VatzPluginManager {
	return &vatzPluginManager{
		home: vatzHome,
	}
}
