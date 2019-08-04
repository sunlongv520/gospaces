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
		putResp *clientv3.PutResponse

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
	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)
	putResp, err = kv.Put(context.TODO(),"/cron/jobs/job2","byehhhhhh",clientv3.WithPrevKV())
	if err != nil{
		fmt.Println(err)
	}else{
		//获取版本信息
		fmt.Println("Revision:",putResp.Header.Revision)
		if putResp.PrevKv != nil{
			fmt.Println("key:",string(putResp.PrevKv.Key))
			fmt.Println("Value:",string(putResp.PrevKv.Value))
			fmt.Println("Version:",string(putResp.PrevKv.Version))
		}
	}
}