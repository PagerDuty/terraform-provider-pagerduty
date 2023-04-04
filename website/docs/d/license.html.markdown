---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_license"
sidebar_current: "docs-pagerduty-datasource-license"
description: |-
  Get information about one of the account's purchased licenses for management of PagerDuty user resources
---

# pagerduty\_license

Use this data source to use a single purchased [license][1] to manage PagerDuty user resources.

## Example Usage

```hcl
locals {
	invalid_roles = ["owner"]
}

data "pagerduty_license" "full_user" {
  name = "Full User"
  description = ""
}

resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"

	license = {
		id = data.pagerduty_license.full_user.id
		type = "license_reference"
	}

  # Role must be included in the assigned license's allowed_roles list.
  # Role may be dynamically referenced from data.pagerduty_license.full_user with the following:
  # tolist(setsubtract(data.pagerduty_license.full_user.valid_roles, local.invalid_roles))[0]
	role = "user"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Used to match the data config *id* with an exact match of a valid license ID assigned to the account.
* `name` - (Optional) Used to determine if the data config *name* is a valid substring of a valid license name assigned to the account.
* `description` - (Optional) Used to determine if the data config *description* is a valid substring of a valid license description assigned to the account.

## Attributes Reference
  * `id` - ID of the license
  * `name` - Name of the license
  * `summary` - Summary of the license
  * `description` - Description of the license
  * `role_group` - The role group for the license that determines the available `valid_roles`
  * `valid_roles` - List of allowed roles that may be assigned to a user with this license
  * `current_value` - The number of allocations already assigned to users
  * `allocations_available` - Available allocations to assign to users

[1]: https://developer.pagerduty.com/api-reference/4c10cb38f7381-list-licenses
