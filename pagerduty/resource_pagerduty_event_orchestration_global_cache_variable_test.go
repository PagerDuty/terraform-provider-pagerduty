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
	resource.AddTestSweepers("pagerduty_event_orchestration_global_cache_variable", &resource.Sweeper{
		Name: "pagerduty_event_orchestration_global_cache_variable",
		F:    testSweepEventOrchestration,
	})
}

func TestAccPagerDutyEventOrchestrationGlobalCacheVariable_Basic(t *testing.T) {
	orch := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))
	cv := "pagerduty_event_orchestration_global_cache_variable.cv_1"

	name1 := fmt.Sprintf("tf_global_cache_variable_%s", acctest.RandString(5))
	orchn1 := "orch_1"
	name2 := fmt.Sprintf("tf_global_cache_variable_updated_%s", acctest.RandString(5))
	orchn2 := "orch_2"

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
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableConfig(orch, name1, orchn1, disabled1, config1, cond1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableID(cv, orchn1),
					resource.TestCheckResourceAttr(cv, "name", name1),
				),
			},
			// update name and disabled state:
			{
				Config: testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableConfig(orch, name2, orchn1, disabled2, config1, cond1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableID(cv, orchn1),
					resource.TestCheckResourceAttr(cv, "name", name2),
					resource.TestCheckResourceAttr(cv, "disabled", disabled2),
				),
			},
			// update config:
			{
				Config: testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableConfig(orch, name1, orchn1, disabled1, config2, cond1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableID(cv, orchn1),
					resource.TestCheckResourceAttr(cv, "configuration.0.type", "recent_value"),
					resource.TestCheckResourceAttr(cv, "configuration.0.source", "event.summary"),
					resource.TestCheckResourceAttr(cv, "configuration.0.regex", ".*"),
				),
			},
			// update condition:
			{
				Config: testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableConfig(orch, name1, orchn1, disabled1, config1, cond2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableID(cv, orchn1),
					resource.TestCheckResourceAttr(cv, "condition.0.expression", "event.source exists"),
				),
			},
			// update parent event orchestration:
			{
				Config: testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableConfig(orch, name1, orchn2, disabled1, config1, cond1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableID(cv, orchn2),
				),
			},
			// delete cache variable:
			{
				Config: testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableDeletedConfig(orch),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableExistsNot(cv),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_event_orchestration_global_cache_variable" {
			continue
		}
		if _, _, err := client.EventOrchestrationCacheVariables.Get(context.Background(), pagerduty.CacheVariableTypeGlobal, r.Primary.Attributes["event_orchestration"], r.Primary.ID); err == nil {
			return fmt.Errorf("Event Orchestration Cache Variables still exist")
		}
	}
	return nil
}

func testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableExistsNot(cv string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[cv]
		if ok {
			return fmt.Errorf("Event Orchestration Cache Variable is not deleted from the state: %s", cv)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableID(cv, orchn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ir, ok := s.RootModule().Resources[cv]
		eor, _ := s.RootModule().Resources[fmt.Sprintf("pagerduty_event_orchestration.%s", orchn)]

		if !ok {
			return fmt.Errorf("Event Orchestration Cache Variable resource not found in the state: %s", cv)
		}

		oid := ir.Primary.Attributes["event_orchestration"]
		id := ir.Primary.ID

		client, _ := testAccProvider.Meta().(*Config).Client()
		i, _, err := client.EventOrchestrationCacheVariables.Get(context.Background(), pagerduty.CacheVariableTypeGlobal, oid, id)
		eo, _, _ := client.EventOrchestrations.Get(eor.Primary.ID)

		if err != nil {
			return err
		}

		if i.ID != id {
			return fmt.Errorf("Event Orchestration Cache Variable ID does not match the resource ID: %v - %v", i.ID, id)
		}

		if eo.ID != oid {
			return fmt.Errorf("Event Orchestration Cache Variable's parent ID does not match the resource event_orchestration attr: %v - %v", eo.ID, oid)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableConfig(orch, name, orchn string, disabled string, config string, cond string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch_1" {
			name = "%s-1"
		}

		resource "pagerduty_event_orchestration" "orch_2" {
			name = "%s-2"
		}

		resource "pagerduty_event_orchestration_global_cache_variable" "cv_1" {
			name = "%s"
			event_orchestration = pagerduty_event_orchestration.%s.id
			disabled = %s

			%s

			%s
		}
	`, orch, orch, name, orchn, disabled, config, cond)
}

func testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableDeletedConfig(orch string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch_1" {
			name = "%s-1"
		}

		resource "pagerduty_event_orchestration" "orch_2" {
			name = "%s-2"
		}
	`, orch, orch)
}
