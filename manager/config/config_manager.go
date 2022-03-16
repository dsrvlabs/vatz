package config

import (
	"fmt"
)

var (
	configInstance Config
	CManager       config_manager
)

func init() {
	configInstance = NewConfig()
}

type config_manager struct {
}

func (c *config_manager) Parse() error {
	fmt.Println("this is a config Parser ")
	return nil
}
