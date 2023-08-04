---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_custom_field_schema_assignment"
sidebar_current: "docs-pagerduty-resource-custom-field-schema-assignment"
description: |-
  Creates and manages a custom field schema assignment in PagerDuty.
---

# pagerduty\_custom\_field\_schema\_assignment

!> This Resource is no longer functional. Documentation is left here for the purpose of documenting migration steps.

A [Custom Field Schema Assignment](https://support.pagerduty.com/docs/custom-fields#associate-schemas-with-services) is a relationship between a Custom Field Schema and a Service.

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

data "pagerduty_service" "first_service" {
  name = "My First Service"
}

resource "pagerduty_custom_field_schema_assignment" "assignment" {
  schema  = pagerduty_custom_field_schema.my_schema.id
  service = data.pagerduty_service.first_service.id
}
```

## Argument Reference

The following arguments are supported:

  * `schema` - (Required) The id of the field schema.
  * `service` - (Required) The id of the service.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the field schema assignment.
