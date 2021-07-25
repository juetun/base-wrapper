package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestDateStandard(t *testing.T) {
	type args struct {
		value string
	}
	ti := time.Now()
	tests := []struct {
		name    string
		args    args
		wantT   string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				value: ti.Format(time.RFC3339),
			},
			wantT: ti.Format("2006-01-02 15:04:05"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := DateStandard(tt.args.value)
			if err != nil {
				t.Errorf("DateStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotT.Format("2006-01-02 15:04:05") != tt.wantT {
				t.Errorf("DateParse() gotStamp = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func TestDateParse(t *testing.T) {
	type args struct {
		value  string
		format []string
	}
	tests := []struct {
		name      string
		args      args
		wantStamp string
		wantErr   bool
	}{
		{
			args: args{
				value: "2019-01-08 13:50:30",
			},
			wantStamp: "2019-01-08 13:50:30",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStamp, err := DateParse(tt.args.value, tt.args.format...)
			if err != nil {
				t.Errorf("DateParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStamp.Format("2006-01-02 15:04:05") != tt.wantStamp {
				t.Errorf("DateParse() gotStamp = %v, want %v", gotStamp, tt.wantStamp)
			}
		})
	}
}

func TestDateDuration(t *testing.T) {
	type args struct {
		value    string
		baseTime *time.Time
		format   []string
	}
	t1, _ := DateParse("2022-01-08 00:00:00", DateTimeGeneral)
 	tests := []struct {
		name    string
		args    args
		wantRes string
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				value:    "2021-01-01 13:50:30",
				baseTime: &t1,
			},
			wantRes: "2021-01-01",
		},
		{
			name: "t2",
			args: args{
				value:    "2022-01-07 23:59:59",
				baseTime: &t1,
			},
			wantRes: "1秒钟前",
		},
		{
			name: "t3",
			args: args{
				value:    "2022-01-07 23:58:59",
				baseTime: &t1,
			},
			wantRes: "1分钟前",
		},
		{
			name: "t4",
			args: args{
				value:    "2022-01-07 22:58:59",
				baseTime: &t1,
			},
			wantRes: "1小时前",
		},
		{
			name: "t5",
			args: args{
				value:    "2022-01-01 00:00:00",
				baseTime: &t1,
			},
			wantRes: "1周前",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := DateDuration(tt.args.value, tt.args.baseTime, tt.args.format...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DateDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(gotRes)
			if gotRes != tt.wantRes {
				t.Errorf("DateDuration() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
