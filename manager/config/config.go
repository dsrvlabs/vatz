package config

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	model "github.com/dsrvlabs/vatz/manager/model"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

type config struct {
}

func (c config) getPingIntervals(pluginInfo interface{}, IntervalKey string) []int {
	var pingIntervals []int
	var defaultValue int
	pluginAPIs := pluginInfo.(map[interface{}]interface{})["plugins"].([]interface{})

	switch IntervalKey {
	case "verify_interval":
		defaultValue = pluginInfo.(map[interface{}]interface{})["default_verify_interval"].(int)
	case "execute_interval":
		defaultValue = pluginInfo.(map[interface{}]interface{})["default_execute_interval"].(int)
	}

	if len(pluginAPIs) > 0 {
		for idx := range pluginAPIs {
			if value, ok := pluginAPIs[idx].(map[interface{}]interface{})[IntervalKey].(int); ok {
				pingIntervals = append(pingIntervals, value)
			} else {
				pingIntervals = append(pingIntervals, defaultValue)
			}
		}
	}
	return pingIntervals
}

func (c config) parse(retrievalInfo model.Type, configData map[interface{}]interface{}) interface{} {
	// FLAG PROTOCOL | PLUGIN
	if retrievalInfo == model.Protocol {
		return configData["vatz_protocol_info"]
	} else {
		return configData["plugins_infos"]
	}
}

func (c config) getYMLData(str string, isDefault bool) map[interface{}]interface{} {
	wd, _ := os.Getwd()
	confPath := str
	if isDefault == true {
		confPath = wd + "/" + confPath
	}

	yamlFile, err := ioutil.ReadFile(confPath)

	if err != nil {
		log.Fatal(err)
	}

	data := make(map[interface{}]interface{})
	err2 := yaml.Unmarshal(yamlFile, &data)

	if err2 != nil {
		log.Fatal(err2)
	}

	return data
}

func (c config) getConfigFromURL() map[interface{}]interface{} {
	var inputArguments = len(os.Args)
	var configFromURL = make(map[interface{}]interface{})

	if inputArguments > 1 {
		url := os.Args[1]
		resp, err := http.Get(url)

		if err != nil {
			log.Fatal("cannot fetch URL %q: %v", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatal("Status error: %v", resp.StatusCode)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Read body: %v", err)
		}

		err2 := yaml.Unmarshal(data, &configFromURL)
		if err2 != nil {
			log.Fatal(err2)
		}

	}
	return configFromURL
}

//TODO: Update to function to create multiple clients
func (c config) getClients(pluginInfo interface{}) []pluginpb.PluginClient {
	var grpcClients []pluginpb.PluginClient

	pluginAPIs := pluginInfo.(map[interface{}]interface{})["plugins"].([]interface{})
	defaultConnectedTarget := "localhost:9091"

	if len(pluginAPIs) > 0 {

		for idx, _ := range pluginAPIs {
			clientAddress := pluginAPIs[idx].(map[interface{}]interface{})["plugin_address"].(string)
			clientPort := pluginAPIs[idx].(map[interface{}]interface{})["plugin_port"].(int)
			connectTarget := clientAddress + ":" + strconv.Itoa(clientPort)
			conn, err := grpc.Dial(connectTarget, grpc.WithInsecure())
			if err != nil {
				log.Fatal(err)
			}
			grpcClients = append(grpcClients, pluginpb.NewPluginClient(conn))
		}

	} else {

		conn, err := grpc.Dial(defaultConnectedTarget, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		//TODO: Please, Create a better client functions with static
		//defer conn.Close()
		grpcClients = append(grpcClients, pluginpb.NewPluginClient(conn))
	}

	return grpcClients
}

type Config interface {
	parse(retrievalInfo model.Type, data map[interface{}]interface{}) interface{}
	getYMLData(str string, isDefault bool) map[interface{}]interface{}
	getConfigFromURL() map[interface{}]interface{}
	getClients(interface{}) []pluginpb.PluginClient
	getPingIntervals(pluginInfo interface{}, IntervalKey string) []int
}

func NewConfig() Config {
	return &config{}
}
