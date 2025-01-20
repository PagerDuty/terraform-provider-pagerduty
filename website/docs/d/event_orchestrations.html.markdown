---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestrations"
sidebar_current: "docs-pagerduty-datasource-event-orchestrations"
description: |-
  Get information about Global Event Orchestrations that you have created.
---

# pagerduty\_event_orchestrations

Use this data source to get information as a list about specific Global [Event Orchestrations][1] filtered by a Regular Expression provided.

## Example Usage
```hcl
resource "pagerduty_event_orchestration" "tf_orch_a" {
  name = "Test Event A Orchestration"
}

resource "pagerduty_event_orchestration" "tf_orch_b" {
  name = "Test Event B Orchestration"
}

data "pagerduty_event_orchestrations" "tf_my_monitor" {
  name_filter = ".*Orchestration$"
}

resource "pagerduty_event_orchestration_global_cache_variable" "cache_var" {
  event_orchestration = data.pagerduty_event_orchestrations.tf_my_monitor.event_orchestrations[0].id
  name = "recent_host"

  condition {
    expression = "event.source exists"
  }

  configuration {
    type = "recent_value"
    source = "event.source"
    regex = ".*"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name_filter` - (Required) The regex name of Global Event Orchestrations to find in the PagerDuty API.

## Attributes Reference

* `name_filter` - The regex supplied to find the list of Global Event Orchestrations
* `event_orchestrations` - The list of the Event Orchestrations with a name that matches the `name_filter` argument.
  * `id` - The ID of the found Event Orchestration.
  * `name` - The name of the found Event Orchestration.
  * `integration` - A list of integrations for the Event Orchestration.
      * `id` - ID of the integration
      * `parameters` - A single-item list containing a parameter object describing the integration
          * `routing_key` - Routing key that routes to this Orchestration.
          * `type` - Type of the routing key. `global` is the default type.


[1]: https://developer.pagerduty.com/api-reference/7ba0fe7bdb26a-list-event-orchestrations
