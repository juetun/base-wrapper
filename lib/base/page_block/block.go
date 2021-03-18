// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package page_block

// 页面缓存对象

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

type Handler func(block *Block) (err error)
type HandlerBlockCache func(block *BlockCache) (err error)

// 页面操作对象 一个html由BLOCK拼凑而成 本结构体设计目的为实现页面局部数据缓存控制
// 每个BLOCK具备独立的缓存对象，独立从数据库、redis、或其他数据源获取数据的能力获取数据的能力
// 缓存规则：
// 1、如果父block的缓存时间大于0，则子BLOCK设置了缓存时间无效，
// 2、如果父BLOCK的缓存时间为0，则子BLOCK的缓存时间有效
type Block struct {
	//Ctx                   context.Context `json:"ctx"`                     // 上下文的操作对象 ，此处主要用来传递上下文参数
	ParentBlockCache      *BlockCache `json:"parent_block_cache"`      // 当前Block的父Block
	Name                  string      `json:"name"`                    // 当前Block的系统唯一名字
	Data                  gin.H       `json:"data"`                    // 当前Block的参数
	TempFile              string      `json:"temp_file"`               // html文件地址
	TemplateBaseDirectory string      `json:"template_base_directory"` // html模板文件所在的公共基础路径
	BlockCache            *BlockCache `json:"cache"`                   // 当前模块缓存的基本参数
	RunChildBefore        []Handler   `json:"-"`                       // 子BLOCK运行之前的动作
	RunBefore             []Handler   `json:"-"`                       // 渲染完数据后执行此方法，主要用来调试数据使用 //渲染完数据后执行此方法，主要用来调试数据使用,返回值为true时跳出
	RunAfter              []Handler   `json:"-"`                       // 渲染完数据前执行此方法，主要用来调试数据使用 //渲染完数据前执行此方法，主要用来调试数据使用,返回值为true时跳出
	ChildBock             []*Block    `json:"child_bock"`              // 当前的子BLOCK
	RefreshForceCache     bool        `json:"refresh_force_cache"`     //是否强制刷新缓存
}

