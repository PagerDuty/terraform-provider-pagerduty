---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_incident_type_custom_field"
sidebar_current: "docs-pagerduty-resource-incident-type-custom-field"
description: |-
  Creates and manages a incident type custom field in PagerDuty.
---

# pagerduty\_incident\_type\_custom\_field

An [incident type custom fields](https://developer.pagerduty.com/api-reference/423b6701f3f1b-create-a-custom-field-for-an-incident-type)
are a feature which will allow customers to extend Incidents with their own
custom data, to provide additional context and support features such as
customized filtering, search and analytics. Custom Fields can be applied to
different incident types. Child types will inherit custom fields from their
parent types.


## Example Usage

```hcl
resource "pagerduty_incident_type_custom_field" "alarm_time" {
  name          = "alarm_time_minutes"
  display_name  = "Alarm Time"
  data_type     = "integer"
  field_type    = "single_value"
  default_value = jsonencode(5)
  incident_type = "incident_default"
}

data "pagerduty_incident_type" "foo" {
    display_name = "Foo"
}

resource "pagerduty_incident_type_custom_field" "level" {
  name          = "level"
  incident_type = data.pagerduty_incident_type.foo.id
  display_name  = "Level"
  data_type     = "string"
  field_type    = "single_value_fixed"
  field_options = ["Trace", "Debug", "Info", "Warn", "Error", "Fatal"]
}

resource "pagerduty_incident_type_custom_field" "cs_impact" {
  name          = "impact"
  incident_type = data.pagerduty_incident_type.foo.id
  display_name  = "Customer Impact"
  data_type     = "string"
  field_type    = "multi_value"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) [Updating causes resource replacement] The name of the custom field.
* `incident_type` - (Required) [Updating causes resource replacement] The id of the incident type the custom field is associated with.
* `display_name` - (Required) The display name of the custom Type.
* `data_type` - (Required) [Updating causes resource replacement] The type of the data of this custom field. Can be one of `string`, `integer`, `float`, `boolean`, `datetime`, or `url` when `data_type` is `single_value`, otherwise must be `string`. Update
* `field_type` - (Required) [Updating causes resource replacement] The field type of the field. Must be one of `single_value`, `single_value_fixed`, `multi_value`, or `multi_value_fixed`.
* `description` - The description of the custom field.
* `default_value` - The default value to set when new incidents are created. Always specified as a string.
* `enabled` - Whether the custom field is enabled. Defaults to true if not provided.
* `field_options` - The options for the custom field. Can only be applied to fields with a type of `single_value_fixed` or `multi_value_fixed`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the field.
* `self` - The API show URL at which the object is accessible.
* `summary` - A short-form, server-generated string that provides succinct, important information about an object suitable for primary labeling of an entity in a client. In many cases, this will be identical to name, though it is not intended to be an identifier.

## Import

Fields can be imported using the combination of `incident_type_id` and `field_id`, e.g.

```
$ terraform import pagerduty_incident_custom_field.cs_impact PT1234:PF1234
```
