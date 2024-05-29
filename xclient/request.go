package xclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// assembleRequest() returns a pointer to a http request instance
// with method, url and params (if method type post) as inputs
func (cli *Client) assembleRequest(method, url string, params interface{}) (*http.Request, error) {
	var body io.Reader

	switch t := params.(type) {
	case io.ReadCloser:
		// assign the raw params as read closer interface to the body as is
		body = t
	case io.Reader:
		body = t
	default:
		if params != nil {
			var b []byte
			var err error

			if cli.config.UseJsoniter {
				b, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(params)
				if err != nil {
					return nil, errors.Wrap(err, "json encode dto jsoniter")
				}
			} else {
				b, err = json.Marshal(params)
				if err != nil {
					return nil, errors.Wrap(err, "json encode dto")
				}
			}
			body = bytes.NewReader(b)
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "%s request initialization failed for %s", method, url)
	}

	req.Header.Add("Accept", cli.config.ContentFormat)
	req.Header.Add("Content-Type", cli.config.ContentFormat)

	for key, val := range cli.config.CustomHeader {
		req.Header.Add(key, val)
	}

	return req, nil
}
