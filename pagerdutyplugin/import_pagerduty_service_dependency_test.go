package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyServiceDependency_import(t *testing.T) {
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	businessService := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyBusinessServiceDependencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceDependencyConfig(service, businessService, username, email, escalationPolicy),
			},

			{
				ResourceName:      "pagerduty_service_dependency.foo",
				ImportStateIdFunc: testAccCheckPagerDutyServiceDependencyID,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				ResourceName:  "pagerduty_service_dependency.foo",
				ImportStateId: "wrongFormatID",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Expecting an importation ID formed as`),
			},
		},
	})
}

func testAccCheckPagerDutyServiceDependencyID(s *terraform.State) (string, error) {
	return fmt.Sprintf("%v.%v.%v", s.RootModule().Resources["pagerduty_business_service.foo"].Primary.ID, "business_service", s.RootModule().Resources["pagerduty_service_dependency.foo"].Primary.ID), nil
}
