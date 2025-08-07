package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_enablement", &resource.Sweeper{
		Name: "pagerduty_enablement",
		F:    testSweepEnablement,
		Dependencies: []string{
			"pagerduty_service",
			"pagerduty_event_orchestration",
		},
	})
}

func testSweepEnablement(region string) error {
	ctx := context.Background()
	client := testAccProvider.client

	// Get all services to sweep their enablements
	servicesResp, err := client.ListServicesWithContext(ctx, pagerduty.ListServiceOptions{})
	if err != nil {
		log.Printf("Error listing services for enablement sweeping: %s", err)
		// Don't fail sweep if we can't list services
		return nil
	}

	for _, service := range servicesResp.Services {
		// Only sweep test services
		if strings.HasPrefix(service.Name, "tf-") {
			enablements, err := client.ListServiceEnablementsWithContext(ctx, service.ID)
			if err != nil {
				log.Printf("Error listing enablements for service %s during sweep: %s", service.ID, err)
				continue
			}

			for _, enablement := range enablements {
				if enablement.Enabled {
					log.Printf("Disabling enablement %s for service %s during sweep", enablement.Feature, service.ID)
					_, err := client.UpdateServiceEnablementWithContext(ctx, service.ID, enablement.Feature, false)
					if err != nil {
						log.Printf("Error disabling enablement %s for service %s during sweep: %s", enablement.Feature, service.ID, err)
					}
				}
			}
		}
	}

	// Get all event orchestrations to sweep their enablements
	orchestrationsResp, err := client.ListOrchestrationsWithContext(ctx, pagerduty.ListOrchestrationsOptions{})
	if err != nil {
		log.Printf("Error listing event orchestrations for enablement sweeping: %s", err)
		// Don't fail sweep if we can't list orchestrations
		return nil
	}

	for _, orchestration := range orchestrationsResp.Orchestrations {
		// Only sweep test orchestrations
		if strings.HasPrefix(orchestration.Name, "tf-") {
			enablements, err := client.ListEventOrchestrationEnablementsWithContext(ctx, orchestration.ID)
			if err != nil {
				log.Printf("Error listing enablements for event orchestration %s during sweep: %s", orchestration.ID, err)
				continue
			}

			for _, enablement := range enablements {
				if enablement.Enabled {
					log.Printf("Disabling enablement %s for event orchestration %s during sweep", enablement.Feature, orchestration.ID)
					_, err := client.UpdateEventOrchestrationEnablementWithContext(ctx, orchestration.ID, enablement.Feature, false)
					if err != nil {
						log.Printf("Error disabling enablement %s for event orchestration %s during sweep: %s", enablement.Feature, orchestration.ID, err)
					}
				}
			}
		}
	}

	return nil
}

