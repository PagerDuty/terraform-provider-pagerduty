---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestration_service"
sidebar_current: "docs-pagerduty-resource-event-orchestration-service"
description: |-
  Creates and manages a Service Orchestration for a Service.
---

# pagerduty_event_orchestration_service

A [Service Orchestration](https://support.pagerduty.com/docs/event-orchestration#service-orchestrations) allows you to create a set of Event Rules. The Service Orchestration evaluates Events sent to this Service against each of its rules, beginning with the rules in the "start" set. When a matching rule is found, it can modify and enhance the event and can route the event to another set of rules within this Service Orchestration for further processing.

**Note:** If you have a Service that uses [Service Event Rules](https://support.pagerduty.com/docs/rulesets#service-event-rules), you can switch to [Service Orchestrations](https://support.pagerduty.com/docs/event-orchestration#service-orchestrations) at any time. Please read the [Switch to Service Orchestrations](https://support.pagerduty.com/docs/event-orchestration#switch-to-service-orchestrations) instructions for more information.

## Example of configuring a Service Orchestration

This example shows creating `Team`, `User`, `Escalation Policy`, and `Service` resources followed by creating a Service Orchestration to handle Events sent to that Service.

This example also shows using `priority` [data source](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/data-sources/priority) to configure `priority` action for a rule. If the Event matches the first rule in set "step-two" the resulting incident will have the Priority `P1`.

This example shows a Service Orchestration that has nested sets: a rule in the "start" set has a `route_to` action pointing at the "step-two" set.

The `catch_all` actions will be applied if an Event reaches the end of any set without matching any rules in that set. In this example the `catch_all` doesn't have any `actions` so it'll leave events as-is.


```hcl
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

data "pagerduty_priority" "p1" {
  name = "P1"
}

resource "pagerduty_event_orchestration_service" "www" {
  service = pagerduty_service.example.id
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
        priority = data.pagerduty_priority.p1.id
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
  catch_all {
    actions { }
  }
}
```
## Argument Reference

The following arguments are supported:

* `service` - (Required) ID of the Service to which this Service Orchestration belongs to.
* `sets` - (Required) A Service Orchestration must contain at least a "start" set, but can contain any number of additional sets that are routed to by other rules to form a directional graph.
* `catch_all` - (Required) the `catch_all` actions will be applied if an Event reaches the end of any set without matching any rules in that set.

### Sets (`sets`) supports the following:
* `id` - (Required) The ID of this set of rules. Rules in other sets can route events into this set using the rule's `route_to` property.
* `rules` - (Optional) The service orchestration evaluates Events against these Rules, one at a time, and applies all the actions for first rule it finds where the event matches the rule's conditions. If no rules are provided as part of Terraform configuration, the API returns empty list of rules.

### Rules (`rules`) supports the following:
* `label` - (Optional) A description of this rule's purpose.
* `conditions` - (Optional) Each of these conditions is evaluated to check if an event matches this rule. The rule is considered a match if any of these conditions match. If none are provided, the event will `always` match against the rule.
* `actions` - (Required) Actions that will be taken to change the resulting alert and incident, when an event matches this rule.
* `disabled` - (Optional) Indicates whether the rule is disabled and would therefore not be evaluated.

### Conditions (`conditions`) supports the following:
* `expression`- (Required) A [PCL condition](https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview) string.

### Actions (`actions`) supports the following:
* `route_to` - (Optional) The ID of a Set from this Service Orchestration whose rules you also want to use with event that match this rule.
* `suppress` - (Optional) Set whether the resulting alert is suppressed. Suppressed alerts will not trigger an incident.
* `suspend` - (Optional) The number of seconds to suspend the resulting alert before triggering. This effectively pauses incident notifications. If a `resolve` event arrives before the alert triggers then PagerDuty won't create an incident for this the resulting alert.
* `priority` - (Optional) The ID of the priority you want to set on resulting incident. Consider using the [`pagerduty_priority`](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/data-sources/priority) data source.
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
  * `path` - (Required) Path to a field in an event, in dot-notation. This supports both PagerDuty Common Event Format [PD-CEF](https://support.pagerduty.com/docs/pd-cef) and non-CEF fields. Eg: Use `event.summary` for the `summary` CEF field. Use `raw_event.fieldname` to read from the original event `fieldname` data. You can use any valid [PCL path](https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview#paths).
  * `type` - (Required) Only `regex` is supported
  * `value` - (Required) The Regex expression to match against. Must use valid [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) syntax.
* `extractions` - (Optional) Replace any CEF field or Custom Details object field using custom variables.
  * `target` - (Required) The PagerDuty Common Event Format [PD-CEF](https://support.pagerduty.com/docs/pd-cef) field that will be set with the value from the `template` or based on `regex` and `source` fields.
  * `template` - (Optional) A string that will be used to populate the `target` field. You can reference variables or event data within your template using double curly braces. For example:
     * Use variables named `ip` and `subnet` with a template like: `{{variables.ip}}/{{variables.subnet}}`
     * Combine the event severity & summary with template like: `{{event.severity}}:{{event.summary}}`
  * `regex` - (Optional) A [RE2 regular expression](https://github.com/google/re2/wiki/Syntax) that will be matched against field specified via the `source` argument. If the regex contains one or more capture groups, their values will be extracted and appended together. If it contains no capture groups, the whole match is used. This field can be ignored for `template` based extractions.
  * `source` - (Optional) The path to the event field where the `regex` will be applied to extract a value. You can use any valid [PCL path](https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview#paths) like `event.summary` and you can reference previously-defined variables using a path like `variables.hostname`. This field can be ignored for `template` based extractions.

### Catch All (`catch_all`) supports the following:
* `actions` - (Required) These are the actions that will be taken to change the resulting alert and incident. `catch_all` supports all actions described above for `rules` _except_ `route_to` action.


## Attributes Reference

The following attributes are exported:
* `self` - The URL at which the Service Orchestration is accessible.
* `rules`
  * `id` - The ID of the rule within the set.

## Import

Service Orchestration can be imported using the `id` of the Service, e.g.

```
$ terraform import pagerduty_event_orchestration_service PFEODA7
```
