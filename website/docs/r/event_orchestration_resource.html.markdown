---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestration"
sidebar_current: "docs-pagerduty-resource-event-orchestration"
description: |-
  Creates and manages a Global Event Orchestration / Service Orchestration in PagerDuty.
---



# pagerduty_event_orchestration

[Event Orchestrations](https://support.pagerduty.com/docs/event-orchestration) allows users to route events to an endpoint and create nested rules, which define sets of actions to take based on event content.

## Example of configuring a Global Event Orchestration

```hcl
resource "pagerduty_team" "engineering" {
  name = "Engineering"
}

resource "pagerduty_event_orchestration" "my_monitor" {
  name = "My Monitoring Orchestration"
  description = "Send events to a pair of services"
  team {
    id = pagerduty_team.engineering.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Event Orchestration.
* `description` - (Optional) A human-friendly description of the Event Orchestration.
* `team` - (Optional) Reference to the team that owns the Event Orchestration. If none is specified, only admins have access.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Event Orchestration.
* `integrations` - Routing keys routed to this Event Orchestration.

## Import

EventOrchestrations can be imported using the `id`, e.g.

```
$ terraform import pagerduty_event_orchestration.main 19acac92-027a-4ea0-b06c-bbf516519601
```

# pagerduty_event_orchestration_service

When integrations exist on a service, [Service Orchestrations](https://support.pagerduty.com/docs/event-orchestration#service-orchestrations) can be used to evaluate incoming events against each of its rules, beginning with the rules in the "start" set. When a matching rule is found, it can modify and enhance the event and can route the event to another set of rules within this Service Orchestration for further processing.

## Example of configuring a Service Orchestration

```hcl
  # user, escalation policy are required for a service.
  # a service orchestration is required to point to an existing service.
  # This example shows creating the prerequisite resources for a Service Orchestration (team, user, escalationpolicy and service)
  resource "pagerduty_team" "engineering" {
  name = "Engineering"
}

  resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
  teams = [pagerduty_team.engineering.id]
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "Engineering Escalation Policy"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user"
      id   = pagerduty_user.example.id
    }
  }
}

resource "pagerduty_service" "example" {
  name                    = "My Web App"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.example.id
  alert_creation          = "create_alerts_and_incidents"
}

