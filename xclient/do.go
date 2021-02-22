package xclient

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/machinebox/progress"
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

	res, err := cli.http.Do(req.WithContext(ctx))
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
	defer res.Body.Close()
	if err != nil {
		return 0, err
	}

	if cli.config.TrackProgress {
		var contentLengthHeader = res.Header.Get("Content-Length")
		if contentLengthHeader == "" {
			return res.StatusCode, errors.New("cannot determine progress without Content-Length")
		}

		size, err := strconv.ParseInt(contentLengthHeader, 10, 64)
		if err != nil {
			return res.StatusCode, errors.Wrapf(err, "bad content-length %q", contentLengthHeader)
		}

		var r = progress.NewReader(res.Body)
		go func() {
			var progressChan = progress.NewTicker(ctx, r, size, 1*time.Second)

			for p := range progressChan {
				cli.log.Printf("%v remaining...", p.Remaining().Round(time.Second))
			}
		}()
	}

	// 	check for nil in case we are not interested in the response body
	if result != nil {
		err = cli.readResponse(res.Body, result)
	}

	return res.StatusCode, err
}
