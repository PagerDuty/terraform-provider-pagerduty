---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_service"
sidebar_current: "docs-pagerduty-datasource-service"
description: |-
  Get information about a service that you have created.
---

# pagerduty\_service

Use this data source to get information about a specific [service][1].

## Example Usage

```hcl
data "pagerduty_service" "example" {
  name = "My Service"
}

data "pagerduty_vendor" "datadog" {
  name = "Datadog"
}

resource "pagerduty_service_integration" "example" {
  name    = "Datadog Integration"
  vendor  = data.pagerduty_vendor.datadog.id
  service = data.pagerduty_service.example.id
  type    = "generic_events_api_inbound_integration"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The service name to use to find a service in the PagerDuty API.

## Attributes Reference

* `id` - The ID of the found service.
* `name` - The short name of the found service.
* `type` - The type of object. The value returned will be `service`. Can be used for passing to a service dependency.
* `auto_resolve_timeout` - Time in seconds that an incident is automatically resolved if left open for that long. Value is null if the feature is disabled. Value must not be negative. Setting this field to 0, null (or unset) will disable the feature.
* `acknowledgement_timeout` - Time in seconds that an incident changes to the Triggered State after being Acknowledged. Value is null if the feature is disabled. Value must not be negative. Setting this field to 0, null (or unset) will disable the feature.
* `alert_creation` - Whether a service creates only incidents, or both alerts and incidents. A service must create alerts in order to enable incident merging.
* `description` - The user-provided description of the service.
* `escalation_policy` - The escalation policy associated with this service.
* `teams` - The set of teams associated with the service.

[1]: https://api-reference.pagerduty.com/#!/Services/get_services
