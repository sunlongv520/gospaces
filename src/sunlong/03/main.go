package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

func main(){
	var (
		config clientv3.Config
		err error
		client *clientv3.Client
		kv clientv3.KV
		delResp *clientv3.DeleteResponse
		idx int
		kvpair *mvccpb.KeyValue

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

	delResp,err = kv.Delete(context.TODO(),"/cron/jobs")
	if err != nil{
		fmt.Println(err)
		return
	}else{
		if len(delResp.PrevKvs) > 0 {
			for idx,kvpair = range delResp.PrevKvs{
				idx = idx
				fmt.Println("删除了",string(kvpair.Key),string(kvpair.Value))
			}
		}
	}
}