package httpClient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"goexamples/kafka/msgMiddleware/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	httpClient http.Client
	retries    int
	timeout    int
}

func CreateHttpClient(info *model.ConsumerGroupInfo) (*Client, error) {
	client := Client{
		retries: info.Retries,
		timeout: info.Timeout,
	}

	// http客户端超时时间
	client.httpClient.Timeout = time.Duration(client.timeout) * time.Millisecond

	return &client, nil
}

func (c *Client) Call(message *model.Message) (err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	for i := 0; i < c.retries; i++ {
		req, err = http.NewRequest("POST", message.Url,
			strings.NewReader(message.Data))
		if err != nil {
			log.Error("(%d) call fail: (%v), error : (%v)",
				i, *message, err)
			continue
		}

		req.Header = message.Headers
		req.Header["Content-Length"] = []string{strconv.Itoa(len(message.Data))}
		req.Header["Content-Type"] = []string{"application/octet-stream"}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			log.Error("(%d) call fail: (%v), error : (%v)",
				i, *message, err)
			continue
		}

		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = errors.New("response status error.")
			log.Error("(%d) response status error: (%v)",
				i, *message)
			continue
		}
	}

	return err
}
