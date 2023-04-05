---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_team_membership"
sidebar_current: "docs-pagerduty-resource-team-membership"
description: |-
  Creates and manages a team membership in PagerDuty.
---

# pagerduty_team_membership

A [team membership](https://developer.pagerduty.com/api-reference/b3A6Mjc0ODIzMg-add-a-user-to-a-team) manages memberships within a team.

-> This resource supports caching to improve performance in use cases when having Teams with 500 or more associations being managed via Terraform and a detrimental of the performance is noticed. So in order to overcome performance issues the **Cache** support can be activated. [Know more here...](https://github.com/PagerDuty/terraform-provider-pagerduty\#caching-support)

## Example Usage

```hcl
resource "pagerduty_user" "foo" {
  name = "foo"
  email = "foo@bar.com"
}

resource "pagerduty_team" "foo" {
  name        = "foo"
  description = "foo"
}

resource "pagerduty_team_membership" "foo" {
  user_id = pagerduty_user.foo.id
  team_id = pagerduty_team.foo.id
  role    = "manager"
}
```

## Argument Reference

The following arguments are supported:

  * `user_id` - (Required) The ID of the user to add to the team.
  * `team_id` - (Required) The ID of the team in which the user will belong.
  * `role`    - (Optional) The role of the user in the team. One of `observer`, `responder`, or `manager`. Defaults to `manager`.  
     These roles match up to user roles in the following ways:
    * User role of `user` is a Team role of `manager`
    * User role of `limited_user` is a Team role of `responder`

## Attributes Reference

The following attributes are exported:

  * `user_id` - The ID of the user belonging to the team.
  * `team_id` - The team ID the user belongs to.
  * `role`    - The role of the user in the team.


## Import

Team memberships can be imported using the `user_id` and `team_id`, e.g.

```
$ terraform import pagerduty_team_membership.main PLBP09X:PLB09Z
```
