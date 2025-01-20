---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestration"
sidebar_current: "docs-pagerduty-resource-event-orchestration"
description: |-
  Creates and manages an Event Orchestration in PagerDuty.
---

# pagerduty_event_orchestration

[Event Orchestrations](https://support.pagerduty.com/docs/event-orchestration) allow you define a set of Event Rules, so that when you ingest events using the Orchestration's Routing Key your events will be routed to the correct Global and/or Service Orchestration, based on the event's content.

## Example of configuring an Event Orchestration

```hcl
resource "pagerduty_team" "engineering" {
  name = "Engineering"
}

resource "pagerduty_event_orchestration" "my_monitor" {
  name = "My Monitoring Orchestration"
  description = "Send events to a pair of services"
  team = pagerduty_team.engineering.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Event Orchestration.
* `description` - (Optional) A human-friendly description of the Event Orchestration.
* `team` - (Optional) ID of the team that owns the Event Orchestration. If none is specified, only admins have access.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Event Orchestration.
* `integration` - A list of integrations for the Event Orchestration.
  * `id` - ID of the integration
  * `parameters` - A single-item list containing a parameter object describing the integration
      * `routing_key` - Routing key that routes to this Orchestration.
      * `type` - Type of the routing key. `global` is the default type.

## Import

EventOrchestrations can be imported using the `id`, e.g.

```
$ terraform import pagerduty_event_orchestration.main 19acac92-027a-4ea0-b06c-bbf516519601
```
