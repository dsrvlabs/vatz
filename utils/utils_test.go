package utils

import (
	"fmt"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/stretchr/testify/assert"
	"syscall"
	"testing"
	"time"
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

func TestInitializeChannel(t *testing.T) {
	sigs := InitializeChannel()
	done := make(chan bool)
	go func() {
		time.Sleep(time.Millisecond * 100) // Wait a bit before sending the signal
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
			fmt.Printf("Failed to kill process: %v\n", err)
		}
	}()

	select {
	case <-sigs:
		assert.True(t, true, "Signal received as expected")
	case <-time.After(time.Second):
		assert.Fail(t, "Expected to receive signal, but did not")
	case <-done:
		// This case ensures the goroutine has completed its execution
	}
}
