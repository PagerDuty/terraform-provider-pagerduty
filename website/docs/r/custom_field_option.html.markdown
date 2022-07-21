---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_custom_field_option"
sidebar_current: "docs-pagerduty-resource-custom-field-option"
description: |-
  Creates and manages a custom field option in PagerDuty.
---

# pagerduty\_custom\_field\_option

A Custom Field Option is a specific value that can be used for [Custom Fields](https://support.pagerduty.com/docs/custom-fields) that only allow values from a set of fixed option. 

-> The Custom Fields feature is currently available in Early Access.

## Example Usage

```hcl
resource "pagerduty_custom_field" "sre_environment" {
  name          = "environment"
  datatype      = "string"
  fixed_options = true
}

resource "pagerduty_custom_field_option" "dev_environment" {
  field    = pagerduty_custom_field.sre_environment.id
  datatype = "string"
  value    = "dev"
}

resource "pagerduty_custom_field_option" "stage_environment" {
  field    = pagerduty_custom_field.sre_environment.id
  datatype = "string"
  value    = "stage"
}

resource "pagerduty_custom_field_option" "prod_environment" {
  field    = pagerduty_custom_field.sre_environment.id
  datatype = "string"
  value    = "prod"
}
```

## Argument Reference

The following arguments are supported:

* `field` - (Required) The ID of the field.
* `datatype` - (Required) The datatype of the field option. Must be one of `string`, `integer`, `float`, `boolean`, `datetime`, or `url`.
* `value` - (Required) The allowed value.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the field option.
