package xclientv2

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type requestInfo struct {
	bodyData []byte
	reader   io.Reader
	header   http.Header
	method   string
	url      string
	ctx      context.Context
}

func (info *requestInfo) request() (*http.Request, error) {
	var reader io.Reader

	if info.reader != nil {
		reader = info.reader
	} else if info.bodyData != nil {
		reader = bytes.NewReader(info.bodyData)
	}

	req, err := http.NewRequestWithContext(info.ctx, info.method, info.url, reader)
	if err != nil {
		return nil, err
	}

	req.Header = info.header

	return req, nil
}

func (c *Client) newRequestInfo(ctx context.Context, method, requestUrl string, body any, header http.Header) (*requestInfo, error) {
	var info requestInfo

	info.url = requestUrl
	info.method = method
	info.ctx = ctx

	reqHeader := http.Header{}

	for k, v := range c.header {
		for _, vv := range v {
			reqHeader.Add(k, vv)
		}
	}

	for k, v := range header {
		for _, vv := range v {
			reqHeader.Add(k, vv)
		}
	}

	info.header = reqHeader

	switch t := body.(type) {
	case io.ReadCloser:
		info.reader = t
	case io.Reader:
		info.reader = t
	case []byte:
		info.bodyData = t
	default:
		if body != nil && c.marshal != nil {
			b, err := c.marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body -> %w", err)
			}
			info.bodyData = b
		}
	}

	return &info, nil
}
