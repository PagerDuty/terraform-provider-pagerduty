---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_custom_field"
sidebar_current: "docs-pagerduty-datasource-custom-field"
description: |-
  Get information about a custom field that you can add to a custom field schema.
---

# pagerduty\_custom\_field

Use this data source to get information about a specific [Custom Field](https://support.pagerduty.com/docs/custom-fields) that you can add to a custom field schema.

-> The Custom Fields feature is currently available in Early Access.

## Example Usage

```hcl
data "pagerduty_custom_field" "sre_environment" {
  name      = "environment"
}

resource "pagerduty_custom_field_schema" "foo" {
  title       = "myschema"
  description = "some description"
  field {
    field = data.pagerduty_custom_field.sre_environment.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the field.

## Attributes Reference

* `id` - The ID of the found field.
