package xclient

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/machinebox/progress"
	"github.com/pkg/errors"
)

// Do sends an HTTP request to the API endpoint and decodes the response body into the specified value.
// If the response is an error, it will be decoded as well.
// If the context is canceled or times out, the request will be aborted.
// If the response has a non-200 status code, an error will be returned.
// If the result parameter is nil, the response body will be discarded.
//
// The base URL for the request is determined based on the Client's configuration.
// To make a request to a different URL, use the full URL as the path parameter.
//
// The request will be retried up to the configured number of times if the request fails due to a temporary error.
// The retry policy will be applied after the backoff delay, so each retry will be delayed for a longer period of time.
//
// The request will be tracked for progress if the Client's TrackProgress option is set to true.
// The progress tracking will be done using the provided context, so the request will be canceled if the context is canceled.
// The Content-Length header of the response will be used to determine the size of the response body,
// and the progress will be updated every second.
//
// The request will be rate limited if the Client's Limiter option is set.
// The request will be blocked until the limit is available, or the context is canceled.
//
// The request method, URL, and any parameters will be logged using the Client's Logger.
//
// This function is intended to be used by the generated clients, and should not be called directly by the user.
func (cli *Client) Do(ctx context.Context, method, path string, params, result interface{}) (actualStatusCode int, err error) {
	url := fmt.Sprintf("%s/%s", cli.baseURL, path)
	if strings.HasPrefix(path, "http") {
		url = path
	}

	cli.log.Debugln(aurora.Cyan(method), aurora.Yellow(url))
	req, err := cli.assembleRequest(method, url, params)
	req.Close = !cli.config.RecycleConnection
	if err != nil {
		return 0, errors.Wrapf(err, "assemble request %s %s", method, url)
	}

	if cli.config.Limiter != nil {
		err = cli.config.Limiter.Wait(ctx) // blocking call to honor the rate limit
		if err != nil {
			return 0, errors.Wrapf(err, "rate limiter %s %s", method, url)
		}
	}

	res, err := cli.http.Do(req.WithContext(ctx))
	if err != nil {
		cli.log.Debugf("error in backoff request: %s", err.Error())
		for i := 0; i < cli.config.MaxRetry; i++ {
			err = cli.handleBackoff(i)
			if err != nil {
				err = errors.Wrapf(err, "%d backoff exhausted", i)
				break
			}
			res, err = cli.http.Do(req)
			if err != nil {
				cli.log.Debugf("error in backoff request: %s", err.Error())
				continue
			}
			defer res.Body.Close()
			break
		}
	}

	defer res.Body.Close()

	if err != nil {
		return 0, errors.Wrapf(err, "%s %s failed", method, url)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if cli.config.TrackProgress {
		contentLengthHeader := res.Header.Get("Content-Length")
		if contentLengthHeader == "" {
			return res.StatusCode, errors.New("cannot determine progress without Content-Length")
		}

		size, err := strconv.ParseInt(contentLengthHeader, 10, 64)
		if err != nil {
			return res.StatusCode, errors.Wrapf(err, "bad content-length %q", contentLengthHeader)
		}

		r := progress.NewReader(res.Body)
		go func() {
			progressChan := progress.NewTicker(ctx, r, size, 1*time.Second)

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
