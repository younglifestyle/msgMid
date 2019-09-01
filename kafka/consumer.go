package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"goexamples/kafka/msgMiddleware/model"
	"sync"
)

type consumer struct {
	*Consumption
	consumer sarama.PartitionConsumer
	closeCh  chan struct{}
}

type Consumption struct {
	cfg *model.Config

	wg            sync.WaitGroup
	consumers     []*consumer
	client        sarama.Client
	offsetManager sarama.OffsetManager
	baseConsumer  sarama.Consumer
}

func NewConsumption(cfg *model.Config) *Consumption { return &Consumption{cfg: cfg} }

func (c *Consumption) CreateConsume() error {

	var err error

	if c.client, err = sarama.NewClient(c.cfg.KafkaServer, nil); err != nil {
		return fmt.Errorf("new kafka client error: %s", err)
	}

	if c.baseConsumer, err = sarama.NewConsumerFromClient(c.client); err != nil {
		return fmt.Errorf("new kafka consumer error: %s", err)
	}

	for _, topicInfo := range c.cfg.KafkaTopics {

		var cp sarama.PartitionConsumer
		if cp, err = c.baseConsumer.ConsumePartition(topicInfo.Name, int32(topicInfo.Partitions),
			sarama.OffsetNewest); err != nil {
			return fmt.Errorf("new consume partition error: %s", err)
		}
		consumer := newConsumer(c, cp)
		c.consumers = append(c.consumers, consumer)
		c.wg.Add(1)
		go consumer.start()
	}

	return nil
}

func (c *Consumption) Close() error {
	for _, c := range c.consumers {
		if err := c.close(); err != nil {
			log.Warn("close consumer error: %s", err)
		}
	}
	c.wg.Wait()
	return nil
}

func newConsumer(c *Consumption, cp sarama.PartitionConsumer) *consumer {
	return &consumer{Consumption: c, consumer: cp, closeCh: make(chan struct{}, 1)}
}

func (c *consumer) close() error {
	c.closeCh <- struct{}{}
	return c.consumer.Close()
}

func (c *consumer) start() {
	defer c.wg.Done()

end_circle:
	for {
		select {
		case msg := <-c.consumer.Messages():
			log.Println("consumer partition :", msg.Partition)
			log.Println(msg.Topic, string(msg.Value))
		case <-c.closeCh:
			log.Println("consumer close...")
			break end_circle
		}
	}
}

//func CreateConsumer(c *model.Config) (*ConsumerS, error) {
//	config := sarama.NewConfig()
//	config.Consumer.Group.Heartbeat.Interval = 1 * time.Second
//	config.Consumer.Group.Session.Timeout = 30 * time.Second
//
//	//注：开发环境中我们使用sarama.OffsetOldest，Kafka将从创建以来第一条消息开始发送。
//	//在生产环境中切换为sarama.OffsetNewest,只会将最新生成的消息发送给我们。
//	config.Consumer.Offsets.Initial = sarama.OffsetOldest
//
//	consumer, err := sarama.NewConsumerGroup([]string{c.KafkaServer}, c.KafkaConsumerGroupId, config)
//	if err != nil {
//		return nil, err
//	}
//
//	err = consumer.Consume(context.Background(), []string{},
//		Consumer(func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
//
//			return nil
//		}))
//
//	return &ConsumerS{}, nil
//}
