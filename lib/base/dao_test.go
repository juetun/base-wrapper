package base

import (
	"testing"
)

func TestServiceDao_getDefaultColumnValue(t *testing.T) {
	type fields struct {
		Context *Context
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes string
	}{
		{
			args: args{
				name: "ServiceDao",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ServiceDao{
				Context: tt.fields.Context,
			}
			if gotRes := r.getDefaultColumnValue(tt.args.name); gotRes != tt.wantRes {
				t.Errorf("getDefaultColumnValue() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
