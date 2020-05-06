package metric

import (
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		options []func(*Client) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			// Returns newrelic.NewApplication error due to invalid license code
			"InvalidClientInitialization",
			args{
				[]func(*Client) error{
					SetAppname("test"),
					SetLicense("12345test"),
					SetLabels(map[string]string{}),
				},
			},
			true,
		},
		{
			"InvalidClientInitialization",
			args{
				[]func(*Client) error{
					SetAppname("test"),
					SetLicense("12345test"),
				},
			},
			true,
		},
		{
			"InvalidClientInitialization",
			args{
				[]func(*Client) error{
					SetAppname("test"),
					SetLabels(map[string]string{}),
				},
			},
			true,
		},
		{
			"InvalidClientInitialization",
			args{
				[]func(*Client) error{
					SetLicense("12345test"),
					SetLabels(map[string]string{}),
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSetAppname(t *testing.T) {
	tests := []struct {
		name    string
		appname string
		want    string
	}{
		{
			"ValidAppnameSetter",
			"test",
			"test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := SetAppname(tt.appname)
			client := new(Client)

			if err := fn(client); err != nil {
				t.Error("want to not have error")
				t.FailNow()
			}

			if !reflect.DeepEqual(client.appname, tt.want) {
				t.Error("invalid response")
				t.FailNow()
			}
		})
	}
}

func TestSetLicense(t *testing.T) {
	tests := []struct {
		name    string
		license string
		want    string
	}{
		{
			"ValidLicenseSetter",
			"12345test",
			"12345test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := SetLicense(tt.license)
			client := new(Client)

			if err := fn(client); err != nil {
				t.Error("want to not have error")
				t.FailNow()
			}

			if !reflect.DeepEqual(client.license, tt.want) {
				t.Error("invalid response")
				t.FailNow()
			}
		})
	}
}

func TestSetLabels(t *testing.T) {
	tests := []struct {
		name   string
		labels map[string]string
		want   map[string]string
	}{
		{
			"ValidLabelsSetter",
			map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := SetLabels(tt.labels)
			client := new(Client)

			if err := fn(client); err != nil {
				t.Error("want to not have error")
				t.FailNow()
			}

			if !reflect.DeepEqual(client.labels, tt.want) {
				t.Error("invalid response")
				t.FailNow()
			}
		})
	}
}

func TestInterface(t *testing.T) {
	// Check if the interface is compatible with type Metricer
	var metricer Metricer
	metricer, _ = NewClient()
	_ = metricer
}
