# 框架目录结构规划

```shell 及说明

├── README.md  																			帮助文档
├── config    
│   ├── app.yaml																		应用配置通用文件
│   ├── dev																					开发环境(dev)配置文件存放目录
│   │   ├── appmap.yml															各服务间访问host配置
│   │   ├── database.yaml														数据库配置
│   │   ├── elasticsearch.yml												es配置文件
│   │   ├── log.yaml																日志配置文件
│   │   └── redis.yaml															redis配置文件
│   ├── release																			生产环境配置						
│   │   ├── database.yaml
│   │   ├── log.yaml
│   │   └── redis.yaml
│   ├── sys_conf																		系统配置如特殊的nginx配置（CI/CD集成环境读取）
│   │   └── nginx
│   │       ├── console.conf
│   │       └── www.conf
│   └── test																				测试环境配置文件
│       ├── database.yaml
│       ├── log.yaml
│       └── redis.yaml
├── docs																						项目文档、swagger生成的目录存放路径
│   └── default.go
├── go.mod																					go mod 相关文件
├── go.sum																					go mod 相关文件
├── log																							日志文件生成目录
├── main.go																					项目启动main.go（main函数所在位置）
├── pkg																							项目调用封装的工具目录（存放常量、一方包的位置）
│   ├── cron																				定时任务实现入口位置
│   │   └── default.go
│   ├── parameter																		常量存储位置
│   │   ├── cache_key.go
│   │   └── chat_arg.go
│   └── redis_mq																		
│       └── mq.go
└── web																							系统应用目录
    ├── cons																				系统controller层结构
    │   ├── admins																	客服后台controller
    │   │   ├── admin_impl													controller具体实现代码目录
    │   │   │   └── default.go											controller具体实现文件
    │   │   └── default.go													controller代码(约定目录)interface目录
    │   ├── intranets																内网访问接口实现目录(rpc接口目录)						
    │   │   ├── default.go
    │   │   └── intranet_impl
    │   │       └── default.go
    │   ├── outernets																外网访问接口目录
    │   │   ├── default.go
    │   │   └── outernet_impl
    │   │       └── default.go
    │   └── pages																		外网访问网页界面实现目录
    │       ├── default.go
    │       └── page_impl
    │           └── default.go
    ├── daos																				数据(数据库、ES、redis等)操作目录
    │   ├── chat.go
    │   ├── dao_impl
    │   │   ├── chat.go
    │   │   └── room.go
    │   └── room.go
    ├── models																			数据库表数据结构关系映射
    │   ├── chat.go
    │   ├── model.go
    │   └── room.go
    ├── router																			路由注册目录
    │   ├── admin																		客服后台路由注册目录
    │   │   └── message.go
    │   ├── intranet																内网访问路由注册目录
    │   │   └── default.go
    │   ├── outernet																外网访问路由注册目录
    │   │   └── default.go
    │   ├── page																		网页访问路由注册目录
    │   │   └── default.go
    │   └── router.go																路由注册公共入口
    ├── srvs																				业务逻辑service层目录
    │   ├── chart.go
    │   ├── chat_act.go
    │   ├── default.go
    │   └── srv_impl
    │       ├── chart.go
    │       ├── chat_act.go
    │       └── default.go
    ├── validate																	其他的参数校验目录（非必须）
    └── wrappers																	常用的参数结构体封装目录，相当于JAVA的(pojo目录)
        ├── chart.go
        ├── wrapper_admin													客服后台所需参数			
        ├── wrapper_intranet											内网所需参数
        ├── wrapper_outernet
        └── wrapper_wrappe


```
