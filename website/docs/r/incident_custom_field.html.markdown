---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_incident_custom_field"
sidebar_current: "docs-pagerduty-resource-incident-custom-field"
description: |-
  Creates and manages an Incident Custom Field in PagerDuty.
---

# pagerduty\_incident\_custom\_field

An [Incident Custom Field](https://support.pagerduty.com/docs/custom-fields-on-incidents) defines a field which can be set on incidents in the target account.

## Example Usage

```hcl
resource "pagerduty_incident_custom_field" "cs_impact" {
  name         = "impact"
  display_name = "Customer Impact"
  data_type    = "string"
  field_type   = "single_value"
}

resource "pagerduty_incident_custom_field" "sre_environment" {
  name         = "environment"
  display_name = "Environment"
  data_type    = "string"
  field_type   = "single_value_fixed"
}

resource "pagerduty_incident_custom_field" "false_alarm" {
  name          = "false_alarm"
  display_name  = "False Alarm"
  data_type     = "boolean"
  field_type    = "single_value"
  default_value = "false"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the field.
  * `display_name` - (Required) The display name of the field.
  * `description` - (Optional) The description of the field.
  * `data_type` - (Required) The data type of the field. Must be one of `string`, `integer`, `float`, `boolean`, `datetime`, or `url`.
  * `field_type` - (Required) The field type of the field. Must be one of `single_value`, `single_value_fixed`, `multi_value`, or `multi_value_fixed`.
  * `default_value` - (Optional) The default value to set when new incidents are created. Always specified as a string.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the field.

## Import

Fields can be imported using the `id`, e.g.

```
$ terraform import pagerduty_incident_custom_field.sre_environment PLBP09X
```
