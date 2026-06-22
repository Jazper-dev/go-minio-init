package storage

import (
	"errors"
	"testing"
	"time"
)

func TestConfigFromEnvPrefix(t *testing.T) {
	vars := map[string]string{
		"TEST_ENDPOINT":       "localhost:9000",
		"TEST_ACCESS_KEY":     "ak",
		"TEST_SECRET_KEY":     "sk",
		"TEST_USE_SSL":        "true",
		"TEST_REGION":         "ap-southeast-1",
		"TEST_BUCKET":         "uploads",
		"TEST_PRESIGN_EXPIRY": "30m",
	}
	for k, v := range vars {
		t.Setenv(k, v)
	}

	c := ConfigFromEnvPrefix("TEST_")
	if c.Endpoint != "localhost:9000" {
		t.Errorf("Endpoint = %q", c.Endpoint)
	}
	if c.AccessKey != "ak" || c.SecretKey != "sk" {
		t.Errorf("creds = %q/%q", c.AccessKey, c.SecretKey)
	}
	if !c.UseSSL {
		t.Error("UseSSL should be true")
	}
	if c.Region != "ap-southeast-1" {
		t.Errorf("Region = %q", c.Region)
	}
	if c.Bucket != "uploads" {
		t.Errorf("Bucket = %q", c.Bucket)
	}
	if c.PresignExpiry != 30*time.Minute {
		t.Errorf("PresignExpiry = %s", c.PresignExpiry)
	}
}

func TestConfigWithDefaults(t *testing.T) {
	c := Config{}.withDefaults()
	if c.PresignExpiry != 15*time.Minute {
		t.Errorf("default PresignExpiry = %s, want 15m", c.PresignExpiry)
	}
}

func TestConfigValidate(t *testing.T) {
	if err := (Config{}).validate(); !errors.Is(err, ErrInvalidConfig) {
		t.Errorf("empty config: got %v, want ErrInvalidConfig", err)
	}
	if err := (Config{Endpoint: "x"}).validate(); !errors.Is(err, ErrInvalidConfig) {
		t.Errorf("missing creds: got %v, want ErrInvalidConfig", err)
	}
	ok := Config{Endpoint: "x", AccessKey: "a", SecretKey: "s"}
	if err := ok.validate(); err != nil {
		t.Errorf("valid config: %v", err)
	}
}
