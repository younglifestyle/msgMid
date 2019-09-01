package main

import (
	"fmt"
	"goexamples/kafka/msgMiddleware/kafka"
	"goexamples/kafka/msgMiddleware/model"
	"goexamples/kafka/msgMiddleware/server"

	log "github.com/Sirupsen/logrus"
)

var (
	cfg *model.Config = &model.Config{
		KafkaServer: []string{"172.16.9.229:9029"},
		KafkaTopics: []model.TopicInfo{
			{Name: "tes2", Partitions: 1},
			{Name: "tes2", Partitions: 2},
			{Name: "tes2", Partitions: 3},
		},

		//KafkaConsumerGroupId: "test-consumer-group",
		ConsumerGroupList: []model.ConsumerGroupInfo{
			{Topic: []string{"tes2"}, GroupId: "test-consumer-group1", Retries: 3, Timeout: 3},
			{Topic: []string{"tes2"}, GroupId: "test-consumer-group2", Retries: 3, Timeout: 3},
		},

		HttpServerPort:               8090,
		HttpServerReadTimeout:        5000,
		HttpServerWriteTimeout:       5000,
		HttpServerHandlerChannelSize: 500,
	}
)

func main() {
	log.SetLevel(log.DebugLevel)

	//cfgFlg := flag.String("c", "./cfgfile/cfg.json", "configuration file")
	//flag.Parse()
	//
	//err := g.ParseConfig(*cfgFlg)
	//if err != nil {
	//	log.Println("parse config is failed :", err)
	//	return
	//}
	//
	//cfg = g.Config()

	fmt.Println("start producer... ")
	producer, err := kafka.CreateProducer(cfg)
	if err != nil {
		log.Error("create produce error :", err)
		return
	}

	fmt.Println("start Consumer... ")
	//consumption := kafka.NewConsumption(cfg)
	//err = consumption.CreateConsume()
	//if err != nil {
	//	log.Error("consumption error :", err)
	//	return
	//}

	_, err = kafka.CreateConsumerGroup(cfg)
	if err != nil {
		fmt.Println("create consumer error.")
		return
	}

	log.Info("server start...")
	runFlg := server.RunHttpServer(producer, cfg)
	if !runFlg {
		log.Panic("server run error !")
	}

	select {}
}
