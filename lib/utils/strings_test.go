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
			args:   args{idCard: "511303198508154557"},
			wantOk: true,
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

func TestHidTel(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
		wantErr bool
	}{
		{
			args: args{phone: "15108352617"},
			wantRes:"151****2617",
		},
		{
			args: args{phone: "+8613888111188"},
			wantRes:"+86138****1188",
		},
		{
			args: args{phone: "010-88111188"},
			wantRes:"010-88***88",
		},
		{
			args: args{phone: "0861088111188"},
			wantRes:"086108***188",
		},
		{
			args: args{phone: "086-010-88111188"},
			wantRes:"086-010-88***88",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := HidTel(tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("HidTel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("HidTel() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
