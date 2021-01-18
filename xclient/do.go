package xclient

import (
	"fmt"
	"io"

	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
)

// Do sends the request to the specified rest path and unmarshals the response into the
// desired results interface{} if not provided as null and under the condition that
// the received status response is as expected
func (cli *Client) Do(method, path string, params io.Reader, result interface{}) (actualStatusCode int, err error) {
	url := fmt.Sprintf("%s/%s", cli.baseURL, path)

	cli.log.Debugln("requesting: ", aurora.Yellow(url))
	req, err := cli.assembleRequest(method, url, params)
	if err != nil {
		return 0, err
	}

	res, err := cli.http.Do(req)
	if err != nil {
		for i := 0; i < cli.config.MaxRetry; i++ {
			err = cli.handleBackoff(i)
			if err != nil {
				err = errors.Wrapf(err, "%d backoff exhausted", i)
				break
			}
			res, err = cli.http.Do(req)
			if err != nil {
				continue
			}
		}
	}

	// 	check in case we are not interested in the response body
	if result != nil && err == nil {
		err = cli.readResponse(res.Body, result)
	}

	res.Body.Close()
	return res.StatusCode, err
}
