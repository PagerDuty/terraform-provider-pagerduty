---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_automation_actions_action_team_association"
sidebar_current: "docs-pagerduty-resource-automation-actions-action"
description: |-
  Creates and manages an Automation Actions action association with a Team in PagerDuty.
---

# pagerduty\_automation\_actions\_action_team_association

An Automation Actions [action association with a team](https://developer.pagerduty.com/api-reference/8f722dd91a4ba-associate-an-automation-action-with-a-team) configures the relation of a specific Action with a Team.

## Example Usage

```hcl
resource "pagerduty_team" "example" {
  name        = "Engineering"
  description = "All engineering"
}

resource "pagerduty_automation_actions_action" "pa_action_example" {
  name        = "PA Action created via TF"
  description = "Description of the PA Action created via TF"
  action_type = "process_automation"
  action_data_reference {
    process_automation_job_id = "P123456"
  }
}

resource "pagerduty_automation_actions_action_team_association" "foo" {
  action_id = pagerduty_automation_actions_action.pa_action_example.id
  team_id   = pagerduty_team.example.id
}

```

## Argument Reference

The following arguments are supported:

  * `action_id` - (Required) Id of the action.
  * `team_id` - (Required) Id of the team associated to the action.

## Import

Action team association can be imported using the `action_id` and `team_id` separated by a colon, e.g.

```
$ terraform import pagerduty_automation_actions_action_team_association.example 01DER7CUUBF7TH4116K0M4WKPU:PLB09Z
```
