package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariable_Basic(t *testing.T) {
	on := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))
	name := fmt.Sprintf("tf_global_cache_variable_%s", acctest.RandString(5))
	irn := "pagerduty_event_orchestration_global_cache_variable.orch_cv"
	n := "data.pagerduty_event_orchestration_global_cache_variable.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// find by id
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableByIdConfig(on, name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariable(irn, n),
				),
			},
			// find by id, ignore name
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableByIdNameConfig(on, name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariable(irn, n),
				),
			},
			// find by name
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableByNameConfig(on, name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariable(irn, n),
				),
			},
			// id and name are both not set
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableIdNameNullConfig(on, name),
				ExpectError: regexp.MustCompile("Invalid Event Orchestration Cache Variable data source configuration: ID and name cannot both be null"),
			},
			// bad event_orchestration
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableBadOrchConfig(on, name),
				ExpectError: regexp.MustCompile("Unable to find a Cache Variable with ID '(.+)' on PagerDuty Event Orchestration 'bad-orchestration-id'"),
			},
			// bad id
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableBadIdConfig(on, name),
				ExpectError: regexp.MustCompile("Unable to find a Cache Variable with ID 'bad-cache-var-id' on PagerDuty Event Orchestration '(.+)'"),
			},
			// bad name
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableBadNameConfig(on, name),
				ExpectError: regexp.MustCompile("Unable to find a Cache Variable on Event Orchestration '(.+)' with name 'bad-cache-var-name'"),
			},
		},
	})
}

func testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariable(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected the Event Orchestration Cache Variable ID to be set")
		}

		testAtts := []string{
			"id", "name", "configuration.0.type", "configuration.0.ttl_seconds",
		}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the Event Orchestration Cache Variable %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableBaseConfig(on, name string) string {
	return fmt.Sprintf(`
    resource "pagerduty_event_orchestration" "orch" {
      name = "%s"
    }

    resource "pagerduty_event_orchestration_global_cache_variable" "orch_cv" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }
    `, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableByIdConfig(on, name string) string {
	return fmt.Sprintf(`
    resource "pagerduty_event_orchestration" "orch" {
      name = "%s"
    }

    resource "pagerduty_event_orchestration_global_cache_variable" "orch_cv" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_global_cache_variable" "by_id" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      id = pagerduty_event_orchestration_global_cache_variable.orch_cv.id
    }
    `, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableByIdNameConfig(on, name string) string {
	return fmt.Sprintf(`
    resource "pagerduty_event_orchestration" "orch" {
      name = "%s"
    }

    resource "pagerduty_event_orchestration_global_cache_variable" "orch_cv" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_global_cache_variable" "by_id" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      id = pagerduty_event_orchestration_global_cache_variable.orch_cv.id
      name = "No such name"
    }
    `, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableByNameConfig(on, name string) string {
	return fmt.Sprintf(`
    resource "pagerduty_event_orchestration" "orch" {
      name = "%[1]s"
    }

    resource "pagerduty_event_orchestration_global_cache_variable" "orch_cv" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      name = "%[2]s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_global_cache_variable" "by_id" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      name = "%[2]s"
    }
    `, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableIdNameNullConfig(on, name string) string {
	return fmt.Sprintf(`
    resource "pagerduty_event_orchestration" "orch" {
      name = "%s"
    }

    resource "pagerduty_event_orchestration_global_cache_variable" "orch_cv" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_global_cache_variable" "by_id" {
      event_orchestration = pagerduty_event_orchestration.orch.id
    }
    `, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableBadOrchConfig(on, name string) string {
	return fmt.Sprintf(`
    resource "pagerduty_event_orchestration" "orch" {
      name = "%s"
    }

    resource "pagerduty_event_orchestration_global_cache_variable" "orch_cv" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_global_cache_variable" "by_id" {
      event_orchestration = "bad-orchestration-id"
      id = pagerduty_event_orchestration_global_cache_variable.orch_cv.id
    }
    `, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableBadIdConfig(on, name string) string {
	return fmt.Sprintf(`
    resource "pagerduty_event_orchestration" "orch" {
      name = "%s"
    }

    resource "pagerduty_event_orchestration_global_cache_variable" "orch_cv" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_global_cache_variable" "by_id" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      id = "bad-cache-var-id"
    }
    `, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationGlobalCacheVariableBadNameConfig(on, name string) string {
	return fmt.Sprintf(`
    resource "pagerduty_event_orchestration" "orch" {
      name = "%s"
    }

    resource "pagerduty_event_orchestration_global_cache_variable" "orch_cv" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_global_cache_variable" "by_id" {
      event_orchestration = pagerduty_event_orchestration.orch.id
      name = "bad-cache-var-name"
    }
    `, on, name)
}
