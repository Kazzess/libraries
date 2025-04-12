package repeat

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type Operation func(context.Context) error

type Config struct {
	minTimeWait  time.Duration
	maxTimeWait  time.Duration
	maxRetries   int
	errorHandler func(err error) bool
}

const (
	defaultMinTimeWait = time.Second
	defaultMaxTimeWait = time.Minute
	defaultMaxRetries  = -1 // Infinite retries
)

func WithMinTimeWait(d time.Duration) OptionSetter {
	return func(c *Config) {
		c.minTimeWait = d
	}
}

func WithMaxTimeWait(d time.Duration) OptionSetter {
	return func(c *Config) {
		c.maxTimeWait = d
	}
}

func WithMaxRetries(retries int) OptionSetter {
	return func(c *Config) {
		c.maxRetries = retries
	}
}

func WithErrorHandler(handler func(err error) bool) OptionSetter {
	return func(c *Config) {
		c.errorHandler = handler
	}
}

type OptionSetter func(*Config)

func Exec(ctx context.Context, op Operation, opts ...OptionSetter) error {
	config := Config{
		minTimeWait:  defaultMinTimeWait,
		maxTimeWait:  defaultMaxTimeWait,
		maxRetries:   defaultMaxRetries,
		errorHandler: func(err error) bool { return true },
	}

	for _, opt := range opts {
		opt(&config)
	}

	if config.minTimeWait > config.maxTimeWait {
		return errors.New("minTimeWait cannot be greater than maxTimeWait")
	}

	var err error
	for retries := 0; config.maxRetries == -1 || retries < config.maxRetries; retries++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err = op(ctx)
			if err == nil {
				return nil
			}
			if !config.errorHandler(err) {
				return err
			}

			waitTime := config.minTimeWait + time.Duration(rand.Float64()*float64((config.maxTimeWait-config.minTimeWait)))
			timer := time.NewTimer(waitTime)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			}
		}
	}

	return err
}
