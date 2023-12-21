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
* `role` - The team role of the found user.
* `type` - The type of object. The value returned will be `user_reference`. Can be used for passing to another object as dependency.
* `summary` - A short-form, server-generated string that provides succinct, important information about an object suitable for primary labeling of an entity in a client. In many cases, this will be identical to name, though it is not intended to be an identifier.

[1]: https://developer.pagerduty.com/api-reference/e35802f3c4ba4-list-members-of-a-team
