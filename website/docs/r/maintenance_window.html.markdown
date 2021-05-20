---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_maintenance_window"
sidebar_current: "docs-pagerduty-resource-maintenance-window"
description: |-
  Creates and manages a maintenance window in PagerDuty.
---

# pagerduty_maintenance_window

A [maintenance window](https://v2.developer.pagerduty.com/v2/page/api-reference#!/Maintenance_Windows/get_maintenance_windows) is used to temporarily disable one or more services for a set period of time. No incidents will be triggered and no notifications will be received while a service is disabled by a maintenance window.

Maintenance windows are specified to start at a certain time and end after they have begun. Once started, a maintenance window cannot be deleted; it can only be ended immediately to re-enable the service.


## Example Usage

```hcl
resource "pagerduty_maintenance_window" "example" {
  start_time  = "2015-11-09T20:00:00-05:00"
  end_time    = "2015-11-09T22:00:00-05:00"
  services    = [pagerduty_service.example.id]
}
```

## Argument Reference

The following arguments are supported:

  * `start_time`  - (Required) The maintenance window's start time. This is when the services will stop creating incidents. If this date is in the past, it will be updated to be the current time.
  * `end_time`    - (Required) The maintenance window's end time. This is when the services will start creating incidents again. This date must be in the future and after the `start_time`.
  * `services`    - (Required) A list of service IDs to include in the maintenance window.
  * `description` - (Optional) A description for the maintenance window.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the maintenance window.


## Import

Maintenance windows can be imported using the `id`, e.g.

```
$ terraform import pagerduty_maintenance_window.main PLBP09X
```
