---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_business_service_subscriber"
sidebar_current: "docs-pagerduty-resource-business-service-subscriber"
description: |-
  Creates and manages a business service subscriber in PagerDuty.
---

# pagerduty\_business\_service_subscriber

A [business service subscriber](https://developer.pagerduty.com/api-reference/b3A6NDUwNDgxOQ-list-business-service-subscribers) allows you to subscribe users or teams to automatically receive updates about key business services.

## Example Usage

```hcl
resource "pagerduty_business_service" "example" {
  name             = "My Web App"
  description      = "A very descriptive description of this business service"
  point_of_contact = "PagerDuty Admin"
  team             = "P37RSRS"
}
resource "pagerduty_team" "engteam" {
  name = "Engineering"
}
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
}
resource "pagerduty_business_service_subscriber" "team_example" {
  subscriber_id = pagerduty_team.engteam.id
  subscriber_type = "team"
  business_service_id = pagerduty_business_service.example.id
}
resource "pagerduty_business_service_subscriber" "user_example" {
  subscriber_id = pagerduty_user.example.id
  subscriber_type = "user"
  business_service_id = pagerduty_business_service.example.id
}
```

## Argument Reference

The following arguments are supported:

  * `subscriber_id` - (Required) The ID of the subscriber entity.
  * `subscriber_type` - (Required) Type of subscriber entity in the subscriber assignment. Possible values can be `user` and `team`.
  * `business_service_id` - (Required) The ID of the business service to subscribe to.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the business service subscriber assignment.

## Import

Services can be imported using the `id` using the related business service ID, the subscriber type and the subscriber ID separated by a dot, e.g.

```
$ terraform import pagerduty_business_service_subscriber.main PLBP09X.team.PLBP09X
```
