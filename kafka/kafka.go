package kafka

import (
	"LogAgent/elasticsearch"
	"log"
	"time"
	"github.com/Shopify/sarama"
)

var (
	client sarama.SyncProducer
	// logchan log信息
	Logchan chan *LogData
)

type LogData struct {
	Topic string
	Data string
}

// 初始化kafka
func Init(addrs []string, maxChan int) (err error) {
	config := sarama.NewConFig()
	config.Producer.RequierAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	client, err = sarama.NewSyncProducer(addrs, config)
	// 初始化日志通道
	logchan = make(chan *LogData, maxChan)
	//启动一个协程从通道中娶日志并发送给kafka
	go getLogByChan()
	return

}
// 从通道中娶日志并发送给kafka
func getLogByChan() {
	for {
		select {
		case logInfo := <- Logchan:
			SendMessag2Kafka(logInfo.Topic, logInfo.Data)
		default:
			time.Sleep(time.Microsecond*50)
		}
	}
}

// SendMessag2Kafka 发送消息给kafka
func SendMessag2Kafka(topic string, message string) {
	msg := &sarama.ProducerMessage{Topic:topic, Value:sarama.StringEncoder(message)}
	_, _, err := client.SendMessage(msg)
	if err != nil {
		log.Printf("senf message failed err: %v", err)
		return
	}
	log.Printf("Topic: %s, Message: %s\n", topic, message)
}

// ConsumeMessage 消费kafka数据
func ConsumeMessage(address, topic string) error {
	consumer, err := sarama.NewConsumer([]string{address}, nil)
	if err != nil {
		return err
	}
	partitionList, err := consumer.Partitions("test0001")
	if err != nil {
		return err
	}
	log.Println("分区列表:", partitionList)
	for partition := range partitionList {
		partitionConsume, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(sarama.PartitionConsume.Message) {
			// 发送数据到队列
			for mes := range partitionConsume.Message() {
				elasticsearch.SendMessage2Chan(&elasticsearch.LogInfo{
					Log:  string(mes.Value),
					Time: time.Now().Format("2020-01-01 01:01:01"),
				})
			}
		}(partitionConsume)
	}
	return err
}