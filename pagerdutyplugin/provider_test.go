package pagerduty

import (
	"context"
	"os"
	"testing"
	"time"

	pd "github.com/PagerDuty/terraform-provider-pagerduty/pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccProvider = New()

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_PARALLEL"); v != "" {
		t.Parallel()
	}

	if v := os.Getenv("PAGERDUTY_TOKEN"); v == "" {
		t.Fatal("PAGERDUTY_TOKEN must be set for acceptance tests")
	}

	if v := os.Getenv("PAGERDUTY_USER_TOKEN"); v == "" {
		t.Fatal("PAGERDUTY_USER_TOKEN must be set for acceptance tests")
	}
}

func testAccCheckAttributes(n string, fn func(map[string]string) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes
		return fn(a)
	}
}

func testAccExternalProviders() map[string]resource.ExternalProvider {
	// Using the latest release before the introduction of
	// Terraform plugin framework
	version := "~> 3.6"
	if v := os.Getenv("PAGERDUTY_ACC_EXTERNAL_PROVIDER_VERSION"); v != "" {
		version = v
	}
	m := map[string]resource.ExternalProvider{
		"pagerduty": {Source: "pagerduty/pagerduty", VersionConstraint: version},
	}
	return m
}

func testAccProtoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"pagerduty": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				pd.Provider(pd.IsMuxed).GRPCProvider,
				providerserver.NewProtocol5(testAccProvider),
			}

			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}

// testAccTimeNow returns the current time in the given location.
// The location defaults to Europe/Dublin but can be controlled by the
// PAGERDUTY_TIME_ZONE environment variable. The location must match the
// PagerDuty account time zone or diff issues might bubble up in tests.
func testAccTimeNow() time.Time {
	name := "Europe/Dublin"
	if v := os.Getenv("PAGERDUTY_TIME_ZONE"); v != "" {
		name = v
	}
	return util.TimeNowInLoc(name)
}
