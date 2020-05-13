---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_schedule"
sidebar_current: "docs-pagerduty-resource-schedule"
description: |-
  Creates and manages a schedule in PagerDuty.
---

# pagerduty\_schedule

A [schedule](https://v2.developer.pagerduty.com/v2/page/api-reference#!/Schedules/get_schedules) determines the time periods that users are on call. Only on-call users are eligible to receive notifications from incidents.


## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
  teams = [pagerduty_team.example.id]
}

resource "pagerduty_schedule" "foo" {
  name      = "Daily Engineering Rotation"
  time_zone = "America/New_York"

  layer {
    name                         = "Night Shift"
    start                        = "2015-11-06T20:00:00-05:00"
    rotation_virtual_start       = "2015-11-06T20:00:00-05:00"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32400
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the schedule.
* `time_zone` - (Required) The time zone of the schedule (e.g Europe/Berlin).
* `description` - (Optional) The description of the schedule
* `layer` - (Required) A schedule layer block. Schedule layers documented below.
* `overflow` - (Optional) Any on-call schedule entries that pass the date range bounds will be truncated at the bounds, unless the parameter `overflow` is passed. For instance, if your schedule is a rotation that changes daily at midnight UTC, and your date range is from `2011-06-01T10:00:00Z` to `2011-06-01T14:00:00Z`:
If you don't pass the overflow=true parameter, you will get one schedule entry returned with a start of `2011-06-01T10:00:00Z` and end of `2011-06-01T14:00:00Z`.
If you do pass the `overflow` parameter, you will get one schedule entry returned with a start of `2011-06-01T00:00:00Z` and end of `2011-06-02T00:00:00Z`.


Schedule layers (`layer`) supports the following:

* `name` - (Optional) The name of the schedule layer.
* `start` - (Required) The start time of the schedule layer. This value will not be read back from the PagerDuty API because the API will always return a new `start` time, which represents the last updated time of the schedule layer.
* `end` - (Optional) The end time of the schedule layer. If not specified, the layer does not end.
* `rotation_virtual_start` - (Required) The effective start time of the schedule layer. This can be before the start time of the schedule.
* `rotation_turn_length_seconds` - (Required) The duration of each on-call shift in `seconds`.
* `users` - (Required) The ordered list of users on this layer. The position of the user on the list determines their order in the layer.
* `restriction` - (Optional) A schedule layer restriction block. Restriction blocks documented below.


Restriction blocks (`restriction`) supports the following:

* `type` - (Required) Can be `daily_restriction` or `weekly_restriction`
* `start_time_of_day` - (Required) The start time in `HH:mm:ss` format.
* `duration_seconds` - (Required) The duration of the restriction in `seconds`.
* `start_day_of_week` - (Required for `weekly_restriction`) Number of the day when restriction starts. From 1 to 7 where 1 is Monday and 7 is Sunday.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the schedule

## Import

Schedules can be imported using the `id`, e.g.

```
$ terraform import pagerduty_schedule.main PLBP09X
```
