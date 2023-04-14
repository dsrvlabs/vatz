package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dsrvlabs/vatz/utils"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/process"
)

const (
	// TODO: Should be configurable?
	pluginDBName = "vatz.db"
)

// VatzPlugin describes plugin information
type VatzPlugin struct {
	PluginID    string    `json:"plugin_id"` // TODO: Handle this after #323
	Name        string    `json:"name"`
	IsEnabled   bool      `json:"is_enabled"`
	Location    string    `json:"location"`
	Repository  string    `json:"repository"`
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installed_at"`
}

// VatzPluginManager provides management functions for plugin.
type VatzPluginManager interface {
	Init() error

	Install(repo, name, version string) error
	List() ([]VatzPlugin, error)

	Update(pluginID string, isEnabled bool) error

	Start(name, args string, logfile *os.File) error
	Stop(name string) error
}

type vatzPluginManager struct {
	home string
}

func (m *vatzPluginManager) Init() error {
	return initDB(fmt.Sprintf("%s/%s", m.home, pluginDBName))
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
		log.Error().Str("module", "plugin").Err(err).Msg("")
		return err
	}

	dbWr, err := newWriter(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err).Msg("")
		return err
	}

	hasValue := utils.UniqueHashValue(fmt.Sprintf("%s%s", repo, version))
	err = dbWr.AddPlugin(pluginEntry{
		PluginID:       hasValue,
		Name:           name,
		IsEnabled:      1,
		Repository:     repo,
		BinaryLocation: newPath,
		Version:        version,
		InstalledAt:    time.Now(),
	})

	if err != nil {
		log.Error().Str("module", "plugin").Err(err).Msg("")
		return err
	}

	return nil
}

func (m *vatzPluginManager) List() ([]VatzPlugin, error) {
	log.Info().Str("module", "plugin").Msgf("List")

	dbRd, err := newReader(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err).Msg("")
		return nil, err
	}

	dbPlugins, err := dbRd.List()
	if err != nil {
		log.Error().Str("module", "plugin").Err(err).Msg("")
		return nil, err
	}

	plugins := make([]VatzPlugin, len(dbPlugins))

	for i, p := range dbPlugins {
		isEnabled := false
		if p.IsEnabled == 1 {
			isEnabled = true
		}
		plugins[i].PluginID = p.PluginID
		plugins[i].Name = p.Name
		plugins[i].IsEnabled = isEnabled
		plugins[i].Repository = p.Repository
		plugins[i].Location = p.BinaryLocation
		plugins[i].Version = p.Version
		plugins[i].InstalledAt = p.InstalledAt
	}

	return plugins, nil
}

func (m *vatzPluginManager) Update(pluginID string, isEnabled bool) error {
	dbWr, err := newWriter(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err).Msg("Get new DB writer")
		return err
	}

	err = dbWr.UpdatePlugin(pluginID, isEnabled)
	if err != nil {
		log.Error().Str("module", "plugin").Err(err).Msg("Update DB plugin")
		return err
	}
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

	f := func(r rune) bool {
		return r == '=' || r == ' '
	}

	splits := strings.FieldsFunc(args, f)
	cmd := exec.Command(e.BinaryLocation, splits...)
	cmd.Stdout = logfile

	return cmd.Start()
}

func (m *vatzPluginManager) Stop(pluginID string) error {
	log.Info().Str("module", "plugin").Msgf("Stop plugin %s", pluginID)

	ps, err := m.findProcessByName(pluginID)
	if err != nil {
		return err
	}

	err = ps.Kill()
	if err != nil {
		log.Info().Str("module", "plugin").Msgf("Stop plugin %s", err)
		return err
	}

	return nil
}

func (m *vatzPluginManager) findProcessByName(name string) (*process.Process, error) {
	log.Info().Str("module", "plugin").Msgf("Find Process %s", name)

	processes, err := process.Processes()
	if err != nil {
		log.Info().Str("module", "plugin").Msgf("Find Process %s", err.Error())
		return nil, err
	}

	for _, p := range processes {
		pName, err := p.Name()
		if err != nil {
			log.Info().Str("module", "plugin").Msgf("Get Name of plugin %s", err.Error())
			continue
		}

		if pName == name {
			return p, nil
		}
	}
	return nil, errors.New("can't find the process")
}

// NewManager creates new plugin manager.
func NewManager(vatzHome string) VatzPluginManager {
	pManager := &vatzPluginManager{
		home: vatzHome,
	}
	dbWr, err := newWriter(fmt.Sprintf("%s/%s", vatzHome, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Err(err).Msg("")
	}
	dbWr.MigratePluginTable()
	return pManager
}
