package model

// Kafka消息体
type Message struct {
	Headers    map[string][]string // 原始http头
	Url        string              // 发送目标
	Data       interface{}         // 发送数据
	Topic      string              // 消息所处的Topic
	Partition  int                 // 消息所处的分区ID
	CreateTime int                 // 消息创建的时间
}

type TopicInfo struct {
	Name       string `json:"name"`
	Partitions int    `json:"partitions"`
}

type Config struct {
	// Kafka地址
	KafkaServer string      `json:"kafka_server"`
	KafkaTopics []TopicInfo `json:"kafka_topics"` // topic信息

	// HTTP服务配置
	HttpServerPort               int `json:"http_server_port"`
	HttpServerReadTimeout        int `json:"http_server_read_timeout"`
	HttpServerWriteTimeout       int `json:"http_server_write_timeout"`
	HttpServerHandlerChannelSize int `json:"http_server_handler_channel_size"`
}
