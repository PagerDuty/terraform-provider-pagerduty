---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_team_members"
sidebar_current: "docs-pagerduty-datasource-team-members"
description: |-
  Get information about a team's members..
---

# pagerduty\_team\_members

Use this data source to get information about a specific [team's members][1].

## Example Usage

```hcl
data "pagerduty_team" "devops" {
  name = "devops"
}

data "pagerduty_team_members" "devops_members" {
  team_id = data.pagerduty_team.devops.id
}
```

## Argument Reference

The following arguments are supported:

* `team_id` - (Required) The ID of the team to find in the PagerDuty API.

## Attributes Reference

* `id` - The ID of the found team.
* `members` - The users of the found team.

### Members (`members`) supports the following:

* `id` - The ID of the found user.
* `name` - The short name of the found user.
* `email` - The email of the found user.
* `role` - The team role of the found user.
* `job_title` - The job title of the found user.
* `time_zone` - The timezone of the found user.
* `description` - The human-friendly description of the found user.

[1]: https://developer.pagerduty.com/api-reference/e35802f3c4ba4-list-members-of-a-team
