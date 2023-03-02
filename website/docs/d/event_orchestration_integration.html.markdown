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

data "pagerduty_event_orchestration_integration" "integration" {
  event_orchestration = "19acac92-027a-4ea0-b06c-bbf516519601"
  id = "1b49abe7-26db-4439-a715-c6d883acfb3e"
  label = "Example integration"
}

```

## Argument Reference

The following arguments are supported:

- `event_orchestration` - (Required) ID of the Event Orchestration to which this Integration belongs.
- `id` - (Optional) ID of the Integration associated with the Event Orchestration. Specify either `id` or `label`. If both are specified `id` takes precedence.
- `label` - (Optional) Name/description of the Integration associated with the Event Orchestration. Specify either `id` or `label`. If both are specified `id` takes precedence. The value of `label` is not unique and potentially there might be multiple Integrations with the same `label` value associated with the Event Orchestration.

## Attributes Reference

- `parameters`
  - `routing_key` - Routing key that routes to this Orchestration.
  - `type` - Type of the routing key. `global` is the default type.

<!-- TODO: Add a link to Integration Page when API docs will be available -->

[1]: https://developer.pagerduty.com/api-reference/<event_orchestration_integration>
