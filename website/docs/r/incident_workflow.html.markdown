---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_incident_workflow"
sidebar_current: "docs-pagerduty-resource-incident-workflow"
description: |-
  Creates and manages an incident workflow in PagerDuty.
---

# pagerduty\_incident\_workflow

An [Incident Workflow](https://support.pagerduty.com/docs/incident-workflows) is a series of steps which can be executed on an incident.

-> The Incident Workflows feature is currently available in Early Access.

## Example Usage

```hcl
resource "pagerduty_incident_workflow" "my_first_workflow" {
  name         = "My First Workflow"
  description  = "Some description"
  step {
    name           = "Example Step"
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the workflow.
* `description` - (Optional) The description of the workflow.
* `step` - (Optional) The steps in the workflow.

Each incident workflow step (`step`) supports the following:

* `name` - (Required) The name of the workflow step.
* `action` - (Required) The action id for the workflow step, including the version.
* `input` - (Optional) The list of inputs for the workflow action.

Each incident workflow step input (`input`) supports the following:

* `name` - (Required) The name of the input.
* `value` - (Required) The value of the input.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the incident workflow.

## Import

Incident workflows can be imported using the `id`, e.g.

```
$ terraform import pagerduty_incident_workflow.major_incident_workflow PLBP09X
```
