---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestration_unrouted"
sidebar_current: "docs-pagerduty-resource-event-orchestration-unrouted"
description: |-
  Creates and manages an Unrouted Orchestration for a Global Event Orchestration in PagerDuty.
---

# pagerduty_event_orchestration_unrouted

An Unrouted Orchestration allows users to create a set of Event Rules that will be evaluated against all events that don't match any rules in the Orchestration's Router.

The Unrouted Orchestration evaluates events sent to it against each of its rules, beginning with the rules in the "start" set. When a matching rule is found, it can modify and enhance the event and can route the event to another set of rules within this Unrouted Orchestration for further processing.

## Example of configuring a Unrouted Rules for an Orchestration

  In this example the user has defined the unrouted orchestration with two rules, each routing to a different service.This example assumes services used in the route_to configuration already exists. So it does not show creation of service resource.As there is a single set, route_to is not defined for the rules.

  If there are more than one set, the rules in start set must define route_to with id of the next set
  The first rule matches only if the condition is matched. The second rule matches always as there are no conditions.

  `catch_all` with empty `actions` block suppresses the alerts that do not match any rules. Using the `catch_all` to set severity, event_action, variables and extractions is also allowed.

```hcl
resource "pagerduty_event_orchestration_unrouted" "unrouted" {
 event_orchestration = pagerduty_event_orchestration.my_monitor.id
  sets {
    id = "start"
    rules {
      label = "Update the summary of un-matched Critical alerts so they're easier to spot"
      conditions {
        expression = "event.severity matches 'critical'"
      }
      actions {
        severity = "info"
        extractions {
            target = "event.summary"
            template = "[Critical Unrouted] {{event.summary}}"
        }
      }
    }
    rules {
      label = "Reduce the severity of all other unrouted events"
      actions {
        severity = "info"
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

* `event_orchestration` - (Required) The Event Orchestration to which this Unrouted Orchestration belongs to.
* `sets` - (Required) An Unrouted Orchestration must contain at least a "start" set, but can contain any number of additional sets that are routed to by other rules to form a directional graph.
* `catch_all` - (Required) When none of the rules match an event, the event will be routed according to the catch_all settings.

### Sets (`sets`) supports the following:
* `id` - (Required) The ID of this set of rules. Rules in other sets can route events into this set using the rule's `route_to` property.
* `rules` - (Optional) The Unrouted Orchestration evaluates Events against these Rules, one at a time, and routes each Event based on the first rule that matches. If no rules are provided as part of Terraform configuration, the API returns empty list of rules.

### Rules (`rules`) supports the following:
* `label` - (Optional) A description of this rule's purpose.
* `conditions` - (Optional) Each of these conditions is evaluated to check if an event matches this rule. The rule is considered a match if any of these conditions match. If none are provided, the event will `always` match against the rule.
* `actions` - (Required) Actions that will be taken to change the resulting alert and incident, when an event matches this rule.
* `disabled` - (Optional) Indicates whether the rule is disabled and would therefore not be evaluated.

### Conditions (`conditions`) supports the following:
* `expression`- (Required) A [PCL condition](https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview) string.

### Actions (`actions`) supports the following:
* `route_to` - (Required) The ID of the target Service for the resulting alert.
* `severity` - (Optional) sets Severity of the resulting alert. Allowed values are: `info`, `error`, `warning`, `critical`
* `event_action` - (Optional) sets whether the resulting alert status is trigger or resolve. Allowed values are: `trigger`, `resolve`
* `variables` - (Optional) Populate variables from event payloads and use those variables in other event actions.
  * `name` - (Required) The name of the variable
  * `path` - (Required) Path to a field in an event, in dot-notation. This supports both [PD-CEF](https://support.pagerduty.com/docs/pd-cef) and non-CEF fields. Eg: Use `event.summary` for the `summary` CEF field. Use raw_event.fieldname to read from the original event `fieldname` data.
  * `type` - (Required) Only `regex` is supported
  * `value` - (Required) The Regex expression to match against. Must use valid [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) syntax.
* `extractions` - (Optional) Replace any CEF field or Custom Details object field using custom variables.
  * `template` - (Optional) A value that will be used to populate the `target` field. The configuration can include variables extracted from the payload by using string interpolation. Eg: If you have defined a variable called `hostname` you can set extraction `template` to `High CPU on variables.hostname server` to use the variable in extraction.  This field can be ignored for `regex` based replacements.
  * `target` - (Required) The PagerDuty Common Event Format [PD-CEF](https://support.pagerduty.com/docs/pd-cef) field that will be set with the value from the `template` or based on `regex` and `source` fields.
  * `regex` - (Optional) The conditions that need to be met for the extraction to happen. Must use valid [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) syntax. This field can be ignored for `template` based replacements.
  * `source` - (Optional) Field where the data is being copied from. Must be a PagerDuty Common Event Format [PD-CEF](https://support.pagerduty.com/docs/pd-cef) field. This field can be ignored for `template` based replacements.

### Catch All (`catch_all`) supports the following:
* `actions` - (Required) These are the actions that will be taken to change the resulting alert and incident. `catch_all` supports all actions described above for `rules` _except_ `route_to` action.

## Attributes Reference

The following attributes are exported:
* `self` - The URL at which the Unrouted Event Orchestration is accessible.
* `rules`
  * `id` - The ID of the rule within the `start` set.

## Import

Unrouted Orchestration can be imported using the `id` of the Event Orchestration, e.g.

```
$ terraform import pagerduty_event_orchestration_unrouted 1b49abe7-26db-4439-a715-c6d883acfb3e
```
