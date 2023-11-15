---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_incident_custom_field_option"
sidebar_current: "docs-pagerduty-resource-incident-custom-field-option"
description: |-
  Creates and manages an field option for an Incident Custom Field in PagerDuty.
---

# pagerduty\_incident\_custom\_field\_option

A Incident Custom Field Option is a specific value that can be used for an [Incident Custom Field](https://support.pagerduty.com/docs/custom-fields-on-incidents) that only allow values from a set of fixed options,
i.e. has the `field_type` of `single_value_fixed` or `multi_value_fixed`.

## Example Usage

```hcl
resource "pagerduty_incident_custom_field" "sre_environment" {
  name         = "environment"
  display_name = "Environment"
  data_type    = "string"
  field_type   = "single_value_fixed"
}

resource "pagerduty_incident_custom_field_option" "dev_environment" {
  field     = pagerduty_incident_custom_field.sre_environment.id
  data_type = "string"
  value     = "dev"
}

resource "pagerduty_incident_custom_field_option" "stage_environment" {
  field    = pagerduty_incident_custom_field.sre_environment.id
  data_type = "string"
  value    = "stage"
}

resource "pagerduty_incident_custom_field_option" "prod_environment" {
  field    = pagerduty_incident_custom_field.sre_environment.id
  data_type = "string"
  value    = "prod"
}
```

## Argument Reference

The following arguments are supported:

* `field` - (Required) The ID of the field.
* `data_type` - (Required) The datatype of the field option. Only `string` is allowed here at present.
* `value` - (Required) The allowed value.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the field option.
