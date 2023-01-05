---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_automation_actions_runner_team_association"
sidebar_current: "docs-pagerduty-resource-automation-actions-runner"
description: |-
  Creates and manages an Automation Actions runner association with a Team in PagerDuty.
---

# pagerduty\_automation\_actions\_runner_team_association

An Automation Actions [runner association with a team](https://developer.pagerduty.com/api-reference/f662de6271a6e-associate-a-runner-with-a-team) configures the relation of a specific Runner with a Team.

## Example Usage

```hcl
resource "pagerduty_team" "team_ent_eng" {
  name        = "Enterprise Engineering"
  description = "Enterprise engineering"
}

resource "pagerduty_automation_actions_runner" "pa_runbook_runner" {
	name = "Runner created via TF"
	description = "Description of the Runner created via TF"
	runner_type = "runbook"
	runbook_base_uri = "cat-cat"
	runbook_api_key = "cat-secret"
}

resource "pagerduty_automation_actions_runner_team_association" "pa_runner_ent_eng_assoc" {
  runner_id = pagerduty_automation_actions_runner.pa_runbook_runner.id
  team_id   = pagerduty_team.team_ent_eng.id
}

```

## Argument Reference

The following arguments are supported:

  * `runner_id` - (Required) Id of the runner.
  * `team_id` - (Required) Id of the team associated with the runner.

## Import

Runner team association can be imported using the `runner_id` and `team_id` separated by a colon, e.g.

```
$ terraform import pagerduty_automation_actions_runner_team_association.example 01DER7CUUBF7TH4116K0M4WKPU:PLB09Z
```
