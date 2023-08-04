---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_custom_field_schema"
sidebar_current: "docs-pagerduty-resource-custom-field-schema"
description: |-
  Creates and manages a custom field schema in PagerDuty.
---

# pagerduty\_custom\_field\_schema

!> This Resource is no longer functional. Documentation is left here for the purpose of documenting migration steps.

A [Custom Field Schema](https://support.pagerduty.com/docs/custom-fields#schemas) is a set of Custom Fields which can be set on an incident.

## Migration

This resource has no currently functional counterpart. Custom Fields on Incidents are now applied globally
to incidents within an account and have no notion of a Field Schema.

## Example Usage

```hcl
resource "pagerduty_custom_field" "cs_impact" {
  name      = "impact"
  datatype  = "string"
}

resource "pagerduty_custom_field_schema" "my_schema" {
  title       = "My Schema"
  description = "Fields used on incidents"
}
```

## Argument Reference

The following arguments are supported:

  * `title` - (Required) The title of the field schema.
  * `description` - (Optional) The description of the field schema.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the field schema.

## Import

Fields schemas can be imported using the `id`, e.g.

```
$ terraform import pagerduty_custom_field_schema.my_schema PLBP09X
```
