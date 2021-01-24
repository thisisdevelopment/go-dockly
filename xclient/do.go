package xclient

import (
	"context"
	"fmt"
	"io"

	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
)

func (cli *Client) Do(ctx context.Context, method, path string, params io.Reader, result interface{}) (actualStatusCode int, err error) {
	url := fmt.Sprintf("%s/%s", cli.baseURL, path)

	cli.log.Debugln("requesting: ", aurora.Yellow(url))
	req, err := cli.assembleRequest(method, url, params)
	if err != nil {
		return 0, err
	}

	err = cli.config.Limiter.Wait(ctx) // blocking call to honor the rate limit
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)

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
	if err != nil {
		return 0, err
	}
	// 	check for nil in case we are not interested in the response body
	if result != nil {
		err = cli.readResponse(res.Body, result)
	}

	res.Body.Close()
	return res.StatusCode, err
}
