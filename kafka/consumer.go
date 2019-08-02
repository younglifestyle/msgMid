package kafka

import (
	"github.com/Shopify/sarama"
	"goexamples/kafka/msgMiddleware/model"
	"time"
)

type Consumer struct {
	c         *model.Config
	consumer  sarama.ConsumerGroup
	consumer1 sarama.Consumer
}

func CreateConsumer(c *model.Config) {
	config := sarama.NewConfig()
	config.Consumer.Group.Heartbeat.Interval = 1*time.Second
	config.Consumer.Group.Session.Timeout = 30*time.Second
}
