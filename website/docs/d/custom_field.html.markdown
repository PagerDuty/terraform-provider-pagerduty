---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_custom_field"
sidebar_current: "docs-pagerduty-datasource-custom-field"
description: |-
  Get information about a custom field that you can add to a custom field schema.
---

# pagerduty\_custom\_field

!> This Data Source is no longer functional. Documentation is left here for the purpose of documenting migration steps.

Use this data source to get information about a specific [Custom Field](https://support.pagerduty.com/docs/custom-fields) that you can add to a custom field schema.

## Migration

The [`incident_custom_field`](./incident_custom_field.html.markdown) data source provides similar functionality
with the same arguments and attributes. The key distinction is that while custom fields returned by this data source
may have only applied to a subset of incidents within the account, custom fields returned by the `incident_custom_field`
data source are applied to all incidents in the account.

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
