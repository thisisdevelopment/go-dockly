package xhttpclient

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type Option func(c *Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func WithMaxRetry(maxRetry int) Option {
	return func(c *Client) {
		c.maxRetry = maxRetry
	}
}

func WithWaitMin(waitMin time.Duration) Option {
	return func(c *Client) {
		c.waitMin = waitMin
	}
}

func WithWaitMax(waitMax time.Duration) Option {
	return func(c *Client) {
		c.waitMax = waitMax
	}
}

func WithLimiter(limiter *rate.Limiter) Option {
	return func(c *Client) {
		c.limiter = limiter
	}
}

func WithRecycleConnection(recycleConnection bool) Option {
	return func(c *Client) {
		c.recycleConnection = recycleConnection
	}
}

func WithContentFormat(format string) Option {
	return func(c *Client) {
		c.contentFormat = format
	}
}

func WithHeader(header http.Header) Option {
	return func(c *Client) {
		c.header = header
	}
}

func WithHeaderMap(headerMap map[string]string) Option {
	return func(c *Client) {
		var header http.Header

		for k, v := range headerMap {
			header.Set(k, v)
		}

		c.header = header
	}
}

func WithTrackProgress(trackProgress bool) Option {
	return func(c *Client) {
		c.trackProgress = trackProgress
	}
}

func WithNeedRetry(needRetry NeedRetryFunc) Option {
	return func(c *Client) {
		c.needRetry = needRetry
	}
}

func WithMarshal(marshal func(any) ([]byte, error)) Option {
	return func(c *Client) {
		c.marshal = marshal
	}
}

func WithUnmarshal(unmarshal func([]byte, any) error) Option {
	return func(c *Client) {
		c.unmarshal = unmarshal
	}
}

func WithLog(logf LogFunc) Option {
	return func(c *Client) {
		c.logf = logf
	}
}
