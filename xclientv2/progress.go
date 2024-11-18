package xclientv2

import (
	"context"
	"errors"
	"fmt"
	"github.com/machinebox/progress"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (c *Client) wrapTrackProgressReader(ctx context.Context, res *http.Response, reader io.Reader) (io.Reader, error) {
	contentLengthHeader := res.Header.Get("content-length")
	if contentLengthHeader == "" {
		return nil, errors.New("track progress: cannot determine progress without Content-Length")
	}

	size, err := strconv.ParseInt(contentLengthHeader, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("track progress: bad Content-Length %s -> %w", contentLengthHeader, err)
	}

	reader = progress.NewReader(reader)

	go func() {
		progressChan := progress.NewTicker(ctx, reader.(progress.Counter), size, 1*time.Second)
		for p := range progressChan {
			c.log("%v remaining...", p.Remaining().Round(time.Second))
		}
	}()

	return reader, nil
}
