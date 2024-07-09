---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_event_orchestration_service"
sidebar_current: "docs-pagerduty-resource-event-orchestration-service"
description: |-
  Creates and manages a Service Orchestration for a Service.
---

# pagerduty_event_orchestration_service

A [Service Orchestration](https://support.pagerduty.com/docs/event-orchestration#service-orchestrations) allows you to create a set of Event Rules. The Service Orchestration evaluates Events sent to this Service against each of its rules, beginning with the rules in the "start" set. When a matching rule is found, it can modify and enhance the event and can route the event to another set of rules within this Service Orchestration for further processing.

-> If you have a Service that uses [Service Event Rules](https://support.pagerduty.com/docs/rulesets#service-event-rules), you can switch to [Service Orchestrations](https://support.pagerduty.com/docs/event-orchestration#service-orchestrations) at any time setting the attribute `enable_event_orchestration_for_service` to `true`. Please read the [Switch to Service Orchestrations](https://support.pagerduty.com/docs/event-orchestration#switch-to-service-orchestrations) instructions for more information.

## Example of configuring a Service Orchestration

This example shows creating `Team`, `User`, `Escalation Policy`, and `Service` resources followed by creating a Service Orchestration to handle Events sent to that Service.

This example also shows using the [pagerduty_priority](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/data-sources/priority) and [pagerduty_escalation_policy](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/data-sources/escalation_policy) data sources to configure `priority` and `escalation_policy` actions for a rule.

This example shows a Service Orchestration that has nested sets: a rule in the "start" set has a `route_to` action pointing at the "step-two" set.

The `catch_all` actions will be applied if an Event reaches the end of any set without matching any rules in that set. In this example the `catch_all` doesn't have any `actions` so it'll leave events as-is.


```hcl
resource "pagerduty_team" "engineering" {
  name = "Engineering"
}

resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
}

resource "pagerduty_team_membership" "foo" {
  user_id = pagerduty_user.example.id
  team_id = pagerduty_team.engineering.id
  role    = "manager"
}

resource "pagerduty_escalation_policy" "example" {
  name      = "Engineering Escalation Policy"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
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

resource "pagerduty_incident_custom_field" "cs_impact" {
  name       = "impact"
  data_type  = "string"
  field_type = "single_value"
}

data "pagerduty_priority" "p1" {
  name = "P1"
}

data "pagerduty_escalation_policy" "sre_esc_policy" {
  name = "SRE Escalation Policy"
}

resource "pagerduty_event_orchestration_service" "www" {
  service = pagerduty_service.example.id
  enable_event_orchestration_for_service = true
  set {
    id = "start"
    rule {
      label = "Always apply some consistent event transformations to all events"
      actions {
        variable {
          name = "hostname"
          path = "event.component"
          value = "hostname: (.*)"
          type = "regex"
        }
        extraction {
          # Demonstrating a template-style extraction
          template = "{{variables.hostname}}"
          target = "event.custom_details.hostname"
        }
        extraction {
          # Demonstrating a regex-style extraction
          source = "event.source"
          regex = "www (.*) service"
          target = "event.source"
        }
        # Id of the next set
        route_to = "step-two"
      }
    }
  }
  set {
    id = "step-two"
    rule {
      label = "All critical alerts should be treated as P1 incident"
      condition {
        expression = "event.severity matches 'critical'"
      }
      actions {
        annotate = "Please use our P1 runbook: https://docs.test/p1-runbook"
        priority = data.pagerduty_priority.p1.id
        incident_custom_field_update {
          id = pagerduty_incident_custom_field.cs_impact.id
          value = "High Impact"
        }
      }
    }
    rule {
      label = "If any of the API apps are unavailable, page the SRE team"
      condition {
        expression = "event.custom_details.service_name matches part '-api' and event.custom_details.status_code matches '502'"
      }
      actions {
        escalation_policy = data.pagerduty_escalation_policy.sre_esc_policy.id
      }
    }
    rule {
      label = "If there's something wrong on the canary let the team know about it in our deployments Slack channel"
      condition {
        expression = "event.custom_details.hostname matches part 'canary'"
      }
      # create webhook action with parameters and headers
      actions {
        automation_action {
          name = "Canary Slack Notification"
          url = "https://our-slack-listerner.test/canary-notification"
          auto_send = true
          parameter {
            key = "channel"
            value = "#my-team-channel"
          }
          parameter {
            key = "message"
            value = "something is wrong with the canary deployment"
          }
          header {
            key = "X-Notification-Source"
            value = "PagerDuty Incident Webhook"
          }
        }
      }
    }
    rule {
      label = "Never bother the on-call for info-level events outside of work hours"
      condition {
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
* `enable_event_orchestration_for_service` - (Optional) Opt-in/out for switching the Service to [Service Orchestrations](https://support.pagerduty.com/docs/event-orchestration#service-orchestrations).
* `set` - (Required) A Service Orchestration must contain at least a "start" set, but can contain any number of additional sets that are routed to by other rules to form a directional graph.
* `catch_all` - (Required) the `catch_all` actions will be applied if an Event reaches the end of any set without matching any rules in that set.

### Set (`set`) supports the following:
* `id` - (Required) The ID of this set of rules. Rules in other sets can route events into this set using the rule's `route_to` property.
* `rule` - (Optional) The service orchestration evaluates Events against these Rules, one at a time, and applies all the actions for first rule it finds where the event matches the rule's conditions. If no rules are provided as part of Terraform configuration, the API returns empty list of rules.

### Rule (`rule`) supports the following:
* `label` - (Optional) A description of this rule's purpose.
* `condition` - (Optional) Each of these conditions is evaluated to check if an event matches this rule. The rule is considered a match if any of these conditions match. If none are provided, the event will `always` match against the rule.
* `actions` - (Required) Actions that will be taken to change the resulting alert and incident, when an event matches this rule.
* `disabled` - (Optional) Indicates whether the rule is disabled and would therefore not be evaluated.

### Condition (`condition`) supports the following:
* `expression`- (Required) A [PCL condition](https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview) string.

### Actions (`actions`) supports the following:
* `route_to` - (Optional) The ID of a Set from this Service Orchestration whose rules you also want to use with events that match this rule.
* `suppress` - (Optional) Set whether the resulting alert is suppressed. Suppressed alerts will not trigger an incident.
* `suspend` - (Optional) The number of seconds to suspend the resulting alert before triggering. This effectively pauses incident notifications. If a `resolve` event arrives before the alert triggers then PagerDuty won't create an incident for this alert.
* `priority` - (Optional) The ID of the priority you want to set on resulting incident. Consider using the [`pagerduty_priority`](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/data-sources/priority) data source.
* `escalation_policy` - (Optional) The ID of the Escalation Policy you want to assign incidents to. Event rules with this action will override the Escalation Policy already set on a Service's settings, with what is configured by this action.
* `annotate` - (Optional) Add this text as a note on the resulting incident.
* `incident_custom_field_update` - (Optional) Assign a custom field to the resulting incident.
  * `id` - (Required) The custom field id
  * `value` - (Required) The value to assign to this custom field
* `pagerduty_automation_action` - (Optional) Configure a [Process Automation](https://support.pagerduty.com/docs/event-orchestration#process-automation) associated with the resulting incident.
  * `action_id` - (Required) Id of the Process Automation action to be triggered.
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
* `self` - The URL at which the Service Orchestration is accessible.
* `rule`
  * `id` - The ID of the rule within the set.

## Import

Service Orchestration can be imported using the `id` of the Service, e.g.

```
$ terraform import pagerduty_event_orchestration_service.service PFEODA7
```
