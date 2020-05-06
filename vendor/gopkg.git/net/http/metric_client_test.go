package http

import (
	"net/http"
	"reflect"
	"testing"
)

func TestSetName(t *testing.T) {
	tests := []struct {
		name        string
		segmentName string
		want        string
	}{
		{
			"ValidNameSetter",
			"test",
			"test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := SetName(tt.segmentName)
			client := new(MetricClient)

			if err := fn(client); err != nil {
				t.Error("want to not have error")
				t.FailNow()
			}

			if !reflect.DeepEqual(client.name, tt.want) {
				t.Error("invalid response")
				t.FailNow()
			}
		})
	}
}

func TestSetHTTPClient(t *testing.T) {
	tests := []struct {
		name   string
		runner *http.Client
		want   *http.Client
	}{
		{
			"ValidHTTPClientSetter",
			nil,
			nil,
		},
		{
			"ValidHTTPClientSetter",
			http.DefaultClient,
			http.DefaultClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := SetHTTPClient(tt.runner)
			client := new(MetricClient)

			if err := fn(client); err != nil {
				t.Error("want to not have error")
				t.FailNow()
			}

			if !reflect.DeepEqual(client.runner, tt.want) {
				t.Error("invalid response")
				t.FailNow()
			}
		})
	}
}
