package xclient

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// assembleRequest() returns a pointer to a http request instance
// with method, url and params (if method type post) as inputs
func (cli *Client) assembleRequest(method, url string, params io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, params)
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
