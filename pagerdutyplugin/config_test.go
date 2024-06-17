package pagerduty

import (
	"context"
	"testing"
)

// Test config with an empty token
func TestConfigEmptyToken(t *testing.T) {
	config := Config{
		Token: "",
	}

	if _, err := config.Client(context.Background()); err == nil {
		t.Fatalf("expected error, but got nil")
	}
}

// Test config with invalid token but with SkipCredsValidation
func TestConfigSkipCredsValidation(t *testing.T) {
	config := Config{
		Token:               "foo",
		SkipCredsValidation: true,
	}

	if _, err := config.Client(context.Background()); err != nil {
		t.Fatalf("error: expected the client to not fail: %v", err)
	}
}

// Test config with a custom ApiUrl
func TestConfigCustomApiUrl(t *testing.T) {
	config := Config{
		Token:               "foo",
		APIURL:              "https://api.domain.tld",
		SkipCredsValidation: true,
	}

	if _, err := config.Client(context.Background()); err != nil {
		t.Fatalf("error: expected the client to not fail: %v", err)
	}
}

// Test config with a custom ApiUrl override
func TestConfigCustomApiUrlOverride(t *testing.T) {
	config := Config{
		Token:               "foo",
		APIURLOverride:      "https://api.domain-override.tld",
		SkipCredsValidation: true,
	}

	if _, err := config.Client(context.Background()); err != nil {
		t.Fatalf("error: expected the client to not fail: %v", err)
	}
}

// Test config with a custom AppUrl
func TestConfigCustomAppUrl(t *testing.T) {
	config := Config{
		Token:               "foo",
		AppURL:              "https://app.domain.tld",
		SkipCredsValidation: true,
	}

	if _, err := config.Client(context.Background()); err != nil {
		t.Fatalf("error: expected the client to not fail: %v", err)
	}
}

// Test config with InsecureTls
func TestConfigInsecureTls(t *testing.T) {
	config := Config{
		Token:               "foo",
		InsecureTls:         true,
		SkipCredsValidation: true,
	}

	if _, err := config.Client(context.Background()); err != nil {
		t.Fatalf("error: expected the client to not fail: %v", err)
	}
}
