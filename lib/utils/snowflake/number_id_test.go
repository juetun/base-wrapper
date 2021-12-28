package snowflake

import (
	"testing"

	"github.com/sony/sonyflake"
)

func TestSnowFlake_GetID(t *testing.T) {
	type fields struct {
		sFlake *sonyflake.Sonyflake
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint64
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
 			for i := 0; i <100; i++ {
				got, err := SFlake.GetID()
				if (err != nil) != tt.wantErr {
					t.Errorf("GetID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("GetID() got = %v", got)
			}

		})
	}
}
