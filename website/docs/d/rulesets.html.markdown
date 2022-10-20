---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_rulesets"
sidebar_current: "docs-pagerduty-datasource-rulesets"
description: |-
  Get information about multiple ruleset that you have created.
---

# pagerduty\_rulesets

Use this data source to get information about multiple [ruleset][1] that you can use for managing and grouping [event rules][2].

## Example Usage

```hcl
resource "pagerduty_ruleset" "first" {
	name = "MyFirstRuleSet"
}

resource "pagerduty_ruleset" "second" {
	name = "MySecondRuleSet"
}

data "pagerduty_rulesets" "example" {
  search = ".*RuleSet$"
}
```

## Argument Reference

The following arguments are supported:

* `search` - (Required) The ruleset regex name to find in the PagerDuty API.

## Attributes Reference

The following attributes are exported:

* `rulesets` - List of rulesets which name match `search` argument.
	* `id` - The ID of the ruleset.
	* `name` - Name of the ruleset.
	* `routing_keys` - Routing keys routed to this ruleset.

[1]: https://developer.pagerduty.com/api-reference/b3A6Mjc0ODE3MQ-list-rulesets
