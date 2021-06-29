// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package redis

// import (
// 	"testing"
// 	"time"
//
// 	"github.com/go-redis/redis"
// 	"github.com/go-redis/redis/v8"
// 	"github.com/juetun/base-wrapper/lib/app/app_start"
// 	_ "github.com/juetun/base-wrapper/lib/app/init" // 加载公共插件项
// )
//
// func Test_utilsRedis_Get(t *testing.T) {
// 	app_start.NewPlugins().LoadPlugins() // 加载插件动作
// 	type fields struct {
// 		Key         string
// 		Type        string
// 		redis       *redis.Client
// 		getDataFunc GetDataFunc
// 		Expiration  time.Duration
// 	}
// 	type args struct {
// 		data interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := NewUtilsRedis().SetKey("test").SetType("object")
//
// 			if err := r.Get(tt.args.data); (err != nil) != tt.wantErr {
// 				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
