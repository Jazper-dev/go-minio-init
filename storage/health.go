package storage

import (
	"context"
	"time"
)

// Ping verifies connectivity and credentials by performing a lightweight
// ListBuckets call. It returns nil when the backend is reachable.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.mc.ListBuckets(ctx)
	return wrapErr(err)
}

// Health runs a connectivity check and returns structured info, including the
// round-trip latency. It never returns an error; inspect Health.OK / Health.Err.
func (c *Client) Health(ctx context.Context) Health {
	start := time.Now()
	err := c.Ping(ctx)
	h := Health{
		Endpoint: c.cfg.Endpoint,
		Latency:  time.Since(start),
		OK:       err == nil,
	}
	if err != nil {
		h.Err = err.Error()
	}
	return h
}
