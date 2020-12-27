---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_escalation_policies"
sidebar_current: "docs-pagerduty-datasource-escalation-policies"
description: |-
  Provides information about all Escalation Policies.

  This data source can be helpful when your escalation policies are handled outside Terraform but you still want to reference them in other resources.
---

# pagerduty\_escalation_policies

Use this data source to get information about all available [escalation policies][1] that you can use for other PagerDuty resources.

## Example Usage

```hcl
data "pagerduty_escalation_policies" "all" {}

resource "pagerduty_service" "test" {
  name                    = "My Web App"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = data.pagerduty_escalation_policies.all.ids[index(data.pagerduty_escalation_policies.all.names, "Office hours")]
}
```

## Argument Reference

The data source doesn't support any arguments.

## Attributes Reference
* `ids` - The list of IDs of all escalation policies.
* `names` - The list of short names of all escalation policies.

[1]: https://v2.developer.pagerduty.com/v2/page/api-reference#!/Escalation_Policies/get_escalation_policies
