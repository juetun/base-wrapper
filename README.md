[TOC]  

# 框架说明

### 序言

​       互联网技术的发展，go语言作为新兴的编程语言在实际应用开发中有着自身的一些优势。随着在市面上应用范围广度增加，市面上的一些主流的开发框架虽然能够完成大部分的开发需求。但是在日志的整合度、系统的集成度和规范性方面有一定的缺失。
​       本代码库的功能:整合市面上常用的go语言开发框架(gin gorm logrus redis elasticsearch kafka clickhouse等)定制适合WEB微服务系统场景。

### 基础使用框架 

| 依赖框架及组件 | 版本号 | 备注 |
| -------------- | ------ | ---- |
| gin            |   v1.7.2      |      |
| gorm           | V2版本         | 源地址 gorm.io/gorm v1.21.11    |
| logrus         | v1.8.1        |      |
| file-rotatelogs | v2.4.0 |  日志文件切割 github.com/lestrrat-go/file-rotatelogs    |
| redis依赖包 | v8.10.0 | github.com/go-redis/redis/v8 |
| yaml文件管理工具 | v2.4.0 | gopkg.in/yaml.v2 |


...       

## 帮助文档
​        使用go mod管理依赖包(详见[go mod](lib/docs/help/other/go_mod.md)帮助文档)。
### 一、框架使用实例
#### 1.1、[框架启动](lib/docs/help/framework/init.md)
##### 1.1.1、 [配置文件](lib/docs/help/framework/config.md)
#### 1.2、[框架插件](lib/docs/help/framework/plugins.md)
#### 1.3、[gin中间件引入](lib/docs/help/framework/gin_middleware.md)
#### 1.4、[框架目录说明](lib/docs/help/framework/dir.md)
#### 1.5、[GORM使用](lib/docs/help/framework/gorm.md)
#### 1.6、[redis使用封装](lib/docs/help/framework/redis.md)
##### 1.6.1、[分布式锁](lib/docs/help/framework/redis/lock.md)
##### 1.6.1、[分布式订阅发布](lib/docs/help/framework/redis/lock.md)
#### 1.7、[WEBSOCKET使用](lib/docs/help/framework/gin_micro.md)
#### 1.9、[微服务注册](lib/docs/help/framework/gin_micro.md)
##### 1.9.1、[服务间的调用实例](lib/docs/help/framework/call_method.md)
##### 1.9.2、[traefik代码改造文档](lib/docs/help/framework/traefik.md)
### 二、其他帮助文档
#### 2.1、[go mod](lib/docs/help/other/go_mod.md)

#### 2.2、[swagger](lib/docs/help/other/swagger.md)




