---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_priority"
sidebar_current: "docs-pagerduty-datasource-priority"
description: |-
  Get information about a priority that you can use with ruleset_rules, etc.
---

# pagerduty\_priority

Use this data source to get information about a specific [priority][1] that you can use for other PagerDuty resources. A priority is a label representing the importance and impact of an incident. This feature is only available on Standard and Enterprise plans.

## Example Usage

```hcl
data "pagerduty_priority" "p1" {
  name = "P1"
}

resource "pagerduty_ruleset" "foo" {
  name = "Primary Ruleset"
}
resource "pagerduty_ruleset_rule" "foo" {
  ruleset = pagerduty_ruleset.foo.id
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
	priority {
      value = data.pagerduty_priority.p1.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the priority to find in the PagerDuty API.

## Attributes Reference
* `id` - The ID of the found priority.
* `name` - The name of the found priority.
* `description` - A description of the found priority.

[1]: https://developer.pagerduty.com/api-reference/reference/REST/openapiv3.json/paths/~1priorities/get
