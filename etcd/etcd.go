package etcd

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"go.etcd.io/etcd/clientv3"
)

var (
	client *clientv3.Client
)

type LogConf struct {
	Path string `json:"path"`
	Topic string `json:"topic"`
}

// Init 初始化etcd
func Init(address string, timeout time.Duration) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:            []string{address},
		DialTimeout:          timeout,
	})
	if err != nil {
		log.Printf("etcd connect failed, err: %v", err)
		return
	}
	return
}

// Get 获取etcd种的配置
func Get(etcK string, timeout time.Duration) (logconf []*LogConf, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	resp, err := client.Get(ctx, etcK)
	cancel()
	if err != nil {
		log.Printf("Get etcd data failed :%v", err)
		return
	}
	for _, ev := range resp.Kvs {
		err = json.Unmarshal(ev.Value, &logconf)
		if err != nil {
			log.Printf("unmarshal failed: %c", err)
			return
		}
	}
	return
}

func Watcher(etcK string, newConfCh chan<- []*LogConf) {
	resp := client.Watch(context.Background(), etcK)
	log.Println("watcher start..")
	for wresp := range resp {
		for _, ev := range wresp.Events {
			var newConf []*LogConf
			//如果不是删除操作 则发送新配置到newConfCh中，否则发送默认配置，即空
			if ev.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(ev.Kv.Value, &newConf)
				if err != nil {
					log.Printf("newconf unmalshal failed %v", err)
				}
			}
			newConfCh <- newConf
		}
	}
}