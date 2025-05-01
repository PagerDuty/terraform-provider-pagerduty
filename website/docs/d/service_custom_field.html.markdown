---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_service_custom_field"
sidebar_current: "docs-pagerduty-datasource-service-custom-field"
description: |-
  Get information about a Service Custom Field in PagerDuty.
---

# pagerduty\_service\_custom\_field

Use this data source to get information about a specific Service Custom Field that has been configured in your PagerDuty account.

## Example Usage

```hcl
data "pagerduty_service_custom_field" "regions" {
  display_name = "AWS Regions"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The human-readable name of the field to look up. This must be unique across an account.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the found field.
* `name` - The name of the field. Contains ASCII characters, specifically lowercase letters, digits, and underscores.
* `type` - API object type.
* `summary` - A short-form, server-generated string that provides succinct, important information about the field.
* `self` - The API show URL at which the object is accessible.
* `description` - A description of the data this field contains.
* `data_type` - The kind of data the custom field is allowed to contain.
* `field_type` - The type of data this field contains. In combination with the data_type field.
* `default_value` - The default value for the custom field, if any.
* `enabled` - Whether the field is enabled.
* `field_options` - The options for the custom field. Only applies to `single_value_fixed` and `multi_value_fixed` field types. Each field option contains:
  * `value` - The value of the field option.
  * `data_type` - The data type of the field option.
