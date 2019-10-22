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
		getResp *clientv3.GetResponse

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

	//读取前缀
	getResp,err = kv.Get(context.TODO(),"/cron/jobs/",clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(getResp.Kvs)
}