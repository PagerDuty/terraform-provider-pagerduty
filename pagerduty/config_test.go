package pagerduty

import (
	"testing"
)

// Test config with an empty token
func TestConfigEmptyToken(t *testing.T) {
	config := Config{
		Token: "",
	}

	if _, err := config.Client(); err == nil {
		t.Fatalf("expected error, but got nil")
	}
}

// Test config with invalid token but with SkipCredsValidation
func TestConfigSkipCredsValidation(t *testing.T) {
	config := Config{
		Token:               "foo",
		SkipCredsValidation: true,
	}

	if _, err := config.Client(); err != nil {
		t.Fatalf("error: expected the client to not fail: %v", err)
	}
}

// Test config with a custom ApiUrl
func TestConfigCustomApiUrl(t *testing.T) {
	config := Config{
		Token:               "foo",
		ApiUrl:              "https://api.domain.tld",
		SkipCredsValidation: true,
	}

	if _, err := config.Client(); err != nil {
		t.Fatalf("error: expected the client to not fail: %v", err)
	}
}

// Test config with a custom ApiUrl override
func TestConfigCustomApiUrlOverride(t *testing.T) {
	config := Config{
		Token:               "foo",
		ApiUrlOverride:      "https://api.domain-override.tld",
		SkipCredsValidation: true,
	}

	if _, err := config.Client(); err != nil {
		t.Fatalf("error: expected the client to not fail: %v", err)
	}
}

// Test config with a custom AppUrl
func TestConfigCustomAppUrl(t *testing.T) {
	config := Config{
		Token:               "foo",
		AppUrl:              "https://app.domain.tld",
		SkipCredsValidation: true,
	}

	if _, err := config.Client(); err != nil {
		t.Fatalf("error: expected the client to not fail: %v", err)
	}
}

func TestConfigXTerraformFunctionHeader(t *testing.T) {
	config := Config{
		Token:               "foo",
		AppUrl:              "https://app.domain.tld",
		SkipCredsValidation: true,
	}

	client, err := config.Client()
	if err != nil {
		t.Fatalf("error: expected the client to not fail: %v", err)
	}
	want := "TestConfigXTerraformFunctionHeader"
	if got := client.Config.XTerraformFunctionHeader; got != want {
		t.Fatalf("error: expected %q as value for \"x-terraform-function-header\", but got %q", want, got)
	}
}

func TestExtractCallerName(t *testing.T) {
	cases := []struct {
		skip int
		want string
	}{
		{0, "extractCallerName"},
		{1, "TestExtractCallerName"},
		{2, "tRunner"},
	}

	for _, c := range cases {
		if got := extractCallerName(c.skip); got != c.want {
			t.Fatalf("error: expected %q, but got %q", c.want, got)
		}
	}
}
