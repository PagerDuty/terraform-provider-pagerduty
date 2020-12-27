---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_vendors"
sidebar_current: "docs-pagerduty-datasource-vendors"
description: |-
  Get information about all vendors that you can use for service integrations (e.g Amazon Cloudwatch, Splunk, Datadog).
---

# pagerduty\_vendor

Use this data source to get information about all [vendors][1] that you can use for service integrations (e.g Amazon Cloudwatch, Splunk, Datadog).

## Example Usage

```hcl
data "pagerduty_vendors" "all" {}

resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
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

resource "pagerduty_service_integration" "example" {
  name    = "Datadog Integration"
  vendor  = data.pagerduty_vendors.all.ids[index(data.pagerduty_vendors.all.names, "Datadog")]
  service = pagerduty_service.example.id
}
```

## Argument Reference

The data source doesn't support any arguments.

## Attributes Reference
* `ids` - The list of IDs of all escalation policies.
* `names` - The list of short names of all escalation policies.
* `types` - The list of generic service types.

[1]: https://v2.developer.pagerduty.com/v2/page/api-reference#!/Vendors/get_vendors
