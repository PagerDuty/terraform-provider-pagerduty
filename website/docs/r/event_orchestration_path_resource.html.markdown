---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestration_path"
sidebar_current: "docs-pagerduty-resource-event-orchestration-path"
description: |-
  Creates and manages an orchestration path (router|unrouted|service) associated with a Global Event Orchestration in PagerDuty.
---

# pagerduty_event_orchestration_router

An Orchestration Router allows users to create a set of Event Rules. The Router evaluates events sent to this Orchestration against each of its rules, one at a time, and routes the event to a specific Service based on the first rule that matches. If an event doesn't match any rules, it'll be sent to service specified in the catch_all or to the "Unrouted" Orchestration if no service is specified. Rules are defined as attributes of the orchestration path resource.

## Example of configuring Router rules for an Orchestration

```hcl
  # In this example the user has defined the router with two rules, each routing to a different service
  # The first rule matches only if either of the two conditions are matched
  # The second rule matches always
  # The catch_all routes to unrouted path.
  # but using the catch_all to route to another service is allowed.
resource "pagerduty_event_orchestration_router" "router" {
  type = "router"
  parent {
    id = pagerduty_event_orchestration.my_monitor.id
  }
  catch_all {
		actions {
			route_to = "unrouted"
		}
	}
  sets {
    rules {
      label = "Events relating to our relational database"
      conditions {
        expression = "event.summary matches part 'database'"
      }
      conditions {
        expression = "event.source matches regex 'db[0-9]+-server'"
      }
      actions {
        route_to = pageduty_service.database.id
      }
    }
    rules {
      actions {
        route_to = pagerduty_service.www.id
      }
      disabled = false
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) Type of the orchestration path. For router path, it is `router`.
* `parent` - (Required) Parent (Event Orchestration) to which this orchestration path belongs to.
* `sets` - (Required) The Router contains a single set of rules  (the "start" set)
* `catch_all` - (Required) When none of the rules match an event, the event will be routed according to the catch_all settings.


### Parent (`parent`) supports the following:
* `id` - (Required) ID of the Event Orchestration to which the Router path belongs to

### Sets (`sets`) supports the following:
* `id` - (Required) ID of the start set. Router path supports only one set and it's id has to be `start`
* `rules` - (Optional) The Router evaluates Events against these Rules, one at a time, and routes each Event to a specific Service based on the first rule that matches. If no rules are provided as part of Terraform configuration, the API returns empty list of rules.

### Rules (`rules`) supports the following:
* `label` - (Optional) A description of this rule's purpose.
* `conditions` - (Optional) Each of these conditions is evaluated to check if an event matches this rule. The rule is considered a match if any of these conditions match. If none are provided, the event will `always` match against the rule.
* `actions` - (Required) Actions that will be taken to change the resulting alert and incident, when an event matches this rule.
* `disabled` - (Optional) Indicates whether the rule is disabled and would therefore not be evaluated.

### Conditions (`conditions`) supports the following:
* `expression`- (Required) A [PCL condition] (https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview) string.

### Actions (`actions`) supports the following:
* `route_to` - (Required) The ID of the target Service for the resulting alert. You can find the service you want to route to by calling the services endpoint.

### Catch All (`catch_all`) supports the following:
* `actions` - (Required) These are the actions that will be taken to change the resulting alert and incident.
  * `route_to' - (Required) With a value of 'unrouted', all events are sent to the Unrouted Orchestration.

## Attributes Reference

The following attributes are exported:
* `parent`
  * `type' - Type of the parent (Event Orchestration) reference for this Event Orchestration Path.
  * `self` - The URL at which the parent object (Event Orchestration) is accessible.
* `self` - The URL at which the Router Event Orchestration path is accessible.
* `rules`
  * `id` - The ID of the rule within the `start` set.

## Import

Router Orchestration path can be imported using the `id` of the event orchestration, e.g.

```
$ terraform import pagerduty_event_orchestration_router 1b49abe7-26db-4439-a715-c6d883acfb3e
```

# pagerduty_event_orchestration_unrouted

An Unrouted Orchestration allows users to create a set of Event Rules that will be evaluated against all events that don't match any rules in the Orchestration's Router.

The Unrouted Orchestration evaluates events sent to it against each of its rules, beginning with the rules in the "start" set. When a matching rule is found, it can modify and enhance the event and can route the event to another set of rules within this Unrouted Orchestration for further processing.

## Example of configuring a Unrouted Rules for an Orchestration

```hcl
  # In this example the user has defined the unrouted orchestration path with two rules, each routing to a different service
  # As there is a single set, route_to is not defined for the rules.
  # If there are more than one set, the rules in start set must define route_to with id of the next set
  # The first rule matches only if the condition is matched
  # The second rule matches always
  # The catch_all without actions as in this example will set suppressed action to true.
  # but using the catch_all to set severity, event_action, variables and extractions is allowed.
