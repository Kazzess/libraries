package psql

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultIntervalCheck = 10 * time.Second
	defaultName          = "postgres"
)

type HealthChecker interface {
	SetStatus(dependencyName string, status bool)
}

type Config struct {
	pgxConfig   *pgxpool.Config
	dsn         string
	maxAttempts int
	maxDelay    time.Duration
	health      struct {
		checker       HealthChecker
		intervalCheck time.Duration
		name          string
	}
}

type Option func(*Config)

func WithBinaryExecMode(binary bool) Option {
	return func(cfg *Config) {
		if binary {
			cfg.pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
		}
	}
}

// WithHealthChecker sets the checker name and health server for the client.
// Empty name value sets the default name.
func WithHealthChecker(name string, hc HealthChecker) Option {
	return func(cfg *Config) {
		cfg.health.checker = hc
		cfg.health.name = name
	}
}

// WithIntervalCheck sets the interval for check availability.
func WithIntervalCheck(interval time.Duration) Option {
	return func(cfg *Config) {
		cfg.health.intervalCheck = interval
	}
}

func NewConfig(dsn string, maxAttempts int, maxDelay time.Duration, options ...Option) (*Config, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Printf("Unable to parse config: %v\n", err)
		return nil, err
	}

	config := &Config{
		pgxConfig:   cfg,
		dsn:         dsn,
		maxAttempts: maxAttempts,
		maxDelay:    maxDelay,
	}

	for _, o := range options {
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

	return config, nil
}

type Client struct {
	*pgxpool.Pool
}

// NewClient creates new postgres client.
func NewClient(ctx context.Context, cfg *Config) (client *Client, err error) {
	pool, configErr := pgxpool.NewWithConfig(ctx, cfg.pgxConfig)
	if configErr != nil {
		log.Printf("Failed to parse PostgreSQL configuration due to error: %v\n", configErr)
		return nil, configErr
	}

	err = checkPostgresAvailability(ctx, pool, cfg)
	if err != nil {
		return nil, err
	}

	err = DoWithAttempts(func() error {
		pingErr := pool.Ping(ctx)
		if pingErr != nil {
			log.Printf("Failed to connect to postgres due to error %v... Going to do the next attempt\n", pingErr)
			return pingErr
		}

		return nil
	}, cfg.maxAttempts, cfg.maxDelay)
	if err != nil {
		log.Fatal("All attempts are exceeded. Unable to connect to PostgreSQL")
	}

	client = &Client{
		Pool: pool,
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

type DSNConfig struct {
	username string
	password string
	host     string
	port     string
	database string
}

func (c *DSNConfig) ToPostgresDSN() string {
	url := strings.Builder{}
	url.WriteString("postgresql://")
	url.WriteString(c.username)
	url.WriteString(":")
	url.WriteString(c.password)
	url.WriteString("@")
	url.WriteString(c.host)
	url.WriteString(":")
	url.WriteString(c.port)
	url.WriteString("/")
	url.WriteString(c.database)

	return url.String()
}
