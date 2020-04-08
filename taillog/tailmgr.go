package taillog

import (
	"LogAgent/etcd"
	"fmt"
	"time"
)

var (
	logMgr *tailLogManger
)

type tailLogManger struct {
	LogConf   []*etcd.LogConf
	taskMap   map[string]*TailTask
	newConfCh chan []*etcd.LogConf
}

// Init 初始化tail管理
func Init(logConf []*etcd.LogConf) {
	logMgr = &tailLogManger{
		LogConf:   logConf,
		taskMap:   make(map[string]*TailTask, 16),
		newConfCh: make(chan []*etcd.LogConf),
	}
	//根据获取的配置创建tail
	for _, conf := range logMgr.LogConf {
		tailObj := NewTail(conf.Path, conf.Topic)
		mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
		logMgr.taskMap[mk] = tailObj
	}
	//读配置通道，有新配置进来则打印
	go logMgr.handleNewConf()
}

//处理配置通道的数据
func (t *tailLogManger) handleNewConf() {
	for {
		select {
		case newConf := <-logMgr.newConfCh:
			//判断旧配置与新配置的差异
			for _, taskMapItem := range t.taskMap {
				//旧配置里没有则新建任务
				isDelete := true
				for _, newConfItem := range newConf {
					if taskMapItem.Path == newConfItem.Path && taskMapItem.Topic == newConfItem.Topic {
						isDelete = false
						continue
					}
					mk := fmt.Sprintf("%s_%s", newConfItem.Path, newConfItem.Topic)
					t.taskMap[mk] = NewTail(newConfItem.Path, newConfItem.Topic)
				}
				//删除新配置里没有旧配置里有的配置
				if isDelete {
					dmk := fmt.Sprintf("%s_%s", taskMapItem.Path, taskMapItem.Topic)
					t.taskMap[dmk].Cancel()
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

// NewConfChan 向外暴露一个只能写入新配置的通道
func NewConfChan() chan<- []*etcd.LogConf {
	return logMgr.newConfCh
}
