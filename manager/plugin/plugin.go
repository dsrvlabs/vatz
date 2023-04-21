package plugin

import (
	"bytes"
	"errors"
	"fmt"
	tp "github.com/dsrvlabs/vatz/manager/types"
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
	ID          string    `json:"id"` // TODO: Handle this after #323
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Repository  string    `json:"repository"`
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installed_at"`
}

// VatzPluginManager provides management functions for plugin.
type VatzPluginManager interface {
	Init(runType tp.Initializer) error

	Install(repo, name, version string) error
	List() ([]VatzPlugin, error)

	Update() error

	Start(name, args string, logfile *os.File) error
	Stop(name string) error
}

type vatzPluginManager struct {
	home string
}

func (m *vatzPluginManager) Init(runType tp.Initializer) error {
	dbName := pluginDBName
	if runType == tp.TEST {
		dbName = "vatz-test.db"
	}
	return initDB(fmt.Sprintf("%s/%s", m.home, dbName))
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

	err = dbWr.AddPlugin(pluginEntry{
		Name:           name,
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

	f := func(r rune) bool {
		return r == '=' || r == ' '
	}

	splits := strings.FieldsFunc(args, f)
	cmd := exec.Command(e.BinaryLocation, splits...)
	cmd.Stdout = logfile

	return cmd.Start()
}

func (m *vatzPluginManager) Stop(name string) error {
	log.Info().Str("module", "plugin").Msgf("Stop plugin %s", name)

	ps, err := m.findProcessByName(name)
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
	return &vatzPluginManager{
		home: vatzHome,
	}
}
