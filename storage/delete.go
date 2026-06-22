package storage

import (
	"context"

	"github.com/minio/minio-go/v7"
)

// Delete removes a single object. Deleting a non-existent object is a no-op
// (S3 semantics) and does not return an error.
func (c *Client) Delete(ctx context.Context, key string, bucket ...string) error {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return err
	}
	return wrapErr(c.mc.RemoveObject(ctx, b, key, minio.RemoveObjectOptions{}))
}

// DeleteMany removes many objects efficiently in a single batch. It returns
// the first error encountered, if any.
func (c *Client) DeleteMany(ctx context.Context, keys []string, bucket ...string) error {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	objCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objCh)
		for _, k := range keys {
			select {
			case objCh <- minio.ObjectInfo{Key: k}:
			case <-ctx.Done():
				return
			}
		}
	}()
	for e := range c.mc.RemoveObjects(ctx, b, objCh, minio.RemoveObjectsOptions{}) {
		if e.Err != nil {
			return wrapErr(e.Err)
		}
	}
	return nil
}
