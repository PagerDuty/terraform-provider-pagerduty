---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_user"
sidebar_current: "docs-pagerduty-datasource-user"
description: |-
  Get information about users of your PagerDuty account as a list, optionally filtered by team ids that you can use for a service integration (e.g Amazon Cloudwatch, Splunk, Datadog).
---

# pagerduty\_users

Use this data source to get information about [list of users][1] that you can use for other PagerDuty resources, optionally filtering by team ids.

## Example Usage

```hcl
data "pagerduty_team" "devops" {
  name = "devops"
}

data "pagerduty_user" "me" {
  email = "me@example.com"
}

resource "pagerduty_user" "example_w_team" {
  name = "user-with-team"
  email = "user-with-team@example.com"
}

resource "pagerduty_team_membership" "example" {
  team_id = pagerduty_team.devops.id
  user_id = pagerduty_user.example_w_team.id
}

data "pagerduty_users" "all_users" {}

data "pagerduty_users" "from_devops_team" {
  depends_on = [pagerduty_team_membership.example]
  team_ids = [pagerduty_team.devops.id]
}
```

## Argument Reference

The following arguments are supported:

* `team_ids` - (Optional) List of team IDs. Only results related to these teams will be returned. Account must have the `teams` ability to use this parameter.

## Attributes Reference
* `id` - The ID of queried list of users.
* `users` - List of users queried.

### Users (`users`) supports the following:

* `id` - The ID of the found user.
* `name` - The short name of the found user.
* `email` - The email of the found user.
* `role` - The role of the found user.
* `job_title` - The job title of the found user.
* `time_zone` - The timezone of the found user.
* `description` - The human-friendly description of the found user.

[1]: https://developer.pagerduty.com/api-reference/b3A6Mjc0ODIzMw-list-users
