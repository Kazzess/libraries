package mynats

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
)

const (
	defaultIntervalCheck = 10 * time.Second
	defaultName          = "nats"
)

type HealthChecker interface {
	SetStatus(dependencyName string, status bool)
}

type SubscribeHandler func(ctx context.Context, msg *nats.Msg) error

type Config struct {
	servers    []string
	consumerID string
	token      string
	username   string
	password   string
	debug      bool
	tracing    bool
	health     struct {
		checker       HealthChecker
		intervalCheck time.Duration
		name          string
	}
}

type OptionSetter func(*Config)

func WithToken(token string) OptionSetter {
	return func(c *Config) { c.token = token }
}

func WithCredentials(username, password string) OptionSetter {
	return func(c *Config) {
		c.username = username
		c.password = password
	}
}

func WithDebug(debug bool) OptionSetter {
	return func(c *Config) { c.debug = debug }
}

func WithTracing(tracing bool) OptionSetter {
	return func(c *Config) { c.tracing = tracing }
}

// WithHealthChecker sets the checker name and health server for the client.
// Empty name value sets the default name.
func WithHealthChecker(name string, hc HealthChecker) OptionSetter {
	return func(cfg *Config) {
		cfg.health.checker = hc
		cfg.health.name = name
	}
}

// WithIntervalCheck sets the interval for check availability.
func WithIntervalCheck(interval time.Duration) OptionSetter {
	return func(cfg *Config) {
		cfg.health.intervalCheck = interval
	}
}

func NewConfig(servers []string, consumerID string, options ...OptionSetter) *Config {
	config := &Config{
		servers:    servers,
		consumerID: consumerID,
	}

	for _, option := range options {
		option(config)
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
	Config *Config
	nc     *nats.Conn
	js     nats.JetStreamContext
}

func NewClient(ctx context.Context, config *Config) (*Client, error) {
	options := nats.GetDefaultOptions()
	options.Servers = config.servers

	if config.token != "" {
		options.Token = config.token
	} else if config.username != "" && config.password != "" {
		options.User = config.username
		options.Password = config.password
	}

	// Debug handlers
	setDebugHandlers(&options, config.debug)

	nc, err := options.Connect()
	if err != nil {
		return nil, err
	}

	checkNatsAvailability(ctx, nc, config)

	js, jetStreamErr := nc.JetStream()
	if jetStreamErr != nil {
		return nil, jetStreamErr
	}

	return &Client{Config: config, nc: nc, js: js}, nil
}

func (c *Client) Fetch(subject, consumerID string, limit int, opts ...SubscribeOption) (_ []*nats.Msg, err error) {
	if consumerID == "" {
		consumerID = c.Config.consumerID
	}

	sub, err := c.js.PullSubscribe(subject, consumerID, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "pull subscribe")
	}

	defer func() {
		errUns := sub.Unsubscribe()
		if errUns != nil {
			err = multierror.Append(err, errUns)
		}
	}()

	messages, err := sub.Fetch(limit)
	if err != nil {
		return nil, errors.Wrap(err, "fetch messages")
	}

	return messages, err
}

func (c *Client) PullSubscribe(subject, consumerID string, opts ...SubscribeOption) (*Subscription, error) {
	if consumerID == "" {
		consumerID = c.Config.consumerID
	}

	sub, err := c.js.PullSubscribe(subject, consumerID, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "pull subscribe")
	}

	return sub, err
}

func (c *Client) Close() error {
	<-c.js.PublishAsyncComplete()

	for info := range c.js.Consumers("FOO", nats.MaxWait(10*time.Second)) {
		fmt.Println("consumer name:", info.Name)
	}

	if err := c.nc.Flush(); err != nil {
		return errors.Wrap(err, "nc.Flush")
	}

	if !c.nc.IsClosed() {
		c.nc.Close()
	}

	return nil
}

func MessageAttributes(msg *nats.Msg) []attribute.KeyValue {
	var header, meta attribute.KeyValue

	// don't handle errors, because it is not important.
	hd, err := json.Marshal(msg.Header)
	if err == nil {
		header = attribute.Key("nats.message.header").String(string(hd))
	}

	metadata, msgErr := msg.Metadata()
	if msgErr == nil {
		md, jsonErr := json.Marshal(metadata)
		if jsonErr == nil {
			meta = attribute.Key("nats.message.metadata").String(string(md))
		}
	}

	return []attribute.KeyValue{
		attribute.Key("subject").String(msg.Subject),
		attribute.Key("data").String(string(msg.Data)),
		attribute.Key("reply").String(msg.Reply),
		header,
		meta,
	}
}

func logDebugNatsConnection(c *nats.Conn, msg string) {
	if c == nil || !c.IsConnected() || c.IsClosed() || c.Status() == nats.CLOSED {
		slog.Error("no nats connection")
		return
	}

	log := slog.With(
		slog.String("status", c.Status().String()),
		slog.Uint64("stats.reconnects", c.Stats().Reconnects),
		slog.Uint64("stats.in_messages", c.Stats().InMsgs),
		slog.Uint64("stats.out_messages", c.Stats().OutMsgs),
		slog.Uint64("stats.in_bytes", c.Stats().InBytes),
		slog.Uint64("stats.out_bytes", c.Stats().OutBytes),
	)

	if lastErr := c.LastError(); lastErr != nil {
		log = log.With(slog.String("error", lastErr.Error()))
	}

	log.Info(msg)
}

func setDebugHandlers(options *nats.Options, debug bool) {
	if options == nil {
		slog.Error("no nats options")
		return
	}

	if debug {
		options.ClosedCB = func(c *nats.Conn) {
			logDebugNatsConnection(c, "NATS closed callback")
		}
		options.ReconnectedCB = func(c *nats.Conn) {
			logDebugNatsConnection(c, "NATS reconnected callback")
		}
		options.DisconnectedErrCB = func(c *nats.Conn, err error) {
			logDebugNatsConnection(c, "NATS disconnected callback")
			if err == nil {
				return
			}

			slog.With(slog.String("error", err.Error())).Info("NATS disconnected")
		}
		options.AsyncErrorCB = func(c *nats.Conn, subscription *nats.Subscription, err error) {
			logDebugNatsConnection(c, "NATS reconnected callback")
			if err == nil {
				return
			}

			slog.With(slog.String("error", err.Error())).
				With(slog.String("subject", subscription.Subject)).
				Info("Async error in subscription")
		}
	}
}
