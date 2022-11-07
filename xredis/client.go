package xredis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/logrusorgru/aurora"
	"github.com/thisisdevelopment/go-dockly/v2/xlogger"
)

// ICache defines and exposes the caching layer
type ICache interface {
	// Set the value to key with TTL overwrite or TTL taken from config
	Set(ctx context.Context, key string, value interface{}, t time.Duration) error
	// Get will return value of key under cancellable context and try unmarshal the result into expected
	Get(ctx context.Context, key string, expected interface{}) error
	// GetCacheKeys returns all the keys in the oartial match pattern
	GetCacheKeys(ctx context.Context, keyPattern string) ([]string, error)
	// GetTTL returns the remaining time to live duration for key or 0 if expired
	GetTTL(ctx context.Context, key string) (time.Duration, error)
}

// Redis implements the ICache interface based on redis
type Redis struct {
	redis  *redis.Client
	config *Config
	log    *xlogger.Logger
}

type Config struct {
	Host         string
	Pass         string
	DB           int
	Expiration   int
	PoolSize     int           `yaml:"pool_size"`
	MaxRetries   int           `yaml:"max_retries"`
	ConnTimeOut  time.Duration `yaml:"conn_timeout"`
	PollInterval time.Duration `yaml:"poll_interval"`
	TLS          bool
}

// New constructs a cache class
func New(config *Config, log *xlogger.Logger) (ICache, error) {

	var opts = defaultOpts(config)
	var client = redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), config.ConnTimeOut)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	log.Printf("connected to redis %s\n", aurora.Cyan(config.Host))

	var o = &Redis{
		config: config,
		log:    log,
		redis:  client,
	}

	go o.checkConnection()

	return o, nil
}

func defaultOpts(config *Config) *redis.Options {
	var opts = &redis.Options{
		Addr:       config.Host,
		Password:   config.Pass,
		DB:         config.DB,
		PoolSize:   config.PoolSize,
		MaxRetries: config.MaxRetries,
	}

	if config.TLS {
		opts.TLSConfig = new(tls.Config)
	}

	return opts
}

// pings the connection at tick interval and tries to continously reconnect on error
func (r *Redis) checkConnection() {

	for range time.Tick(r.config.PollInterval) {
		ctx, cancel := context.WithTimeout(context.Background(), r.config.PollInterval)

		err := r.redis.Ping(ctx).Err()
		if err != nil {
			// redis disconnected
			opts := defaultOpts(r.config)
			client := redis.NewClient(opts)
			r.redis = client

			if err := r.redis.Ping(ctx).Err(); err != nil {
				r.log.Warningln("redis connnection failed, failed to set new one")
				cancel()
				continue
			}
		}
		cancel()
	}
}
