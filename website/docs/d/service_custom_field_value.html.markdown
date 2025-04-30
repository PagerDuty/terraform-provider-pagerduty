---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_service_custom_field_value"
sidebar_current: "docs-pagerduty-datasource-service-custom-field-value"
description: |-
  Get information about service custom field values in PagerDuty.
---

# pagerduty\_service\_custom\_field\_value

Use this data source to get information about service custom field values in PagerDuty.

## Example Usage

```hcl
data "pagerduty_service_custom_field_value" "example" {
  service_id = pagerduty_service.example.id
}

output "environment_value" {
  value = [
    for field in data.pagerduty_service_custom_field_value.example.custom_fields : 
    field.value if field.name == "environment"
  ][0]
}

# Create a service
resource "pagerduty_service" "example" {
  name                    = "Example Service"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.example.id
}

# Set custom field values on the service
resource "pagerduty_service_custom_field_value" "example" {
  service_id = pagerduty_service.example.id
  
  custom_fields {
    name  = "environment"
    value = jsonencode("production")
  }
  
  custom_fields {
    name  = "region"
    value = jsonencode("us-east-1")
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) The ID of the service to get custom field values for.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service.
* `custom_fields` - A list of custom field values associated with the service. Each element contains:
  * `id` - The ID of the custom field.
  * `name` - The name of the custom field.
  * `display_name` - The human-readable name of the custom field.
  * `description` - A description of the data this field contains.
  * `data_type` - The kind of data the custom field is allowed to contain. Can be one of: `string`, `integer`, `float`, `boolean`, `datetime`, or `url`.
  * `field_type` - The type of field. Can be one of: `single_value`, `single_value_fixed`, `multi_value`, or `multi_value_fixed`.
  * `type` - The type of the reference, typically "field_value".
  * `value` - The value of the custom field. This is a JSON-encoded string matching the field's data type.
