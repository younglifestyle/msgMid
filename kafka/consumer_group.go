package kafka

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/Shopify/sarama"
	"golang.org/x/net/context"
)

type ConsumptionMember struct {
	Ready    chan bool
	Messages chan *sarama.ConsumerMessage
	Errs     chan error
	Close    chan bool
	Name     string
	Number   int
}

func NewConsumptionMember(groupID string, topics []string, offset int64, name string, number int) *ConsumptionMember {

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = offset

	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	client, err := sarama.NewConsumerGroup(kafkaBrokers, groupID, config)
	if err != nil {
		log.Error("error :", err)
		return nil
	}

	readyChan := make(chan bool)
	messagesChan := make(chan *sarama.ConsumerMessage)
	errsChan := make(chan error)
	closeChan := make(chan bool)

	member := &ConsumptionMember{
		Close:    closeChan,
		Errs:     errsChan,
		Messages: messagesChan,
		Ready:    readyChan,
		Name:     name,
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
		err := client.Consume(context.Background(), topics, member)
		if err != nil {
			log.Debug("error :", err)
			return
		}
	}()

	return member
}

func (cm *ConsumptionMember) Setup(session sarama.ConsumerGroupSession) error {
	close(cm.Ready)
	log.Printf("%s consumer member no = %d ready...", cm.Name, cm.Number)
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
