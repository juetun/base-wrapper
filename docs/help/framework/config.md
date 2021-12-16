#### 一、数据库连接配置
```yaml
app:
  alias: ""                                           #如果要访问当前服务器的别名 如IP地址，或者反向代理的名称别名
  system_name: "xxx"                                  #系统名称
  name: "base-wrapper"                                #应用名称 (URL路径 微服务注册与发现使用)
  port: 8193                                          #应用启动端口
  version: "2.0"                                      #应用版本
  grace_reload: 1                                       #是否支持优雅重启  
  app_need_p_prof: false                                #是否开启性能分析
  administrator: "此参数会出现在每个接口的response中"        #代码、接口负责人信息
 

```

#### 一、数据库连接配置

```yaml
default: # 数据库配置名称（使用数据库操作时初始化和获取数据库操作连接会使用）
  maxidleconns: 1   # 连接池空闲连接数
  maxopenconns: 2   # 连接池最大连接数
  addr: "账号:密码@tcp(127.0.0.1:3306)/数据库名?charset=utf8mb4&parseTime=True&loc=Local" #连接配置
juetun_user0:
  maxidleconns: 1
  maxopenconns: 2
  addr: "账号:密码@tcp(127.0.0.1:3306)/数据库名?charset=utf8mb4&parseTime=True&loc=Local"

```

#### 二、redis连接配置

```yaml
default: # redis配置名称（使用redis操作时初始化和获取redis操作连接会使用）
  namespace: "default" 
  addr: "127.0.0.1:6379"                  # 连接地址
  db: 0                                   # 访问的redisDB
  password: "123456"                      # 密码
  poolsize: 1                             # 连接池信息
  maxidleconns: 1                         # 连接池信息
default1:
  namespace: "default"
  addr: "127.0.0.1:6379"
  db: 0
  password: "123456"
  poolsize: 1
  maxidleconns: 1
```

#### 三、ES连接配置
```yaml

default:
  addresses:
    - "http://localhost:9200"
  username: ""
  password: ""
  cloudid: ""
  apikey: ""
  retryonstatus:
    - 502
    - 503
    - 504
  disableretry: false
  enableretryontimeout : false
  maxretries: 3
default1:
  addresses:
    - "http://localhost:9200"
  username: ""
  password: ""
  cloudid: ""
  apikey: ""
  retryonstatus:
    - 502
    - 503
    - 504
  disableretry: false
  enableretryontimeout : false
  maxretries: 3
```

#### 三、日志配置
```yaml
outputs:                    # 日志输出的位置
  - "stdout"                # 日志输出到控制台
  - "file"                   # 日志输出到本地文件
format: "json"              # 日志格式采用 json
logcollectlevel: 4          # 日志采集级别  
logfilepath: "/tmp"          # 日志文件所在目录 
logfilename: ""              # 日志文件名称前缀
logiscut: false             # 日志文件是否切割

```


#### 三、服务注册发现配置（待实现完整逻辑）
```yaml

endpoints:
  #  etcd链接地址
  - http://localhost:2379
etcdendpoints:
  - api-domain-common
  - api-domain-websecure
dir: traefik
lockkey: etcdserver
host: api.test.com
```
### 四、阿里云OSS使用配置
```yaml
default:
  endpoint: "https://testupload.xxx.com"
  endsrcpoint: "https://testuploadsrc.xxx.com"
  accesskeyid: "xxx"
  accesskeysecret: "xxx"
  bucketname: "test-md-fileupload"
  bucketurl: "oss-cn-beijing.aliyuncs.com"
  sessionname: "xxx"
  dirname: "xx"
  expiredtime: 7200
  #  cdn有效时长
  cdnexpiredtime: 300
  # cdn 验证key
  cdnauthkey: "xxx"
  rolearn: "acs:ram::xx75022225676xx:role/aliyunosstokengenerxxx"
  videobucketlocation: "oss-cn-beijing"
video:
  endpoint: "https://testupload.xxx.com"
  endsrcpoint: "https://testuploadsrc.xxx.com"
  accesskeyid: "xxx"
  accesskeysecret: "xxx"
  bucketname: "xxx"
  bucketurl: "oss-cn-beijing.aliyuncs.com"
  sessionname: "xxx"
  dirname: ""
  expiredtime: 7200
  #  cdn有效时长
  cdnexpiredtime: 300
  # cdn 验证key
  cdnauthkey: "xxx"
  rolearn: "acs:ram::xxx:role/xxx"
  pipelineidvideo: "xxx"
  videobucketlocation: "oss-cn-beijing"
  parsecodetemp:
    #高清
    - templateid: "xxx"
      extname: "mp4"
      label: "高清"
      key: "HD"
    #标清
    - templateid: "xxx"
      extname: "mp4"
      label: "标清"
      key: "SD"
    #普清
    - templateid: "xxx"
      extname: "mp4"
      label: "普清"
      key: "LD"

```


### 七、nginx相关配置（此功能配置需要特殊支持、普通场景请忽略）

本配置需要nginx高于1.11.1版本
```conf
server {
        listen 80 default;
        server_name juetun.com 47.99.219.36;
        return 301 http://www.test.com$request_uri;
}
server {

	listen       80;
	server_name  www.test.com;
	root /var/www/html/;
	if ($host = 'test.com') {
		rewrite	^/(.*)$	http://www.test.com/$1	permanent;
	}

	# nginx 1.11.1版本以后有效
	set $trace_id "${request_id}";
	if ($http_x_atrace_id != "" ){
		set $trace_id "${http_x_atrace_id}";
	}
	add_header trace_id $trace_id;
	proxy_set_header X-Request-Id $trace_id;
	location / {
		proxy_pass			http://127.0.0.1:8092/;
		proxy_redirect			off;
		proxy_set_header		Host	$http_host;
		proxy_set_header		X-Real-IP	$remote_addr;
		proxy_set_header		X-Forwarded-For	$proxy_add_x_forwarded_for;
		proxy_set_header		Cookie	$http_cookie;
		chunked_transfer_encoding	off;
	}
	gzip on;
       	gzip_buffers 32 4K;
       	gzip_comp_level 6;
       	gzip_min_length 100;
       	gzip_types application/javascript text/css text/xml;
       	gzip_disable "MSIE [1-6]\."; #配置禁用gzip条件，支持正则。此处表示ie6及以下不启用gzip（因为ie低版本不支持）
       	gzip_vary on;
}



```