// 判断文件目录是否存在
func (r *Block) Exists(path string) bool {
	_, err := os.Stat(path) // os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 模板文件路径
func (r *Block) tempFilePath() {
	r.TempFile = r.TemplateBaseDirectory + r.TempFile
}

// 将HTML模板文件绑定参数
func (r *Block) ParseHtml() (res string, err error) {
	var tmp *template.Template
	r.tempFilePath()

	if !r.Exists(r.TempFile) {
		if err = fmt.Errorf("the template file(%s) is not exists",
			r.TempFile); err != nil {
			return
		}
		return
	}

	buf := new(bytes.Buffer)

	// 拼接TemplateFile path
	if tmp, err = template.New(r.Name).
		Funcs(app_obj.FuncMap).
		ParseFiles(r.TempFile); err != nil {
		return
	} else {
		tmp.Execute(buf, r.Data)
		//tmp.ExecuteTemplate(buf, r.Name, r.Data)
	}

	res = buf.String()
	return
}

// 传递上下文参数
func (r *Block) setChildContext(item *Block) {

	data := gin.H{} // 合并页面数据
	for key, value := range r.Data {
		data[key] = value
	}
	for key, value := range item.Data {
		data[key] = value
	}

	item.Data = data
	//item.Ctx = r.Ctx
	item.ParentBlockCache = r.BlockCache
}

// 从缓存中获取数据
func (r *Block) getCache() (res string, err error) {
	if r.BlockCache.Cache == nil {
		return
	}
	res, err = r.BlockCache.Cache.Get(r.getCacheKey())
	return
}

// 将数据写入缓存
func (r *Block) writeToCache(data string) (err error) {

	if r.BlockCache.ExpireTime.IsZero() || r.BlockCache.ExpireTime.Unix() < time.Now().Unix() {
		return
	}

	// 缓存时间
	if lt := r.BlockCache.ExpireTime.Unix() - time.Now().Unix(); lt > 0 {
		lTime := time.Duration(r.BlockCache.ExpireTime.Unix() - time.Now().Unix())
		r.BlockCache.Cache.Set(r.getCacheKey(), data, lTime*time.Second)
	}

	return
}

// 获取缓存数据的Key
func (r *Block) getCacheKey() (res string) {

	if r.BlockCache.CacheKey != "" {
		res = r.BlockCache.CacheKey
		return
	}

	uniqueKey := fmt.Sprintf("p:%s:%s", app_obj.App.AppName, r.Name)

	switch r.BlockCache.CacheType {

	case CacheRedis: // 如果缓存类型为redis
		res = uniqueKey
	case CacheFile, CacheDatabase:
		res = base64.StdEncoding.EncodeToString([]byte(uniqueKey))
	default:
		res = base64.StdEncoding.EncodeToString([]byte(uniqueKey))
	}

	return
}

func (r *Block) before() (err error) {
	// 如果配置了运行后执行
	for _, runBefore := range r.RunBefore {
		if err = runBefore(r); err != nil {
			return
		}
	}
	return
}
func (r *Block) after() (err error) {
	// 如果配置了运行后执行
	for _, runAfter := range r.RunAfter {
		if err = runAfter(r); err != nil {
			return
		}
	}
	return
}

func (r *Block) haveCacheDo() (res template.HTML, err error) {

	//如果不是强制刷新缓存
	if !r.RefreshForceCache {
		// 从缓存中拿数据
		if r.BlockCache.CacheData, err = r.getCache(); err != nil {
			if res != "" {
				err = r.after()
				return
			}
		}
	}

	for _, item := range r.ChildBock {
		r.setChildContext(item) // 传递上下文参数
		if r.Data[item.Name], err = item.Run(); err != nil {
			return
		}
	}

	// 解析HTML模板代码
	if r.BlockCache.CacheData, err = r.ParseHtml(); err != nil {
		return
	}

	// 返回值赋值在after后的目的是，可以通过后边的注入修改缓存值
	if err = r.after(); err != nil {
		return
	}
	// 将数据写入缓存
	if err = r.writeToCache(r.BlockCache.CacheData); err != nil {
		return
	}

	return
}

func (r *Block) hasNotCacheDo() (res template.HTML, err error) {
	for _, item := range r.ChildBock {
		r.setChildContext(item) // 传递上下文参数
		if r.Data[item.Name], err = item.Run(); err != nil {
			return
		}
	}

	// 解析HTML模板代码
	if r.BlockCache.CacheData, err = r.ParseHtml(); err != nil {
		return
	}

	// 返回值赋值在after后的目的是，可以通过后边的注入修改缓存值
	if err = r.after(); err != nil {
		return
	}
	return

}

// 解析模板数据
func (r *Block) Run() (res template.HTML, err error) {

	defer func() {
		res = template.HTML(r.BlockCache.CacheData)
	}()

	// 初始化默认值
	if err = r.defaultOption(); err != nil {
		return
	}

	// 获取缓存数据或者解析Block之前的动作
	if err = r.before(); err != nil {
		return
	}

	// 如果没有缓存,获取过期时间已过期
	if r.BlockCache == nil || r.BlockCache.ExpireTime.Unix() < time.Now().Unix() {
		if res, err = r.hasNotCacheDo(); err != nil {
			return
		}
		return
	}

	if res, err = r.haveCacheDo(); err != nil {
		return
	}
	return
}

// 初始化页面缓存对象
func (r *Block) defaultCacheBlock() {

	if r.BlockCache == nil {
		r.BlockCache = NewBlockCache()
	}

	return
}

// BLOCK 默认数据逻辑处理
func (r *Block) defaultOption() (err error) {

	if r.Data == nil && len(r.Data) == 0 {
		r.Data = gin.H{}
	}

	// 初始化页面缓存对象
	r.defaultCacheBlock()

	// 如果名称没定义
	if r.Name == "" {
		err = fmt.Errorf("您没有定义当前BLOCK的name(%T)", r)
		return
	}

	// 默认初始化当前模板文件所在位置
	r.defaultTemplateBaseDirectory()

	// 初始化过期时间
	r.initExpireTime()
	return
}

// 缓存时间处理
func (r *Block) initExpireTime() {

	if r.ParentBlockCache == nil || r.BlockCache == nil {
		return
	}
	// 如果当前BLOCK的缓存时间为0 则不缓存
	// 如果当前的父BLOCK缓存为0,则指定使用当前缓存时间
	if r.ParentBlockCache.ExpireTime.Unix() == 0 || r.BlockCache.ExpireTime.Unix() == 0 {
		return
	}

	// 如果当前BLOCK的父block不等于0,则本次缓存就为不缓存，（设置的过期时间是昨天当前时间）
	if r.ParentBlockCache.ExpireTime.Unix()-time.Now().Unix() > 0 {
		r.BlockCache.ExpireTime = time.Now().AddDate(0, 0, -1)
	}

}

func NewBlock(option ...BlockOption) (block *Block) {

	block = &Block{}
	for _, handler := range option {
		handler(block)
	}
	return
}

// 默认初始化当前模板文件所在位置
func (r *Block) defaultTemplateBaseDirectory() {
	if r.TemplateBaseDirectory == "" {
		r.TemplateBaseDirectory = app_obj.App.AppTemplateDirectory
	}
	return
}

type BlockOption func(block *Block)

// 是否强制刷新缓存
func RefreshForceCache(refreshForceCache bool) BlockOption {
	return func(block *Block) {
		block.RefreshForceCache = refreshForceCache
	}
}

// 当前BLOCK的子Block
func ChildBock(childBock ...*Block) BlockOption {
	return func(block *Block) {
		block.ChildBock = childBock
	}
}

// BLOCK的Run方法运行主要逻辑后执行此方法
func RunAfter(runAfter ...Handler) BlockOption {
	return func(block *Block) {
		block.RunAfter = runAfter
	}
}

func RunBefore(runBefore ...Handler) BlockOption {
	return func(block *Block) {
		block.RunBefore = runBefore
	}
}
func RunChildBefore(runChildBefore ...Handler) BlockOption {
	return func(block *Block) {
		block.RunChildBefore = runChildBefore
	}
}

//
func CacheBlock(cacheBlock *BlockCache) BlockOption {
	return func(block *Block) {
		block.BlockCache = cacheBlock
	}
}

func CacheBlockOption(cacheBlock ...BlockCacheOption) BlockOption {
	return func(block *Block) {
		block.BlockCache = NewBlockCache(cacheBlock...)
	}
}

func TempFile(tempFile string) BlockOption {
	return func(block *Block) {
		block.TempFile = tempFile
	}
}

// html模板文件所在的基础路径
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

func ParentBlockCacheOption(value *BlockCache) BlockOption {
	return func(block *Block) {
		block.ParentBlockCache = &BlockCache{}
	}
}

func ParentBlockCache(value *BlockCache) BlockOption {
	return func(block *Block) {
		block.ParentBlockCache = value
	}
}

//func Ctx(value context.Context) BlockOption {
//	return func(block *Block) {
//		block.Ctx = value
//	}
//}
