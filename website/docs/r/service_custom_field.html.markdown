---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_service_custom_field"
sidebar_current: "docs-pagerduty-resource-service-custom-field"
description: |-
  Creates and manages a service custom field in PagerDuty.
---

# pagerduty\_service\_custom\_field

A [service custom field](https://developer.pagerduty.com/api-reference/6075929031f7d-create-a-field)
allows you to extend PagerDuty Services with custom data fields to provide
additional context and support features such as customized filtering, search,
and analytics.

## Example Usage

```hcl
# Simple string field
resource "pagerduty_service_custom_field" "environment" {
  name         = "environment"
  display_name = "Environment"
  data_type    = "string"
  field_type   = "single_value"
  description  = "The environment this service runs in"
}

# Field with fixed options
resource "pagerduty_service_custom_field" "deployment_tier" {
  name          = "deployment_tier"
  display_name  = "Deployment Tier"
  data_type     = "string"
  field_type    = "single_value_fixed"
  description   = "The deployment tier of the service"

  field_option {
    value     = "production"
    data_type = "string"
  }

  field_option {
    value     = "staging"
    data_type = "string"
  }

  field_option {
    value     = "development"
    data_type = "string"
  }
}

# Multi-value field with fixed options
resource "pagerduty_service_custom_field" "regions" {
  name         = "regions"
  display_name = "AWS Regions"
  data_type    = "string"
  field_type   = "multi_value_fixed"
  description  = "AWS regions where this service is deployed"

  field_option {
    value     = "us-east-1"
    data_type = "string"
  }

  field_option {
    value     = "us-west-1"
    data_type = "string"
  }
}

# Boolean field
resource "pagerduty_service_custom_field" "critical" {
  name          = "is_critical"
  display_name  = "Is Critical"
  data_type     = "boolean"
  field_type    = "single_value"
  description   = "Whether this is a critical service"
}

# Integer field
resource "pagerduty_service_custom_field" "priority" {
  name          = "priority_level"
  display_name  = "Priority Level"
  data_type     = "integer"
  field_type    = "single_value"
  description   = "Service priority level"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the field. May include ASCII characters, specifically lowercase letters, digits, and underscores. Must be unique and cannot be changed once created.
* `display_name` - (Required) The human-readable name of the field. Must be unique across an account.
* `data_type` - (Required) The kind of data the custom field is allowed to contain. Can be one of: `string`, `integer`, `float`, `boolean`, `datetime`, or `url`.
* `field_type` - (Required) The type of field. Must be one of: `single_value`, `single_value_fixed`, `multi_value`, or `multi_value_fixed`.
* `description` - (Optional) A description of the data this field contains.
* `enabled` - (Optional) Whether the field is enabled. Defaults to `true`.
* `field_option` - (Optional) Configuration block for defining options for `single_value_fixed` or `multi_value_fixed` field types. Can be specified multiple times for multiple options.
  * `value` - (Required) The value of the option.
  * `data_type` - (Required) Must be `string`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service custom field.

## Import

Service custom fields can be imported using the field ID, e.g.

```
$ terraform import pagerduty_service_custom_field.example P123456
```
