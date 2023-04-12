---
layout: 'pagerduty'
page_title: 'PagerDuty: pagerduty_event_orchestration_integration'
sidebar_current: 'docs-pagerduty-datasource-event-orchestration-integration'
description: |-
  Get information about an Integration for an Event Orchestration that you have created.
---

# pagerduty_event_orchestration_integration

Use this data source to get information about a specific [Integration][1] for an Event Orchestration.

## Example Usage

```hcl

resource "pagerduty_event_orchestration" "event_orchestration" {
  name = "Test Event Orchestration"
}

data "pagerduty_event_orchestration_integration" "integration" {
  event_orchestration = pagerduty_event_orchestration.event_orchestration.id
  label = "Test Event Orchestration Default Integration"
}

```

## Argument Reference

The following arguments are supported:

- `event_orchestration` - (Required) ID of the Event Orchestration to which this Integration belongs.
- `id` - (Optional) ID of the Integration associated with the Event Orchestration. Specify either `id` or `label`. If both are specified `id` takes precedence.
- `label` - (Optional) Name/description of the Integration associated with the Event Orchestration. Specify either `id` or `label`. If both are specified `id` takes precedence. The value of `label` is not unique. Potentially there might be multiple Integrations with the same `label` value associated with the Event Orchestration and retrieving data by `label` attribute will result in an error during the planning step.

## Attributes Reference

- `parameters`
  - `routing_key` - Routing key that routes to this Orchestration.
  - `type` - Type of the routing key. `global` is the default type.

[1]: https://developer.pagerduty.com/api-reference/1c6607db389a8-get-an-integration-for-an-event-orchestration
