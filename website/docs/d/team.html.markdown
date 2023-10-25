---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_team"
sidebar_current: "docs-pagerduty-datasource-team"
description: |-
  Get information about a team that you can use with escalation_policies, schedules etc.
---

# pagerduty\_team

Use this data source to get information about a specific [team][1] that you can use for other PagerDuty resources.

## Example Usage

```hcl
data "pagerduty_user" "me" {
  email = "me@example.com"
}

data "pagerduty_team" "devops" {
  name = "devops"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "DevOps Escalation Policy"
  num_loops = 2

  teams = [data.pagerduty_team.devops.id]

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user"
      id   = data.pagerduty_user.me.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the team to find in the PagerDuty API.

## Attributes Reference
* `id` - The ID of the found team.
* `name` - The name of the found team.
* `description` - A description of the found team.
* `parent` - ID of the parent team. This is available to accounts with the Team Hierarchy feature enabled. Please contact your account manager for more information.
* `default_role` - (Optional) The team is private if the value is "none", or public if it is "manager" (the default permissions for a non-member of the team are either "none", or their base role up until "manager").

[1]: https://developer.pagerduty.com/api-reference/b3A6Mjc0ODIyMw-list-teams
