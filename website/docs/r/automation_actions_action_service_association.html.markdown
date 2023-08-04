---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_automation_actions_action_service_association"
sidebar_current: "docs-pagerduty-resource-automation-actions-action"
description: |-
  Creates and manages an Automation Actions action association with a Service in PagerDuty.
---

# pagerduty\_automation\_actions\_action_service_association

An Automation Actions [action association with a service](https://developer.pagerduty.com/api-reference/5d2f051f3fb43-associate-an-automation-action-with-a-service) configures the relation of a specific Action with a Service.

## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
}

resource "pagerduty_escalation_policy" "foo" {
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
  escalation_policy       = pagerduty_escalation_policy.foo.id
  alert_creation          = "create_alerts_and_incidents"

  auto_pause_notifications_parameters {
    enabled = true
    timeout = 300
  }
}

resource "pagerduty_automation_actions_action" "pa_action_example" {
  name        = "PA Action created via TF"
  description = "Description of the PA Action created via TF"
  action_type = "process_automation"
  action_data_reference {
    process_automation_job_id = "P123456"
  }
}

resource "pagerduty_automation_actions_action_service_association" "foo" {
  action_id = pagerduty_automation_actions_action.pa_action_example.id
  service_id   = pagerduty_service.example.id
}

```

## Argument Reference

The following arguments are supported:

  * `action_id` - (Required) Id of the action.
  * `service_id` - (Required) Id of the service associated to the action.

## Import

Action service association can be imported using the `action_id` and `service_id` separated by a colon, e.g.

```
$ terraform import pagerduty_automation_actions_action_service_association.example 01DER7CUUBF7TH4116K0M4WKPU:PLB09Z
```
