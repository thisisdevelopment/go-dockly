package xclient

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// assembleRequest() returns a pointer to a http request instance
// with method, url and params (if method type post) as inputs
func (cli *Client) assembleRequest(method, url string, params interface{}) (*http.Request, error) {
	var body *bytes.Reader
	switch params.(type) {
	case bytes.Reader:
		// assign the raw params as read closer interface to the body as is
		body = params.(*bytes.Reader)
	default:
		if params == nil {
			body = bytes.NewReader(nil)
		} else {
			b, err := json.Marshal(params)
			if err != nil {
				return nil, errors.Wrap(err, "json encode dto")
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
