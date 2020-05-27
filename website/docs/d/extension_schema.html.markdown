---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_extension_schema"
sidebar_current: "docs-pagerduty-datasource-extension-schema"
description: |-
  Get information about an extension vendor that you can use for a service (e.g: Slack, Generic Webhook, ServiceNow).
---

# pagerduty\_extension\_schema

Use this data source to get information about a specific [extension][1] vendor that you can use for a service (e.g: Slack, Generic Webhook, ServiceNow).

## Example Usage

```hcl
data "pagerduty_extension_schema" "webhook" {
  name = "Generic V2 Webhook"
}

resource "pagerduty_user" "example" {
  name  = "Howard James"
  email = "howard.james@example.domain"
  teams = [pagerduty_team.example.id]
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "Engineering Escalation Policy"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user"
      id   = pagerduty_user.example.id
    }
  }
}

resource "pagerduty_service" "example" {
  name                    = "My Web App"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.example.id
}


resource "pagerduty_extension" "slack"{
  name = "My Web App Extension"
  endpoint_url = "https://generic_webhook_url/XXXXXX/BBBBBB"
  extension_schema = data.pagerduty_extension_schema.webhook.id
  extension_objects    = [pagerduty_service.example.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The extension name to use to find an extension vendor in the PagerDuty API.

## Attributes Reference
* `id` - The ID of the found extension vendor.
* `name` - The short name of the found extension vendor.
* `type` - The generic service type for this extension vendor.

[1]: https://v2.developer.pagerduty.com/v2/page/api-reference#!/Extension_Schemas/get_extension_schemas
