package kafka

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"goexamples/kafka/msgMiddleware/httpClient"
	"goexamples/kafka/msgMiddleware/model"

	"github.com/Shopify/sarama"
	"golang.org/x/net/context"
)

type ConsumptionMember struct {
	Ready      chan struct{}
	Messages   chan *sarama.ConsumerMessage
	HttpClient *httpClient.Client
	Errs       chan error
	Close      chan bool
	Number     int
}

func CreateConsumerGroup(cfg *model.Config) (consumer []*ConsumptionMember, err error) {

	for index, gcInfo := range cfg.ConsumerGroupList {
		member := newConsumptionMember(cfg.KafkaServer, &gcInfo, sarama.OffsetNewest, index)
		if member == nil {
			return nil, errors.New("create consumer error.")
		}

		member.HttpClient, err = httpClient.CreateHttpClient(&gcInfo)
		if err != nil {
			log.Error("create http client error :", err)
			return nil, err
		}

		// 消费
		go member.Start()

		consumer = append(consumer, member)
	}

	return
}

func newConsumptionMember(kafkaAddr []string,
	dcInfo *model.ConsumerGroupInfo,
	offset int64, number int) *ConsumptionMember {

	config := sarama.NewConfig()
	// sarama.OffsetNewest sarama.OffsetOldest
	config.Consumer.Offsets.Initial = offset
	config.Version = sarama.V0_10_2_1

	client, err := sarama.NewConsumerGroup(kafkaAddr, dcInfo.GroupId, config)
	if err != nil {
		log.Error("NewConsumerGroup error :", err)
		return nil
	}

	readyChan := make(chan struct{})
	messagesChan := make(chan *sarama.ConsumerMessage)
	errsChan := make(chan error)
	closeChan := make(chan bool)

	member := &ConsumptionMember{
		Close:    closeChan,
		Errs:     errsChan,
		Messages: messagesChan,
		Ready:    readyChan,
		Number:   number,
	}

	go func() {
		member.Errs <- <-client.Errors()
	}()

	go func() {
		<-member.Close
		err := client.Close()
		if err != nil {
			log.Debug("error :", err)
			return
		}
	}()

	go func() {
		err := client.Consume(context.Background(),
			dcInfo.Topic, member)
		if err != nil {
			log.Error("configure consume error :", err)
			return
		}
	}()

	return member
}

func (cm *ConsumptionMember) Setup(session sarama.ConsumerGroupSession) error {
	close(cm.Ready) // 该标志可用于外围进行测试consumer是否都已就位
	log.Printf("consumer member no = %d ready...", cm.Number)
	return nil
}

func (cm *ConsumptionMember) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (cm *ConsumptionMember) ConsumeClaim(session sarama.ConsumerGroupSession, consClaim sarama.ConsumerGroupClaim) error {

	for message := range consClaim.Messages() {
		cm.Messages <- message
		session.MarkMessage(message, "consumed")
	}

	return nil
}

// 可以设置并发参数，多go几个协程加速消费
func (cm *ConsumptionMember) Start() {

	for {
		select {
		case msg := <-cm.Messages:
			// get message
			fmt.Println("get message : ",
				msg.Topic, msg.Partition, string(msg.Value))

			cm.handleMsg(msg.Value, cm.Number)
		case errs := <-cm.Errs:
			log.Error("consumer get err : ", errs)
		}
	}
}

func (cm *ConsumptionMember) handleMsg(val []byte, idx int) error {

	var msg model.Message

	err := json.Unmarshal(val, &msg)
	if err != nil {
		log.Error("unmarshal value error :", err)
		return err
	}

	err = cm.HttpClient.Call(&msg)
	if err != nil {
		return err
	}

	return nil
}
