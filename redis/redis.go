package redis

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	defaultIntervalCheck = 10 * time.Second
	defaultName          = "redis"
)

type HealthChecker interface {
	SetStatus(dependencyName string, status bool)
}

type Config struct {
	address  string
	password string
	db       int
	isTLS    bool
	health   struct {
		checker       HealthChecker
		intervalCheck time.Duration
		name          string
	}
}

type Options func(*Config)

// WithHealthChecker sets the checker name and health server for the client.
// Empty name value sets the default name.
func WithHealthChecker(name string, hc HealthChecker) Options {
	return func(cfg *Config) {
		cfg.health.checker = hc
		cfg.health.name = name
	}
}

// WithIntervalCheck sets the interval for check availability.
func WithIntervalCheck(interval time.Duration) Options {
	return func(cfg *Config) {
		cfg.health.intervalCheck = interval
	}
}

func NewRedisConfig(address, password string, db int, isTLS bool, opt ...Options) *Config {
	config := &Config{
		address:  address,
		password: password,
		db:       db,
		isTLS:    isTLS,
	}

	for _, o := range opt {
		o(config)
	}

	if config.health.checker != nil {
		if config.health.name == "" {
			config.health.name = defaultName
		}
	}

	if config.health.intervalCheck <= 0 {
		config.health.intervalCheck = defaultIntervalCheck
	}

	return config
}

type Client struct {
	*redis.Client
}

// NewClient Returns new redis client
func NewClient(ctx context.Context, maxAttempts int, maxDelay time.Duration, cfg *Config) (*Client, error) {
	options := &redis.Options{
		Addr:     cfg.address,
		Password: cfg.password,
		DB:       cfg.db,
	}

	if cfg.isTLS {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := &Client{Client: redis.NewClient(options)}

	checkRedisAvailability(ctx, client, cfg)

	err := DoWithAttempts(func() error {
		pingErr := client.Ping(ctx).Err()
		if pingErr != nil {
			log.Printf("Failed to ping redis server due to error: %v\n", pingErr)
			return pingErr
		}

		return nil
	}, maxAttempts, maxDelay)
	if err != nil {
		log.Fatal("All attempts are exceeded. Unable to connect to redis server")
	}

	return client, nil
}

func DoWithAttempts(fn func() error, maxAttempts int, delay time.Duration) error {
	var err error

	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			maxAttempts--

			continue
		}

		return nil
	}

	return err
}
