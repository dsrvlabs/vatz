package utils

import (
	"encoding/json"
	"fmt"
	tp "github.com/dsrvlabs/vatz/types"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func GetPluginStatus(vatzRPC string) (tp.PluginState, error) {
	var newPluginStatus tp.PluginState
	statusRequest := fmt.Sprintf("%s/v1/plugin_status", vatzRPC)

	req, err := http.NewRequest(http.MethodGet, statusRequest, nil)
	if err != nil {
		return newPluginStatus, err
	}

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
		return newPluginStatus, err
	}

	log.Debug().Str("module", "plugin").Msgf("Plugin(s) status is requested to  %s.", statusRequest)

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
		return newPluginStatus, err
	}

	err = json.Unmarshal(respData, &newPluginStatus)
	if err != nil {
		log.Error().Str("module", "plugin").Err(err)
		return newPluginStatus, err
	}

	return newPluginStatus, nil
}
