// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package etcd

import (
	"context"
	"fmt"
	"github.com/etcd-io/etcd/clientv3"
	"sync"
	"time"
)

type MutexEtcd struct {
	Ttl     int64              //租约时间
	Conf    clientv3.Config    //etcd集群配置
	Key     string             //etcd的key
	cancel  context.CancelFunc //关闭续租的func
	lease   clientv3.Lease
	leaseID clientv3.LeaseID
	txn     clientv3.Txn
	m       sync.Locker
}

//初始化锁
func (em *MutexEtcd) init() error {
	em.m.Lock()
	if em.Key == "" {
		em.Key = "platLock"
	}
	var err error
	var ctx context.Context
	client, err := clientv3.New(em.Conf)
	if err != nil {
		return err
	}
	em.txn = clientv3.NewKV(client).Txn(context.TODO())
	em.lease = clientv3.NewLease(client)
	leaseResp, err := em.lease.Grant(context.TODO(), em.Ttl)
	if err != nil {
		return err
	}
	ctx, em.cancel = context.WithCancel(context.TODO())
	em.leaseID = leaseResp.ID
	_, err = em.lease.KeepAlive(ctx, em.leaseID)
	return err
}

//获取锁:
func (em *MutexEtcd) Lock() (err error) {

	if err = em.init(); err != nil {
		return
	}
	//LOCK:
	em.txn.If(clientv3.Compare(clientv3.CreateRevision(em.Key), "=", 0)).
		Then(clientv3.OpPut(em.Key, "", clientv3.WithLease(em.leaseID))).
		Else()
	var txnResp *clientv3.TxnResponse
	if txnResp, err = em.txn.Commit(); err != nil {
		return
	}
	if !txnResp.Succeeded { //判断txn.if条件是否成立
		err = fmt.Errorf("抢锁失败")
		return
		//goto LOCK    //txn如果没有succeeded就不能重复提交，不然会panic,不知道怎么解决，求大佬告知
	}
	return
}

//释放锁：
func (em *MutexEtcd) UnLock() {
	em.cancel()
	em.lease.Revoke(context.TODO(), em.leaseID)
	em.m.Unlock()
	fmt.Println("释放了锁")
}

func (em *MutexEtcd) callExample() {
	var conf = clientv3.Config{
		Endpoints:   []string{"172.16.196.129:2380", "192.168.50.250:2380"},
		DialTimeout: 5 * time.Second,
	}
	eMutex1 := &MutexEtcd{
		Conf: conf,
		Ttl:  10,
		Key:  "lock",
	}
	eMutex2 := &MutexEtcd{
		Conf: conf,
		Ttl:  10,
		Key:  "lock",
	}
	//groutine1
	go func() {
		err := eMutex1.Lock()
		if err != nil {
			fmt.Println("groutine1抢锁失败")
			fmt.Println(err)
			return
		}
		//可以做点其他事，比如访问和操作分布式资源
		fmt.Println("groutine1抢锁成功")
		time.Sleep(10 * time.Second)
		defer eMutex1.UnLock()
	}()

	//groutine2
	go func() {
		err := eMutex2.Lock()
		if err != nil {
			fmt.Println("groutine2抢锁失败")
			fmt.Println(err)
			return
		}
		//可以做点其他事，比如访问和操作分布式资源
		fmt.Println("groutine2抢锁成功")
		defer eMutex2.UnLock()
	}()
	time.Sleep(30 * time.Second)
}
