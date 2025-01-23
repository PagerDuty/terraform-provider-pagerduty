---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestration"
sidebar_current: "docs-pagerduty-datasource-event-orchestration"
description: |-
  Get information about a Global Event Orchestration that you have created.
---

# pagerduty\_event_orchestration

Use this data source to get information about a specific Global [Event Orchestration][1]

## Example Usage
```hcl
resource "pagerduty_event_orchestration" "tf_orch_a" {
  name = "Test Event Orchestration"
}

data "pagerduty_event_orchestration" "tf_my_monitor" {
  name = pagerduty_event_orchestration.tf_orch_a.name
}

resource "pagerduty_event_orchestration_unrouted" "unrouted" {
  event_orchestration = data.pagerduty_event_orchestration.tf_my_monitor.id
  catch_all {
    actions {
      severity = "info"
    }
  }
  set {
    id = "start"
    rule {
      actions {
        extraction {
          target = "event.custom_details.integration_type"
          template = data.pagerduty_event_orchestration.tf_my_monitor.integration[0].parameters[0].type
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Global Event Orchestration to find in the PagerDuty API.

## Attributes Reference

* `id` - The ID of the found Event Orchestration.
* `name` - The name of the found Event Orchestration.
* `integration` - A list of integrations for the Event Orchestration.
  * `id` - ID of the integration
  * `parameters` - A single-item list containing a parameter object describing the integration
      * `routing_key` - Routing key that routes to this Orchestration.
      * `type` - Type of the routing key. `global` is the default type.


[1]: https://developer.pagerduty.com/api-reference/7ba0fe7bdb26a-list-event-orchestrations
