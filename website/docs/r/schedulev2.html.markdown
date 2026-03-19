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

### Rotating member assignment

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
        type             = "rotating_member_assignment_strategy"
        shifts_per_member = 1

        member {
          type    = "user_member"
          user_id = pagerduty_user.example.id
        }
      }
    }
  }
}
```

### Every-member assignment (all members on-call simultaneously)

```hcl
resource "pagerduty_user" "primary" {
  name  = "Alice"
  email = "alice@example.com"
}

resource "pagerduty_user" "secondary" {
  name  = "Bob"
  email = "bob@example.com"
}

resource "pagerduty_schedulev2" "all_hands" {
  name      = "Weekend All-Hands On-Call"
  time_zone = "UTC"

  rotation {
    event {
      name            = "Weekend Coverage"
      start_time      = "2026-06-06T00:00:00Z"
      end_time        = "2026-06-07T23:59:00Z"
      effective_since = "2026-06-06T00:00:00Z"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=SA,SU"]

      assignment_strategy {
        type = "every_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.primary.id
        }

        member {
          type    = "user_member"
          user_id = pagerduty_user.secondary.id
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the schedule. Maximum 255 characters.
* `time_zone` - (Required) The time zone of the schedule (IANA format, e.g. `America/New_York`).
* `description` - (Optional) A description of the schedule. Maximum 1024 characters.
* `teams` - (Optional) List of team IDs to associate with this schedule.
* `rotation` - (Required) One or more rotation blocks. Rotations documented below.

---

Rotation blocks (`rotation`) support the following:

* `event` - (Required) One or more event blocks defining on-call periods within this rotation. Events documented below.

---

Event blocks (`event`) support the following:

* `name` - (Required) The name of the event. Maximum 255 characters.
* `start_time` - (Required) The shift start time in ISO-8601 format (e.g. `2026-06-01T09:00:00Z`). The v3 API normalizes this to UTC.
* `end_time` - (Required) The shift end time in ISO-8601 format. The v3 API normalizes this to UTC.
* `effective_since` - (Required) When this event configuration begins producing shifts (ISO-8601 UTC). The API adjusts past values to the current time.
* `effective_until` - (Optional) When this event configuration stops producing shifts (ISO-8601 UTC). Omit for an indefinite schedule.
* `recurrence` - (Required) List of RFC 5545 recurrence rule strings. Must contain exactly one `RRULE` entry. May optionally include one or more `EXDATE` entries (dates to exclude) and one or more `RDATE` entries (additional dates to include). Example: `["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"]`. You can generate RRULE strings interactively using tools like [RRULE Tool](https://icalendar.org/rrule-tool.html).
* `assignment_strategy` - (Required) A block defining how on-call responsibility is assigned. Assignment strategy documented below.

---

Assignment strategy blocks (`assignment_strategy`) support the following:

* `type` - (Required) The assignment strategy type. Supported values:
  * `"rotating_member_assignment_strategy"` — listed members rotate in sequence. Each member covers `shifts_per_member` consecutive shift periods before the next member takes over.
  * `"every_member_assignment_strategy"` — all listed members are on-call simultaneously for every occurrence.

  ~> **Breaking change:** The previous value `"user_assignment_strategy"` is no longer valid. Use `"rotating_member_assignment_strategy"` instead.

* `shifts_per_member` - (Optional) Number of consecutive shift occurrences each member covers before rotating. Minimum value: `1`. Required when `type` is `"rotating_member_assignment_strategy"`.
* `member` - (Required) One or more member blocks identifying who is on call. Required for both strategy types. Maximum 20 members. Members documented below.

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
