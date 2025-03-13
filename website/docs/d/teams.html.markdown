---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_teams"
sidebar_current: "docs-pagerduty-datasource-teams"
description: |-
  Get information about a team that you can use with escalation_policies, schedules etc.
---

# pagerduty\_teams

Use this data source to [list teams][1] in your PagerDuty account.

## Example Usage

```hcl
data "pagerduty_teams" "all_teams" {}

# Fetch only teams whose name matches "devops"
data "pagerduty_teams" "devops" {
  query = "devops"
}
```

## Argument Reference

The following arguments are supported:

* `query` - (Optional) Filters the result, showing only the records whose name matches the query.

## Attributes Reference

* `teams` - The teams found.

### Teams (`teams`) supports the following:

* `id` - The ID of the team.
* `name` - The name of the team.
* `summary` - A short-form, server-generated string that provides succinct, important information about an object suitable for primary labeling of an entity in a client. In many cases, this will be identical to name, though it is not intended to be an identifier.
* `description` - The description of the team.

[1]: https://developer.pagerduty.com/api-reference/0138639504311-list-teams