func TestAccPagerDutyEnablement_Service_Basic(t *testing.T) {
	serviceName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyEnablementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEnablementServiceConfig(serviceName, username, email, escalationPolicy, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEnablementExists("pagerduty_enablement.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "entity_type", "service"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "feature", "aiops"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_enablement.test", "entity_id"),
				),
			},
			{
				ResourceName:      "pagerduty_enablement.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPagerDutyEnablement_Service_Update(t *testing.T) {
	serviceName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyEnablementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEnablementServiceConfig(serviceName, username, email, escalationPolicy, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEnablementExists("pagerduty_enablement.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "enabled", "false"),
				),
			},
			{
				Config: testAccCheckPagerDutyEnablementServiceConfig(serviceName, username, email, escalationPolicy, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEnablementExists("pagerduty_enablement.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "enabled", "false"),
				),
			},
			{
				Config: testAccCheckPagerDutyEnablementServiceConfig(serviceName, username, email, escalationPolicy, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEnablementExists("pagerduty_enablement.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccPagerDutyEnablement_EventOrchestration_Basic(t *testing.T) {
	orchestrationName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyEnablementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEnablementEventOrchestrationConfig(orchestrationName, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEnablementExists("pagerduty_enablement.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "entity_type", "event_orchestration"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "feature", "aiops"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_enablement.test", "entity_id"),
				),
			},
			{
				ResourceName:      "pagerduty_enablement.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPagerDutyEnablement_EventOrchestration_Update(t *testing.T) {
	orchestrationName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyEnablementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEnablementEventOrchestrationConfig(orchestrationName, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEnablementExists("pagerduty_enablement.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "enabled", "false"),
				),
			},
			{
				Config: testAccCheckPagerDutyEnablementEventOrchestrationConfig(orchestrationName, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEnablementExists("pagerduty_enablement.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "enabled", "false"),
				),
			},
			{
				Config: testAccCheckPagerDutyEnablementEventOrchestrationConfig(orchestrationName, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEnablementExists("pagerduty_enablement.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccPagerDutyEnablement_MultipleFeatures(t *testing.T) {
	serviceName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyEnablementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEnablementMultipleFeaturesConfig(serviceName, username, email, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEnablementExists("pagerduty_enablement.test_aiops"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test_aiops", "entity_type", "service"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test_aiops", "feature", "aiops"),
					resource.TestCheckResourceAttr(
						"pagerduty_enablement.test_aiops", "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEnablementDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := testAccProvider.client

	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_enablement" {
			continue
		}

		// Parse the resource ID which should be in format "entity_type.entity_id.feature"
		idParts := strings.Split(r.Primary.ID, ".")
		if len(idParts) != 3 {
			return fmt.Errorf("Invalid enablement ID format: %s", r.Primary.ID)
		}

		entityType := idParts[0]
		entityID := idParts[1]
		feature := idParts[2]

		// Get all enablements for the entity to check if our specific enablement still exists and is enabled
		var enablements []pagerduty.Enablement
		var err error
		switch entityType {
		case "service":
			enablements, err = client.ListServiceEnablementsWithContext(ctx, entityID)
		case "event_orchestration":
			enablements, err = client.ListEventOrchestrationEnablementsWithContext(ctx, entityID)
		default:
			err = fmt.Errorf("unsupported entity type: %s", entityType)
		}
		if err != nil {
			// If we get a 404 or similar error, the entity itself might be deleted, which is fine
			if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Error checking enablements during destroy validation for %s %s: %s", entityType, entityID, err)
		}

		// Check if our specific enablement still exists and is enabled
		for _, enablement := range enablements {
			if enablement.Feature == feature && enablement.Enabled {
				return fmt.Errorf("Enablement %s for %s %s is still enabled after destroy", feature, entityType, entityID)
			}
		}
	}

	return nil
}

func testAccCheckPagerDutyEnablementExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Enablement ID is set")
		}

		// Parse the resource ID which should be in format "entity_type.entity_id.feature"
		idParts := strings.Split(rs.Primary.ID, ".")
		if len(idParts) != 3 {
			return fmt.Errorf("Invalid enablement ID format: %s", rs.Primary.ID)
		}

		entityType := idParts[0]
		entityID := idParts[1]
		feature := idParts[2]

		client := testAccProvider.client

		// Get all enablements for the entity to check if our specific enablement exists
		var enablements []pagerduty.Enablement
		var err error
		switch entityType {
		case "service":
			enablements, err = client.ListServiceEnablementsWithContext(ctx, entityID)
		case "event_orchestration":
			enablements, err = client.ListEventOrchestrationEnablementsWithContext(ctx, entityID)
		default:
			err = fmt.Errorf("unsupported entity type: %s", entityType)
		}
		if err != nil {
			return fmt.Errorf("Error checking if enablement exists for %s %s: %s", entityType, entityID, err)
		}

		// Look for our specific enablement
		found := false
		for _, enablement := range enablements {
			if enablement.Feature == feature {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Enablement %s not found for %s %s", feature, entityType, entityID)
		}

		return nil
	}
}

func testAccCheckPagerDutyEnablementServiceConfig(serviceName, username, email, escalationPolicy, enabled string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_escalation_policy" "test" {
  name      = "%s"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.test.id
    }
  }
}

resource "pagerduty_service" "test" {
  name                    = "%s"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.test.id
}

resource "pagerduty_enablement" "test" {
  entity_type = "service"
  entity_id   = pagerduty_service.test.id
  feature     = "aiops"
  enabled     = %s
}
`, username, email, escalationPolicy, serviceName, enabled)
}

func testAccCheckPagerDutyEnablementEventOrchestrationConfig(orchestrationName, enabled string) string {
	return fmt.Sprintf(`
resource "pagerduty_event_orchestration" "test" {
  name = "%s"
}

resource "pagerduty_enablement" "test" {
  entity_type = "event_orchestration"
  entity_id   = pagerduty_event_orchestration.test.id
  feature     = "aiops"
  enabled     = %s
}
`, orchestrationName, enabled)
}

func testAccCheckPagerDutyEnablementMultipleFeaturesConfig(serviceName, username, email, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_escalation_policy" "test" {
  name      = "%s"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.test.id
    }
  }
}

resource "pagerduty_service" "test" {
  name                    = "%s"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.test.id
}

resource "pagerduty_enablement" "test_aiops" {
  entity_type = "service"
  entity_id   = pagerduty_service.test.id
  feature     = "aiops"
  enabled     = false
}
`, username, email, escalationPolicy, serviceName)
}
