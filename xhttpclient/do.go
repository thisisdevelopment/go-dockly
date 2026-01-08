package xhttpclient

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Response struct {
	StatusCode    int
	Status        string
	Header        http.Header
	ContentLength int64
}

/*
Do executes a request on this client. First arg is a context which is passed to the http request, if a deadline is not set, we
will add one for request timeout. Query parameters (as url.Values) and headers (as http.Header) for this specific request
can be passed as optional args in any order.
*/
func (c *Client) Do(ctx context.Context, method, path string, body any, result any, args ...any) (int, error) {
	resp, err := c.DoWithResponse(ctx, method, path, body, result, args...)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

/*
DoWithResponse sames as Do but returns a more informative response besides status code
*/
func (c *Client) DoWithResponse(ctx context.Context, method, path string, body any, result any, args ...any) (*Response, error) {
	var (
		query  url.Values
		header http.Header
	)

	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		// no deadline is set, we will set one ourselves
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
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

	requestUrl, err := c.buildUrl(path, query)
	if err != nil {
		return nil, err
	}

	c.log("%s %s", method, requestUrl)
	for k, v := range header {
		c.log("request header %s: [%s]", k, strings.Join(v, ","))
	}

	info, err := c.newRequestInfo(ctx, method, requestUrl, body, header)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request info -> %w", err)
	}

	return c.do(info, result)
}

func (c *Client) mergeQueryParams(inputParams url.Values) url.Values {
	merged := make(url.Values)

	// first add global query params
	for k, v := range c.queryParams {
		merged[k] = v
	}

	// local query params, overwrite existing from global
	for k, v := range inputParams {
		merged[k] = v
	}

	return merged
}

func (c *Client) buildUrl(path string, inputParams url.Values) (string, error) {
	var rawURL string

	if strings.HasPrefix(path, "http") {
		rawURL = path
	} else {
		rawURL = fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(path, "/"))
	}

	// parse url to validate and get query
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("error parsing raw url %s: %w", rawURL, err)
	}

	// get query params from url
	q := u.Query()

	// set global and request based query params, overwrite path query params
	for k, vs := range c.mergeQueryParams(inputParams) {
		q[k] = vs
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (c *Client) do(info *requestInfo, result any) (*Response, error) {
	var (
		numRetries int
		resp       *http.Response
		req        *http.Request
		err        error
	)

	cleanupResponse := func(res *http.Response) {
		// make sure we read everything even if we do nothing with the
		// see discussion: https://www.reddit.com/r/golang/comments/13fphyz/til_go_response_body_must_be_closed_even_if_you/
		_, _ = io.Copy(io.Discard, res.Body)
		// close the body, Body is guaranteed to be there even if the response does not have return data (see docs)
		_ = res.Body.Close()
	}

	retry := true
	for retry {
		if c.limiter != nil {
			// blocking call to honor the rate limit
			err = c.limiter.Wait(info.ctx)
			if err != nil {
				return nil, fmt.Errorf("rate limiter %s %s: %w", info.method, info.url, err)
			}
		}

		// create fresh request from info each retry round, if user passed in io.Reader or io.ReaderCloser we can not
		// quarantee correct behaviour (see discussion https://github.com/golang/go/issues/19653)
		req, err = info.request()
		if err != nil {
			return nil, fmt.Errorf("%s %s failed to create request from info -> %w", info.method, info.url, err)
		}

		// recycle connections or not
		req.Close = !c.recycleConnection

		// make actual request
		resp, err = c.httpClient.Do(req)
		if c.needRetry(resp, err) {
			// handle retry
			var statusCode int

			if resp != nil {
				statusCode = resp.StatusCode
			}

			c.log("need retry request -> error = %v, status code = %d", err, statusCode)

			if err == nil {
				// clean up response if not empty
				cleanupResponse(resp)
			}

			numRetries++

			if numRetries < c.maxRetry {
				// try again after some time
				sleep := time.Duration(math.Pow(2, float64(numRetries)) * float64(c.waitMin))
				if sleep > c.waitMax {
					sleep = c.waitMax
				}
				time.Sleep(sleep)
				continue
			}
		}

		retry = false

		if err != nil {
			return nil, err
		}
	}

	if resp == nil {
		return nil, fmt.Errorf("empty response should not occur")
	}

	c.log("response status: %s...", resp.Status)

	defer cleanupResponse(resp)

	var bodyReader io.Reader = resp.Body

	if result != nil {
		if c.trackProgress {
			bodyReader, err = c.wrapTrackProgressReader(info.ctx, resp, bodyReader)
			if err != nil {
				return nil, err
			}
		}

		err = c.readResponse(bodyReader, result)
	}

	for k, v := range resp.Header {
		c.log("response header %s: [%s]", k, strings.Join(v, ","))
	}

	return &Response{
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		Header:        resp.Header,
		ContentLength: resp.ContentLength,
	}, err
}
