---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_schedulev2"
sidebar_current: "docs-pagerduty-resource-schedulev2"
description: |-
  Creates and manages an on-call schedule using the PagerDuty v3 Schedules API.
---

# pagerduty\_schedulev2

A [v3 schedule](https://developer.pagerduty.com/api-reference/d90c4c94e3ce2-create-a-schedule) determines the time periods that users are on call using flexible rotation configurations. This resource uses the PagerDuty v3 Schedules API, which supports per-event assignment strategies and RFC 5545 recurrence rules.

~> **Note:** This resource requires the `flexible-schedules-early-access` early access flag on your PagerDuty account. The required `X-Early-Access` header is sent automatically by the provider.

## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "earline@example.com"
}

resource "pagerduty_schedulev2" "example" {
  name        = "Engineering On-Call"
  time_zone   = "America/New_York"
  description = "Managed by Terraform"

  rotation {
    event {
      name            = "Weekday Business Hours"
      start_time      = "2026-06-01T09:00:00Z"
      end_time        = "2026-06-01T17:00:00Z"
      effective_since = "2026-06-01T09:00:00Z"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"]

      assignment_strategy {
        type = "user_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.example.id
        }
      }
    }
  }

  rotation {
    event {
      name            = "Weekend On-Call"
      start_time      = "2026-06-06T00:00:00Z"
      end_time        = "2026-06-07T23:59:00Z"
      effective_since = "2026-06-06T00:00:00Z"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=SA,SU"]

      assignment_strategy {
        type = "user_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.example.id
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the schedule.
* `time_zone` - (Required) The time zone of the schedule (IANA format, e.g. `America/New_York`).
* `description` - (Optional) A description of the schedule.
* `rotation` - (Required) One or more rotation blocks. Rotations documented below.

---

Rotation blocks (`rotation`) support the following:

* `event` - (Required) One or more event blocks defining on-call periods within this rotation. Events documented below.

---

Event blocks (`event`) support the following:

* `name` - (Required) The name of the event.
* `start_time` - (Required) The shift start time in ISO-8601 format (e.g. `2026-06-01T09:00:00Z`). The v3 API normalizes this to UTC.
* `end_time` - (Required) The shift end time in ISO-8601 format. The v3 API normalizes this to UTC.
* `effective_since` - (Required) When this event configuration begins producing shifts (ISO-8601 UTC). The API adjusts past values to the current time.
* `effective_until` - (Optional) When this event configuration stops producing shifts (ISO-8601 UTC). Omit for an indefinite schedule.
* `recurrence` - (Required) List of recurrence rule strings in RFC 5545 RRULE format (e.g. `"RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"`). At least one rule is required. You can generate RRULE strings interactively using tools like [RRULE Tool](https://icalendar.org/rrule-tool.html).
* `assignment_strategy` - (Required) A block defining how on-call responsibility is assigned. Assignment strategy documented below.

---

Assignment strategy blocks (`assignment_strategy`) support the following:

* `type` - (Required) The assignment strategy type. Currently only `"user_assignment_strategy"` is supported.
* `member` - (Required) One or more member blocks identifying who is on call. Members documented below.

---

Member blocks (`member`) support the following:

* `type` - (Required) The member type. Supported values: `"user_member"`, `"empty_member"`.
* `user_id` - (Optional) The ID of the user to assign. Required when `type` is `"user_member"`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the schedule.
* `rotation.*.id` - The ID of each rotation.
* `rotation.*.event.*.id` - The ID of each event within a rotation.

## Import

Schedules can be imported using the schedule `id`, e.g.

```
$ terraform import pagerduty_schedulev2.example P1234AB
```
