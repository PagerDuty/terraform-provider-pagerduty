---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_tag"
sidebar_current: "docs-pagerduty-datasource-tag"
description: |-
  Get information about a tag that you can use to assign to users, teams, and escalation_policies.
---

# pagerduty\_tag

Use this data source to get information about a specific [tag][1] that you can use to assign to users, teams, and escalation_policies.

## Example Usage

```hcl
data "pagerduty_user" "me" {
  email = "me@example.com"
}

data "pagerduty_tag" "devops" {
  label = "devops"
}

resource "pagerduty_tag_assignment" "foo" {
  tag_id      = data.pagerduty_tag.devops.id
  entity_id   = data.pagerduty_user.me.id
  entity_type = "users"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the tag to find in the PagerDuty API.

## Attributes Reference

* `id` - The ID of the found team.

[1]: https://developer.pagerduty.com/api-reference/b3A6Mjc0ODIxNw-list-tags
