// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

import (
	"testing"
)

func TestTraefikConfig_AppendToFile(t *testing.T) {
	type fields struct {
		RouterHttpConfig []TraefikDynamic
		RouterTcpConfig  []TcpTraefikRouters
		MapValue         []KeyValue
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				filename: "a.yml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewTraefikConfigTest()
			if err := r.AppendToFile(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("AppendToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTraefikConfig_KVShow(t *testing.T) {
	type fields struct {
		TraefikDynamic TraefikDynamic
		MapValue       []KeyValue
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "traefik_to_kv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewTraefikConfig()
			r.KVShow()
		})
	}
}
