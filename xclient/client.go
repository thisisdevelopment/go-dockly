package xclient

import (
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thisisdevelopment/go-dockly/xlogger"
)

// todo https://medium.com/mflow/rate-limiting-in-golang-http-client-a22fba15861a

// IAPIClient interface definition
type IAPIClient interface {
	Do(method, path string, params io.Reader, result interface{}) (actualStatusCode int, err error)
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
	// eg map[string]string{"Authorization": "bearer " + cfg.Token},
	CustomHeader  map[string]string
	ContentFormat string
	WaitMin       time.Duration
	WaitMax       time.Duration
	MaxRetry      int
}

// New returns an initiliazed API client
func New(log *xlogger.Logger,
	baseURL string,
	customHTTP *http.Client,
	customConfig *Config) (IAPIClient, error) {

	if baseURL == "" {
		return nil, errors.New("api needs a base URL")
	}

	var config *Config
	if customConfig != nil {
		config = customConfig
	} else {
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
		WaitMin:       500 * time.Millisecond,
		WaitMax:       2 * time.Second,
		MaxRetry:      5,
		ContentFormat: "application/json",
	}
}
