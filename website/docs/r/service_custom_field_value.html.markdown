---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_service_custom_field_value"
sidebar_current: "docs-pagerduty-resource-service-custom-field-value"
description: |-
  Creates and manages custom field values for a service in PagerDuty.
---

# pagerduty\_service\_custom\_field\_value

A [service custom field value](https://developer.pagerduty.com/api-reference/6075929031f7d-update-custom-field-values)
allows you to set values for custom fields on a PagerDuty service. These values
provide additional context for services and can be used for filtering, search,
and analytics.

## Example Usage

```hcl
# First, create a service custom field
resource "pagerduty_service_custom_field" "environment" {
  name         = "environment"
  display_name = "Environment"
  data_type    = "string"
  field_type   = "single_value"
  description  = "The environment this service runs in"
}

# Create a service
resource "pagerduty_service" "example" {
  name                    = "Example Service"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.example.id
}

# Set a custom field value on the service
resource "pagerduty_service_custom_field_value" "example" {
  service_id = pagerduty_service.example.id
  
  custom_fields {
    name  = pagerduty_service_custom_field.environment.name
    value = jsonencode("production")
  }
}

# Set multiple custom field values on a service
resource "pagerduty_service_custom_field" "region" {
  name         = "region"
  display_name = "Region"
  data_type    = "string"
  field_type   = "single_value"
  description  = "The region this service is deployed in"
}

resource "pagerduty_service_custom_field_value" "multiple_example" {
  service_id = pagerduty_service.example.id
  
  custom_fields {
    name  = pagerduty_service_custom_field.environment.name
    value = jsonencode("production")
  }
  
  custom_fields {
    name  = pagerduty_service_custom_field.region.name
    value = jsonencode("us-east-1")
  }
}

# Example with a boolean field
resource "pagerduty_service_custom_field" "is_critical" {
  name         = "is_critical"
  display_name = "Is Critical"
  data_type    = "boolean"
  field_type   = "single_value"
  description  = "Whether this service is critical"
}

resource "pagerduty_service_custom_field_value" "boolean_example" {
  service_id = pagerduty_service.example.id
  
  custom_fields {
    name  = pagerduty_service_custom_field.is_critical.name
    value = jsonencode(true)
  }
}

# Example with a multi-value field
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

resource "pagerduty_service_custom_field_value" "multi_value_example" {
  service_id = pagerduty_service.example.id
  
  custom_fields {
    name  = pagerduty_service_custom_field.regions.name
    value = jsonencode(["us-east-1", "us-west-1"])
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) The ID of the service to set custom field values for.
* `custom_fields` - (Required) A list of custom field values to set on the service. Each block supports the following:
  * `id` - (Optional) The ID of the custom field. Either `id` or `name` must be provided.
  * `name` - (Optional) The name of the custom field. Either `id` or `name` must be provided.
  * `value` - (Required) The value to set for the custom field. Must be provided as a JSON-encoded string matching the field's data type. Use the `jsonencode()` function to ensure proper formatting.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service custom field value, which is the same as the service ID.

## Import

Service custom field values can be imported using the service ID, e.g.

```
$ terraform import pagerduty_service_custom_field_value.example PXYZ123
```
