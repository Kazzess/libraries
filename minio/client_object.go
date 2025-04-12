package minio

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/minio/minio-go/v7"
)

// Upload upload object to minio. if bucketName set - upload object to specific bucket.
// if bucketName not set - generate bucket, add prefix (if set).
// return objectID and bucket.
func (c *Client) Upload(
	ctx context.Context,
	object UploadObject,
	bucketName string,
) (string, string, error) {
	info, err := c.minioClient.PutObject(ctx, bucketName, object.Name, object.Reader, object.Size,
		minio.PutObjectOptions{
			UserMetadata: object.Meta,
			ContentType:  object.ContentType.String(),
		})
	if err != nil {
		return "", "", realError(err)
	}

	return info.Key, info.Bucket, nil
}

func (c *Client) Object(
	ctx context.Context,
	bucketName, fileName string,
	opts GetObjectOptions,
) (*Object, ObjectInfo, error) {
	obj, err := c.minioClient.GetObject(ctx, bucketName, fileName, opts)
	if err != nil {
		err = realError(err)

		return nil, ObjectInfo{}, err
	}

	stat, err := obj.Stat()
	if err != nil {
		err = realError(err)

		return nil, ObjectInfo{}, err
	}

	return obj, stat, nil
}

func (c *Client) Delete(ctx context.Context, bucketName, fileName string, opts RemoveObjectOptions) error {
	err := c.minioClient.RemoveObject(ctx, bucketName, fileName, opts)
	if err != nil {
		return realError(err)
	}

	return nil
}

func fileNameWithPrefix(fileName string) string {
	hashRaw := sha256.Sum256([]byte(fileName))
	hash := hex.EncodeToString(hashRaw[:])
	return hash[0:2] + "/" + hash[2:4] + "/" + hash[4:6] + "/" + fileName
}
