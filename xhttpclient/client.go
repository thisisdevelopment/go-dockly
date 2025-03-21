package xhttpclient

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type NeedRetryFunc func(*http.Response, error) bool

type LogFunc func(format string, v ...any)

type Client struct {
	baseURL           string
	httpClient        *http.Client
	timeout           time.Duration
	maxRetry          int
	waitMin           time.Duration
	waitMax           time.Duration
	limiter           *rate.Limiter
	recycleConnection bool
	header            http.Header
	trackProgress     bool
	needRetry         NeedRetryFunc
	marshal           func(any) ([]byte, error)
	unmarshal         func([]byte, any) error
	logf              LogFunc
	contentFormat     string
}

func New(baseURL string, options ...Option) *Client {
	c := &Client{
		baseURL:           baseURL,
		timeout:           DefaultTimeout,
		maxRetry:          DefaultMaxRetry,
		waitMin:           DefaultWaitMin,
		waitMax:           DefaultWaitMax,
		recycleConnection: DefaultRecycleConnection,
		trackProgress:     DefaultTrackProgress,
		needRetry:         defaultNeedRetry,
		marshal:           json.Marshal,
		unmarshal:         json.Unmarshal,
		logf:              log.Printf,
		contentFormat:     DefaultContentFormat,
	}

	for _, opt := range options {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}

	if c.header == nil && c.contentFormat != "" {
		c.header = make(http.Header)
		c.header.Add("accept", c.contentFormat)
		c.header.Add("content-type", c.contentFormat)
	}

	return c
}

func (c *Client) log(format string, v ...interface{}) {
	if c.logf != nil {
		c.logf(format, v...)
	}
}
