---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_business_service"
sidebar_current: "docs-pagerduty-resource-business-service"
description: |-
  Creates and manages a business service in PagerDuty.
---

# pagerduty\_business\_service

A [business service](https://v2.developer.pagerduty.com/v2/page/api-reference#!/Business_Services/get_business_services) allows you to model capabilities that span multiple technical services and that may be owned by several different teams. 


## Example Usage

```hcl
resource "pagerduty_business_service" "example" {
  name             = "My Web App"
  description      = "A very descriptive description of this business service"
  point_of_contact = "PagerDuty Admin"
  team             = "P37RSRS"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the business service.
  * `description` - (Optional) A human-friendly description of the service.
    If not set, a placeholder of "Managed by Terraform" will be set.
  * `point_of_contact` - (Optional) The owner of the business service. 
  * `type` - (Optional) Default value is `business_service`. Can also be set as `business_service_reference`.
  * `team` - (Optional) ID of the team that owns the business service.
  
## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the service.
  * `summary`- A short-form, server-generated string that provides succinct, important information about an object suitable for primary labeling of an entity in a client. In many cases, this will be identical to `name`, though it is not intended to be an identifier.
  * `html_url`- A URL at which the entity is uniquely displayed in the Web app.
  * `self`- The API show URL at which the object is accessible.

## Import

Services can be imported using the `id`, e.g.

```
$ terraform import pagerduty_business_service.main PLBP09X
```
