# base-wrapper

### 关于GO MOD使用
go mod 执行步骤
####第一步修改环境变量

//查看环境变量
```cassandraql
go env

GO111MODULE="auto"
GOARCH="amd64"
GOBIN="/Users/xxx/go/bin"
GOCACHE="/Users/xxx/Library/Caches/go-build"
GOENV="/Users/xxx/Library/Application Support/go/env"
GOEXE=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="darwin"
GOINSECURE=""
GOMODCACHE="/Users/xxx/go/pkg/mod"
GONOPROXY="gitlab.xxx.com/md-go/library"
GONOSUMDB="gitlab.xxx.com/md-go/library"
GOOS="darwin"
GOPATH="/Users/xxx/go"
GOPRIVATE="gitlab.xxx.com/md-go/library"
GOPROXY="https://mirrors.aliyun.com/goproxy/,direct"
GOROOT="/usr/local/go"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/usr/local/go/pkg/tool/darwin_amd64"
GCCGO="gccgo"
AR="ar"
CC="clang"
CXX="clang++"
CGO_ENABLED="1"
GOMOD="/Users/xxx/go/src/example/go.mod"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -pthread -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -fdebug-prefix-map=/var/folders/jq/_r0ghj8d099dyshx7v8284p80000gn/T/go-build133838908=/tmp/go-build -gno-record-gcc-switches -fno-common"
```
####配置代理
GO111MODULE  GOPROXY

| 命令 | 描述 | 功能 |
| :----: | :----: | :----: |
| download | download modules to local cache | 下载依赖的module到本地cache（gopath\pkg\mod\cache） |
| edit |  edit go.mod from tools or scripts | 编辑go.mod文件 |
| graph |  print module requirement graph | 打印模块依赖图 |
| init |  initialize new module in current directory | 在当前文件夹下初始化一个module, 生成go.mod文件 |
| tidy |  add missing and remove unused modules | 增加丢失的module，去掉未用的module |
| vendor |  make vendored copy of dependencies | 将依赖复制到当前模块vendor下 |
| verify |  verify dependencies have expected content | 校验依赖 |
| why |  explain why packages or modules are needed | 解释为什么需要依赖 |


####1、配置 GOPRIVATE 
```cassandraql
// go env 查看GOPRIVATE 的值
export GOPRIVATE="gitlab.xxx.com/md-go/library"
```

####2、拉取私有依赖包（二方包  insecure参数是切换 http包）

```cassandraql
go get --insecure gitlab.xxx.com/md-go/library
```

####3、拉取三方包
```cassandraql
go get ./...
```

####备注：
拉取指定版本命令
```cassandraql
//d3f30cfc81f9109850fa9043a19783cbbde68a5 git提交的SHA-1 散列值
go get github.com/gin-gonic/gin@5d3f30cfc81f9109850fa9043a19783cbbde68a5
```


###swagger:使用

#####刷新swagger配置的接口信息
```cassandraql
swag init

//使用swagger必须先引入 /docs包
//访问地址/swagger/index.html
```
