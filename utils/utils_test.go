package utils

import (
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetClients(t *testing.T) {
	type args struct {
		plugins []config.Plugin
	}
	var tests []struct {
		name string
		args args
		want []pluginpb.PluginClient
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetClients(tt.args.plugins), "GetClients(%v)", tt.args.plugins)
		})
	}
}

func TestMakeUniqueValue(t *testing.T) {
	var testUnique1 = MakeUniqueValue("aa", "bb", 8080)
	var testUnique2 = MakeUniqueValue("GetCPU", "localhost", 9090)
	var testUnique3 = MakeUniqueValue("GetCPU", "128.97.26.11", 9090)

	assert.Equal(t, "aabb8080", testUnique1)
	assert.Equal(t, "GetCPUlocalhost9090", testUnique2)
	assert.Equal(t, "GetCPU128.97.26.119090", testUnique3)
}
