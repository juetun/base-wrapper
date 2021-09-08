package common

import (
	"testing"

	"github.com/juetun/base-wrapper/lib/base"
)

const (
	ChatEncryptionKey = "wumansgygoaesctr" // 秘钥长度为16的倍数

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
				text:   "asdfasdf",
				aesKey: ChatEncryptionKey,
			},
			wantRes: "b0pM44LPTAc=",
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
				encrypted: "b0pM44LPTAc=",
				aesKey:    ChatEncryptionKey,
			},
			wantRes: "asdfasdf",
		},
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
			if gotRes != tt.wantRes {
				t.Errorf("Decrypt() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
