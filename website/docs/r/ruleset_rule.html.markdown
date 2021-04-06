---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_ruleset_rule"
sidebar_current: "docs-pagerduty-resource-ruleset-rule"
description: |-
  Creates and manages a ruleset rule in PagerDuty.
---

# pagerduty\_ruleset_rule

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
  time_frame {
    scheduled_weekly {
	  weekdays = [3,7]
	  timezone = "America/Los_Angeles"
	  start_time = "1000000"
	  duration = "3600000"
	}
  }
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
	severity  {
	  value = "warning"
	}
	annotate {
	  value = "From Terraform"
	}
	extractions {
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
* `time_frame` - (Optional) Settings for [scheduling the rule](https://support.pagerduty.com/docs/rulesets#section-scheduled-event-rules). 
* `actions` - (Optional) Actions to apply to an event if the conditions match.

### Conditions (`conditions`) supports the following:
* `operator` - Operator to combine sub-conditions. Can be `and` or `or`.
* `subconditions` - List of sub-conditions that define the the condition. 

### Sub-Conditions (`subconditions`) supports the following:
* `operator` - Type of operator to apply to the sub-condition. Can be `exists`,`nexists`,`equals`,`nequals`,`contains`,`ncontains`,`matches`, or `nmatches`.
* `parameter` - Parameter for the sub-condition. It requires both a `path` and `value` to be set.

### Action (`actions`) supports the following:
* `route` (Optional) - The ID of the service where the event will be routed.
* `priority` (Optional) - The ID of the priority applied to the event.
* `severity` (Optional)  - The [severity level](https://support.pagerduty.com/docs/rulesets#section-set-severity-with-event-rules) of the event. Can be either `info`,`error`,`warning`, or `critical`.
* `annotate` (Optional) - Note added to the event.
* `extractions` (Optional) - Allows you to copy important data from one event field to another. Extraction rules must use valid [RE2 regular expression syntax](https://github.com/google/re2/wiki/Syntax). Extraction objects consist of the following fields:
	* `source` - Field where the data is being copied from.
	* `target` - Field where the data is being copied to.
	* `regex` - The conditions that need to be met for the extraction to happen.
	* *NOTE: A rule can have multiple `extraction` objects attributed to it.*

* `suppress` (Optional) - Controls whether an alert is [suppressed](https://support.pagerduty.com/docs/rulesets#section-suppress-but-create-triggering-thresholds-with-event-rules) (does not create an incident).
	* `value` - Boolean value that indicates if the alert should be suppressed before the indicated threshold values are met.
	* `threshold_value` - The number of alerts that should be suppressed.
	* `threshold_time_amount` - The number value of the `threshold_time_unit` before an incident is created.
	* `threshold_time_unit` - The `minutes`,`hours`, or `days` that the `threshold_time_amount` should be measured. 

### Time Frame (`time_frame`) supports the following:
* `scheduled_weekly` (Optional) - Values for executing the rule on a recurring schedule.
	* `weekdays` - An integer array representing which days during the week the rule executes. For example `weekdays = [1,3,7]` would execute on Monday, Wednesday and Sunday.
	* `timezone` - Timezone for the given schedule.
	* `start_time` - Time when the schedule will start. Unix timestamp in milliseconds. For example, if you have a rule with a `start_time` of `0` and a `duration` of `60,000` then that rule would be active from `00:00` to `00:01`. If the `start_time` was `3,600,000` the it would be active starting at `01:00`.
	* `duration` - Length of time the schedule will be active.  Unix timestamp in milliseconds.
* `active_between` (Optional) - Values for executing the rule during a specific time period.
	* `start_time` - Beginning of the scheduled time when the rule should execute.  Unix timestamp in milliseconds.
	* `end_time` - Ending of the scheduled time when the rule should execute.  Unix timestamp in milliseconds.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the rule.
  * `catch_all` - Indicates whether the rule is the last rule of the ruleset that serves as a catch-all. It has limited functionality compared to other rules.

## Import

Ruleset rules can be imported using using the related `ruleset` id and the `ruleset_rule` id separated by a dot, e.g.

```
$ terraform import pagerduty_ruleset_rule.main a19cdca1-3d5e-4b52-bfea-8c8de04da243.19acac92-027a-4ea0-b06c-bbf516519601
```
