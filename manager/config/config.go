package config

import (
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	model "vatz/manager/model"
)

type config struct {
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
func (c config) getClient(pluginInfo interface{}) pluginpb.PluginClient {
	pluginAPIs := pluginInfo.(map[interface{}]interface{})["plugins"].([]interface{})
	connectTarget := ":9091"

	if len(pluginAPIs) > 0 {
		clientAddress := pluginAPIs[0].(map[interface{}]interface{})["plugin_address"].(string)
		clientPort := pluginAPIs[0].(map[interface{}]interface{})["plugin_port"].(int)
		connectTarget = clientAddress + ":" + strconv.Itoa(clientPort)
	}

	conn, err := grpc.Dial(connectTarget, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	//TODO: Please, Create a better client functions with static
	//defer conn.Close()
	return pluginpb.NewPluginClient(conn)
}

type Config interface {
	parse(retrievalInfo model.Type, data map[interface{}]interface{}) interface{}
	getYMLData(str string, isDefault bool) map[interface{}]interface{}
	getConfigFromURL() map[interface{}]interface{}
	getClient(interface{}) pluginpb.PluginClient
}

func NewConfig() Config {
	return &config{}
}
