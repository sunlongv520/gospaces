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

	kv = clientv3.NewKV(client)

	go func() {
		for{
			kv.Put(context.TODO(),"/cron/jobs/job7","i am job7")
			kv.Delete(context.TODO(),"/cron/jobs/job7")
			time.Sleep(time.Second*1)
		}
	}()

	//先get到当前的值  并监听后续变化
	if getResp,err = kv.Get(context.TODO(),"/cron/jobs/job7");err != nil{
		fmt.Println(err)
		return
	}

	//现在key是存在的
	if len(getResp.Kvs) != 0{
		fmt.Println("当前值：",string(getResp.Kvs[0].Value))
	}



}