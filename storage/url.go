package storage

import (
	"context"
	"net/url"
	"strings"
	"time"
)

// PresignGet returns a temporary, signed URL that grants read access to the
// object. If expiry <= 0 the client's configured PresignExpiry is used.
func (c *Client) PresignGet(ctx context.Context, key string, expiry time.Duration, bucket ...string) (string, error) {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return "", err
	}
	if expiry <= 0 {
		expiry = c.cfg.PresignExpiry
	}
	u, err := c.mc.PresignedGetObject(ctx, b, key, expiry, nil)
	if err != nil {
		return "", wrapErr(err)
	}
	return u.String(), nil
}

// PresignPut returns a temporary, signed URL that grants write access, so a
// client can upload directly to storage. If expiry <= 0 the configured
// PresignExpiry is used.
func (c *Client) PresignPut(ctx context.Context, key string, expiry time.Duration, bucket ...string) (string, error) {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return "", err
	}
	if expiry <= 0 {
		expiry = c.cfg.PresignExpiry
	}
	u, err := c.mc.PresignedPutObject(ctx, b, key, expiry)
	if err != nil {
		return "", wrapErr(err)
	}
	return u.String(), nil
}

// PublicURL builds an unsigned URL for an object. Use it for public buckets or
// objects fronted by a CDN. If Config.PublicBaseURL is set it is used as the
// base; otherwise the endpoint + bucket are used.
//
// Note: the URL only works if the object is actually publicly readable (see
// SetPublicReadPolicy) or served through a public CDN.
func (c *Client) PublicURL(key string, bucket ...string) (string, error) {
	b, err := c.resolveBucket(bucket)
	if err != nil {
		return "", err
	}
	key = strings.TrimLeft(key, "/")
	if c.cfg.PublicBaseURL != "" {
		return strings.TrimRight(c.cfg.PublicBaseURL, "/") + "/" + key, nil
	}
	ep := c.mc.EndpointURL() // *url.URL, never nil for a valid client
	u := &url.URL{
		Scheme: ep.Scheme,
		Host:   ep.Host,
		Path:   "/" + b + "/" + key,
	}
	return u.String(), nil
}
