package main

import (
	"LogAgent/common"
	"LogAgent/config"
	"LogAgent/etcd"
	"LogAgent/kafka"
	"LogAgent/taillog"
	"fmt"
	"log"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

var (
	cfg = new(config.Conf)
)

//logagent 入口程序
func main() {
	//  加载配置文件
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		log.Printf("load config failed err: %v", err)
		return
	}
	log.Println("配置文件加载成功")

	//	初始化kafka链接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.MaxChan)
	if err != nil {
		log.Printf("kafka init failed,err:%v\n", err)
		return
	}
	log.Println("连接kakfa成功")

	//  初始化etcd
	err = etcd.Init(cfg.Etcd.Address, time.Duration(cfg.Etcd.Timeout)*time.Millisecond)
	if err != nil {
		log.Printf("etcd init failed,err:%v\n", err)
		return
	}
	log.Println("etcd初始化成功")

	//测试连接
	// etcd.Put("/etcd", `[{"path":"/data/webroot/go/src/LogAgent/my.log","topic":"test01"},{"path":"/data/webroot/go/src/LogAgent/my.log","topic":"test02"}]`)
	ip, err := common.GetOutboundIP()
	if err != nil {
		log.Printf("获取本机ip失败", err)
		return
	}
	localConf := fmt.Sprintf(cfg.Etcd.Key, ip)
	log.Println(localConf)
	logconf, err := etcd.Get(localConf, time.Duration(cfg.Timeout)*time.Millisecond)
	if err != nil {
		log.Printf("get etcd data failed %v", err)
		return
	}
	//打印数据
	// for _, v := range logconf {
	// 	log.Printf("log conf %v", v)
	// }
	//启动日志收集管理器
	taillog.Init(logconf)
	//初始化新配置通道
	newConfCh := taillog.NewConfChan()
	//起一个wg防止主程序直接结束
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.Watcher(localConf, newConfCh)
	wg.Wait()
}
