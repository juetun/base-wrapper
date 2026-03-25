package micro_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 全局变量
var (
//etcdEndpoints = []string{"127.0.0.1:2379"} // etcd 地址
//serviceName   = "gin-demo-service"         // 服务名称
//serviceIP     = "127.0.0.1"                // 服务IP
//servicePort   = 8080                       // 服务端口
//leaseID       clientv3.LeaseID             // etcd 租约ID
//etcdClient    *clientv3.Client             // etcd 客户端
)

type (
	ETCDRegisterAndUnRegister struct {
		ctx        context.Context  `json:"-"`
		LeaseID    clientv3.LeaseID `json:"lease_id"`    // etcd 租约ID
		EtcdClient *clientv3.Client `json:"etcd_client"` // etcd 客户端
	}

	// MicroServiceInfo 服务注册信息结构体
	MicroServiceInfo struct {
		ID      string `json:"id"`      // 服务实例唯一ID
		Name    string `json:"name"`    // 服务名称
		Address string `json:"address"` // 服务地址 (ip:port)
		Port    int    `json:"port"`    // 服务端口
	}
)

func (r *ETCDRegisterAndUnRegister) orgServerInfo() (serviceId, serviceName string, serviceInfoBytes []byte, err error) {
	// 3. 准备服务注册信息
	var ip string
	if ip, err = utils.GetLocalIP(); err != nil {
		return
	}
	serviceId = uuid.New().String()
	serviceInfo := MicroServiceInfo{
		ID:      serviceId, // 生成唯一实例ID
		Name:    app_obj.App.AppName,
		Address: fmt.Sprintf("%s:%d", ip, app_obj.App.AppPort),
		Port:    app_obj.App.AppPort,
	}

	serviceName = fmt.Sprintf("%v(%v)", serviceInfo.Name, serviceInfo.Address)

	if serviceInfoBytes, err = json.Marshal(serviceInfo); err != nil {
		err = fmt.Errorf("序列化服务信息失败: %v", err)
		return
	}
	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("当前serviceInfo信息: %v\n", string(serviceInfoBytes))
	return
}

// 注册服务到etcd
func (r *ETCDRegisterAndUnRegister) RegisterMicro(c *gin.Engine, cTxs ...context.Context) (ok bool, err error) {
	// 4. 序列化服务信息
	var serviceInfoBytes []byte
	var serviceId, serviceName string
	var eTCDConfig = clientv3.Config{
		Endpoints:   app_obj.RegistryServiceConfig.Consul.Endpoints, // etcd 地址
		DialTimeout: 5 * time.Second,
	}

	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("开始注册服务到 ETCD....")
	serviceInfoBytes, _ = json.Marshal(eTCDConfig)
	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("注册服务内容: %v \n", string(serviceInfoBytes))

	// 1. 创建etcd客户端
	if r.EtcdClient, err = clientv3.New(eTCDConfig); err != nil {
		base.Io.SetInfoType(base.LogLevelError).SystemOutFatalf("创建etcd客户端失败: %v", err)
		return
	}

	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("创建注册服务连接完成 \n")

	//组织服务信息
	if serviceId, serviceName, serviceInfoBytes, err = r.orgServerInfo(); err != nil {
		return
	}

	// 2. 创建租约（10秒过期，需要心跳保活）
	var leaseResp *clientv3.LeaseGrantResponse
	if leaseResp, err = r.EtcdClient.Grant(r.ctx, 10); err != nil {
		base.Io.SetInfoType(base.LogLevelError).SystemOutFatalf("创建租约失败: %v", err)
		return
	}

	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("注册完的LeaseID: %v\n", leaseResp.ID)
	r.LeaseID = leaseResp.ID

	// 5. 注册服务（将服务信息写入etcd，绑定租约）
	eTCDKey := fmt.Sprintf("/services/%s/%s", app_obj.App.AppName, serviceId)
	if _, err = r.EtcdClient.Put(r.ctx, eTCDKey, string(serviceInfoBytes), clientv3.WithLease(r.LeaseID)); err != nil {
		base.Io.SetInfoType(base.LogLevelError).SystemOutFatalf("注册服务失败: %v", err)
		return
	}

	// 6. 启动心跳保活
	var ch <-chan *clientv3.LeaseKeepAliveResponse
	if ch, err = r.EtcdClient.KeepAlive(r.ctx, r.LeaseID); err != nil {
		base.Io.SetInfoType(base.LogLevelError).SystemOutFatalf("启动心跳保活失败: %v", err)
		return
	}

	// 监听心跳响应
	go func() {
		for range ch {
			base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("心跳保活成功 \n")
		}
		base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("心跳保活通道关闭\n")
	}()
	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("服务注册成功: %v -> %v\n", eTCDKey, serviceName)
	return
}

// 从etcd注销服务
func (r *ETCDRegisterAndUnRegister) UnRegisterMicro() {
	if r.EtcdClient == nil {
		return
	}
	var err error
	// 撤销租约（自动删除注册信息）
	if _, err = r.EtcdClient.Revoke(r.ctx, r.LeaseID); err != nil {
		base.Io.SetInfoType(base.LogLevelError).SystemOutFatalf("注销服务失败: %v \n", err)
	} else {
		base.Io.SetInfoType(base.LogLevelError).SystemOutFatalf("服务注销成功 \n")
	}

	// 关闭etcd客户端
	err = r.EtcdClient.Close()
	return
}

func NewETCDRegisterAndUnRegister() (r app_start.MicroOperateInterface) {
	res := &ETCDRegisterAndUnRegister{}
	res.ctx = context.Background()
	return res
}

//func main() {
//	// 1. 注册服务到etcd
//	if err := registerService(); err != nil {
//		fmt.Printf("服务注册失败: %v\n", err)
//		return
//	}
//
//	// 2. 创建Gin引擎
//	r := gin.Default()
//
//	// 3. 定义路由
//	r.GET("/hello", func(c *gin.Context) {
//		c.JSON(http.StatusOK, gin.H{
//			"message":  "Hello, Gin + etcd + Traefik!",
//			"service":  serviceName,
//			"instance": fmt.Sprintf("%s:%d", serviceIP, servicePort),
//		})
//	})
//
//	// 4. 启动HTTP服务（非阻塞）
//	go func() {
//		addr := fmt.Sprintf("%s:%d", serviceIP, servicePort)
//		if err := r.Run(addr); err != nil {
//			fmt.Printf("HTTP服务启动失败: %v\n", err)
//			unregisterService()
//		}
//	}()
//
//	fmt.Printf("Gin服务已启动: http://%s:%d\n", serviceIP, servicePort)
//
//	// 5. 监听退出信号（优雅关闭）
//	quit := make(chan os.Signal, 1)
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//	<-quit
//
//	// 6. 优雅关闭
//	fmt.Println("\n开始优雅关闭服务...")
//	unregisterService()
//	fmt.Println("服务已完全关闭")
//}
