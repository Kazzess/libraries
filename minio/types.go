package minio

import (
	"fmt"
	"io"
)

type ContentType string

func (c ContentType) String() string {
	return string(c)
}

const (
	OctetContentType     ContentType = "application/octet-stream"
	ImagePNGContentType  ContentType = "image/png"
	ImageWebPContentType ContentType = "image/webp"
)

type UploadObject struct {
	Name        string
	Size        int64
	Meta        map[string]string
	ContentType ContentType
	Reader      io.Reader
}

func NewUploadObject(
	name string,
	size int64,
	meta map[string]string,
	contentType ContentType,
	reader io.Reader,
) *UploadObject {
	return &UploadObject{
		Name:        name,
		Size:        size,
		Meta:        meta,
		ContentType: contentType,
		Reader:      reader,
	}
}

func (o UploadObject) String() string {
	return fmt.Sprintf("Name: %s, Size: %d, ContentType: %s", o.Name, o.Size, o.ContentType)
}
