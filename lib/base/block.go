package base

//页面缓存对象

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_obj"
)

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
}

type Handler func(block *Block) (err error)

//页面操作对象
type Block struct {
	Ctx                   context.Context `json:"ctx"`                     //上下文的操作对象 ，此处主要用来传递上下文参数
	ParentBlockCache      *BlockCache     `json:"parent_block_cache"`      //当前Block的父Block
	Name                  string          `json:"name"`                    //当前Block的系统唯一名字
	Data                  gin.H           `json:"data"`                    //当前Block的参数
	TempFile              string          `json:"temp_file"`               //html文件地址
	TemplateBaseDirectory string          `json:"template_base_directory"` //html模板文件所在的公共基础路径
	CacheBlock            *BlockCache     `json:"cache"`                   //当前模块缓存的基本参数
	RunChildBefore        []Handler       `json:"-"`
	RunBefore             []Handler       `json:"-"`          //渲染完数据后执行此方法，主要用来调试数据使用 //渲染完数据后执行此方法，主要用来调试数据使用,返回值为true时跳出
	RunAfter              []Handler       `json:"-"`          //渲染完数据前执行此方法，主要用来调试数据使用 //渲染完数据前执行此方法，主要用来调试数据使用,返回值为true时跳出
	ChildBock             []*Block        `json:"child_bock"` //当前的子BLOCK
}

//缓存信息对象
type BlockCache struct {
	ExpireTime time.Time           `json:"expire_time"` //静态化时间周期(单位秒)，设置当前BLOCK的生命周期，如果父Block>0时以父Block的值为准。
	CacheType  string              `json:"cache_type"`  //当前界面缓存类型 如 file:文件缓存,redis:缓存，database:数据库缓存
	Cache      BlockCacheInterface `json:"cache"`       //当前界面缓存的相关信息
	CacheKey   string              `json:"cache_key"`
}

//判断文件目录是否存在
func (r *Block) Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//将HTML模板文件绑定参数
func (r *Block) ParseHtml() (res string, err error) {
	var tmp *template.Template

	if !r.Exists(r.TempFile) {
		if err = fmt.Errorf("the template file(%s) is not exists",
			r.TempFile); err != nil {
			return
		}
		return
	}

	buf := new(bytes.Buffer)

	//拼接TemplateFile path
	if tmp, err = template.ParseFiles(r.TemplateBaseDirectory + r.TempFile); err != nil {
		return
	} else {
		tmp.Execute(buf, r.Data)
	}

	res = buf.String()
	return
}

//传递上下文参数
func (r *Block) setChildContext(item *Block) {

	data := gin.H{} //合并页面数据
	for key, value := range r.Data {
		data[key] = value
	}
	for key, value := range item.Data {
		data[key] = value
	}

	item.Data = data
	item.Ctx = r.Ctx
	item.ParentBlockCache = r.CacheBlock
}

//从缓存中获取数据
func (r *Block) getCache() (res string, err error) {
	if r.CacheBlock.Cache == nil {
		return
	}
	res, err = r.CacheBlock.Cache.Get(r.getKey())
	return
}

//将数据写入缓存
func (r *Block) writeToCache(data string) (err error) {

	if r.CacheBlock.ExpireTime.IsZero() {
		return
	}

	//缓存时间
	if lt := r.CacheBlock.ExpireTime.Unix() - time.Now().Unix(); lt > 0 {
		lTime := time.Duration(r.CacheBlock.ExpireTime.Unix() - time.Now().Unix())
		r.CacheBlock.Cache.Set(r.getKey(), data, lTime*time.Second)
	}

	return
}

//获取缓存数据的Key
func (r *Block) getKey() (res string) {
	if r.CacheBlock.CacheKey != "" {
		res = r.CacheBlock.CacheKey
		return
	}
	res = fmt.Sprintf("%s:%s", app_obj.App.AppName, r.Name)
	return
}

//解析模板数据
func (r *Block) Run() (res string, err error) {

	r.defaultValue() //初始化默认值

	if res, err = r.getCache(); err != nil {
		return
	}

	for _, item := range r.ChildBock {
		r.setChildContext(item)
		if r.Data[item.Name], err = item.Run(); err != nil {
			return
		}
	}

	//如果配置了运行后执行
	for _, runBefore := range r.RunBefore {
		if err = runBefore(r); err != nil {
			return
		}
	}
	if res, err = r.ParseHtml(); err != nil {
		return
	}

	//将数据写入缓存
	if err = r.writeToCache(res); err != nil {
		return
	}
	//如果配置了运行后执行
	for _, runAfter := range r.RunAfter {
		if err = runAfter(r); err != nil {
			return
		}
	}
	return
}

//BLOCK 默认数据逻辑处理
func (r *Block) defaultValue() {
	if r.Data == nil && len(r.Data) == 0 {
		r.Data = gin.H{}
	}
	if r.CacheBlock == nil {
		r.CacheBlock = &BlockCache{}
	}
	//如果名称没定义
	if r.Name == "" {
		r.Name = fmt.Sprintf("%T", r)
	}
	//初始化过期时间
	r.initExpireTime()
}

//缓存时间处理
func (r *Block) initExpireTime() {

	if r.ParentBlockCache == nil || r.CacheBlock == nil {
		return
	}
	//如果当前BLOCK的缓存时间为0 则不缓存
	//如果当前的父BLOCK缓存为0,则指定使用当前缓存时间
	if r.ParentBlockCache.ExpireTime.Unix() == 0 || r.CacheBlock.ExpireTime.Unix() == 0 {
		return
	}

	//如果当前BLOCK的父block不等于0,则本次缓存就为不缓存，（设置的过期时间是昨天当前时间）
	if r.ParentBlockCache.ExpireTime.Unix()-time.Now().Unix() > 0 {
		r.CacheBlock.ExpireTime = time.Now().AddDate(0, 0, -1)
	}

}

func NewBlock(option ...BlockOption) (block *Block) {
	block = &Block{}
	for _, handler := range option {
		handler(block)
	}
	return
}

type BlockOption func(block *Block)

func ChildBock(childBock []*Block) BlockOption {
	return func(block *Block) {
		block.ChildBock = childBock
	}
}
func RunAfter(runAfter []Handler) BlockOption {
	return func(block *Block) {
		block.RunAfter = runAfter
	}
}

func RunBefore(runBefore []Handler) BlockOption {
	return func(block *Block) {
		block.RunBefore = runBefore
	}
}
func RunChildBefore(runChildBefore []Handler) BlockOption {
	return func(block *Block) {
		block.RunChildBefore = runChildBefore
	}
}

func CacheBlock(cacheBlock *BlockCache) BlockOption {
	return func(block *Block) {
		block.CacheBlock = cacheBlock
	}
}

func TempFile(tempFile string) BlockOption {
	return func(block *Block) {
		block.TempFile = tempFile
	}
}

//html模板文件所在的基础路径
func TemplateBaseDirectory(templateBaseDirectory string) BlockOption {
	return func(block *Block) {
		block.TemplateBaseDirectory = templateBaseDirectory
	}
}

func Data(data gin.H) BlockOption {
	return func(block *Block) {
		block.Data = data
	}
}
func Name(value string) BlockOption {
	return func(block *Block) {
		block.Name = value
	}
}
func ParentBlockCache(value *BlockCache) BlockOption {
	return func(block *Block) {
		block.ParentBlockCache = value
	}
}
func Ctx(value context.Context) BlockOption {
	return func(block *Block) {
		block.Ctx = value
	}
}
