package xclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/pkg/errors"
)

// assembleRequest() returns a pointer to a http request instance
// with method, url and params (if method type post) as inputs
func (cli *Client) assembleRequest(method, url string, params interface{}) (*http.Request, error) {
	var body io.ReadCloser
	switch params.(type) {
	case io.ReadCloser:
		// assign the raw read closer of params to the results interface as is
		reflect.ValueOf(body).Elem().Set(reflect.ValueOf(params))
	default:
		b, err := json.Marshal(params)
		if err != nil {
			return nil, errors.Wrap(err, "json encode dto")
		}
		body = io.NopCloser(bytes.NewReader(b))
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
