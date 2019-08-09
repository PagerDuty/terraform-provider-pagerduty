---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_rule"
sidebar_current: "docs-pagerduty-resource-event-rule"
description: |-
  Creates and manages an event rule in PagerDuty.
---

# pagerduty\_event_rule

An [event rule](https://v2.developer.pagerduty.com/docs/global-event-rules-api) determines what happens to an event that is sent to PagerDuty by monitoring tools and other integrations.


## Example Usage

```hcl
variable "action_list" {
    default = [
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
            "Managed by terraform"
        ],
        [
            "priority",
            "PL451DT"
        ]
    ]
}
variable "condition_list" {
    default = [
        "and",
        ["contains",["path","payload","source"],"website"],
        ["contains",["path","headers","from","0","address"],"homer"]
    ]
}
variable "advanced_condition_list" {
    default = [
      [
          "scheduled-weekly",
          1565392127032,
          3600000,
          "America/Los_Angeles",
          [
              1,
              3,
              5,
              7
          ]
      ]
    ]
}
resource "pagerduty_event_rule" "example" {
    action_json = jsonencode(var.action_list)
    condition_json = jsonencode(var.condition_list)
    advanced_condition_json = jsonencode(var.advanced_condition_list)
}
```

## Argument Reference

The following arguments are supported:

* `action_json` - (Required) A list of one or more actions for each rule. Each action within the list is itself a list.
* `condition_json` - (Required) Contains a list of conditions. The first field in the list is `and` or `or`, followed by a list of operators and values.
* `advanced_condition_json` - (Required) Contains a list of specific conditions including `active-between`,`scheduled-weekly`, and `frequency-over`. The first element in the list is the label for the condition, followed by a list of values for the specific condition. For more details on these conditions see [Advanced Condition](https://v2.developer.pagerduty.com/docs/global-event-rules-api#section-advanced-condition) in the PagerDuty API documentation.
* `catch_all` - (Optional) A boolean that indicates whether the rule is a catch all for the account. 

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the event rule.

## Import

Escalation policies can be imported using the `id`, e.g.

```
$ terraform import pagerduty_event_rule.main 19acac92-027a-4ea0-b06c-bbf516519601
```
