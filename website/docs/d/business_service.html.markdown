---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_business_service"
sidebar_current: "docs-pagerduty-datasource-business-service"
description: |-
  Get information about a business service that you have created.
---

# pagerduty\_business\_service

Use this data source to get information about a specific [business service][1].

## Example Usage

```hcl
data "pagerduty_business_service" "example" {
  name = "My Service"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The business service name to use to find a business service in the PagerDuty API.

## Attributes Reference
* `id` - The ID of the found business service.
* `name` - The short name of the found business service.
* `type` - The type of object. The value returned will be `business_service`. Can be used for passing to a service dependency.

[1]: https://api-reference.pagerduty.com/#!/Business_Services/get_business_services
