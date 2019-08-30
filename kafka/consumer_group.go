package kafka

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"goexamples/kafka/msgMiddleware/model"

	"github.com/Shopify/sarama"
	"golang.org/x/net/context"
)

type ConsumptionMember struct {
	//Ready    chan bool
	Messages chan *sarama.ConsumerMessage
	Errs     chan error
	Close    chan bool
	Number   int
}

func NewConsumptionMember(kafkaAddr []string,
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

	//readyChan := make(chan bool)
	messagesChan := make(chan *sarama.ConsumerMessage)
	errsChan := make(chan error)
	closeChan := make(chan bool)

	member := &ConsumptionMember{
		Close:    closeChan,
		Errs:     errsChan,
		Messages: messagesChan,
		//Ready:    readyChan,
		Number: number,
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
	//close(cm.Ready)
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
		case errs := <-cm.Errs:
			fmt.Println("get err : ", errs)
		}
	}
}
