package xclientv2

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const (
	DefaultTimeout           = 60 * time.Minute
	DefaultMaxRetry          = 3
	DefaultWaitMin           = 500 * time.Millisecond
	DefaultWaitMax           = 2 * time.Second
	DefaultRecycleConnection = true
	DefaultTrackProgress     = false
	DefaultContentFormat     = "application/json"
)

func defaultNeedRetry(resp *http.Response, err error) bool {
	return (err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled)) ||
		(resp != nil && resp.StatusCode >= http.StatusInternalServerError)
}
