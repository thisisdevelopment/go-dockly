package xclient

import (
	"context"
	"errors"
	"math"
	"net/http"
	"time"
)

func needRetry(resp *http.Response, err error) bool {
	return (err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled)) ||
		(resp != nil && resp.StatusCode >= http.StatusInternalServerError)
}

func (cli *Client) retry(i int) bool {
	if i > cli.config.MaxRetry {
		return false
	}

	backoff := cli.backoff(i)
	cli.log.Info("retry ", i, backoff)
	time.Sleep(backoff)

	return true
}

// performs exponential backoff based on attempts and limited by waitMax
func (cli *Client) backoff(attempts int) time.Duration {

	mul := math.Pow(2, float64(attempts)) * float64(cli.config.WaitMin)
	sleep := time.Duration(mul)

	if sleep > cli.config.WaitMax {
		sleep = cli.config.WaitMax
	}

	return sleep
}
