---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_ruleset_rule"
sidebar_current: "docs-pagerduty-resource-ruleset-rule"
description: |-
  Creates and manages a ruleset rule in PagerDuty.
---

# pagerduty\_event_rule

An [event rule](https://support.pagerduty.com/docs/rulesets#section-create-event-rules) allows you to set actions that should be taken on events that meet your designated rule criteria.

## Example Usage

```hcl
resource "pagerduty_team" "foo" {
    name = "Engineering (Seattle)"
}

resource "pagerduty_ruleset" "foo" {
	name = "Primary Ruleset"
	team { 
		id = pagerduty_team.foo.id
	}
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
			parameter {
				value = "db"
				path = "payload.source"
			}
		}
	}
	action {
		action = "route"
		parameters {
			value = "P5DTL0K"
		}
	}
	action {
		action = "severity"
		parameters {
			value = "warning"
		}
	}
	action {
		action = "extract"
		parameters {
			target = "dedup_key"
			source = "details.host"
			regex = "(.*)"
		}
	}
}
```

## Argument Reference

The following arguments are supported:

* `ruleset` - (Required) The ID of the ruleset that the rule belongs to.
* `conditions` - (Required) Conditions evaluated to check if an event matches this event rule. Is always empty for the catch all rule, though.
* `position` - (Optional) Position/index of the rule within the ruleset.
* `disabled` - (Optional) Indicates whether the rule is disabled and would therefore not be evaluated.
* `advanced_conditions` - (Optional) Advanced conditions evaluated to check if an event matches this event rule.
* `action` - (Optional) Actions to apply to an event if the conditions match.

### Conditions (`conditions`) supports the following:
* `operator` - Operator to combine sub-conditions. Can be `and` or `or`.
* `subconditions` - List of sub-conditions that define the the condition. 

### Sub-Conditions (`subconditions`) supports the following:
* `operator` - Type of operator to apply to the sub-condition. Can be `exists`,`nexists`,`equals`,`nequals`,`contains`,`ncontains`,`matches`, or `nmatches`.
* `parameter` - Parameter for the sub-condition. It requires both a `path` and `value` to be set.
### Action (`action`) supports the following:
* `action` - Type of action to apply. Can be `route`, `suppress`, `priority`, `severity`, `annotate` or `extract`.
* `parameters` - Parameters for the given action. For actions such as `priority`,`route`,`severity`, and `annotate` only a single `value` field is needed. For the `extract` action, use `target`, `source` and `regex` parameter fields.  

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the rule.
  * `catch_all` - Indicates whether the rule is the last rule of the ruleset that serves as a catch-all. It has limited functionality compared to other rules.

## Import

Ruleset rules can be imported using the `id`, e.g.

```
$ terraform import pagerduty_ruleset_rule.main 19acac92-027a-4ea0-b06c-bbf516519601
```
