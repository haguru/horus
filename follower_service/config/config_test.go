package config

import (
	"reflect"
	"testing"
)

func TestReadLocalConfig(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *ServiceConfig
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				configPath: "../res/config.yaml",
			},
			want: &ServiceConfig{
				ServiceName: "follower_service",
				Consul: Consul{
					Host: "consul",
					Port: 8500,
				},
				Port:     50055,
				LogLevel: "DEBUG",
				Database: Database{
					Host:         "followerdb",
					Port:         27017,
					DatabaseName: "horus",
					Timeout:      "5s",
					PingInterval: "5s",
					Collection:   "users",
					Options: ServerOptions{
						SetStrict:            true,
						SetDeprecationErrors: true,
					},
				},
				Metrics: Metrics{
					Port: 52112,
				},
			},
			wantErr: false,
		},
		{
			name: "file does not exist",
			args: args{
				configPath: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadLocalConfig(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLocalConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadLocalConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
