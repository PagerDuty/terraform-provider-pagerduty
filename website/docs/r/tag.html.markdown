---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_tag"
sidebar_current: "docs-pagerduty-resource-tag"
description: |-
  Creates and manages a tag in PagerDuty.
---

# pagerduty\_tag

A [tag](https://developer.pagerduty.com/api-reference/b3A6Mjc0ODIxNw-list-tags) is applied to Escalation Policies, Teams or Users and can be used to filter them.

## Example Usage

```hcl
resource "pagerduty_tag" "example" {
  label = "Product"
  type = "tag"
}
```

## Argument Reference

The following arguments are supported:

  * `label` - (Required) The label of the tag.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the tag.
  * `summary`- A short-form, server-generated string that provides succinct, important information about an object suitable for primary labeling of an entity in a client. In many cases, this will be identical to name, though it is not intended to be an identifier.
  * `html_url` - URL at which the entity is uniquely displayed in the Web app

## Import

Tags can be imported using the `id`, e.g.

```
$ terraform import pagerduty_tag.main PLBP09X
```
