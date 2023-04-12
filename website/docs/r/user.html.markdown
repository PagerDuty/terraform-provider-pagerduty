---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_user"
sidebar_current: "docs-pagerduty-resource-user"
description: |-
  Creates and manages a user in PagerDuty.
---

# pagerduty\_user

A [user][1] is a member of a PagerDuty account that have the ability to interact with incidents and other data on the account.


## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the user.
  * `email` - (Required) The user's email address.
  * `color` - (Optional) The schedule color for the user. Valid options are purple, red, green, blue, teal, orange, brown, turquoise, dark-slate-blue, cayenne, orange-red, dark-orchid, dark-slate-grey, lime, dark-magenta, lime-green, midnight-blue, deep-pink, dark-green, dark-orange, dark-cyan, darkolive-green, dark-slate-gray, grey20, firebrick, maroon, crimson, dark-red, dark-goldenrod, chocolate, medium-violet-red, sea-green, olivedrab, forest-green, dark-olive-green, blue-violet, royal-blue, indigo, slate-blue, saddle-brown, or steel-blue.
  * `role` - (Optional) The user role. Can be `admin`, `limited_user`, `observer`, `owner`, `read_only_user`, `read_only_limited_user`, `restricted_access`, or `user`.
     Notes:
    * Account must have the `read_only_users` ability to set a user as a `read_only_user` or a `read_only_limited_user`, and must have advanced permissions abilities to set a user as `observer` or `restricted_access`.
    * With advanced permissions, users can have both a user role (base role) and a team role. The team role can be configured in the `pagerduty_team_membership` resource.
    * Mapping of `role` values to Web UI user role names available in the [user roles support page](https://support.pagerduty.com/docs/advanced-permissions#roles-in-the-rest-api-and-saml).
  * `job_title` - (Optional) The user's title.
  * `teams` - (Optional, **DEPRECATED**) A list of teams the user should belong to. Please use `pagerduty_team_membership` instead.
  * `time_zone` - (Optional) The time zone of the user. Default is account default timezone.
  * `description` - (Optional) A human-friendly description of the user.
    If not set, a placeholder of "Managed by Terraform" will be set.
  * `license` - (Optional) The license id assigned to the user. If provided the user's role must exist in the assigned license's `valid_roles` list. To reference purchased licenses' ids see data source `pagerduty_licenses` [data source][1].

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the user.
  * `avatar_url` - The URL of the user's avatar.
  * `time_zone` - The timezone of the user.
  * `html_url` - URL at which the entity is uniquely displayed in the Web app
  * `invitation_sent` - If true, the user has an outstanding invitation.

## Import

Users can be imported using the `id`, e.g.

```
$ terraform import pagerduty_user.main PLBP09X
```

[1]: https://developer.pagerduty.com/api-reference/b3A6Mjc0ODIzNA-create-a-user
[2]: https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/data-sources/pagerduty_license
