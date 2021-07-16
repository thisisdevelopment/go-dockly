package xredis

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
)

// Set the value to key with TTL overwrite or TTL taken from config
func (c *Redis) Set(ctx context.Context, key string, value interface{}, t time.Duration) (err error) {
	var g []byte
	// 0 will default to config controlled cache ttl
	if t == 0 {
		t = time.Duration(c.config.Expiration) * time.Minute
	}

	switch value.(type) {
	// if set value is byte slice it is assumed you know what you do
	// as in the value is already marshalled from some format eg yaml
	case []byte:
		// gzipped storing in redis yields x10 size reduction
		g, err = c.gzip(value.([]byte))
		if err != nil {
			return errors.Wrapf(err, "gzip byte %s", aurora.Yellow(key))
		}

	default:
		// if set value is interface it will be json marshalled
		b, err := json.Marshal(value)
		if err != nil {
			return errors.Wrapf(err, "marshal %s", aurora.Yellow(key))
		}

		g, err = c.gzip(b)
		if err != nil {
			return errors.Wrapf(err, "gzip interface %s", aurora.Yellow(key))
		}
	}

	return c.redis.Set(ctx, key, g, t).Err()
}

// Get will return value of key under cancellable context and try unmarshal the result into expected
func (c *Redis) Get(ctx context.Context, key string, expected interface{}) error {

	val, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		return errors.Wrapf(err, "get %s", aurora.Yellow(key))
	}

	// each value this application controls is assumed to be gzipped compressed
	b, err := c.gunzip(val)
	if err != nil {
		// or an error will be thrown in case it is eg json marshalled byte slice
		return errors.Wrapf(err, "gunzip %s", aurora.Yellow(key))
	}

	switch expected.(type) {
	// the special case bypassing an unmarshal indicating a different format otherwise json is assumed
	case *[]byte:
		// assign the raw byte slice to expected interface as is (caller handles payload)
		reflect.ValueOf(expected).Elem().Set(reflect.ValueOf(b))
	default:
		// we handle the payload and unmarshal into the expected interface directly
		if err = json.Unmarshal(b, expected); err != nil {
			return errors.Wrapf(err, "unmarshal %s ", aurora.Yellow(key))
		}
	}

	return nil
}

// GetCacheKeys returns all the keys in the oartial match pattern
func (c *Redis) GetCacheKeys(ctx context.Context, keyPattern string) (ks []string, err error) {

	var cursor uint64

	for {
		var keys []string
		var err error
		// only string keys are returned no payloads
		keys, cursor, err = c.redis.Scan(ctx, cursor, keyPattern, 512).Result()
		if err != nil {
			return nil, err
		}

		ks = append(ks, keys...)

		if cursor == 0 {
			break
		}
	}

	sort.Slice(ks, func(i, j int) bool {
		return ks[i] < ks[j]
	})

	return ks, err
}

// GetTTL returns the remaining time to live duration for key or 0 if expired
func (c *Redis) GetTTL(ctx context.Context, key string) (time.Duration, error) {

	ttl, err := c.redis.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if ttl < 0 {
		return 0, nil
	}

	return ttl, err
}
