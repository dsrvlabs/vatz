package plugin

import (
	"bytes"
	"errors"
	"fmt"
	tp "github.com/dsrvlabs/vatz/types"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/process"

	"github.com/dsrvlabs/vatz/manager/config"
)

const (
	// TODO: Should be configurable?
	pluginDBName = "vatz.db"
)

// VatzPlugin describes plugin information
type VatzPlugin struct {
	Name        string    `json:"name"`
	IsEnabled   bool      `json:"is_enabled"`
	Location    string    `json:"location"`
	Repository  string    `json:"repository"`
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installed_at"`
}

// VatzPluginManager provides management functions for plugin.
type VatzPluginManager interface {
	Init(runType tp.Initializer) error
	Get(name string) (VatzPlugin, error)
	Install(repo, name, version string) error
	Uninstall(name string) error
	List() ([]VatzPlugin, error)

	SetEnabled(pluginID string, isEnabled bool) error

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

	path, err := config.GetConfig().Vatz.AbsoluteHomePath()
	if err != nil {
		return err
	}

	return initDB(fmt.Sprintf("%s/%s", path, dbName))
}

func (m *vatzPluginManager) Install(repo, name, version string) error {
	var fixedRepo string
	strTokens := strings.Split(repo, "://")
	if len(strTokens) >= 2 {
		fixedRepo = strTokens[1]
	} else {
		fixedRepo = strTokens[0]
	}

	log.Debug().Str("module", "plugin").Msgf("Install new plugin %s", fixedRepo)

	var stdout, stderr bytes.Buffer

	os.Setenv("GOBIN", m.home)

	exeCmd := exec.Command("go", "install", fixedRepo+"@"+version)
	exeCmd.Stdout = &stdout
	exeCmd.Stderr = &stderr

	err := exeCmd.Run()
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Install > exeCmd.Run Error: %s", stderr.String())
		return err
	}

	dirTokens := strings.Split(repo, "/")
	binName := dirTokens[len(dirTokens)-1]

	origPath := fmt.Sprintf("%s/%s", m.home, binName)
	newPath := fmt.Sprintf("%s/%s", m.home, name)

	// Binary name should be changed.
	err = os.Rename(origPath, newPath)
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Install > os.Rename Error: %s", err)
		return err
	}

	dbWr, err := newWriter(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Install > newWriter Error: %s", err)
		return err
	}

	err = dbWr.AddPlugin(pluginEntry{
		Name:           name,
		IsEnabled:      1,
		Repository:     repo,
		BinaryLocation: newPath,
		Version:        version,
		InstalledAt:    time.Now(),
	})

	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Install > dbWr.AddPlugin Error: %s", err)
		return err
	}

	log.Debug().Str("module", "plugin").Msgf("A new plugin %s from %s is installed at %s.", name, repo, newPath)
	return nil
}

func (m *vatzPluginManager) List() ([]VatzPlugin, error) {
	log.Debug().Str("module", "plugin").Msgf("List")

	dbRd, err := newReader(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Install > newReader Error: %s", err)
		return nil, err
	}

	dbPlugins, err := dbRd.List()
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Install > dbRd.List Error: %s", err)
		return nil, err
	}

	plugins := make([]VatzPlugin, len(dbPlugins))

	for i, p := range dbPlugins {
		isEnabled := false
		if p.IsEnabled == 1 {
			isEnabled = true
		}
		plugins[i].Name = p.Name
		plugins[i].IsEnabled = isEnabled
		plugins[i].Repository = p.Repository
		plugins[i].Location = p.BinaryLocation
		plugins[i].Version = p.Version
		plugins[i].InstalledAt = p.InstalledAt
	}

	return plugins, nil
}

