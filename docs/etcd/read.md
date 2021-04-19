```
查看版本
etcdctl version


//查看所有的key
etcdctl get --prefix ""

//添加值 foo=bar
etcdctl put foo bar

//创建租约
etcdctl put foo1 bar1 --lease=1234abcd（创建的租约ID）

查看数据
etcdctl get foo --hex
etcdctl get --from-key b
etcdctl get --prefix --limit=2 foo

// 删除数据
etcdctl  del foo1 foo9      --范围删除
etcdctl del --prev-kv zoo   --根据前缀删除数据
etcdctl del --from-key b




作者：半兽人
链接：https://www.orchome.com/620
来源：OrcHome
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

# 列表
etcdctl ls /kube-centos/network/config

# 查看
etcdctl get /kube-centos/network/config

# v2移除
etcdctl rm /kube-centos/network/config

# v3移除
ETCDCTL_API=3 etcdctl del /kube-centos/network/config

# 递归移除
etcdctl rm --recursive registry

# 修改
etcdctl mk /kube-centos/network/config "{ \"Network\": \"172.30.0.0/16\", \"Backend\": { \"Type\": \"vxlan\" } }"

# 命令将数据存到指定位置。这部分数据可以用来灾难恢复
etcdctl backup

# 健康检查
etcdctl endpoint health



作者：半兽人
链接：https://www.orchome.com/620
来源：OrcHome
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

存储:
    curl http://127.0.0.1:4001/v2/keys/testkey -XPUT -d value='testvalue'
    curl -s http://127.0.0.1:4001/v2/keys/message2 -XPUT -d value='hello etcd' -d ttl=5

获取:
    curl http://127.0.0.1:4001/v2/keys/testkey

查看版本:
    curl  http://127.0.0.1:4001/version

删除:
    curl -s http://127.0.0.1:4001/v2/keys/testkey -XDELETE

监视:
    窗口1：curl -s http://127.0.0.1:4001/v2/keys/message2 -XPUT -d value='hello etcd 1'
          curl -s http://127.0.0.1:4001/v2/keys/message2?wait=true
    窗口2：
          curl -s http://127.0.0.1:4001/v2/keys/message2 -XPUT -d value='hello etcd 2'

自动创建key:
    curl -s http://127.0.0.1:4001/v2/keys/message3 -XPOST -d value='hello etcd 1'
    curl -s 'http://127.0.0.1:4001/v2/keys/message3?recursive=true&sorted=true'

创建目录：
    curl -s http://127.0.0.1:4001/v2/keys/message8 -XPUT -d dir=true

删除目录：
    curl -s 'http://127.0.0.1:4001/v2/keys/message7?dir=true' -XDELETE
    curl -s 'http://127.0.0.1:4001/v2/keys/message7?recursive=true' -XDELETE

查看所有key:
    curl -s http://127.0.0.1:4001/v2/keys/?recursive=true

存储数据：
    curl -s http://127.0.0.1:4001/v2/keys/file -XPUT --data-urlencode value@upfile


使用etcdctl客户端：

存储:
    etcdctl set /liuyiling/testkey "610" --ttl '100'
                                         --swap-with-value value

获取：
    etcdctl get /liuyiling/testkey

更新：
    etcdctl update /liuyiling/testkey "world" --ttl '100'

删除：
    etcdctl rm /liuyiling/testkey

使用ca获取：
etcdctl --cert-file=/etc/etcd/ssl/etcd.pem   --key-file=/etc/etcd/ssl/etcd-key.pem  --ca-file=/etc/etcd/ssl/ca.pem get /message

目录管理：

    etcdctl mk /liuyiling/testkey "hello"    类似set,但是如果key已经存在，报错

    etcdctl mkdir /liuyiling 

    etcdctl setdir /liuyiling  

    etcdctl updatedir /liuyiling      

    etcdctl rmdir /liuyiling    

查看：
    etcdctl ls --recursive

监视：
    etcdctl watch mykey  --forever         +    etcdctl update mykey "hehe"

    #监视目录下所有节点的改变

    etcdctl exec-watch --recursive /foo -- sh -c "echo hi"

    etcdctl exec-watch mykey -- sh -c 'ls -al'    +    etcdctl update mykey "hehe"
    
    #列出集群内的成员以及他们当前的角色是不是leader
    etcdctl member list
```
```
$ etcdctl put flag 1
OK
$ etcdctl txn -i
compares:
value("flag") = "1"
success requests (get, put, delete):
put result true
failure requests (get, put, delete):
put result false
SUCCESS
OK
$ etcdctl get result
result
true


```

####解释：

1、etcdctl put flag 1设置flag为1

2、etcdctl txn -i开启事务（-i表示交互模式）

3、第2步输入命令后回车，终端显示出compares：

4、输入value("flag") = "1"，此命令是比较flag的值与1是否相等

5、第4步完成后输入回车，终端会换行显示，此时可以继续输入判断条件（前面说过事务由条件列表组成），再次输入回车表示判断条件输入完毕

6、第5步连续输入两个回车后，终端显示出success requests (get, put, delete):，表示下面输入判断条件为真时要执行的命令

7、与输入判断条件相同，连续两个回车表示成功时的执行列表输入完成

8、终端显示failure requests (get, put, delete):后输入条件判断失败时的执行列表

9、为了看起来简洁，此实例中条件列表和执行列表只写了一行命令，实际可以输入多行

10、总结上面的事务，要做的事情就是flag为1时设置result为true，否则设置result为false

11、事务执行完成后查看result值为true

####watch监听：

watch后etcdctl阻塞，在另一个终端中执行etcdctl put flag 2后，watch会打印出相关信息
```bigquery
$ etcdctl watch flag
PUT
flag
2

```

####租约
```bigquery
$ etcdctl lease grant 100
lease 38015a3c00490513 granted with TTL(100s)

$ etcdctl put k1 v1 --lease=38015a3c00490513
OK

$ etcdctl lease timetolive 38015a3c00490513
lease 38015a3c00490513 granted with TTL(100s), remaining(67s)

$ etcdctl lease timetolive 38015a3c00490513
lease 38015a3c00490513 granted with TTL(100s), remaining(64s)

$ etcdctl lease timetolive 38015a3c00490513 --keys
lease 38015a3c00490513 granted with TTL(100s), remaining(59s), attached keys([k1])

$ etcdctl put k2 v2 --lease=38015a3c00490513
OK

$ etcdctl lease timetolive 38015a3c00490513 --keys
lease 38015a3c00490513 granted with TTL(100s), remaining(46s), attached keys([k1 k2])

$ etcdctl lease revoke 38015a3c00490513 
lease 38015a3c00490513 revoked

$ etcdctl get k1
$ etcdctl get k2
$ 
$ etcdctl lease grant 10
lease 38015a3c0049051d granted with TTL(10s)

$ etcdctl lease keep-alive 38015a3c0049051d
lease 38015a3c0049051d keepalived with TTL(10)
lease 38015a3c0049051d keepalived with TTL(10)
lease 38015a3c0049051d keepalived with TTL(10)

```