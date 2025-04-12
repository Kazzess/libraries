package minio

import (
	"errors"

	"github.com/minio/minio-go/v7"
)

type (
	Object              = minio.Object
	ObjectInfo          = minio.ObjectInfo
	GetObjectOptions    = minio.GetObjectOptions
	ListObjectsOptions  = minio.ListObjectsOptions
	RemoveObjectOptions = minio.RemoveObjectOptions
	MakeBucketOptions   = minio.MakeBucketOptions
)

var (
	ErrBucketDoesNotExist = errors.New("bucket does not exist")
	ErrUploadNoBucketName = errors.New(
		"bucket name must be specified or specify autp create bucket option",
	)
	ErrBucketAlreadyOwnedByYou = errors.New("bucket already owned by you")
	ErrRemoveBucketNotEmpty    = errors.New("bucket not empty")
	ErrObjectDoesNotExist      = errors.New("object does not exist")
)
