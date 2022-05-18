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
data "pagerduty_event_orchestration" "example" {
  name = "My Global Event Orchestration"
}

resource "pagerduty_event_orchestration" "foo" {
  name = data.pagerduty_event_orchestration.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Global Event orchestration to find in the PagerDuty API.

## Attributes Reference

* `id` - The ID of the found Event Orchestration.
* `name` - The name of the found Event Orchestration.
* `integrations` - Routing keys routed to this Event Orchestration.


[1]: https://developer.pagerduty.com/api-reference/7ba0fe7bdb26a-list-event-orchestrations
