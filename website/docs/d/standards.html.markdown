---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_standards"
sidebar_current: "docs-pagerduty-datasource-standards"
description: |-
  Get information about all standards of an account.
---

# pagerduty\_standards

Use this data source to get information about the [standards][1] applicable to
the PagerDuty account.

## Example Usage

```hcl
data "pagerduty_standards" "standards" {}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Optional) Filter by `resource_type` the received standards. Allowed values are `technical_service`.

## Attributes Reference

* `standards` - The list of standards defined.

### Standards (`standards`) is a list of objects that support the following:

  * `id` - A unique identifier for the standard.
  * `name` - The human-readable name of the standard.
  * `active` - Indicates whether the standard is currently active and applicable to the resource.
  * `description` - Provides a textual description of the standard.
  * `type` - The type of the standard.
  * `resource_type` - Specifies the type of resource to which the standard applies.
  * `exclusions` - A list of exceptions for the application of this standard.
    * `id` - The unique identifier for the resource being excluded.
    * `type` - Specifies the type of resource this exclusion applies to.
  * `inclusions` - A list of explict instances this standard applies to.
    * `id` - The unique identifier for the resource being included.
    * `type` - Specifies the type of resource this inclusion applies to.

[1]: https://developer.pagerduty.com/api-reference/dbed9a0ff9355-list-standards
