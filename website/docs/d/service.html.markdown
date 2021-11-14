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

[1]: https://api-reference.pagerduty.com/#!/Services/get_services
