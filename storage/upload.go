package storage

import (
	"bytes"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

// Upload streams r into the store under key. Pass size = -1 if the size is
// unknown (minio-go will use multipart uploads). opts may be nil.
func (c *Client) Upload(ctx context.Context, key string, r io.Reader, size int64, opts *UploadOptions) (*UploadResult, error) {
	if opts == nil {
		opts = &UploadOptions{}
	}
	b, err := c.bucket(opts.Bucket)
	if err != nil {
		return nil, err
	}
	put := minio.PutObjectOptions{
		ContentType:  opts.ContentType,
		UserMetadata: opts.Metadata,
	}
	if put.ContentType == "" {
		put.ContentType = "application/octet-stream"
	}
	info, err := c.mc.PutObject(ctx, b, key, r, size, put)
	if err != nil {
		return nil, wrapErr(err)
	}
	return &UploadResult{
		Bucket:      info.Bucket,
		Key:         info.Key,
		ETag:        info.ETag,
		Size:        info.Size,
		VersionID:   info.VersionID,
		ContentType: put.ContentType,
	}, nil
}

// UploadBytes uploads an in-memory byte slice.
func (c *Client) UploadBytes(ctx context.Context, key string, data []byte, opts *UploadOptions) (*UploadResult, error) {
	return c.Upload(ctx, key, bytes.NewReader(data), int64(len(data)), opts)
}

// UploadFile uploads a local file by path. The content type is auto-detected
// from the file extension when opts.ContentType is empty.
func (c *Client) UploadFile(ctx context.Context, key, filePath string, opts *UploadOptions) (*UploadResult, error) {
	if opts == nil {
		opts = &UploadOptions{}
	}
	b, err := c.bucket(opts.Bucket)
	if err != nil {
		return nil, err
	}
	put := minio.PutObjectOptions{
		ContentType:  opts.ContentType,
		UserMetadata: opts.Metadata,
	}
	info, err := c.mc.FPutObject(ctx, b, key, filePath, put)
	if err != nil {
		return nil, wrapErr(err)
	}
	return &UploadResult{
		Bucket:      info.Bucket,
		Key:         info.Key,
		ETag:        info.ETag,
		Size:        info.Size,
		VersionID:   info.VersionID,
		ContentType: put.ContentType,
	}, nil
}