# In this example the user has defined the service orchestration with nested rulesets.
# The start set rule in this example routes to the second set
resource "pagerduty_event_orchestration_service" "www" {
  type = "service"
  parent {
    # id of the Service
    id = pagerduty_service.example.id
  }
  sets {
    id = "start"
    rules {
      label = "Always apply some consistent event transformations to all events"
      actions {
        variables {
          name = "hostname"
          path = "event.component"
          value = "hostname: (.*)"
          type = "regex"
        }
        extractions {
          # Demonstrating a template-style extraction
          template = "{{variables.hostname}}"
          target = "event.custom_details.hostname"
        }
        extractions {
          # Demonstrating a regex-style extractions
          source = "event.source"
          regex = "www (.*) service"
          target = "event.source"
        }
        # Id of the next set
        route_to = "step-two"
      }
    }
  }
  sets {
     id = "step-two"
     rules {
       label = "All critical alerts should be treated as P1 incident"
       conditions {
         expression = "event.severity matches 'critical'"
       }
       actions {
         annotate = "Please use our P1 runbook: https://docs.test/p1-runbook"
         priority = "P0IN2KQ"
         suppress = false
       }
     }
     rules {
       label = "If there's something wrong on the canary let the team know about it in our deployments Slack channel"
       conditions {
         expression = "event.custom_details.hostname matches part 'canary'"
       }
       # create webhook action with parameters and headers
       actions {
         automation_actions {
           name = "Canary Slack Notification"
           url = "https://our-slack-listerner.test/canary-notification"
           auto_send = true
           parameters {
             key = "channel"
             value = "#my-team-channel"
           }
           parameters {
             key = "message"
             value = "something is wrong with the canary deployment"
           }
           headers {
             key = "X-Notification-Source"
             value = "PagerDuty Incident Webhook"
           }
         }
       }
     }
     rules {
       label = "Never bother the on-call for info-level events outside of work hours"
       conditions {
         expression = "event.severity matches 'info' and not (now in Mon,Tue,Wed,Thu,Fri 09:00:00 to 17:00:00 America/Los_Angeles)"
       }
       actions {
         suppress = true
       }
     }
  }
  # catch_all always sets suppressed action to true. Other actions like annotate, severity, priority, variables and extractions, webhooks can be set as well
  catch_all {
    actions { }
  }
}
```
## Argument Reference

The following arguments are supported:

* `type` - (Required) Type of the orchestration. For service orchestrations, it is `service`.
* `parent` - (Required) Parent (Service) to which this orchestration belongs to.
* `sets` - (Required) A Service Orchestration must contain at least a "start" set, but can contain any number of additional sets that are routed to by other rules to form a directional graph.
* `catch_all` - (Required) When none of the rules match an event, the event will be routed according to the catch_all settings.


### Parent (`parent`) supports the following:
* `id` - (Required) ID of the Service to which this service orchestration belongs to.

### Sets (`sets`) supports the following:
* `id` - (Required) The ID of this set of rules. Rules in other sets can route events into this set using the rule's `route_to` property.
* `rules` - (Optional) The service orchestration evaluates Events against these Rules, one at a time, and routes each Event based on the first rule that matches. If no rules are provided as part of Terraform configuration, the API returns empty list of rules.

### Rules (`rules`) supports the following:
* `label` - (Optional) A description of this rule's purpose.
* `conditions` - (Optional) Each of these conditions is evaluated to check if an event matches this rule. The rule is considered a match if any of these conditions match. If none are provided, the event will `always` match against the rule.
* `actions` - (Required) Actions that will be taken to change the resulting alert and incident, when an event matches this rule.
* `disabled` - (Optional) Indicates whether the rule is disabled and would therefore not be evaluated.

### Conditions (`conditions`) supports the following:
* `expression`- (Required) A [PCL condition] (https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview) string.

### Actions (`actions`) supports the following:
* `route_to` - (Required) The ID of the target Service for the resulting alert. You can find the service you want to route to by calling the [Services API endpoint](https://developer.pagerduty.com/api-reference/e960cca205c0f-list-services).
* `suppress` - (Optional) Set whether the resulting alert is suppressed. Suppressed alerts will not trigger an incident.
* `suspend` - (Optional) The number of seconds to suspend the resulting alert before triggering. This effectively pauses incident notifications. If a `resolve` event arrives before the alert triggers then PagerDuty won't create an incident for this the resulting alert.
* `priority` - (Optional) The ID of the priority you want to set on resulting incident. You can find the list of priority IDs for your account by calling the priorities endpoint.
* `annotate` - (Optional) Add this text as a note on the resulting incident.
* `pagerduty_automation_actions` - (Optional) Configure a [Process Automation](https://support.pagerduty.com/docs/event-orchestration#process-automation) associated with the resulting incident.
  * `action_id` - (Required) Id of the Process Automation action to be triggered.
* `automation_actions` - (Optional) Create a [Webhook](https://support.pagerduty.com/docs/event-orchestration#webhooks) associated with the resulting incident.
  * `name` - (Required) Name of this Webhook.
  * `url` - (Required) The API endpoint where PagerDuty's servers will send the webhook request.
  * `auto_send` - (Optional) When true, PagerDuty's servers will automatically send this webhook request as soon as the resulting incident is created. When false, your incident responder will be able to manually trigger the Webhook via the PagerDuty website and mobile app.
  * `headers` - (Optional) Specify custom key/value pairs that'll be sent with the webhook request as request headers.
    * `key` - (Required) Name to identify the header
    * `value` - (Required) Value of this header
  * `parameters` - (Optional) Specify custom key/value pairs that'll be included in the webhook request's JSON payload.
    * `key` - (Required) Name to identify the parameter
    * `value` - (Required) Value of this parameter
* `severity` - (Optional) sets Severity of the resulting alert. Allowed values are: `info`, `error`, `warning`, `critical`
* `event_action` - (Optional) sets whether the resulting alert status is trigger or resolve. Allowed values are: `trigger`, `resolve`
* `variables` - (Optional) Populate variables from event payloads and use those variables in other event actions.
  * `name` - (Required) The name of the variable
  * `path` - (Required) Path to a field in an event, in dot-notation. This supports both PagerDuty Common Event Format [PD-CEF](https://support.pagerduty.com/docs/pd-cef) and non-CEF fields. Eg: Use `event.summary` for the `summary` CEF field. Use raw_event.fieldname to read from the original event `fieldname` data.
  * `type` - (Required) Only `regex` is supported
  * `value` - (Required) The Regex expression to match against. Must use valid [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) syntax.
* `extractions` - (Optional) Replace any CEF field or Custom Details object field using custom variables.
  * `template` - (Optional) A value that will be used to populate the `target` field. The configuration can include variables extracted from the payload by using string interpolation. Eg: If you have defined a variable called `hostname` you can set extraction `template` to `High CPU on variables.hostname server` to use the variable in extraction.  This field can be ignored for `regex` based replacements.
  * `target` - (Required) The PagerDuty Common Event Format [PD-CEF](https://support.pagerduty.com/docs/pd-cef) field that will be set with the value from the `template` or based on `regex` and `source` fields.
  * `regex` - (Optional) The conditions that need to be met for the extraction to happen. Must use valid [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) syntax. This field can be ignored for `template` based replacements.
  * `source` - (Optional) Field where the data is being copied from. Must be a PagerDuty Common Event Format [PD-CEF](https://support.pagerduty.com/docs/pd-cef) field. This field can be ignored for `template` based replacements.

### Catch All (`catch_all`) supports the following:
* `actions` - (Required) These are the actions that will be taken to change the resulting alert and incident. `catch_all` supports all actions described above for rules except `route_to` action.


## Attributes Reference

The following attributes are exported:
* `parent`
  * `type` - Type of the parent (Event Orchestration) reference for this Event Orchestration Path
  * `self` - The URL at which the parent object (Event Orchestration) is accessible
* `self` - The URL at which the Service Orchestration path is accessible
* `rules`
  * `id` - The ID of the rule within the `start` set.
* `catch_all`
  * `actions`
    * `suppress` - The suppress action for catch_all rule. This is always True.

## Import

Service Orchestrations can be imported using the `id` of the service, e.g.

```
$ terraform import pagerduty_event_orchestration_service PFEODA7
```
