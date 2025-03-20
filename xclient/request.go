package xclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
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

func (cli *Client) newRequestInfo(ctx context.Context, method, requestUrl string, body any, header http.Header) (*requestInfo, error) {
	var info requestInfo

	info.url = requestUrl
	info.method = method
	info.ctx = ctx

	reqHeader := http.Header{}

	reqHeader.Add("Accept", cli.config.ContentFormat)
	reqHeader.Add("Content-Type", cli.config.ContentFormat)

	for key, val := range cli.config.CustomHeader {
		reqHeader.Add(key, val)
	}

	for key, val := range cli.perRequestHeader {
		reqHeader.Add(key, val)
	}

	for key, vals := range header {
		for _, val := range vals {
			reqHeader.Add(key, val)
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
		if body != nil {
			var b []byte
			var err error

			if cli.config.UseJsoniter {
				b, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(body)
				if err != nil {
					return nil, errors.Wrap(err, "json encode dto jsoniter")
				}
			} else {
				b, err = json.Marshal(body)
				if err != nil {
					return nil, errors.Wrap(err, "json encode dto")
				}
			}
			info.bodyData = b
		}
	}

	return &info, nil
}
