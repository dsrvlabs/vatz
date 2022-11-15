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

	//Installed() error
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
		return err
	}

	dirTokens := strings.Split(repo, "/")
	binName := dirTokens[len(dirTokens)-1]

	origPath := fmt.Sprintf("%s/%s", m.home, binName)
	newPath := fmt.Sprintf("%s/%s", m.home, name)

	// Binary name should be changed.
	err = os.Rename(origPath, newPath)
	if err != nil {
		log.Info().Str("module", "plugin").Err(err)
		return err
	}

	dbWr, err := newWriter(fmt.Sprintf("%s/%s", m.home, pluginDBName))
	if err != nil {
		log.Info().Str("module", "plugin").Err(err)
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
		log.Info().Str("module", "plugin").Err(err)
		return err
	}

	// Required parameters: github repository.
	//  - But how can I confirm the repository is the implementation of Vatz plugin?

	return nil
}

func (m *vatzPluginManager) Update() error {
	return nil
}

//func (*vatzPluginManager) Installed() error {
//	return nil
//}

// NewManager creates new plugin manager.
func NewManager(vatzHome string) VatzPluginManager {
	return &vatzPluginManager{
		home: vatzHome,
	}
}
