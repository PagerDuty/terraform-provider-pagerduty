---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_custom_field_schema"
sidebar_current: "docs-pagerduty-datasource-custom-field_schema"
description: |-
  Get information about a custom field schema that you can assign to services.
---

# pagerduty\_custom\_field\_schema

Use this data source to get information about a specific [Custom Field Schema](https://support.pagerduty.com/docs/custom-fields#schemas) that you can assign to a service.

-> The Custom Fields feature is currently available in Early Access.

## Example Usage

```hcl
data "pagerduty_custom_field_schema" "myschema" {
  title = "myschema title"
}

data "pagerduty_service" "first_service" {
  name = "My Service"
}

resource "pagerduty_custom_field_schema_assignment" "foo" {
  schema  = data.pagerduty_custom_field_schema.myschema.id
  service = data.pagerduty_service.first_service.id
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) The title of the field schema.

## Attributes Reference

* `id` - The ID of the found field schema.
