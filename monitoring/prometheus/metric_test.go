package prometheus

import (
	"reflect"
	"testing"

	"github.com/dsrvlabs/vatz/utils"
)

func Test_prometheusManager_getPluginUp(t *testing.T) {
	type fields struct {
		Protocol string
	}
	type args struct {
		hostName string
	}
	var tests []struct {
		name         string
		fields       fields
		args         args
		wantPluginUp map[int]*prometheusValue
	}
	var grpcClientWithPlugins []utils.GClientWithPlugin
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &prometheusManager{
				Protocol: tt.fields.Protocol,
			}
			if gotPluginUp := c.getPluginUp(tt.args.hostName, grpcClientWithPlugins); !reflect.DeepEqual(gotPluginUp, tt.wantPluginUp) {
				t.Errorf("getPluginUp() = %v, want %v", gotPluginUp, tt.wantPluginUp)
			}
		})
	}
}
