package xclient

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thisisdevelopment/go-dockly/v2/xlogger"
	"golang.org/x/time/rate"
)

// IAPIClient interface definition
type IAPIClient interface {
	Do(ctx context.Context, method, path string, params any, header map[string]string, result any) (actualStatusCode int, err error)
}

// Client defines the class implementation for this package
type Client struct {
	config  *Config
	log     *xlogger.Logger
	http    *http.Client
	baseURL string
}

// Config defines the config properties of the package
type Config struct {
	CustomHeader      map[string]string
	ContentFormat     string
	TrackProgress     bool
	RecycleConnection bool
	Limiter           *rate.Limiter // nil here will use default rate limit
	MaxRetry          int
	WaitMin           time.Duration
	WaitMax           time.Duration
	UseJsoniter       bool
}

// New returns an initiliazed API client
func New(log *xlogger.Logger,
	baseURL string,
	customHTTP *http.Client,
	customConfig *Config,
) (IAPIClient, error) {
	if baseURL == "" {
		return nil, errors.New("api needs a base URL")
	}

	config := customConfig
	if customConfig == nil {
		config = GetDefaultConfig()
	}

	client := &Client{
		log:     log,
		config:  config,
		baseURL: baseURL,
	}

	if customHTTP != nil {
		client.http = customHTTP
	} else {
		client.http = new(http.Client)
	}

	return client, nil
}

// GetDefaultConfig returns the default config for this package
func GetDefaultConfig() *Config {
	return &Config{
		WaitMin:           500 * time.Millisecond,
		WaitMax:           2 * time.Second,
		MaxRetry:          5,
		TrackProgress:     false,
		ContentFormat:     "application/json",
		RecycleConnection: true,
	}
}
