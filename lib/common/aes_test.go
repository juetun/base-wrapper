package common

import (
	"testing"

	"github.com/juetun/base-wrapper/lib/base"
)


func TestAes_Encryption(t *testing.T) {
	type fields struct {
		Context *base.Context
	}
	type args struct {
		text   string
		aesKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes string
		wantErr bool
	}{
		{
			args: args{
				text:   "eyJoX2FwcCI6Imp1ZXR1biIsImhfdGVybWluYWwiOiJ3ZWJzaXRlIiwiaF9jaGFubmVsIjoidXNyIiwiaF92ZXJzaW9uIjoiMS4wIiwiaF9kZWJ1ZyI6dHJ1ZX0=",
				aesKey: "jueTungygoaesctr",
			},
			wantRes: "vPXu4GKAc2B2V8CXasiQ59LnZzSfdEBeFIx2zNLSk8l7o9j2cZmwd/ZYhChaS9e2oY6a2Ur0X4WdSjI3RR/9DADFgozX5lnQK3ZNjKrV/cdwl1ENdBrtp/9E5mFYe9SXQAC51I0RImimZfr7W9B/C7cP6CA/07NBE/U1YcuPsZ8=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Aes{
				Context: tt.fields.Context,
			}
			gotRes, err := r.EncryptionCtr(tt.args.text, tt.args.aesKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encryption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v",gotRes)
			if gotRes != tt.wantRes {
				t.Errorf("Encryption() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestAes_Decrypt(t *testing.T) {
	type fields struct {
		Context *base.Context
	}
	type args struct {
		encrypted string
		aesKey    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes string
		wantErr bool
	}{
		{
			args: args{
				//encrypted: "vPXu51n6dH5afcOSefKqrOnnZ2m0dH5EBNN5w/n8j8p4s430Yaewd91Yhy1Jdti5jZWBxkrweoClYl0zbTLLSQHGgpP43QuE",
				//encrypted: "vPXu51n6dH5afcOSefKqrOnnZ2m0dH5EBNN5w/n8j8p4s430Yaewd91Yhy1Jdti5jZWBxkrweoClYl0zbTLLSQHGgpP43QuERyYL/eu0gqY=",
				encrypted: "vPXu4GKAc2B2V8CXasiQ59LnZzSfdEBeFIx2zNLSk8l7o9j2cZmwd/ZYhChaS9e2oY6a2Ur0X4WdSjI3RR/9DADFgozX5lnQK3ZNjKrV/cdwl1ENdBrtp/9E5mFYe9SXQAC51I0RImimZfr7W9B/C7cP6CA/07NBE/U1YcuPsZ8=",
				//encrypted: "vPXu4GKAc2B2V8CXasiQ59LnZzSfdEBeFIx2zNLSk8l7o9j2cZmwd/ZYhChaS9e2oY6a2Ur0X4WdSjI3RR/9DADFgozX5lnQFnlRgYLrvsddkiJQdnH9sf19kWRzI4+XQjrE24gSFGiLYIn/WbVnVrQY93gS8YlAKsNTMKnaiKbA/QH/kHwkcTWD77j2u0tr",
				//encrypted: "vPXu4GKAc2B2V8CXasiQ59LnZzSfdEBeFIx2zNLSk8l7o9j2cZmwd/ZYhChaS9e2oY6a2Ur0X4WdSjI3RR/9DADFgozX5lnQFnlRgYLrvsddkiJQdnH9sf19kWRzI4+XQjrE24gSFGiLYIn/WbVnVrQY93gS8YlAKsNTMKnaiKbA/QH/kHwkcTWD77j2u0tr",
				aesKey:    "jueTungygoaesctr",
			},
			wantRes: "eyJoX2FwcCI6Imp1ZXR1biIsImhfdGVybWluYWwiOiJ3ZWJzaXRlIiwiaF9jaGFubmVsIjoidXNyIiwiaF92ZXJzaW9uIjoiMS4wIiwiaF9kZWJ1ZyI6dHJ1ZX0=",
		},
		//{
		//	args: args{
		//		encrypted: "b0pM44LPTAc=",
		//		aesKey:    ChatEncryptionKey,
		//	},
		//	wantRes: "asdfasdf",
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Aes{
				Context: tt.fields.Context,
			}
			gotRes, err := r.DecryptCtr(tt.args.encrypted, tt.args.aesKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v",gotRes)
			if gotRes != tt.wantRes {
				t.Errorf("Decrypt() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
