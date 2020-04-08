package kafka

import (
	"LogAgent/elasticsearch"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

var (
	client sarama.SyncProducer
	// Logchan log信息
	Logchan chan *LogData
)

// LogData 日志信息结构体
type LogData struct {
	Topic string
	Data  string
}

// Init 初始化kafka
func Init(addrs []string, maxChan int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	client, err = sarama.NewSyncProducer(addrs, config)
	//初始化日志通道
	Logchan = make(chan *LogData, maxChan)
	//启动一个协程从通道中取日志并发送给kafka
	go getLogByChan()
	return
}

// 从通道中取出log并发送到kafka
func getLogByChan() {
	for {
		select {
		case logInfo := <-Logchan:
			SendMessag2Kafka(logInfo.Topic, logInfo.Data)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

// SendMessag2Kafka 发送信息到kafka
func SendMessag2Kafka(topic string, message string) {
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(message)}
	_, _, err := client.SendMessage(msg)
	if err != nil {
		log.Printf("send message failed err : %v", err)
		return
	}
	log.Printf("Topic: %s , Message : %s\n", topic, message)
}

// ConsumeMessage 消费kafka数据
func ConsumeMessage(address, topic string) error {
	consumer, err := sarama.NewConsumer([]string{address}, nil)
	if err != nil {
		return err
	}
	partitionList, err := consumer.Partitions("test01")
	if err != nil {
		return err
	}
	log.Println("分区列表：", partitionList)
	for partion := range partitionList {
		partitionConsume, err := consumer.ConsumePartition(topic, int32(partion), sarama.OffsetNewest)
		if err != nil {
			return err
		}
		// defer partitionConsume.AsyncClose()
		go func(sarama.PartitionConsumer) {
			for msg := range partitionConsume.Messages() {
				//发送数据到队列
				elasticsearch.SendMessage2Chan(&elasticsearch.LogInfo{
					Log:  string(msg.Value),
					Time: time.Now().Format("2006-01-02 15:04:05"),
				})
			}
		}(partitionConsume)
	}
	return err
}
