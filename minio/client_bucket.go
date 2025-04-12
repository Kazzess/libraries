package minio

import (
	"context"
	"errors"

	"github.com/minio/minio-go/v7"
)

// ListBucket list all objects in a bucket. If any error why any object - function return list of success objects
// and errors with information failed objects.
func (c *Client) ListBucket(
	ctx context.Context,
	bucketName string,
	opts ListObjectsOptions,
) (res []ObjectInfo, err error) {
	exists, err := c.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrBucketDoesNotExist
	}

	for lobj := range c.minioClient.ListObjects(ctx, bucketName, opts) {
		if lobj.Err != nil {
			err = errors.Join(err, lobj.Err)
			continue
		}

		res = append(res, lobj)
	}

	return res, err
}

func (c *Client) CreateBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	err := c.minioClient.MakeBucket(ctx, bucketName, opts)
	if err != nil {
		return realError(err)
	}

	return nil
}

func (c *Client) RemoveBucket(ctx context.Context, bucketName string) error {
	err := c.minioClient.RemoveBucket(ctx, bucketName)
	if err != nil {
		return realError(err)
	}

	return nil
}

func (c *Client) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	exists, err := c.minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return false, realError(err)
	}

	if !exists {
		return false, nil
	}

	return true, nil
}
