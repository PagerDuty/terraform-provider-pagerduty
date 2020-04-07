---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_rule"
sidebar_current: "docs-pagerduty-resource-event-rule"
description: |-
  Creates and manages an event rule in PagerDuty.
---

# pagerduty\_event_rule

*NOTE: The `pagerduty_event_rule` resource has been deprecated in favor the the [pagerduty_ruleset](ruleset.html) and [pagerduty_ruleset_rule](ruleset_rule.html) resources. Please use the `ruleset` based resources for working with Event Rules.*


An [event rule](https://v2.developer.pagerduty.com/docs/global-event-rules-api) determines what happens to an event that is sent to PagerDuty by monitoring tools and other integrations.


## Example Usage

```hcl
resource "pagerduty_event_rule" "second" {
    action_json = jsonencode([
        [
            "route",
            "P5DTL0K"
        ],
        [
            "severity",
            "warning"
        ],
        [
            "annotate",
            "2 Managed by terraform"
        ],
        [
            "priority",
            "PL451DT"
        ]
    ])
    condition_json = jsonencode([
        "and",
        ["contains",["path","payload","source"],"website"],
        ["contains",["path","headers","from","0","address"],"homer"]
    ])
    advanced_condition_json = jsonencode([
        [
            "scheduled-weekly",
            1565392127032,
            3600000,
            "America/Los_Angeles",
            [
                1,
                2,
                3,
                5,
                7
            ]
        ]
    ])
}
resource "pagerduty_event_rule" "third" {
    action_json = jsonencode([
        [
            "route",
            "P5DTL0K"
        ],
        [
            "severity",
            "warning"
        ],
        [
            "annotate",
            "3 Managed by terraform"
        ],
        [
            "priority",
            "PL451DT"
        ]
    ])
    condition_json = jsonencode([
        "and",
        ["contains",["path","payload","source"],"website"],
        ["contains",["path","headers","from","0","address"],"homer"]
    ])
    depends_on = [pagerduty_event_rule.two]
}
```

## Argument Reference

The following arguments are supported:

* `action_json` - (Required) A list of one or more actions for each rule. Each action within the list is itself a list.
* `condition_json` - (Required) Contains a list of conditions. The first field in the list is `and` or `or`, followed by a list of operators and values.
* `advanced_condition_json` - (Optional) Contains a list of specific conditions including `active-between`,`scheduled-weekly`, and `frequency-over`. The first element in the list is the label for the condition, followed by a list of values for the specific condition. For more details on these conditions see [Advanced Condition](https://v2.developer.pagerduty.com/docs/global-event-rules-api#section-advanced-condition) in the PagerDuty API documentation.
* `depends_on` - (Optional) A [Terraform meta-parameter](https://www.terraform.io/docs/configuration-0-11/resources.html#depends_on) that ensures that the `event_rule` specified is created before the current rule. This is important because Event Rules in PagerDuty are executed in order. `depends_on` ensures that  the rules are created in the order specified.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the event rule.
  * `catch_all` - A boolean that indicates whether the rule is a catch all for the account. This field is read-only through the PagerDuty API.

## Import

Event rules can be imported using the `id`, e.g.

```
$ terraform import pagerduty_event_rule.main 19acac92-027a-4ea0-b06c-bbf516519601
```