func (m *vatzPluginManager) Uninstall(name string) error {
	log.Debug().Str("module", "plugin").Msgf("Uninstall")
	ps, err := m.findProcessByName(name)
	if err != nil {
		if !strings.Contains(err.Error(), "can't find the process") {
			return err
		}
	} else {
		running, err := ps.IsRunning()
		if err != nil {
			log.Error().Str("module", "plugin").Msgf("Uninstall > ps.IsRunning Error: %s", err)
			return err
		}
		if running {
			log.Error().Str("module", "plugin").Err(err).Msgf("Plugin %s is currently running, Please, stop plugin first.", name)
			return fmt.Errorf("Please, stop plugin: %s first.", name)
		}
	}

	pluginInfo, err := m.Get(name)
	if err != nil {
		return err
	}

	var stdout, stderr bytes.Buffer

	os.Setenv("GOBIN", m.home)

	exeCmd := exec.Command("rm", "-rf", pluginInfo.Location)
	exeCmd.Stdout = &stdout
	exeCmd.Stderr = &stderr

	err = exeCmd.Run()
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Uninstall > exeCmd.Run Error: %s", err)
		return err
	}

	dbWr, err := newWriter(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Uninstall > newWriter Error: %s", err)
		return err
	}

	err = dbWr.DeletePlugin(name)
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Uninstall > dbWr.DeletePlugin Error: %s", err)
		return err
	}

	return nil
}

func (m *vatzPluginManager) Get(name string) (VatzPlugin, error) {
	log.Debug().Str("module", "plugin").Msgf("Get %s", name)

	dbRd, err := newReader(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Get > newReader Error: %s", err)
		return VatzPlugin{}, err
	}

	dbPlugin, err := dbRd.Get(name)
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Get > dbRd.Get Error: %s", err)
		return VatzPlugin{}, err
	}

	isEnabled := false
	if dbPlugin.IsEnabled > 0 {
		isEnabled = true
	}

	return VatzPlugin{
		Name:        dbPlugin.Name,
		IsEnabled:   isEnabled,
		Repository:  dbPlugin.Repository,
		Location:    dbPlugin.BinaryLocation,
		Version:     dbPlugin.Version,
		InstalledAt: dbPlugin.InstalledAt,
	}, nil
}

func (m *vatzPluginManager) SetEnabled(pluginID string, isEnabled bool) error {
	dbWr, err := newWriter(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Update > newWriter Error: %s", err)
		return err
	}

	err = dbWr.UpdatePluginEnabling(pluginID, isEnabled)
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Update > dbWr.UpdatePlugin Error: %s", err)
		return err
	}
	return nil
}

func (m *vatzPluginManager) Start(name, args string, logfile *os.File) error {
	log.Debug().Str("module", "plugin").Msgf("Start plugin %s", name)

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
	log.Info().Str("module", "plugin").Msgf("Plugin %s is successfully started.", name)
	return cmd.Start()
}

func (m *vatzPluginManager) Stop(name string) error {
	log.Debug().Str("module", "plugin").Msgf("Stop plugin %s", name)

	ps, err := m.findProcessByName(name)
	if err != nil {
		return err
	}

	err = ps.Kill()
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("Stop > ps.Kill Error: %s", err)
		return err
	}
	log.Info().Str("module", "plugin").Msgf("Plugin %s is successfully stopped.", name)
	return nil
}

func (m *vatzPluginManager) findProcessByName(name string) (*process.Process, error) {
	log.Debug().Str("module", "plugin").Msgf("Find Process %s", name)

	processes, err := process.Processes()
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("findProcessByName > process.Processes Error: %s", err)
		return nil, err
	}

	for _, p := range processes {
		pName, err := p.Name()
		if err != nil {
			log.Error().Str("module", "plugin").Msgf("findProcessByName > p.Name Error: %s", err)
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
	fullPath := fmt.Sprintf("%s/%s", vatzHome, pluginDBName)
	pManager := &vatzPluginManager{
		home: vatzHome,
	}

	if _, err := os.Stat(fullPath); err == nil {
		dbWr, err := newWriter(fullPath)
		if err != nil {
			log.Error().Str("module", "plugin").Msgf("NewManager > newWriter Error: %s", err)
		}
		err = dbWr.MigratePluginTable()
		if err != nil {
			log.Error().Str("module", "plugin").Msgf("NewManager > MigratePluginTable Error: %s", err)
		}

	}
	return pManager
}
