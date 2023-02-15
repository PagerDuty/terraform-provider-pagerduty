---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_custom_field_schema_assignment"
sidebar_current: "docs-pagerduty-resource-custom-field-schema-assignment"
description: |-
  Creates and manages a custom field schema assignment in PagerDuty.
---

# pagerduty\_custom\_field\_schema\_assignment

A [Custom Field Schema Assignment](https://support.pagerduty.com/docs/custom-fields#associate-schemas-with-services) is a relationship between a Custom Field Schema and a Service.

-> The Custom Fields feature is currently available in Early Access.

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
