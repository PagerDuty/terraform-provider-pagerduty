---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestrations"
sidebar_current: "docs-pagerduty-datasource-event-orchestrations"
description: |-
  Get information about Global Event Orchestrations that you have created.
---

# pagerduty\_event_orchestrations

Use this data source to get informations about specific Global [Event Orchestrations][1]

## Example Usage
```hcl
resource "pagerduty_event_orchestration" "tf_orch_a" {
  name = "Test Event A Orchestration"
}

resource "pagerduty_event_orchestration" "tf_orch_b" {
  name = "Test Event B Orchestration"
}

data "pagerduty_event_orchestrations" "tf_my_monitor" {
  search = ".*Orchestration$"
}

```

## Argument Reference

The following arguments are supported:

* `search` - (Required) The regex name of Global Event orchestrations to find in the PagerDuty API.

## Attributes Reference

* `event_orchestrations` - The list of the Event Orchestrations which name match `search` argument.
  * `id` - The ID of the found Event Orchestration.
  * `name` - The name of the found Event Orchestration.
  * `integration` - An integration for the Event Orchestration.
    * `id` - ID of the integration
    * `parameters`
      * `routing_key` - Routing key that routes to this Orchestration.
      * `type` - Type of the routing key. `global` is the default type.


[1]: https://developer.pagerduty.com/api-reference/7ba0fe7bdb26a-list-event-orchestrations
