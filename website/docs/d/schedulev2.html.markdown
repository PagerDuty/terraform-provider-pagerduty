---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_schedulev2"
sidebar_current: "docs-pagerduty-datasource-schedulev2"
description: |-
  Provides information about a PagerDuty v3 schedule.
---

# pagerduty\_schedulev2

Use this data source to look up a specific [v3 schedule](https://developer.pagerduty.com/api-reference/d90c4c94e3ce2-create-a-schedule) by name so you can reference its ID in other resources such as escalation policies.

~> **Note:** This data source requires the `flexible-schedules-early-access` early access flag on your PagerDuty account. The required `X-Early-Access` header is sent automatically by the provider.

## Example Usage

```hcl
data "pagerduty_schedulev2" "oncall" {
  name = "Engineering On-Call"
}

resource "pagerduty_escalation_policy" "example" {
  name      = "Engineering Escalation Policy"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "schedule_reference"
      id   = data.pagerduty_schedulev2.oncall.id
    }
  }
}
```

## Argument Reference

* `name` - (Required) The exact name of the v3 schedule to look up.

## Attributes Reference

* `id` - The ID of the found schedule.
* `name` - The name of the found schedule.
