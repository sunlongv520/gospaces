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
		getResp *clientv3.GetResponse
		watchStartRevision int64
		watcher clientv3.Watcher
		watchRespChan <-chan clientv3.WatchResponse
		watchResp clientv3.WatchResponse
		event *clientv3.Event

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

	// 当前etcd集群事务ID, 单调递增的
	watchStartRevision = getResp.Header.Revision + 1

	// 创建一个watcher
	watcher = clientv3.NewWatcher(client)

	// 启动监听
	fmt.Println("从该版本向后监听:", watchStartRevision)

	watchRespChan = watcher.Watch(context.TODO(),"/cron/jobs/job7",clientv3.WithRev(watchStartRevision))

	// 处理kv变化事件
	for watchResp = range watchRespChan{
		for _,event = range watchResp.Events{
			switch event.Type {
				case mvccpb.PUT:
					fmt.Println("修改为:",event.Kv)
				case mvccpb.DELETE:
					fmt.Println("删除了",event.Kv)
			}
		}
	}

}