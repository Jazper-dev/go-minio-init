package storage

import (
	"errors"
	"testing"
)

func testClient(t *testing.T, cfg Config) *Client {
	t.Helper()
	cfg.Endpoint = "localhost:9000"
	cfg.AccessKey = "ak"
	cfg.SecretKey = "sk"
	c, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func TestPublicURLFromEndpoint(t *testing.T) {
	c := testClient(t, Config{Bucket: "uploads"})
	got, err := c.PublicURL("a/b.png")
	if err != nil {
		t.Fatalf("PublicURL: %v", err)
	}
	want := "http://localhost:9000/uploads/a/b.png"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPublicURLFromBaseURL(t *testing.T) {
	c := testClient(t, Config{Bucket: "uploads", PublicBaseURL: "https://cdn.example.com/"})
	got, err := c.PublicURL("/a/b.png") // leading slash should be trimmed
	if err != nil {
		t.Fatalf("PublicURL: %v", err)
	}
	want := "https://cdn.example.com/a/b.png"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPublicURLBucketOverride(t *testing.T) {
	c := testClient(t, Config{Bucket: "uploads"})
	got, _ := c.PublicURL("x.png", "other")
	want := "http://localhost:9000/other/x.png"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPublicURLNoBucket(t *testing.T) {
	c := testClient(t, Config{}) // no default bucket, no override
	if _, err := c.PublicURL("x.png"); !errors.Is(err, ErrEmptyBucket) {
		t.Errorf("got %v, want ErrEmptyBucket", err)
	}
}
