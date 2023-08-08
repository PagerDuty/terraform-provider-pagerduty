---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_service_status"
sidebar_current: "docs-pagerduty-datasource-service-status"
description: |-
  Get status information about a service that you have created.
---

# pagerduty\_service\_status

Use this data source to get information about a specific [service][1]'s status.

## Example Usage

```hcl
data "pagerduty_service_status" "example" {
  name = "My Service"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The service name to use to find a service in the PagerDuty API.

## Attributes Reference

* `id` - The ID of the found service.
* `name` - The short name of the found service.
* `last_incident_timestamp` - Last incident timestamp of the service.
* `status` - The status of the service.



[1]: https://api-reference.pagerduty.com/#!/Services/get_services
