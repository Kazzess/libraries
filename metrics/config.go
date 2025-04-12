package metrics

import (
	"errors"
	"fmt"
	"time"
)

const path = "/metrics"

var (
	ErrBadHost = errors.New("unknown host")
	ErrBadPort = errors.New("unknown port")
)

// Config - metrics server configuration.
type Config struct {
	address           string
	host              string
	port              int
	readTimeout       time.Duration
	writeTimeout      time.Duration
	readHeaderTimeout time.Duration
}

func (c *Config) Validate() error {
	if c.port == 0 {
		return ErrBadPort
	}

	if c.host == "" {
		return ErrBadHost
	}

	return nil
}

// Option - Defines an option type for configuring the Config.
type Option func(*Config)

// WithHost - HTTP port to listen on.
func WithHost(host string) Option {
	return func(c *Config) {
		c.host = host
	}
}

// WithPort - HTTP port to listen on.
func WithPort(port int) Option {
	return func(c *Config) {
		c.port = port
	}
}

// WithReadTimeout - Configures the readTimeout
func WithReadTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.readTimeout = timeout
	}
}

// WithWriteTimeout - Configures the writeTimeout
func WithWriteTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.writeTimeout = timeout
	}
}

// WithReadHeaderTimeout - Configures the readHeaderTimeout
func WithReadHeaderTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.readHeaderTimeout = timeout
	}
}

// NewConfig - Creates a new Config with the provided options
func NewConfig(opts ...Option) *Config {
	config := &Config{} // Could set default values here
	for _, opt := range opts {
		opt(config)
	}

	config.address = fmt.Sprintf("%s:%d", config.host, config.port)

	return config
}
