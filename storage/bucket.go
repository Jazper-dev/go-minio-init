package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

// BucketExists reports whether the named bucket exists.
func (c *Client) BucketExists(ctx context.Context, name string) (bool, error) {
	if name == "" {
		return false, ErrEmptyBucket
	}
	ok, err := c.mc.BucketExists(ctx, name)
	if err != nil {
		return false, wrapErr(err)
	}
	return ok, nil
}

// MakeBucket creates a bucket using the client's configured region.
func (c *Client) MakeBucket(ctx context.Context, name string) error {
	if name == "" {
		return ErrEmptyBucket
	}
	return wrapErr(c.mc.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: c.cfg.Region}))
}

// RemoveBucket deletes an empty bucket.
func (c *Client) RemoveBucket(ctx context.Context, name string) error {
	if name == "" {
		return ErrEmptyBucket
	}
	return wrapErr(c.mc.RemoveBucket(ctx, name))
}

// EnsureBucket creates the bucket if it does not already exist. Without an
// override it operates on the client's default bucket. It is safe to call on
// every startup.
func (c *Client) EnsureBucket(ctx context.Context, bucket ...string) error {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return err
	}
	exists, err := c.mc.BucketExists(ctx, b)
	if err != nil {
		return wrapErr(err)
	}
	if exists {
		return nil
	}
	return wrapErr(c.mc.MakeBucket(ctx, b, minio.MakeBucketOptions{Region: c.cfg.Region}))
}

// SetPublicReadPolicy applies an anonymous read-only policy to the bucket so
// objects can be fetched without signing (combine with PublicURL).
func (c *Client) SetPublicReadPolicy(ctx context.Context, bucket ...string) error {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return err
	}
	policy := fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {"AWS": ["*"]},
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::%s/*"]
    }
  ]
}`, b)
	return wrapErr(c.mc.SetBucketPolicy(ctx, b, policy))
}
