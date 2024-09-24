package tests

import (
	"context"
	"github.com/thisisdevelopment/go-dockly/v2/xclient"
	"github.com/thisisdevelopment/go-dockly/v2/xlogger"
	"testing"
)

func TestTrackProgress(t *testing.T) {
	l, err := xlogger.New(&xlogger.Config{
		Level:  "debug",
		Format: "text",
	})

	if err != nil {
		t.Fatalf("fatal err: %v", err.Error())
	}

	cfg := xclient.GetDefaultConfig()
	cfg.TrackProgress = true

	cl, err := xclient.New(l, "test", nil, cfg)
	if err != nil {
		t.Fatalf("fatal err: %v", err.Error())
	}

	var data []byte

	statusCode, err := cl.Do(context.Background(), "GET", "https://ash-speed.hetzner.com/100MB.bin", nil, nil, &data)
	if err != nil {
		t.Fatalf("fatal err: %v", err.Error())
	}

	t.Logf("statusCode: %v", statusCode)
}
