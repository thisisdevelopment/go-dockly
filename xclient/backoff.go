package xclient

import (
	"math"
	"time"

	"github.com/pkg/errors"
)

func (cli *Client) handleBackoff(i int) error {
	if i > cli.config.MaxRetry {
		return errors.New("backoff exhausted")
	}

	backoff := cli.backoff(i)
	cli.log.Info("retry ", i, backoff)
	time.Sleep(backoff)

	return nil
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
