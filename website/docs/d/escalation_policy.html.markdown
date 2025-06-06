---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_escalation_policy"
sidebar_current: "docs-pagerduty-datasource-escalation-policy"
description: |-
  Provides information about a Escalation Policy.

  This data source can be helpful when an escalation policy is handled outside Terraform, but you still want to reference it in other resources.
---

# pagerduty\_escalation_policy

Use this data source to get information about a specific [escalation policy][1] that you can use for other PagerDuty resources.

## Example Usage

```hcl
data "pagerduty_escalation_policy" "test" {
  name = "Engineering Escalation Policy"
}

resource "pagerduty_service" "test" {
  name                    = "My Web App"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = data.pagerduty_escalation_policy.test.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use to find an escalation policy in the PagerDuty API.

## Attributes Reference
* `id` - The ID of the found escalation policy.
* `name` - The short name of the found escalation policy.
* `description` - The description of the found escalation policy.
* `teams` - The IDs of the teams associated with the found escalation policy.

[1]: https://developer.pagerduty.com/api-reference/b3A6Mjc0ODEyNA-list-escalation-policies
