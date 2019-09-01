package model

// Kafka消息体
type Message struct {
	Headers    map[string][]string `json:"headers"`     // 原始http头
	Url        string              `json:"url"`         // 发送目标
	Data       string              `json:"data"`        // json数据 发送数据
	Topic      string              `json:"topic"`       // 消息所处的Topic
	Partition  int                 `json:"partition"`   // 消息所处的分区ID
	CreateTime int                 `json:"create_time"` // 消息创建的时间
}

type TopicInfo struct {
	Name       string `json:"name"`
	Partitions int    `json:"partitions"`
}

type ConsumerGroupInfo struct {
	Topic   []string `json:"topic"`
	GroupId string   `json:"group_id"`
	Retries int      `json:"retries"`
	Timeout int      `json:"timeout"`
}

type Config struct {
	// Kafka地址
	KafkaServer []string    `json:"kafka_server"`
	KafkaTopics []TopicInfo `json:"kafka_topics"` // topic信息 producer往哪个partition生成

	ConsumerGroupList []ConsumerGroupInfo `json:"consumer_group_list"`

	// HTTP服务配置
	HttpServerPort               int `json:"http_server_port"`
	HttpServerReadTimeout        int `json:"http_server_read_timeout"`
	HttpServerWriteTimeout       int `json:"http_server_write_timeout"`
	HttpServerHandlerChannelSize int `json:"http_server_handler_channel_size"`
}
