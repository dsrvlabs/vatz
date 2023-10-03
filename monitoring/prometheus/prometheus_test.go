package prometheus

import (
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
	"testing"
)

func TestInitPrometheusServer(t *testing.T) {
	type args struct {
		addr     string
		port     string
		protocol string
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	var grpcClients []pluginpb.PluginClient
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitPrometheusServer(tt.args.addr, tt.args.port, tt.args.protocol, grpcClients); (err != nil) != tt.wantErr {
				t.Errorf("InitPrometheusServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newPrometheusManager(t *testing.T) {
	type args struct {
		protocol string
		reg      prometheus.Registerer
	}
	var tests []struct {
		name string
		args args
		want *prometheusManager
	}
	var grpcClients []pluginpb.PluginClient
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPrometheusManager(tt.args.protocol, tt.args.reg, grpcClients); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPrometheusManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prometheusManagerCollector_Collect(t *testing.T) {
	type fields struct {
		prometheusManager *prometheusManager
	}
	type args struct {
		ch chan<- prometheus.Metric
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = prometheusManagerCollector{
				prometheusManager: tt.fields.prometheusManager,
			}
		})
	}
}

func Test_prometheusManagerCollector_Describe(t *testing.T) {
	type fields struct {
		prometheusManager *prometheusManager
	}
	type args struct {
		ch chan<- *prometheus.Desc
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = prometheusManagerCollector{
				prometheusManager: tt.fields.prometheusManager,
			}
		})
	}
}
