---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_service_integration"
sidebar_current: "docs-pagerduty-datasource-service_integration"
description: |-
  Get information about a service integration.
---

# pagerduty\_service\_integration

Use this data source to get information about a specific service_integration.

## Example Usage

```hcl
data "pagerduty_service_integration" "example" {
  service_name = "My Service"
  integration_summary = "Datadog"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service name to use to find a service in the PagerDuty API.
* `integration_summary` - (Required) The integration summary used to find the desired integration on the service.

## Attributes Reference

* `integration_key` - The integration key for the integration. This can be used to configure alerts.
