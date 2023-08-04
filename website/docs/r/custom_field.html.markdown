---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_custom_field"
sidebar_current: "docs-pagerduty-resource-custom-field"
description: |-
  Creates and manages a custom field in PagerDuty.
---

# pagerduty\_custom\_field

!> This Resource is no longer functional. Documentation is left here for the purpose of documenting migration steps.

A [Custom Field](https://support.pagerduty.com/docs/custom-fields) is a resuable element which can be added to Custom Field Schemas.

## Migration

The [`incident_custom_field`](./incident_custom_field.html.markdown) resource provides similar functionality
with largely the same arguments and attributes. The key distinction is that while custom fields created by this data source
may have only applied to a subset of incidents within the account after being added to a schema and assigned to a service,
custom fields managed by the `incident_custom_field` resource are applied to all incidents in the account.

Additionally:
* The separate `multi_value` and `fixed_options` arguments have been merged into a single argument
named `field_type`.
* The `datatype` argument has been renamed `data_type` to match the Public API for the Custom Fields on Incidents feature.

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
