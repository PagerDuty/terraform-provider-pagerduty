---
layout: 'pagerduty'
page_title: 'PagerDuty: pagerduty_event_orchestration_integration'
sidebar_current: 'docs-pagerduty-resource-event-orchestration-integration'
description: |-
  Creates and manages an Integration for an Event Orchestration.
---

# pagerduty_event_orchestration_integration

An Event Orchestration Integration allows you to create and manage multiple Integrations (and Routing Keys) per Event Orchestration _and_ will allow you to move (migrate) Integrations _between_ two Event Orchestrations.

## Example of configuring an Integration for an Event Orchestration

This example shows creating `Event Orchestration` and `Team` resources followed by creating an Event Orchestration Integration to handle Events sent to that Event Orchestration.

-> When a new Event Orchestration is created there will be one Integration (and Routing Key) included by default. Example below shows how to create an extra Integration associated with this Event Orchestration.

```hcl
resource "pagerduty_team" "database_team" {
  name = "Database Team"
}

resource "pagerduty_event_orchestration" "event_orchestration" {
  name = "Example Orchestration"
  team = pagerduty_team.database_team.id
}

resource "pagerduty_event_orchestration_integration" "integration" {
  event_orchestration = pagerduty_event_orchestration.event_orchestration.id
  label = "Example integration"
}
```

## Argument Reference

-> Modifying `event_orchestration` property will cause Integration migration process and as a result all future events sent with this Integrations's Routing Key will be evaluated against the new Event Orchestration.

The following arguments are supported:

- `event_orchestration` - (Required) ID of the Event Orchestration to which this Integration belongs to. If value is changed, current Integration is associated with a newly provided ID.
- `label` - (Required) Name/description of the Integration.

## Attributes Reference

The following attributes are exported:

- `id` - ID of this Integration.
- `parameters`
  - `routing_key` - Routing key that routes to this Orchestration.
  - `type` - Type of the routing key. `global` is the default type.

## Import

Event Orchestration Integration can be imported using colon-separated IDs, which is the combination of the Event Orchestration ID followed by the Event Orchestration Integration ID, e.g.

```
$ terraform import pagerduty_event_orchestration_integration.integration 19acac92-027a-4ea0-b06c-bbf516519601:1b49abe7-26db-4439-a715-c6d883acfb3e
```
