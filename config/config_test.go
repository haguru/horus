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
				Name: "crumbdb",
				Port: 50051,
				LogLevel: "DEBUG",
				Database: Database{
					Host: "localhost",
					Port: 27017,
					Name: "horus",
					Collection: "crumbs",
					Options: ServerOptions{
						SetStrict: true,
						SetDeprecationErrors: true,
					},
				},
			},
			wantErr: false,

		},
		{
			name: "file does not exist",
			args: args{
				configPath: "",
			},
			want: nil,
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
