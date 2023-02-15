---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_custom_field"
sidebar_current: "docs-pagerduty-resource-custom-field"
description: |-
  Creates and manages a custom field in PagerDuty.
---

# pagerduty\_custom\_field

A [Custom Field](https://support.pagerduty.com/docs/custom-fields) is a resuable element which can be added to Custom Field Schemas.

-> The Custom Fields feature is currently available in Early Access.

## Example Usage

```hcl
resource "pagerduty_custom_field" "cs_impact" {
  name      = "impact"
  datatype  = "string"
}

resource "pagerduty_custom_field" "sre_environment" {
  name          = "environment"
  datatype      = "string"
  fixed_options = true
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the field.
  * `datatype` - (Required) The datatype of the field. Must be one of `string`, `integer`, `float`, `boolean`, `datetime`, or `url`.
  * `multi_value` - (Optional) True if the field can accept multiple values.
  * `fixed_options` - (Optional) True if the field can only accept values from a set of options.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the field.

## Import

Fields can be imported using the `id`, e.g.

```
$ terraform import pagerduty_custom_field.sre_environment PLBP09X
```
