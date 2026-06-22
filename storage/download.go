package storage

import (
	"context"
	"errors"
	"io"

	"github.com/minio/minio-go/v7"
)

// Get returns a reader for the object plus its metadata. The caller MUST close
// the returned ReadCloser. An optional bucket override may be supplied.
func (c *Client) Get(ctx context.Context, key string, bucket ...string) (io.ReadCloser, *ObjectInfo, error) {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return nil, nil, err
	}
	obj, err := c.mc.GetObject(ctx, b, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, wrapErr(err)
	}
	// GetObject is lazy; Stat forces the request so a missing object errors here.
	st, err := obj.Stat()
	if err != nil {
		_ = obj.Close()
		return nil, nil, wrapErr(err)
	}
	return obj, objInfo(b, st), nil
}

// GetBytes reads the whole object into memory.
func (c *Client) GetBytes(ctx context.Context, key string, bucket ...string) ([]byte, *ObjectInfo, error) {
	rc, info, err := c.Get(ctx, key, bucket...)
	if err != nil {
		return nil, nil, err
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, nil, err
	}
	return data, info, nil
}

// Download writes the object to a local file path.
func (c *Client) Download(ctx context.Context, key, filePath string, bucket ...string) error {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return err
	}
	return wrapErr(c.mc.FGetObject(ctx, b, key, filePath, minio.GetObjectOptions{}))
}

// Stat returns metadata for an object without downloading it.
func (c *Client) Stat(ctx context.Context, key string, bucket ...string) (*ObjectInfo, error) {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return nil, err
	}
	st, err := c.mc.StatObject(ctx, b, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, wrapErr(err)
	}
	return objInfo(b, st), nil
}

// Exists reports whether an object exists.
func (c *Client) Exists(ctx context.Context, key string, bucket ...string) (bool, error) {
	_, err := c.Stat(ctx, key, bucket...)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ListObjects lists objects under prefix. Set recursive to true to descend
// into "directories" (prefixes).
func (c *Client) ListObjects(ctx context.Context, prefix string, recursive bool, bucket ...string) ([]ObjectInfo, error) {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return nil, err
	}
	var out []ObjectInfo
	for o := range c.mc.ListObjects(ctx, b, minio.ListObjectsOptions{Prefix: prefix, Recursive: recursive}) {
		if o.Err != nil {
			return nil, wrapErr(o.Err)
		}
		out = append(out, *objInfo(b, o))
	}
	return out, nil
}
