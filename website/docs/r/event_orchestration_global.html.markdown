---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestration_global"
sidebar_current: "docs-pagerduty-resource-event-orchestration-global"
description: |-
  Creates and manages a Global Orchestration for an Event Orchestration.
---

# pagerduty_event_orchestration_global

A [Global Orchestration](https://support.pagerduty.com/docs/event-orchestration#global-orchestrations) allows you to create a set of Event Rules. The Global Orchestration evaluates Events sent to it against each of its rules, beginning with the rules in the "start" set. When a matching rule is found, it can modify and enhance the event and can route the event to another set of rules within this Global Orchestration for further processing.

## Example of configuring a Global Orchestration

This example shows creating `Team`, and `Event Orchestration` resources followed by creating a Global Orchestration to handle Events sent to that Event Orchestration.

This example also shows using the [pagerduty_priority](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/data-sources/priority) and [pagerduty_escalation_policy](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/data-sources/escalation_policy) data sources to configure `priority` and `escalation_policy` actions for a rule.

This example shows a Global Orchestration that has nested sets: a rule in the "start" set has a `route_to` action pointing at the "step-two" set.

The `catch_all` actions will be applied if an Event reaches the end of any set without matching any rules in that set. In this example the `catch_all` doesn't have any `actions` so it'll leave events as-is.


```hcl
resource "pagerduty_team" "database_team" {
  name = "Database Team"
}

resource "pagerduty_event_orchestration" "event_orchestration" {
  name = "Example Orchestration"
  team = pagerduty_team.database_team.id
}

data "pagerduty_priority" "p1" {
  name = "P1"
}

data "pagerduty_escalation_policy" "sre_esc_policy" {
  name = "SRE Escalation Policy"
}

resource "pagerduty_event_orchestration_global" "global" {
  event_orchestration = pagerduty_event_orchestration.event_orchestration.id
  set {
    id = "start"
    rule {
      label = "Always annotate a note to all events"
      actions {
        annotate = "This incident was created by the Database Team via a Global Orchestration"
        # Id of the next set
        route_to = "step-two"
      }
    }
  }
  set {
    id = "step-two"
    rule {
      label = "Drop events that are marked as no-op"
      condition {
        expression = "event.summary matches 'no-op'"
      }
      actions {
        drop_event = true
      }
    }
    rule {
      label = "If the DB host is running out of space, then page the SRE team"
      condition {
        expression = "event.summary matches part 'running out of space'"
      }
      actions {
        escalation_policy = data.pagerduty_escalation_policy.sre_esc_policy.id
      }
    }
    rule {
      label = "If there's something wrong on the replica, then mark the alert as a warning"
      condition {
        expression = "event.custom_details.hostname matches part 'replica'"
      }
      actions {
        severity = "warning"
      }
    }
    rule {
      label = "Otherwise, set the incident to P1 and run a diagnostic"
      actions {
        priority = data.pagerduty_priority.p1.id
        automation_action {
          name = "db-diagnostic"
          url = "https://example.com/run-diagnostic"
          auto_send = true
        }
      }
    }
  }
  catch_all {
    actions { }
  }
}
```
## Argument Reference

The following arguments are supported:

* `event_orchestration` - (Required) ID of the Event Orchestration to which this Global Orchestration belongs to.
* `set` - (Required) A Global Orchestration must contain at least a "start" set, but can contain any number of additional sets that are routed to by other rules to form a directional graph.
* `catch_all` - (Required) the `catch_all` actions will be applied if an Event reaches the end of any set without matching any rules in that set.

### Set (`set`) supports the following:
* `id` - (Required) The ID of this set of rules. Rules in other sets can route events into this set using the rule's `route_to` property.
* `rule` - (Optional) The Global Orchestration evaluates Events against these Rules, one at a time, and applies all the actions for first rule it finds where the event matches the rule's conditions. If no rules are provided as part of Terraform configuration, the API returns empty list of rules.

### Rule (`rule`) supports the following:
* `label` - (Optional) A description of this rule's purpose.
* `condition` - (Optional) Each of these conditions is evaluated to check if an event matches this rule. The rule is considered a match if any of these conditions match. If none are provided, the event will `always` match against the rule.
* `actions` - (Required) Actions that will be taken to change the resulting alert and incident, when an event matches this rule.
* `disabled` - (Optional) Indicates whether the rule is disabled and would therefore not be evaluated.

### Condition (`condition`) supports the following:
* `expression`- (Required) A [PCL condition](https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview) string.

### Actions (`actions`) supports the following:
* `route_to` - (Optional) The ID of a Set from this Global Orchestration whose rules you also want to use with events that match this rule.
* `drop_event` - (Optional) When true, this event will be dropped. Dropped events will not trigger or resolve an alert or an incident. Dropped events will not be evaluated against router rules.
* `suppress` - (Optional) Set whether the resulting alert is suppressed. Suppressed alerts will not trigger an incident.
* `suspend` - (Optional) The number of seconds to suspend the resulting alert before triggering. This effectively pauses incident notifications. If a `resolve` event arrives before the alert triggers then PagerDuty won't create an incident for this alert.
* `priority` - (Optional) The ID of the priority you want to set on resulting incident. Consider using the [`pagerduty_priority`](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/data-sources/priority) data source.
* `escalation_policy` - (Optional) The ID of the Escalation Policy you want to assign incidents to. Event rules with this action will override the Escalation Policy already set on a Service's settings, with what is configured by this action.
* `annotate` - (Optional) Add this text as a note on the resulting incident.
* `incident_custom_field_update` - (Optional) Assign a custom field to the resulting incident.
  * `id` - (Required) The custom field id
  * `value` - (Required) The value to assign to this custom field
* `automation_action` - (Optional) Create a [Webhook](https://support.pagerduty.com/docs/event-orchestration#webhooks) associated with the resulting incident.
  * `name` - (Required) Name of this Webhook.
  * `url` - (Required) The API endpoint where PagerDuty's servers will send the webhook request.
  * `auto_send` - (Optional) When true, PagerDuty's servers will automatically send this webhook request as soon as the resulting incident is created. When false, your incident responder will be able to manually trigger the Webhook via the PagerDuty website and mobile app.
  * `header` - (Optional) Specify custom key/value pairs that'll be sent with the webhook request as request headers.
    * `key` - (Required) Name to identify the header
    * `value` - (Required) Value of this header
  * `parameter` - (Optional) Specify custom key/value pairs that'll be included in the webhook request's JSON payload.
    * `key` - (Required) Name to identify the parameter
    * `value` - (Required) Value of this parameter
* `severity` - (Optional) sets Severity of the resulting alert. Allowed values are: `info`, `error`, `warning`, `critical`
* `event_action` - (Optional) sets whether the resulting alert status is trigger or resolve. Allowed values are: `trigger`, `resolve`
* `variable` - (Optional) Populate variables from event payloads and use those variables in other event actions.
  * `name` - (Required) The name of the variable
  * `path` - (Required) Path to a field in an event, in dot-notation. This supports both PagerDuty Common Event Format [PD-CEF](https://support.pagerduty.com/docs/pd-cef) and non-CEF fields. Eg: Use `event.summary` for the `summary` CEF field. Use `raw_event.fieldname` to read from the original event `fieldname` data. You can use any valid [PCL path](https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview#paths).
  * `type` - (Required) Only `regex` is supported
  * `value` - (Required) The Regex expression to match against. Must use valid [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) syntax.
* `extraction` - (Optional) Replace any CEF field or Custom Details object field using custom variables.
  * `target` - (Required) The PagerDuty Common Event Format [PD-CEF](https://support.pagerduty.com/docs/pd-cef) field that will be set with the value from the `template` or based on `regex` and `source` fields.
  * `template` - (Optional) A string that will be used to populate the `target` field. You can reference variables or event data within your template using double curly braces. For example:
     * Use variables named `ip` and `subnet` with a template like: `{{variables.ip}}/{{variables.subnet}}`
     * Combine the event severity & summary with template like: `{{event.severity}}:{{event.summary}}`
  * `regex` - (Optional) A [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) that will be matched against field specified via the `source` argument. If the regex contains one or more capture groups, their values will be extracted and appended together. If it contains no capture groups, the whole match is used. This field can be ignored for `template` based extractions.
  * `source` - (Optional) The path to the event field where the `regex` will be applied to extract a value. You can use any valid [PCL path](https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview#paths) like `event.summary` and you can reference previously-defined variables using a path like `variables.hostname`. This field can be ignored for `template` based extractions.

### Catch All (`catch_all`) supports the following:
* `actions` - (Required) These are the actions that will be taken to change the resulting alert and incident. `catch_all` supports all actions described above for `rule` _except_ `route_to` action.


## Attributes Reference

The following attributes are exported:
* `self` - The URL at which the Global Orchestration is accessible.
* `rule`
  * `id` - The ID of the rule within the set.

## Import

Global Orchestration can be imported using the `id` of the Event Orchestration, e.g.


```
$ terraform import pagerduty_event_orchestration_global.global 1b49abe7-26db-4439-a715-c6d883acfb3e
```
