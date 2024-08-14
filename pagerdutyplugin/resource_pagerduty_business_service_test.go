package pagerduty

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyBusinessService_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	description := fmt.Sprintf("tf-%s", acctest.RandString(5))
	pointOfContact := fmt.Sprintf("tf-%s", acctest.RandString(5))

	nameUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	descriptionUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	pointOfContactUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyBusinessServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceConfig(name, description, pointOfContact),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceExists("pagerduty_business_service.foo"),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "name", name),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "description", description),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "point_of_contact", pointOfContact),
					resource.TestCheckResourceAttrSet("pagerduty_business_service.foo", "self"),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "type", "business_service"),
				),
			},
			{
				Config: testAccCheckPagerDutyBusinessServiceConfig(nameUpdated, descriptionUpdated, pointOfContactUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceExists("pagerduty_business_service.foo"),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "name", nameUpdated),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "description", descriptionUpdated),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "point_of_contact", pointOfContactUpdated),
					resource.TestCheckResourceAttrSet("pagerduty_business_service.foo", "self"),
				),
			},
			// Validating that externally removed business services are detected and
			// planed for re-creation
			{
				Config: testAccCheckPagerDutyBusinessServiceConfig(nameUpdated, descriptionUpdated, pointOfContactUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccExternallyDestroyBusinessService("pagerduty_business_service.foo"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPagerDutyBusinessService_WithTeam(t *testing.T) {
	businessService := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	description := fmt.Sprintf("tf-%s", acctest.RandString(5))
	pointOfContact := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyBusinessServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceWithTeamConfig(businessService, teamName, description, pointOfContact),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceExists("pagerduty_business_service.bar"),
					resource.TestCheckResourceAttr("pagerduty_business_service.bar", "name", businessService),
					resource.TestCheckResourceAttr("pagerduty_business_service.bar", "description", description),
					resource.TestCheckResourceAttr("pagerduty_business_service.bar", "point_of_contact", pointOfContact),
					resource.TestCheckResourceAttrSet("pagerduty_business_service.bar", "self"),
				),
			},
		},
	})
}

func TestAccPagerDutyBusinessService_SDKv2Compatibility(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	description := fmt.Sprintf("tf-%s", acctest.RandString(5))
	pointOfContact := fmt.Sprintf("tf-%s", acctest.RandString(5))
	commonConfig := testAccCheckPagerDutyBusinessServiceConfig(name, description, pointOfContact)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: testAccExternalProviders(),
				Config:            commonConfig,
				Check: resource.ComposeTestCheckFunc(
					// Can't call `testAccCheckPagerDutyBusinessServiceExists` because the external
					// provider doesn't call testAccProvider's Configure method, and its client is
					// left empty.
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "name", name),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "description", description),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "point_of_contact", pointOfContact),
					resource.TestCheckResourceAttrSet("pagerduty_business_service.foo", "self"),
				),
			},
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
				Config:                   commonConfig,
				ConfigPlanChecks:         resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
		},
	})
}

func testAccCheckPagerDutyBusinessServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Business Service ID is set")
		}

		businessService, err := testAccProvider.client.GetBusinessServiceWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if businessService.ID != rs.Primary.ID {
			return fmt.Errorf("Business Service not found: %v - %v", rs.Primary.ID, businessService)
		}

		return nil
	}
}

func testAccCheckPagerDutyBusinessServiceDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_business_service" {
			continue
		}
		ctx := context.Background()
		_, err := testAccProvider.client.GetBusinessServiceWithContext(ctx, r.Primary.ID)
		if err == nil {
			return fmt.Errorf("Business service still exists")
		}

	}
	return nil
}

func testAccExternallyDestroyBusinessService(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.client
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Business Service ID is set")
		}

		if err := client.DeleteBusinessServiceWithContext(ctx, rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckPagerDutyBusinessServiceConfig(name, description, poc string) string {
	return fmt.Sprintf(`
resource "pagerduty_business_service" "foo" {
	name = "%s"
	description = "%s"
	point_of_contact = "%s"
}
`, name, description, poc)
}

func testAccCheckPagerDutyBusinessServiceWithTeamConfig(businessServiceName, teamName, description, poc string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "bar" {
	name = "%s"
}

resource "pagerduty_business_service" "bar" {
	name = "%s"
	description = "%s"
	point_of_contact = "%s"
	team = pagerduty_team.bar.id
}
`, teamName, businessServiceName, description, poc)
}
