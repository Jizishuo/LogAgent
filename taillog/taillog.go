package taillog

import (
	"LogAgent/kafka"
	"context"
	"log"
	"time"

	"github.com/hpcloud/tail"
)

var (
	tailsObj *tail.Tail
)

type TailTask struct {
	Path string
	Topic sting
	Instance *tail.Tail
	Ctx context.Context
	Cancel context.CancelFunc
}

// newtail 创建一个tail
func NewTail(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		Path:     path,
		Topic:    topic,
		Ctx:      ctx,
		Cancel:   cancel,
	}
	tailObj.init()
	return
}

func (t *TailTask)init() {
	config := tail.Config{
		ReOpen: true,
		Follow:true,
		location: &tail.SeekInfo{Offset:0, Whence:2},
		MustExist:false,
		Poll: true,
	}
	var err error
	t.Instance, err = tail.TailFile(t.Path, config)
	if err != nil {
		log.Printf("init tail failed %v", err)
		return
	}
	log.Printf("new tailed create success path: %s, topic: %s", t.Path, t.Topic)
	go t.getTailLines()
}
//获取要收集的日志并发送到kafka
func (t *TailTask)getTailLines() {
	for {
		select {
		case line := <- t.Instance.Lines:
			kafka.Logchan <-&kafka.LogData{
				Topic:t.Topic,
				Data:line.Text,
			}
		default:
			time.Sleep(time.Millisecond*50)
		}
	}
}