package minio

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	defaultIntervalCheck = 10 * time.Second
	defaultName          = "minio"
)

type HealthChecker interface {
	SetStatus(dependencyName string, status bool)
}

type Config struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	health          struct {
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

func NewConfig(endpoint, accessKeyID, secretAccessKey string, options ...Options) *Config {
	config := &Config{
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
	}

	for _, opt := range options {
		opt(config)
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
	minioClient *minio.Client
}

func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	minioClient, err := minio.New(cfg.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.accessKeyID, cfg.secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}

	client := &Client{
		minioClient: minioClient,
	}

	err = checkMinioAvailability(ctx, client, cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func realError(err error) error {
	responseErr := minio.ToErrorResponse(err)
	switch responseErr.Code {
	case "BucketNotEmpty":
		return ErrRemoveBucketNotEmpty
	case "NoSuchBucket":
		return ErrBucketDoesNotExist
	case "BucketAlreadyOwnedByYou":
		return ErrBucketAlreadyOwnedByYou
	case "NoSuchKey":
		return ErrObjectDoesNotExist
	}

	return err
}
