---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_custom_field_schema_field_configuration"
sidebar_current: "docs-pagerduty-resource-custom-field-schema-field-configuration"
description: |-
  Creates and manages a custom field configuration in PagerDuty.
---

# pagerduty\_custom\_field\_schema\_field\_configuration

A [Custom Field Configuration](https://support.pagerduty.com/docs/custom-fields#associate-schemas-with-services) is a declaration of a specific Custom Field in a specific Custom Field Schema.

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

resource "pagerduty_custom_field_schema_field_configuration" "first_field_configuration" {
  schema                 = pagerduty_custom_field_schema.my_schema.id
  field                  = pagerduty_custom_field.cs_impact.id
  required               = true
  default_value          = "none"
  default_value_datatype = "string"
}
```

## Argument Reference

The following arguments are supported:

* `field` - (Required) The ID of the field.
* `schema` - (Required) The ID of the schema.
* `required` - (Optional) True if the field is required
* `default_value` - (Optional) The default value for the field.
* `default_value_datatype` - (Optional) The datatype of the default value.
* `default_value_multi_value` - (Optional) Whether or not the default value is multi-valued.

#### Required and Default Value

Although `required`, `default_value`, and `default_value_datatype` are all
technically optional, they must be provided together, i.e. if `required` is `true`,
there **must** be a `default_value` and a `default_value_datatype`.

The `default_value_multi_value` attribute is only required if it is `true`, i.e. these two fragments
are semantically identical:

```hcl
    field                     = "ID1"
    schema                    = "ID2"
    required                  = true
    default_value             = "foo"
    default_value_datatype    = "string"
```

```hcl
    field                     = "ID1"
    schema                    = "ID2"
    required                  = true
    default_value             = "foo"
    default_value_datatype    = "string"
    default_value_multi_value = false
```

Default values are always strings. When providing a default value for a multi-valued field, the default value
needs to be a JSON-encoded array, e.g.

```hcl
    field                     = "ID1"
    required                  = true
    default_value             = jsonencode(["foo", "bar"])
    default_value_datatype    = "string"
    default_value_multi_value = true
```


```hcl
    field                     = "ID1"
    required                  = true
    default_value             = jsonencode([50, 60])
    default_value_datatype    = "string"
    default_value_multi_value = true
```

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the field configuration.

