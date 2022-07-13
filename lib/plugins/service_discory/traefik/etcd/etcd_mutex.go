// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
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
func (r *MutexEtcd) init() error {
	r.m.Lock()
	if r.Key == "" {
		r.Key = "platLock"
	}
	var err error
	var ctx context.Context
	client, err := clientv3.New(r.Conf)
	if err != nil {
		return err
	}
	r.txn = clientv3.NewKV(client).Txn(context.TODO())
	r.lease = clientv3.NewLease(client)
	leaseResp, err := r.lease.Grant(context.TODO(), r.Ttl)
	if err != nil {
		return err
	}
	ctx, r.cancel = context.WithCancel(context.TODO())
	r.leaseID = leaseResp.ID
	_, err = r.lease.KeepAlive(ctx, r.leaseID)
	return err
}

//获取锁:
func (r *MutexEtcd) Lock() (err error) {

	if err = r.init(); err != nil {
		return
	}
	//LOCK:
	r.txn.If(clientv3.Compare(clientv3.CreateRevision(r.Key), "=", 0)).
		Then(clientv3.OpPut(r.Key, "", clientv3.WithLease(r.leaseID))).
		Else()
	var txnResp *clientv3.TxnResponse
	if txnResp, err = r.txn.Commit(); err != nil {
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
func (r *MutexEtcd) UnLock() {
	r.cancel()
	r.lease.Revoke(context.TODO(), r.leaseID)
	r.m.Unlock()
	fmt.Println("释放了锁")
}

func (r *MutexEtcd) callExample() {
	var conf = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	rutex1 := &MutexEtcd{
		Conf: conf,
		Ttl:  10,
		Key:  "lock",
	}
	rutex2 := &MutexEtcd{
		Conf: conf,
		Ttl:  10,
		Key:  "lock",
	}
	//groutine1
	go func() {
		err := rutex1.Lock()
		if err != nil {
			fmt.Println("groutine1抢锁失败")
			fmt.Println(err)
			return
		}
		//可以做点其他事，比如访问和操作分布式资源
		fmt.Println("groutine1抢锁成功")
		time.Sleep(10 * time.Second)
		defer rutex1.UnLock()
	}()

	//groutine2
	go func() {
		err := rutex2.Lock()
		if err != nil {
			fmt.Println("groutine2抢锁失败")
			fmt.Println(err)
			return
		}
		//可以做点其他事，比如访问和操作分布式资源
		fmt.Println("groutine2抢锁成功")
		defer rutex2.UnLock()
	}()
	time.Sleep(30 * time.Second)
}
