---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_incident_workflow_trigger"
sidebar_current: "docs-pagerduty-resource-incident-workflow-trigger"
description: |-
  Creates and manages an incident workflow trigger in PagerDuty.
---

# pagerduty\_incident\_workflow\_trigger

An [Incident Workflow Trigger](https://support.pagerduty.com/docs/incident-workflows#triggers) defines when and if an [Incident Workflow](https://support.pagerduty.com/docs/incident-workflows) will be triggered.

-> The Incident Workflows feature is currently available in Early Access.

## Example Usage

```hcl
resource "pagerduty_incident_workflow" "my_first_workflow" {
  name         = "My First Workflow"
  description  = "Some description"
  step {
    name           = "Example Step"
    description    = "An example workflow step"
    action         = "something"
    input {
      name = "name"
      value = "value"
    }
  }
  step {
    name          = "Another Step"
    action        = "something_else"
    input {
      name  = "name"
      value = "value"
    }
  }
}

data "pagerduty_service" "first_service" {
  name = "My First Service"
}

resource "pagerduty_incident_workflow_trigger" "automatic_trigger" {
  type                       = "conditional"
  workflow                   = pagerduty_incident_workflow.my_first_workflow.id
  services                   = [pagerduty_service.first_service.id]
  condition                  = "incident.priority matches 'P1'"
  subscribed_to_all_services = false
}

data "pagerduty_team" "devops" {
  name = "devops"
}

resource "pagerduty_incident_workflow_trigger" "manual_trigger" {
  type       = "manual"
  workflow   = pagerduty_incident_workflow.my_first_workflow.id
  services   = [pagerduty_service.first_service.id]
}

```

## Argument Reference

The following arguments are supported:

* `type` - (Required) May be either `manual` or `conditional`.
* `workflow` - (Required) The workflow ID for the workflow to trigger.
* `services` - (Optional) A list of service IDs. Incidents in any of the listed services are eligible to fire this trigger.
* `subscribed_to_all_services` - (Required) Set to `true` if the trigger should be eligible for firing on all services. Only allowed to be `true` if the services list is not defined or empty.
* `condition` - (Required for `conditional`-type triggers) A [PCL](https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview) condition string which must be satisfied for the trigger to fire.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the incident workflow.

## Import

Incident workflows can be imported using the `id`, e.g.

```
$ terraform import pagerduty_incident_workflow.pagerduty_incident_workflow_trigger PLBP09X
```
