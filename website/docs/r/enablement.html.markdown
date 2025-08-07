---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_enablement"
sidebar_current: "docs-pagerduty-resource-enablement"
description: |-
  Creates and manages feature enablements for PagerDuty entities.
---

# pagerduty\_enablement

[Enablements](https://developer.pagerduty.com/api-reference/b3A6Mjc0ODE5Nw-list-enablements) allow you to enable or disable specific features for PagerDuty entities such as services and event orchestrations.

## Example Usage

```hcl
data "pagerduty_service" "example" {
  name = "My Web Service"
}

resource "pagerduty_enablement" "example" {
  entity_type = "service"
  entity_id   = data.pagerduty_service.example.id
  feature     = "aiops"
  enabled     = true
}
```

## Argument Reference

The following arguments are supported:

* `entity_type` - (Required) The type of entity for which to manage the enablement. Possible values can be `service` and `event_orchestration`.
* `entity_id` - (Required) The ID of the entity for which to manage the enablement.
* `feature` - (Required) The name of the feature to enable or disable. Possible values can be `aiops`.
* `enabled` - (Required) Whether the feature should be enabled (`true`) or disabled (`false`) for the specified entity.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the enablement, constructed as `entity_type.entity_id.feature`.

## Import

Enablements can be imported using the `id`, which is constructed by concatenating the `entity_type`, `entity_id`, and `feature` with dots, e.g.

```
$ terraform import pagerduty_enablement.example service.P7HHMVK.aiops
```