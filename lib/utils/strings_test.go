package utils

import "testing"

func TestIsIdCard(t *testing.T) {
	type args struct {
		idCard string
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			args:args{idCard:"511303198508154557"},
			wantOk:true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, err := IsIdCard(tt.args.idCard)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsIdCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("IsIdCard() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}