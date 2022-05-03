package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"testing"
)

func TestDaoUserImpl_BatchData(t *testing.T) {
	type fields struct {
		ServiceDao base.ServiceDao
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := &DaoUserImpl{}
			if err := r.BatchData(); (err != nil) != tt.wantErr {
				t.Errorf("BatchData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
