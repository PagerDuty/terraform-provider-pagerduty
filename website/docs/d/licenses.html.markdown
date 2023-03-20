---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_licenses"
sidebar_current: "docs-pagerduty-datasource-licenses"
description: |-
  Get information about the account's purchased licenses for management of PagerDuty user resources
---

# pagerduty\_licenses

Use this data source to get information about the purchased [licenses][1] that you can use for other managing PagerDuty user resources.

## Example Usage

```hcl
locals {
	invalid_roles = ["owner"]
}

data "pagerduty_licenses" "licenses" {
  name = "licenses"
}

resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"

	license = {
		id = data.pagerduty_licenses.licenses.licenses[0].id
		type = "license_reference"
	}

  # Role must be included in the assigned license's allowed_roles list.
  # Role may be dynamically referenced from data.pagerduty_licenses.licenses with the following:
  # tolist(setsubtract(data.pagerduty_licenses.licenses.licenses[0].valid_roles, local.invalid_roles))[0]
	role = "user"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Used for referencing the data source.

## Attributes Reference
* `licenses` - The list of purchased licenses.

### Licenses (`licenses`) is a list of objects that support the following:
  * `id` - ID of the license
  * `name` - Name of the license
  * `summary` - Summary of the license
  * `description` - Description of the license
  * `role_group` - The role group for the license that determines the available `valid_roles`
  * `valid_roles` - List of allowed roles that may be assigned to a user with this license
  * `current_value` - The number of allocations already assigned to users
  * `allocations_available` - Available allocations to assign to users

[1]: https://developer.pagerduty.com/api-reference/4c10cb38f7381-list-licenses
