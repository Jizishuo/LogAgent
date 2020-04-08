package main

import (
	"LogAgent/common"
	"LogAgent/config"
	"LogAgent/kafka"
	"LogAgent/taillog"
	"fmt"
	"log"
	"LogAgent/etcd"
	"sync"
	"time"
	"gopkg.in/ini.v1"
)

var (
	cfg = new(config.Conf)
)

func main() {
	// 加载配置
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		log.Printf("load config faild err: %v", err)
		return
	}
	log.Println("配置文件加载成功")
	// c初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.MaxChan)
	if err != nil {
		log.Printf("kafka init failed err:%v\n", err)
		return
	}
	log.Println("连接kakfa成功")

	// 初始化etcd
	err = etcd.Init(cfg.Etcd.Address, time.Duration(cfg.Etcd.Timeout)*time.Millisecond)
	if err != nil {
		log.Printf("etcd init faild, err: %v\n", err)
		return
	}

}
