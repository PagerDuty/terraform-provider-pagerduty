---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_ruleset"
sidebar_current: "docs-pagerduty-datasource-ruleset"
description: |-
  Get information about a ruleset that you have created.
---

# pagerduty\_ruleset

Use this data source to get information about a specific [ruleset][1] that you can use for managing and grouping [event rules][2].

## Example Usage

```hcl
data "pagerduty_ruleset" "example" {
  name = "My Ruleset"
}

resource "pagerduty_ruleset_rule" "foo" {
  ruleset = data.pagerduty_ruleset.example.id
  position = 0
  disabled = "false"
  conditions {
    operator = "and"
	subconditions {
	  operator = "contains"
	  parameter {
	    value = "disk space"
		path = "payload.summary"
	  }
	}
	subconditions {
	  operator = "contains"
	  parameter {
	    value = "db"
	    path = "payload.source"
	  }
	}
  }
  actions {
    route {
	  value = "P5DTL0K"
	}
  }
}
```

### Default Global Ruleset

```hcl
data "pagerduty_ruleset" "default_global" {
  name = "Default Global"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ruleset to find in the PagerDuty API.

## Attributes Reference
* `id` - The ID of the found ruleset.
* `name` - The name of the found ruleset.

[1]: https://developer.pagerduty.com/api-reference/reference/REST/openapiv3.json/paths/~1rulesets/get
[2]: https://developer.pagerduty.com/api-reference/reference/REST/openapiv3.json/paths/~1rulesets~1%7Bid%7D~1rules/get
