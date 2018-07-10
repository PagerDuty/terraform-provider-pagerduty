---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_extension"
sidebar_current: "docs-pagerduty-resource-extension"
description: |-
  Creates and manages a service extension in PagerDuty.
---

# pagerduty\_extension

An [extension](https://v2.developer.pagerduty.com/v2/page/api-reference#!/Extensions/post_extensions) can be associated with a service.

## Example Usage

```hcl
data "pagerduty_extension_schema" "webhook" {
  name = "Generic V2 Webhook"
}

resource "pagerduty_user" "example" {
  name  = "Howard James"
  email = "howard.james@example.domain"
  teams = ["${pagerduty_team.example.id}"]
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "Engineering Escalation Policy"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user"
      id   = "${pagerduty_user.example.id}"
    }
  }
}

resource "pagerduty_service" "example" {
  name                    = "My Web App"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = "${pagerduty_escalation_policy.example.id}"
}


resource "pagerduty_extension" "slack"{
  name = "My Web App Extension"
  endpoint_url = "https://generic_webhook_url/XXXXXX/BBBBBB"
  extension_schema = "${data.pagerduty_extension_schema.webhook.id}"
  extension_objects    = ["${pagerduty_service.example.id}"]
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Optional) The name of the service extension.
  * `endpoint_url` - (Optional) The url of the extension.
  * `extension_schema` - (Required) This is the schema for this extension.
  * `extension_objects` - (Required) This is the objects for which the extension applies (An array of service ids).

    **Note:** You can use the `pagerduty_extension_schema` data source to locate the appropriate extension vendor ID.
## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the extension.
  * `html_url` - a URL at which the entity is uniquely displayed in the Web app

## Import

Extensions can be imported using the id.e.g.

```
$ terraform import pagerduty_extension.main PLBP09X
```

