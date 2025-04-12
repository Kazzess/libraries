package minio

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/Kazzess/libraries/core/rnd"
	"github.com/stretchr/testify/assert"
)

var minioConfig = NewConfig("localhost:9000", "minio", "minio123")

func minioClient() (*Client, error) {
	return NewClient(context.Background(), minioConfig)
}

func TestClient_NotExistObject(t *testing.T) {
	ctx := context.Background()

	client, err := minioClient()
	if err != nil {
		t.Fatalf("Error creating Minio client: %v", err)
	}

	key := rnd.RandomString(10)
	bucket := rnd.RandomString(10)

	o, stat, err := client.Object(ctx, bucket, key, GetObjectOptions{})
	assert.Error(t, err, ErrBucketDoesNotExist, "Expected error")
	assert.Nil(t, o, "Object should not be nil")
	assert.True(t, stat.Size == 0, "Size should be 0")
	assert.True(t, stat.Key == "", "Key should be empty")
}

func TestClient_UploadAndGet(t *testing.T) {
	ctx := context.Background()

	client, err := minioClient()
	if err != nil {
		t.Fatalf("Error creating Minio client: %v", err)
	}

	fileName := rnd.RandomString(20)
	bucketName := "bucket-" + strconv.Itoa(rnd.RandInt(0, 1000))

	fileContent := []byte(rnd.RandomString(10003))
	reader := ioutil.NopCloser(bytes.NewReader(fileContent))
	defer reader.Close()

	object := UploadObject{
		Name:        fileName,
		Size:        int64(len(fileContent)),
		Reader:      reader,
		ContentType: "application/octet-stream",
	}

	err = client.CreateBucket(ctx, bucketName, MakeBucketOptions{})
	assert.NoError(t, err, "CreateBucket should not return error")

	key, bucket, err := client.Upload(ctx, object, bucketName)
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, key, "Key should not be nil")
	assert.NotNil(t, bucket, "Bucket should not be nil")

	res, err := client.ListBucket(ctx, bucketName, ListObjectsOptions{})
	assert.True(t, len(res) == 1)
	assert.NoError(t, err)

	o, stat, err := client.Object(ctx, bucket, object.Name, GetObjectOptions{})
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, o, "Object should not be nil")
	assert.NoError(t, err, "Expected no error statting object")
	assert.NotNil(t, stat, "Stat should not be nil")
	assert.NoError(t, stat.Err, "Stat should not contain error")
	assert.Equal(t, key, stat.Key, "Key should match")

	o, stat, err = client.Object(ctx, bucket, fileName, GetObjectOptions{})
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, o, "Object should not be nil")
	assert.NoError(t, err, "Expected no error statting object")
	assert.NotNil(t, stat, "Stat should not be nil")
	assert.NoError(t, stat.Err, "Stat should not contain error")
	assert.Equal(t, key, stat.Key, "Key should match")

	err = client.Delete(ctx, bucket, key, RemoveObjectOptions{})
	assert.NoError(t, err, "Delete should not return error")

	err = client.RemoveBucket(ctx, bucket)
	assert.NoError(t, err, "RemoveBucket should not return error")

}

func TestClient_UploadExistBucket(t *testing.T) {
	ctx := context.Background()

	client, err := minioClient()
	if err != nil {
		t.Fatalf("Error creating Minio client: %v", err)
	}

	fileName := rnd.RandomString(10)
	bucketName := "bucket-" + strconv.Itoa(rnd.RandInt(0, 1000))

	fileContent := []byte(rnd.RandomString(1023))
	reader := ioutil.NopCloser(bytes.NewReader(fileContent))
	defer reader.Close()

	object := UploadObject{
		Name:        fileName,
		Size:        int64(len(fileContent)),
		Reader:      reader,
		ContentType: "application/octet-stream",
	}

	err = client.CreateBucket(ctx, bucketName, MakeBucketOptions{})
	assert.NoError(t, err, "CreateBucket should not return error")

	// Test successful upload to a custom bucket
	key, bucket, err := client.Upload(ctx, object, bucketName)
	assert.NoError(t, err, "Upload should not return error")
	assert.NotNil(t, key, "Upload key should not be nil")
	assert.NotNil(t, bucket, "Bucket should not be nil")

	err = client.RemoveBucket(ctx, bucket)
	assert.True(t, errors.Is(err, ErrRemoveBucketNotEmpty))

	err = client.Delete(ctx, bucket, key, RemoveObjectOptions{})
	assert.NoError(t, err, "Delete should not return error")

	err = client.RemoveBucket(ctx, bucket)
	assert.NoError(t, err, "RemoveBucket should not return error")
}

func TestClient_UploadBucketDoesNotExist(t *testing.T) {
	ctx := context.Background()

	client, err := minioClient()
	if err != nil {
		t.Fatalf("Error creating Minio client: %v", err)
	}

	fileName := rnd.RandomString(10)
	bucketName := "bucket-" + strconv.Itoa(rnd.RandInt(0, 1000))

	// Create a test object to upload
	fileContent := []byte(rnd.RandomString(10124))
	reader := ioutil.NopCloser(bytes.NewReader(fileContent))
	defer reader.Close()

	object := UploadObject{
		Name:        fileName,
		Size:        int64(len(fileContent)),
		Reader:      reader,
		ContentType: "application/octet-stream",
	}

	// Test error when the target bucket does not exist
	key, bucket, err := client.Upload(ctx, object, bucketName)
	assert.Error(t, err, ErrBucketDoesNotExist)
	assert.Empty(t, key, "Key should be empty")
	assert.Empty(t, bucket, "Bucket should be empty")
}

func TestClient_DeleteBucket_NotExist(t *testing.T) {
	client, err := minioClient()
	if err != nil {
		t.Fatalf("Error creating Minio client: %v", err)
	}

	ctx := context.Background()
	bucketName := "test-bucket"

	// Test deleting a non-existent bucket
	err = client.RemoveBucket(ctx, bucketName)
	assert.True(t, errors.Is(err, ErrBucketDoesNotExist), err)
}

func TestClient_DeleteObject_NotExist(t *testing.T) {
	client, err := minioClient()
	if err != nil {
		t.Fatalf("Error creating Minio client: %v", err)
	}

	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := "nonexistent-file"
	opts := RemoveObjectOptions{}

	// Test deleting a non-existent file
	err = client.Delete(ctx, bucketName, fileName, opts)
	assert.True(t, errors.Is(err, ErrBucketDoesNotExist), err)
}
