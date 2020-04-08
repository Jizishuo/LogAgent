package elasticsearch

import (
	"context"
	"go-common/library/database/elastic"
	"log"
	"strconv"
	"time"
	"github.com/olivere/elastic/v7"
)

var (
	client *elastic.Client
	elasticCh chan *LogInfo
	indexStr string
	typeStr string
)

type LogInfo struct {
	Log string `json:"log"`
	Time string `json:"time"`
}

func Init(address, indexStr, typeStr string, size int) (err error) {
	client, err = elastic.NewClient(elastic.SetURL(address))
	//indexStr = indexStrParam
	//typeStr = typeStrParam
	elasticCh = make(chan *LogInfo, size)
	if err != nil {
		log.Printf("indexexists failed %v", err)
		return
	}
	if !isExists {
		err = createIndex(indexStr)
		if err != nil {
			log.Printf("create index faild %v", err)
			return
		}
	}
	go SendMessag2Elastic()
	return
}

// 创建索引
func createIndex(indexStr string) (err error) {
	mapping := `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"properties":{
				"log":{
					"type":"text"
				},
				"time":{
					"type": "date",
          			"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				}
			}
		}
	}`
	// defer client.Stop()
	ctx := context.Background()
	_, err = client.CreateIndex(indexStr).BodyString(mapping).DO(ctx)
	return
}

// 确认索引是否存在
func indexExists(indexStr string) (resp bool, err error) {
	ctx := context.Background()
	resp, err = client.IndexExists(indexStr).Do(ctx)
	return
}

// SendMessag2Elastic 发送数据到es
func SendMessag2Elastic() {
	for {
		select {
		case logitem := <-elasticCh:
			id := time.Now().UnixNano()
			resp, err := client.Index().Index(indexStr).Type(typeStr).Id(strconv.Itoa(int(id))).BodyJson(logitem).Do(context.Background())
			if err != nil {
				log.Printf("insert log failed %v err", err)
				return
			}
			log.Printf("id: %v , index ： %v, type : %v", resp.Id, resp.Index, resp.Type)
		default:
			time.Sleep(time.Second)
		}
	}
}

// SendMessage2Chan 发送消息到通道
func SendMessage2Chan(message *LogInfo) {
	elasticCh <- message
}