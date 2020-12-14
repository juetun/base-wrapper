// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package inte

import "time"

//缓存数据的接口，开发中需自定义实现逻辑
type BlockCacheInterface interface {
	//存储缓存数据
	//@param  name 缓存的kEY
	//@param  val  缓存的值
	//@param  cacheTime 缓存的时间
	//@return error
	Set(name string, val string, cacheTime time.Duration) (err error)

	//获取缓存数据
	//@param name 缓存的key
	//@return res 获取的数据值
	Get(name string) (res string, err error)

	//属性默认操作
	DefaultValue()
}