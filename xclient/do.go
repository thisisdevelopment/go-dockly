package xclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/machinebox/progress"
	"github.com/pkg/errors"
)

func clearResponse(resp *http.Response) {
	// make sure we read everything even if we do nothing with the
	// see discussion: https://www.reddit.com/r/golang/comments/13fphyz/til_go_response_body_must_be_closed_even_if_you/
	_, _ = io.Copy(io.Discard, resp.Body)
	// close the body, Body is guaranteed to be there even if the response does not have return data (see docs)
	_ = resp.Body.Close()
}

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
// The optional args are checked and if it is a url.Values it will be set as query parameters and if it is http.Header
// it will be set as per request header.
//
// This function is intended to be used by the generated clients, and should not be called directly by the user.
func (cli *Client) Do(ctx context.Context, method, path string, body any, result any, args ...any) (actualStatusCode int, err error) {
	var (
		query  url.Values
		header http.Header
		reqUrl string
		req    *http.Request
		res    *http.Response
	)

	if strings.HasPrefix(path, "http") {
		reqUrl = path
	} else {
		reqUrl, err = url.JoinPath(cli.baseURL, path)
		if err != nil {
			return 0, fmt.Errorf("failed to join url path: %w", err)
		}
	}

	for _, arg := range args {
		// assign query values and possible header if found
		switch t := arg.(type) {
		case url.Values:
			query = t
		case http.Header:
			header = t
		}
	}

	if query != nil {
		reqUrl = fmt.Sprintf("%s?%s", reqUrl, query.Encode())
	}

	cli.log.Debugln(aurora.Cyan(method), aurora.Yellow(reqUrl))

	info, err := cli.newRequestInfo(ctx, method, reqUrl, body, header)
	if err != nil {
		return 0, fmt.Errorf("assemble request %s %s: %w", method, reqUrl, err)
	}

	var numRetries int

	retry := true
	for retry {
		if cli.config.Limiter != nil {
			err = cli.config.Limiter.Wait(ctx) // blocking call to honor the rate limit
			if err != nil {
				return 0, errors.Wrapf(err, "rate limiter %s %s", method, reqUrl)
			}
		}

		req, err = info.request()
		if err != nil {
			return 0, errors.Wrapf(err, "assemble request %s %s", method, reqUrl)
		}

		req.Close = !cli.config.RecycleConnection

		res, err = cli.http.Do(req.WithContext(ctx))
		if needRetry(res, err) {
			if err == nil {
				clearResponse(res)
			}

			if cli.retry(numRetries) {
				numRetries++
				continue
			}
		}

		retry = false
	}

	if err != nil {
		return 0, errors.Wrapf(err, "%s %s failed", method, reqUrl)
	}

	defer clearResponse(res)

	var bodyReader io.Reader = res.Body

	if cli.config.TrackProgress {
		contentLengthHeader := res.Header.Get("Content-Length")
		if contentLengthHeader == "" {
			return res.StatusCode, errors.New("cannot determine progress without Content-Length")
		}

		size, err := strconv.ParseInt(contentLengthHeader, 10, 64)
		if err != nil {
			return res.StatusCode, errors.Wrapf(err, "bad content-length %q", contentLengthHeader)
		}

		// wrap reader in progress reader that implements counter interface
		bodyReader = progress.NewReader(bodyReader)
		go func() {
			progressChan := progress.NewTicker(ctx, bodyReader.(progress.Counter), size, 1*time.Second)

			for p := range progressChan {
				cli.log.Printf("%v remaining...", p.Remaining().Round(time.Second))
			}
		}()
	}

	// check for nil in case we are not interested in the response body
	if result != nil {
		err = cli.readResponse(bodyReader, result)
	}

	return res.StatusCode, err
}
