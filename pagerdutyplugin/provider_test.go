package pagerduty

import (
	"context"
	"os"
	"testing"

	pd "github.com/PagerDuty/terraform-provider-pagerduty/pagerduty"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
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

func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"pagerduty": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			upgradedSdkServer, err := tf5to6server.UpgradeServer(ctx, pd.Provider(pd.IsMuxed).GRPCProvider)
			if err != nil {
				return nil, err
			}

			muxServer, err := tf6muxserver.NewMuxServer(
				ctx,
				providerserver.NewProtocol6(testAccProvider),
				func() tfprotov6.ProviderServer { return upgradedSdkServer },
			)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}
