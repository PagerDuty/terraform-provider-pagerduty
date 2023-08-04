---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_incident_workflow"
sidebar_current: "docs-pagerduty-datasource-incident-workflow"
description: |-
  Get information about an incident workflow.
---

# pagerduty\_incident\_workflow

Use this data source to get information about a specific [Incident Workflow](https://support.pagerduty.com/docs/incident-workflows) so that you can create a trigger for it.

## Example Usage

```hcl
data "pagerduty_incident_workflow" "my_workflow" {
  name = "Some Workflow Name"
}

data "pagerduty_service" "first_service" {
  name = "My First Service"
}

resource "pagerduty_incident_workflow_trigger" "automatic_trigger" {
  type       = "conditional"
  workflow   = data.pagerduty_incident_workflow.my_workflow.id
  services   = [data.pagerduty_service.first_service.id]
  condition  = "incident.priority matches 'P1'"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the workflow.

## Attributes Reference

* `id` - The ID of the found workflow.
