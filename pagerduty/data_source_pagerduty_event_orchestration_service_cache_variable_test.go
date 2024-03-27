package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyEventOrchestrationServiceCacheVariable_Basic(t *testing.T) {
	on := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))
	name := fmt.Sprintf("tf_service_cache_variable_%s", acctest.RandString(5))
	irn := "pagerduty_event_orchestration_service_cache_variable.orch_cv"
	n := "data.pagerduty_event_orchestration_service_cache_variable.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// find by id
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableByIdConfig(on, name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariable(irn, n),
				),
			},
			// find by id, ignore name
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableByIdNameConfig(on, name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariable(irn, n),
				),
			},
			// find by name
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableByNameConfig(on, name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariable(irn, n),
				),
			},
			// id and name are both not set
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableIdNameNullConfig(on, name),
				ExpectError: regexp.MustCompile("Invalid Event Orchestration Cache Variable data source configuration: ID and name cannot both be null"),
			},
			// bad event_orchestration
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableBadOrchConfig(on, name),
				ExpectError: regexp.MustCompile("Unable to find a Cache Variable with ID '(.+)' on PagerDuty Event Orchestration 'bad-orchestration-id'"),
			},
			// bad id
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableBadIdConfig(on, name),
				ExpectError: regexp.MustCompile("Unable to find a Cache Variable with ID 'bad-cache-var-id' on PagerDuty Event Orchestration '(.+)'"),
			},
			// bad name
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableBadNameConfig(on, name),
				ExpectError: regexp.MustCompile("Unable to find a Cache Variable on Event Orchestration '(.+)' with name 'bad-cache-var-name'"),
			},
		},
	})
}

const EPResources = `
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
`

func testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariable(src, n string) resource.TestCheckFunc {
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

func testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableBaseConfig(on, name string) string {
	return fmt.Sprintf(`
		%s

		resource "pagerduty_service" "svc" {
		  name = "%s"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

    resource "pagerduty_event_orchestration_service_cache_variable" "orch_cv" {
      service = pagerduty_service.svc.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }
    `, EPResources, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableByIdConfig(on, name string) string {
	return fmt.Sprintf(`
		%s

		resource "pagerduty_service" "svc" {
		  name = "%s"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

    resource "pagerduty_event_orchestration_service_cache_variable" "orch_cv" {
      service = pagerduty_service.svc.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_service_cache_variable" "by_id" {
      service = pagerduty_service.svc.id
      id = pagerduty_event_orchestration_service_cache_variable.orch_cv.id
    }
    `, EPResources, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableByIdNameConfig(on, name string) string {
	return fmt.Sprintf(`
		%s

		resource "pagerduty_service" "svc" {
		  name = "%s"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

    resource "pagerduty_event_orchestration_service_cache_variable" "orch_cv" {
      service = pagerduty_service.svc.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_service_cache_variable" "by_id" {
      service = pagerduty_service.svc.id
      id = pagerduty_event_orchestration_service_cache_variable.orch_cv.id
      name = "No such name"
    }
    `, EPResources, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableByNameConfig(on, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pagerduty_service" "svc" {
		  name = "%[2]s"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

    resource "pagerduty_event_orchestration_service_cache_variable" "orch_cv" {
      service = pagerduty_service.svc.id
      name = "%[3]s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_service_cache_variable" "by_id" {
      service = pagerduty_service.svc.id
      name = "%[3]s"
    }
    `, EPResources, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableIdNameNullConfig(on, name string) string {
	return fmt.Sprintf(`
		%s

		resource "pagerduty_service" "svc" {
		  name = "%s"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

    resource "pagerduty_event_orchestration_service_cache_variable" "orch_cv" {
      service = pagerduty_service.svc.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_service_cache_variable" "by_id" {
      service = pagerduty_service.svc.id
    }
    `, EPResources, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableBadOrchConfig(on, name string) string {
	return fmt.Sprintf(`
		%s

		resource "pagerduty_service" "svc" {
		  name = "%s"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

    resource "pagerduty_event_orchestration_service_cache_variable" "orch_cv" {
      service = pagerduty_service.svc.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_service_cache_variable" "by_id" {
      service = "bad-orchestration-id"
      id = pagerduty_event_orchestration_service_cache_variable.orch_cv.id
    }
    `, EPResources, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableBadIdConfig(on, name string) string {
	return fmt.Sprintf(`
		%s

		resource "pagerduty_service" "svc" {
		  name = "%s"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

    resource "pagerduty_event_orchestration_service_cache_variable" "orch_cv" {
      service = pagerduty_service.svc.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_service_cache_variable" "by_id" {
      service = pagerduty_service.svc.id
      id = "bad-cache-var-id"
    }
    `, EPResources, on, name)
}

func testAccDataSourcePagerDutyEventOrchestrationServiceCacheVariableBadNameConfig(on, name string) string {
	return fmt.Sprintf(`
  	%s

		resource "pagerduty_service" "svc" {
		  name = "%s"
		  escalation_policy = pagerduty_escalation_policy.ep.id
		}

    resource "pagerduty_event_orchestration_service_cache_variable" "orch_cv" {
      service = pagerduty_service.svc.id
      name = "%s"

      configuration {
    		type = "trigger_event_count"
    		ttl_seconds = 60
      }
    }

    data "pagerduty_event_orchestration_service_cache_variable" "by_id" {
      service = pagerduty_service.svc.id
      name = "bad-cache-var-name"
    }
    `, EPResources, on, name)
}