resource "pagerduty_event_orchestration_unrouted" "unrouted" {
  type = "unrouted"
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
      disabled = false
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

* `type` - (Required) Type of the orchestration path. For unrouted path, it is `unrouted`.
* `parent` - (Required) Parent (Event Orchestration) to which this orchestration path belongs to.
* `sets` - (Required) An Unrouted Orchestration must contain at least a "start" set, but can contain any number of additional sets that are routed to by other rules to form a directional graph.
* `catch_all` - (Required) When none of the rules match an event, the event will be routed according to the catch_all settings.


### Parent (`parent`) supports the following:
* `id` - (Required) ID of the Event Orchestration to which the Unrouted path belongs to.

### Sets (`sets`) supports the following:
* `id` - (Required) The ID of this set of rules. Rules in other sets can route events into this set using the rule's `route_to` property.
* `rules` - (Optional) The unrouted path evaluates Events against these Rules, one at a time, and routes each Event based on the first rule that matches. If no rules are provided as part of Terraform configuration, the API returns empty list of rules.

### Rules (`rules`) supports the following:
* `label` - (Optional) A description of this rule's purpose.
* `conditions` - (Optional) Each of these conditions is evaluated to check if an event matches this rule. The rule is considered a match if any of these conditions match. If none are provided, the event will `always` match against the rule.
* `actions` - (Required) Actions that will be taken to change the resulting alert and incident, when an event matches this rule.
* `disabled` - (Optional) Indicates whether the rule is disabled and would therefore not be evaluated.

### Conditions (`conditions`) supports the following:
* `expression`- (Required) A [PCL condition] (https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview) string.

### Actions (`actions`) supports the following:
* `route_to` - (Required) The ID of the target Service for the resulting alert. You can find the service you want to route to by calling the services endpoint.
* `severity` - (Optional) sets Severity of the resulting alert. Allowed values are: `info`, `error`, `warning`, `critical`
* `event_action` - (Optional) sets whether the resulting alert status is trigger or resolve. Allowed values are: `trigger`, `resolve`
* `variables` - (Optional) Populate variables from event payloads and use those variables in other event actions.
  * `name` - (Required) The name of the variable
  * `path` - (Required) Path to a field in an event, in dot-notation. This supports both [PD-CEF](https://support.pagerduty.com/docs/pd-cef) and non-CEF fields. Eg: Use `event.summary` for the `summary` CEF field. Use raw_event.fieldname to read from the original event`s `fieldname` data.
  * `type` - (Required) Only `regex` is supported
  * `value` - The Regex expression to match against.Must use valid [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) syntax.
* `extractions` - (Optional) Replace any CEF field or Custom Details object field using custom variables.
  * `template` - (Optional) A value that will be used to populate the `target` field. The configuration can include variables extracted from the payload by using string interpolation. Eg: If you have defined a variable called `hostname` you can set extraction`s `template` to `High CPU on variables.hostname server` to use the variable in extraction.  This field can be ignored for `regex` based replacements.
  * `target` - (Required) The PagerDuty Common Event Format [PD-CEF](https://support.pagerduty.com/docs/pd-cef) field that will be set with the value from the `template` or based on `regex` and `source` fields.
  * `regex` - (Optional) The conditions that need to be met for the extraction to happen. Must use valid [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) syntax. This field can be ignored for `template` based replacements.
  * `source` - (Optional) Field where the data is being copied from. Must be a PagerDuty Common Event Format [PD-CEF] (https://support.pagerduty.com/docs/pd-cef) field. This field can be ignored for `template` based replacements.

### Catch All (`catch_all`) supports the following:
* `actions` - (Required) These are the actions that will be taken to change the resulting alert and incident. `catch_all` supports all actions described above for rules except `route_to` action.

## Attributes Reference

The following attributes are exported:
* `parent`
  * `type' - Type of the parent (Event Orchestration) reference for this Event .Orchestration Path
  * `self` - The URL at which the parent object (Event Orchestration) is accessible.
* `self` - The URL at which the Unrouted Event Orchestration path is accessible.
* `rules`
  * `id` - The ID of the rule within the `start` set.
* `catch_all`
  * `actions`
    * `suppress` - The suppress action for catch_all rule. This is always True.

## Import

Unrouted Orchestration path can be imported using the `id` of the event orchestration, e.g.

```
$ terraform import pagerduty_event_orchestration_unrouted 1b49abe7-26db-4439-a715-c6d883acfb3e
```
