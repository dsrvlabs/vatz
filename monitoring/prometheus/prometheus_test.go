package prometheus

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitPrometheusServer(tt.args.addr, tt.args.port, tt.args.protocol); (err != nil) != tt.wantErr {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPrometheusManager(tt.args.protocol, tt.args.reg); !reflect.DeepEqual(got, tt.want) {
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
		prometheus.Metric
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
		*prometheus.Desc
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
