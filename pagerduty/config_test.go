package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func TestAccXTerraformFunctionCustomHeader_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	usernameSpaces := " " + username + " "
	email := fmt.Sprintf("%s@foo.test", username)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserConfig(usernameSpaces, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserExists("pagerduty_user.foo"),
					testAccCheckXTerraformFunctionCustomHeaderUserCreate("pagerduty_user.foo"),
				),
			},
		},
	})
}

func testAccCheckXTerraformFunctionCustomHeaderUserCreate(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No user ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		want := "resourcePagerDutyUserCreate"
		if got := client.Config.XTerraformFunctionHeader; got != want {
			return fmt.Errorf("error: expected %q as value for \"client.Config.XTerraformFunctionHeader\", but got %q", want, got)
		}

		return nil
	}
}
