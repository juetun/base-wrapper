package base

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeNormal_MarshalJSON(t1 *testing.T) {
	type fields struct {
		Time time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:   "123",
			fields: fields{Time: time.Now()},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TimeNormal{
				Time: tt.fields.Time,
			}
			res, err := json.Marshal(t)
			if err != nil {
				t1.Fatal(err)
				return
			}
			t1.Log(string(res))
		})
	}
}

func TestTimeNormal_UnmarshalJSON(t1 *testing.T) {
	type fields struct {
		Time time.Time
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "123",
			fields: fields{Time: time.Now()},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TimeNormal{
				Time: tt.fields.Time,
			}
			bt, err := json.Marshal(t)
			if err != nil {
				t1.Fatal(err)
				return
			}
			t1.Log(string(bt))
			var dt TimeNormal
			err=json.Unmarshal(bt,&dt)
			if err != nil {
				t1.Fatal(err)
				return
			}
			t1.Log(dt)
		})
	}
}
