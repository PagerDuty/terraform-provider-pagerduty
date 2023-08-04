---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_licenses"
sidebar_current: "docs-pagerduty-datasource-licenses"
description: |-
  Get information about the account's purchased licenses for management of PagerDuty user resources
---

# pagerduty\_licenses

Use this data source to get information about the purchased [licenses][1] that you can use for other managing PagerDuty user resources. To reference a unique license, see `pagerduty_license` [data source][2]. After applying changes to users' licenses, the `current_value` and `allocations_available` attributes of licenses will change.

## Example Usage

```hcl
locals {
  invalid_roles = ["owner"]
}

data "pagerduty_licenses" "licenses" {}

resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"

  license = data.pagerduty_licenses.licenses.licenses[0].id

  # Role must be included in the assigned license's allowed_roles list.
  # Role may be dynamically referenced from data.pagerduty_licenses.licenses with the following:
  # tolist(setsubtract(data.pagerduty_licenses.licenses.licenses[0].valid_roles, local.invalid_roles))[0]
  role = "user"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Allows to override the default behavior for setting the `id` attribute that is required for data sources.

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
[2]: https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/data-sources/pagerduty_license
