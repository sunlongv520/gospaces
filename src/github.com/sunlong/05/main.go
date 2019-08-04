package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"context"
)

func main(){
	var (
		config clientv3.Config
		err error
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
		leaseid clientv3.LeaseID
		leaseGrantResp *clientv3.LeaseGrantResponse
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		keepresp *clientv3.LeaseKeepAliveResponse
		keepRestChan <-chan *clientv3.LeaseKeepAliveResponse

	)
	//配置
	config = clientv3.Config{
		Endpoints:[]string{"192.168.1.188:2379"},
		DialTimeout:time.Second*5,
	}
	//连接 床见一个客户端
	if client,err = clientv3.New(config);err != nil{
		fmt.Println(err)
		return
	}




	//申请一个lease 租约
	lease = clientv3.NewLease(client)

	//申请一个10秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(),10);err != nil{
		fmt.Println(err)
		return
	}

	//拿到租约id
	leaseid = leaseGrantResp.ID

	//获得kv api子集
	kv = clientv3.NewKV(client)


	//自动续租
	if keepRestChan,err = lease.KeepAlive(context.TODO(),leaseid);err != nil{
		fmt.Println(err)
		return
	}
	//处理续租应答的协程
	go func() {
		for {
			select {
				case keepresp = <-keepRestChan:
					if keepRestChan == nil{
						fmt.Println("租约已失效了")
						goto END
					}else{//每秒会续租一次，所以就会收到一次应答
						fmt.Println("收到自动续租的应答")
					}
			}
		}
		END:
	}()





	//put一个kv 让它与租约关联起来 从而实现10秒自动过期
	if putResp,err = kv.Put(context.TODO(),"cron/lock/job1","v5",clientv3.WithLease(leaseid));err != nil{
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功",putResp.Header.Revision)

	//定时的看一下key过期了没有
	for{
		if getResp,err = kv.Get(context.TODO(),"cron/lock/job1");err != nil{
			fmt.Println(err)
			return
		}
		if getResp.Count == 0{
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期：",getResp.Kvs)
		time.Sleep(time.Second*2)
	}
}