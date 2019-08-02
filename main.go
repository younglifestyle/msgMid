package main

import (
	log "github.com/Sirupsen/logrus"
	"goexamples/kafka/msgMiddleware/kafka"
	"goexamples/kafka/msgMiddleware/model"
	"goexamples/kafka/msgMiddleware/server"
)

var (
	cfg = &model.Config{
		KafkaServer: "172.16.9.229:9029",
		KafkaTopics: []model.TopicInfo{
			{Name: "tes2", Partitions: 5},
		},
		HttpServerPort:               8090,
		HttpServerReadTimeout:        5000,
		HttpServerWriteTimeout:       5000,
		HttpServerHandlerChannelSize: 100,
	}
)

func main() {

	log.SetLevel(log.DebugLevel)

	producer, err := kafka.CreateProducer(cfg)
	if err != nil {
		log.Error("create produce error :", err)
		return
	}

	log.Info("server start...")
	runFlg := server.RunHttpServer(producer, cfg)
	if !runFlg {
		log.Panic("server run error !")
	}

	select {}
}

//func main() {
//	producer, err := CreateProducer(cfg)
//	if err != nil {
//		fmt.Println("create producer error:", err)
//		return
//	}
//	defer producer.AsyncClose()
//
//	go func(p sarama.AsyncProducer) {
//		errors := p.Errors()
//		success := p.Successes()
//		for {
//			select {
//			case err := <-errors:
//				if err != nil {
//					fmt.Println("return error :", err)
//				}
//			case sus := <-success:
//				fmt.Println("success...",
//					sus.Topic,
//					sus.Partition,
//					sus.Offset)
//			}
//		}
//	}(producer)
//
//	fmt.Println("Partitions :", int32(cfg.KafkaTopics[0].Partitions))
//
//	msg := &sarama.ProducerMessage{
//		Topic:     cfg.KafkaTopics[0].Name,
//		Partition: int32(cfg.KafkaTopics[0].Partitions),
//		Value:     sarama.StringEncoder("this is a good test"),
//	}
//
//	producer.Input() <- msg
//
//	select {}
//}
