package tracing

import (
	"errors"
)

var (
	ErrHostIsEmpty = errors.New("host is empty")
	ErrPortIsEmpty = errors.New("port is empty")
)

type config struct {
	host           string
	port           string
	serviceID      string
	serviceName    string
	serviceVersion string
	envName        string
}

func (c *config) Validate() error {
	if c.host == "" {
		return ErrHostIsEmpty
	}

	if c.port == "" {
		return ErrPortIsEmpty
	}

	return nil
}

type ConfigParam func(c *config)

func WithHost(val string) ConfigParam {
	return func(c *config) {
		c.host = val
	}
}

func WithPort(val string) ConfigParam {
	return func(c *config) {
		c.port = val
	}
}

func WithServiceID(val string) ConfigParam {
	return func(c *config) {
		c.serviceID = val
	}
}

func WithServiceName(val string) ConfigParam {
	return func(c *config) {
		c.serviceName = val
	}
}

func WithServiceVersion(val string) ConfigParam {
	return func(c *config) {
		c.serviceVersion = val
	}
}

func WithEnvName(val string) ConfigParam {
	return func(c *config) {
		c.envName = val
	}
}
