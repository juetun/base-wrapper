package base

import "testing"

func TestSignUtils_Encrypt(t *testing.T) {
	type fields struct {
		mapExtend *MapExtend
	}
	type args struct {
		argJoin             string
		secret              string
		listenHandlerStruct ListenHandlerStruct
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes string
	}{
		{
			args: args{
				argJoin:             "Z2V0L2FwaS11c2VyL291dC91c2VyL2dldF91c2VyX2luZm8xNjg0NjM3NTg3NTYyanVldHVuZ3lnb2Flc2N0cnVzZXJfaGlkMA==",
				secret:              "jueTungygoaesctr",
				listenHandlerStruct: ListenHandlerStruct{},
			},
			wantRes: "c47a7fd1d3ef9e893d1f7117e264ad8bb7eb9b89",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SignUtils{
				mapExtend: tt.fields.mapExtend,
			}
			if gotRes := s.Encrypt(tt.args.argJoin, tt.args.secret, tt.args.listenHandlerStruct); gotRes != tt.wantRes {
				t.Errorf("Encrypt() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
