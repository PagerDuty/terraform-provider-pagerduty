---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_incident_workflow"
sidebar_current: "docs-pagerduty-resource-incident-workflow"
description: |-
  Creates and manages an incident workflow in PagerDuty.
---

# pagerduty\_incident\_workflow

An [Incident Workflow](https://support.pagerduty.com/docs/incident-workflows) is a series of steps which can be executed on an incident.

## Example Usage

```hcl
resource "pagerduty_incident_workflow" "my_first_workflow" {
  name         = "Example Incident Workflow"
  description  = "This Incident Workflow is an example"
  step {
    name           = "Send Status Update"
    action         = "pagerduty.com:incident-workflows:send-status-update:1"
    input {
      name = "Message"
      value = "Example status message sent on {{current_date}}"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the workflow.
* `description` - (Optional) The description of the workflow.
* `team` - (Optional) A team ID. If specified then workflow edit permissions will be scoped to members of this team.
* `step` - (Optional) The steps in the workflow.

Each incident workflow step (`step`) supports the following:

* `name` - (Required) The name of the workflow step.
* `action` - (Required) The action id for the workflow step, including the version. A list of actions available can be retrieved using the [PagerDuty API](https://developer.pagerduty.com/api-reference/aa192a25fac39-list-actions).
* `input` - (Optional) The list of standard inputs for the workflow action.
* `inline_steps_input` - (Optional) The list of inputs that contain a series of inline steps for the workflow action.

Each incident workflow step standard input (`input`) supports the following:

* `name` - (Required) The name of the input.
* `value` - (Required) The value of the input.

Each incident workflow step inline steps input (`inline_steps_input`) points to an input whose metadata describes the `format` as `inlineSteps` and supports the following:

* `name` - (Required) The name of the input.
* `step` - (Required) The inline steps of the input. An inline step adheres to the step schema described above.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the incident workflow.

## Import

Incident workflows can be imported using the `id`, e.g.

```
$ terraform import pagerduty_incident_workflow.major_incident_workflow PLBP09X
```
