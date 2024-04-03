package pagerduty

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_event_orchestration_service_cache_variable", &resource.Sweeper{
		Name: "pagerduty_event_orchestration_service_cache_variable",
		F:    testSweepService,
	})
}

func TestAccPagerDutyEventOrchestrationServiceCacheVariable_Basic(t *testing.T) {
	svc := fmt.Sprintf("tf-service-%s", acctest.RandString(5))
	cv := "pagerduty_event_orchestration_service_cache_variable.cv_1"

	name1 := fmt.Sprintf("tf_service_cache_variable_%s", acctest.RandString(5))
	svcn1 := "svc_1"
	name2 := fmt.Sprintf("tf_service_cache_variable_updated_%s", acctest.RandString(5))
	svcn2 := "svc_2"

	config1 := `
		configuration {
  		type = "trigger_event_count"
  		ttl_seconds = 60
    }
  `
	config2 := `
		configuration {
  		type = "recent_value"
  		source = "event.summary"
  		regex = ".*"
    }
  `
	cond1 := ``
	cond2 := `
		condition {
			expression = "event.source exists"
		}
	`
	disabled1 := "false"
	disabled2 := "true"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationServiceCacheVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceCacheVariableConfig(svc, name1, svcn1, disabled1, config1, cond1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServiceCacheVariableID(cv, svcn1),
					resource.TestCheckResourceAttr(cv, "name", name1),
				),
			},
			// update name and disabled state:
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceCacheVariableConfig(svc, name2, svcn1, disabled2, config1, cond1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServiceCacheVariableID(cv, svcn1),
					resource.TestCheckResourceAttr(cv, "name", name2),
					resource.TestCheckResourceAttr(cv, "disabled", disabled2),
				),
			},
			// update config:
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceCacheVariableConfig(svc, name1, svcn1, disabled1, config2, cond1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServiceCacheVariableID(cv, svcn1),
					resource.TestCheckResourceAttr(cv, "configuration.0.type", "recent_value"),
					resource.TestCheckResourceAttr(cv, "configuration.0.source", "event.summary"),
					resource.TestCheckResourceAttr(cv, "configuration.0.regex", ".*"),
				),
			},
			// update condition:
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceCacheVariableConfig(svc, name1, svcn1, disabled1, config1, cond2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServiceCacheVariableID(cv, svcn1),
					resource.TestCheckResourceAttr(cv, "condition.0.expression", "event.source exists"),
				),
			},
			// update parent event orchestration:
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceCacheVariableConfig(svc, name1, svcn2, disabled1, config1, cond1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServiceCacheVariableID(cv, svcn2),
				),
			},
			// delete cache variable:
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceCacheVariableDeletedConfig(svc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServiceCacheVariableExistsNot(cv),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationServiceCacheVariableDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_event_orchestration_service_cache_variable" {
			continue
		}
		if _, _, err := client.EventOrchestrationCacheVariables.Get(context.Background(), pagerduty.CacheVariableTypeService, r.Primary.Attributes["event_orchestration"], r.Primary.ID); err == nil {
			return fmt.Errorf("Event Orchestration Cache Variables still exist")
		}
	}
	return nil
}

func testAccCheckPagerDutyEventOrchestrationServiceCacheVariableExistsNot(cv string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[cv]
		if ok {
			return fmt.Errorf("Event Orchestration Cache Variable is not deleted from the state: %s", cv)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationServiceCacheVariableID(cv, svcn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ir, ok := s.RootModule().Resources[cv]
		svcr, _ := s.RootModule().Resources[fmt.Sprintf("pagerduty_service.%s", svcn)]

		if !ok {
			return fmt.Errorf("Event Orchestration Cache Variable resource not found in the state: %s", cv)
		}

		sid := ir.Primary.Attributes["service"]
		id := ir.Primary.ID

		client, _ := testAccProvider.Meta().(*Config).Client()
		i, _, err := client.EventOrchestrationCacheVariables.Get(context.Background(), pagerduty.CacheVariableTypeService, sid, id)
		svc, _, _ := client.Services.Get(svcr.Primary.ID, &pagerduty.GetServiceOptions{})

		if err != nil {
			return err
		}

		if i.ID != id {
			return fmt.Errorf("Event Orchestration Cache Variable ID does not match the resource ID: %v - %v", i.ID, id)
		}

		if svc.ID != sid {
			return fmt.Errorf("Event Orchestration Cache Variable's parent ID does not match the resource service attr: %v - %v", svc.ID, sid)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationServiceCacheVariableConfig(svc, name, svcn string, disabled string, config string, cond string) string {
	return fmt.Sprintf(`
	  resource "pagerduty_user" "user" {
		  email = "user@pagerduty.com"
		  name = "test user"
		}

		resource "pagerduty_escalation_policy" "ep" {
		  name = "Test EP"
		  rule {
		    escalation_delay_in_minutes = 5
		    target {
		      type = "user_reference"
		      id = pagerduty_user.user.id
		    }
		  }
		}

		resource "pagerduty_service" "svc_1" {
		  name = "%s-1"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

		resource "pagerduty_service" "svc_2" {
		  name = "%s-2"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

		resource "pagerduty_event_orchestration_service_cache_variable" "cv_1" {
			name = "%s"
			service = pagerduty_service.%s.id
			disabled = %s

			%s

			%s
		}
	`, svc, svc, name, svcn, disabled, config, cond)
}

func testAccCheckPagerDutyEventOrchestrationServiceCacheVariableDeletedConfig(svc string) string {
	return fmt.Sprintf(`
		resource "pagerduty_user" "user" {
		  email = "user@pagerduty.com"
		  name = "test user"
		}

		resource "pagerduty_escalation_policy" "ep" {
		  name = "Test EP"
		  rule {
		    escalation_delay_in_minutes = 5
		    target {
		      type = "user_reference"
		      id = pagerduty_user.user.id
		    }
		  }
		}

		resource "pagerduty_service" "svc_1" {
		  name = "%s-1"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

		resource "pagerduty_service" "svc_2" {
		  name = "%s-2"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}
	`, svc, svc)
}
