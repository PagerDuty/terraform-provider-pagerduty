---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestration_router"
sidebar_current: "docs-pagerduty-resource-event-orchestration-router"
description: |-
  Creates and manages a Router for Global Event Orchestration in PagerDuty.
---

# pagerduty_event_orchestration_router

An Orchestration Router allows users to create a set of Event Rules. The Router evaluates events sent to this Orchestration against each of its rules, one at a time, and routes the event to a specific Service based on the first rule that matches. If an event doesn't match any rules, it'll be sent to service specified in the `catch_all` or to the "Unrouted" Orchestration if no service is specified.

## Example of configuring Router rules for an Orchestration

In this example the user has defined the Router with three rules. The first rule configures a dynamic route: any event containing a value in its `pd_service_id` custom detail will be routed to the Service with the ID specified by that value. The other rules route events matching a condition to specific services.

```hcl
data "pagerduty_service" "database" {
  name = "Primary Data Store"
}

data "pagerduty_service" "www" {
  name = "Web Server App"
}

resource "pagerduty_event_orchestration_router" "router" {
  event_orchestration = pagerduty_event_orchestration.my_monitor.id
  set {
    id = "start"
    rule {
      label = "Dynamically route events related to specific PagerDuty services"
      actions {
        dynamic_route_to {
          lookup_by = "service_id"
          source = "event.custom_details.pd_service_id"
          regex = "(.*)"
        }
      }
    }
    rule {
      label = "Events relating to our relational database"
      condition {
        expression = "event.summary matches part 'database'"
      }
      condition {
        expression = "event.source matches regex 'db[0-9]+-server'"
      }
      actions {
        route_to = data.pagerduty_service.database.id
      }
    }
    rule {
      condition {
        expression = "event.summary matches part 'www'"
      }
      actions {
        route_to = data.pagerduty_service.www.id
      }
    }
  }
  catch_all {
    actions {
      route_to = "unrouted"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `event_orchestration` - (Required) ID of the Event Orchestration to which the Router belongs.
* `set` - (Required) The Router contains a single set of rules  (the "start" set).
* `catch_all` - (Required) When none of the rules match an event, the event will be routed according to the catch_all settings.

### Set (`set`) supports the following:
* `id` - (Required) ID of the `start` set. Router supports only one set and it's id has to be `start`
* `rule` - (Optional) The Router evaluates Events against these Rules, one at a time, and routes each Event to a specific Service based on the first rule that matches. If no rules are provided as part of Terraform configuration, the API returns empty list of rules.

### Rule (`rule`) supports the following:
* `label` - (Optional) A description of this rule's purpose.
* `condition` - (Optional) Each of these conditions is evaluated to check if an event matches this rule. The rule is considered a match if any of these conditions match. If none are provided, the event will _always_ match against the rule.
* `actions` - (Required) Actions that will be taken to change the resulting alert and incident, when an event matches this rule.
* `disabled` - (Optional) Indicates whether the rule is disabled and would therefore not be evaluated.

### Condition (`condition`) supports the following:
* `expression`- (Required) A [PCL condition](https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview) string.

### Actions (`actions`) supports the following:

#### Dynamic Routing

Use the contents of an event payload to dynamically route an event to the target service. Note these setting can only be used once in the Router, and only in the first rule. The dynamic routing rule cannot have `conditions` nor a `route_to` action defined.

* `dynamic_route_to` - (Required) supports the following:
    * `source` - (Required) The path to a field in an event.
    * `regex` - (Required) The regular expression, used to extract a value from the source field. Must use valid [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) syntax.
    * `lookup_by` - (Required) Indicates whether the extracted value from the source is a service's name or ID. Allowed values are: `service_name`, `service_id`

If an event has a value at the specified `source`, and if the `regex` successfully matches the value, and if the matching portion is valid Service ID or Name, then the event will be routed to that service. Otherwise the event will be checked against any subsequent router rules.

#### Service Route

If an event matches this rule's conditions, then route it to the specified Service.

* `route_to` - (Required) The ID of the target Service for the resulting alert.

### Catch All (`catch_all`) supports the following:
* `actions` - (Required) These are the actions that will be taken to change the resulting alert and incident.
  * `route_to` - (Required) Defines where an alert will be sent if doesn't match any rules. Can either be the ID of a Service _or_ the string `"unrouted"` to send events to the Unrouted Orchestration.

## Attributes Reference

The following attributes are exported:
* `self` - The URL at which the Router Orchestration is accessible.
* `rule`
  * `id` - The ID of the rule within the `start` set.

## Import

Router can be imported using the `id` of the Event Orchestration, e.g.

```
$ terraform import pagerduty_event_orchestration_router.router 1b49abe7-26db-4439-a715-c6d883acfb3e
```
