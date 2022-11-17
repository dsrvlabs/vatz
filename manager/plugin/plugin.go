package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	// TODO: Should be configurable?
	pluginDBName = "vatz.db"
)

// VatzPluginManager provides management functions for plugin.
type VatzPluginManager interface {
	Install(repo, name, version string) error
	Update() error

	Start(name, args string) error
}

type vatzPluginManager struct {
	home string
}

func (m *vatzPluginManager) Install(repo, name, version string) error {
	log.Info().Str("module", "plugin").Msgf("Install new plugin %s", repo)

	os.Setenv("GOBIN", m.home)
	exeCmd := exec.Command("go", "install", repo+"@"+version)
	err := exeCmd.Run()
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
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

	// TODO
	// Do we need additional info?
	// version
	// binary path
	// install date
	// etc.
	err = dbWr.AddPlugin(pluginEntry{
		Name:       name,
		Repository: repo,
	})
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
		return err
	}

	// Required parameters: github repository.
	//  - But how can I confirm the repository is the implementation of Vatz plugin?

	return nil
}

func (m *vatzPluginManager) Update() error {
	return nil
}

func (m *vatzPluginManager) Start(name, args string) error {
	log.Info().Str("module", "plugin").Msgf("Start plugin %s", name)

	// TODO: How to handle log?

	dbRd, err := newReader(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		return err
	}

	e, err := dbRd.Get(name)
	if err != nil {
		return err
	}

	cmd := exec.Command(m.home+"/"+e.Name, args)
	return cmd.Start()
}

// NewManager creates new plugin manager.
func NewManager(vatzHome string) VatzPluginManager {
	return &vatzPluginManager{
		home: vatzHome,
	}
}
