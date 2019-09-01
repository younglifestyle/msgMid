package kafka

import (
	"encoding/json"
	"fmt"
	"goexamples/kafka/msgMiddleware/model"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	c        *model.Config
	producer sarama.AsyncProducer
}

// syncProducer 同步生产者
// 并发量小时，可以用这种方式

// asyncProducer 异步生产者
// 并发量大时，必须采用这种方式

// 创造生产者
func CreateProducer(c *model.Config) (*Producer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true //必须有这个选项
	config.Producer.Timeout = 5 * time.Second
	config.Producer.Partitioner = sarama.NewManualPartitioner

	producer, err := sarama.NewAsyncProducer(c.KafkaServer, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer, c: c}, err
}

// 发送消息
func (p *Producer) SendMessage(topic string, partitionNum int32,
	message *model.Message) bool {

	bytes, err := json.Marshal(message)
	if err != nil {
		return false
	}

	msg := sarama.ProducerMessage{
		Topic:     topic,
		Partition: partitionNum,
		Value:     sarama.ByteEncoder(bytes),
		Timestamp: time.Time{},
	}

	p.producer.Input() <- &msg

	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					fmt.Println("return error :", err)
				}
			case <-success:
				fmt.Println("success...")
			}
		}
	}(p.producer)

	return true
}
