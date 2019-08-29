package server

import (
	"goexamples/kafka/msgMiddleware/kafka"
	"goexamples/kafka/msgMiddleware/model"
	"net/http"
	"runtime"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	producer *kafka.Producer

	c       *model.Config
	msgChan chan *model.Message
}

type msgContext struct {
	message      *model.Message
	topic        string
	partitionNum int
}

type msgRequest struct {
	Url          string      `json:"url" binding:"required"`
	Data         interface{} `json:"data" binding:"required"`
	Topic        string      `json:"topic" binding:"required"`
	PartitionNum int         `json:"partition_num" binding:"required"`
}

type ServerCfg struct {
	httpServer  *http.Server // HTTP服务器
	httpHandler *Handler     // 路由
}

func handleSendMsg(c *gin.Context, h *Handler) {
	var (
		request msgRequest
	)

	if err := c.BindJSON(&request); err != nil {
		log.Debug("error :", err)
		c.JSON(http.StatusBadRequest, "parsing error")
		return
	}

	msg := model.Message{
		Headers:    c.Request.Header,
		Url:        request.Url,
		Data:       request.Data,
		Topic:      request.Topic,
		Partition:  request.PartitionNum,
		CreateTime: int(time.Now().Unix()),
	}

	// 发送消息
	select {
	case h.msgChan <- &msg:
		c.JSON(http.StatusOK, "send success...")
	default:
		c.JSON(http.StatusOK, "default...")
	}
}

func RunHttpServer(p *kafka.Producer, cfg *model.Config) bool {
	handler := &Handler{
		producer: p,
		c:        cfg,
	}

	r := gin.Default()

	r.POST("/msg", func(c *gin.Context) {
		handleSendMsg(c, handler)
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})

	s := &http.Server{
		Addr:         ":" + strconv.Itoa(handler.c.HttpServerPort),
		ReadTimeout:  time.Duration(handler.c.HttpServerReadTimeout),
		WriteTimeout: time.Duration(handler.c.HttpServerWriteTimeout),
		Handler:      r,
	}

	go s.ListenAndServe()

	//time.Sleep(1 * time.Second)
	//runFlg := checkRun("http://127.0.0.1" + ":" + strconv.Itoa(handler.c.HttpServerPort))
	//if !runFlg {
	//	return false
	//}

	// 启动处理器
	handler.msgChan = make(chan *model.Message,
		handler.c.HttpServerHandlerChannelSize)

	// 运行推送消息的协程
	for i := 0; i < runtime.NumCPU(); i++ {
		go handler.workProduceMsg()
	}

	return true
}

func (h *Handler) workProduceMsg() {
	for {
		select {
		case ctx := <-h.msgChan:
			log.Debug("send...")
			h.producer.SendMessage(ctx.Topic,
				int32(ctx.Partition),
				ctx)
		}
	}
}

//func (h *Handler) CheckRun(addr string) bool {
//	resp, err := resty.R().
//		Get(fmt.Sprintf("%s/health", addr))
//	if err != nil {
//		return false
//	}
//	defer resp.RawBody().Close()
//
//	if resp.StatusCode() != http.StatusOK {
//		return false
//	}
//
//	return true
//}
