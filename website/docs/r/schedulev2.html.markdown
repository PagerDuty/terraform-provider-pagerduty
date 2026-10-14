---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_schedulev2"
sidebar_current: "docs-pagerduty-resource-schedulev2"
description: |-
  Creates and manages an on-call schedule using the PagerDuty v3 Schedules API.
---

# pagerduty\_schedulev2

A [v3 schedule](https://developer.pagerduty.com/api-reference/d90c4c94e3ce2-create-a-schedule) determines the time periods that users are on call using flexible rotation configurations. This resource uses the PagerDuty v3 Schedules API, which supports per-event assignment strategies and RFC 5545 recurrence rules.

## Schedule versions and resource naming

The Terraform resource names do not line up one-to-one with the API version numbers, which is a common source of confusion. The table below maps the two:

| Schedule type        | Terraform resource     | API version            |
| -------------------- | ---------------------- | ---------------------- |
| Legacy schedule      | `pagerduty_schedule`   | v2 (Accept header)     |
| Shift-based schedule | `pagerduty_schedulev2` | v3 (in the URL path)   |

`pagerduty_schedule` manages the legacy schedule (now deprecated); `pagerduty_schedulev2` manages the newer shift-based schedule. In the UI, legacy schedules are marked with a `legacy` tag on the schedules list page.

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
  time_zone = "Etc/UTC"

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

## Migrating from `pagerduty_schedule`

The legacy `pagerduty_schedule` resource models on-call coverage with `layer` blocks (a list of `users` rotating on a fixed `rotation_turn_length_seconds`, optionally constrained by `restriction` blocks). The shift-based `pagerduty_schedulev2` resource models the same coverage with `rotation` → `event` blocks, where an `assignment_strategy` decides how the listed `member`s cover each occurrence of an RFC 5545 `recurrence`.

The most common legacy shape is a single layer with several users handing off once per week (a 24/7 weekly rotation) and one or more team associations. The example below shows that shape before and after migration.

**Before** — legacy `pagerduty_schedule` (single `OnCall` layer, six users, weekly handoff, two teams):

```hcl
resource "pagerduty_schedule" "example_oncall" {
  name        = "Example OnCall Schedule"
  description = "A Example OnCall Schedule."
  time_zone   = "Europe/Amsterdam"

  layer {
    name                         = "OnCall"
    start                        = "2024-06-24T00:00:00-00:00"
    rotation_virtual_start       = "2025-03-17T07:00:00+02:00"
    rotation_turn_length_seconds = 604800 # one week
    users = [
      data.pagerduty_user.user1.id,
      data.pagerduty_user.user2.id,
      data.pagerduty_user.user3.id,
      data.pagerduty_user.user4.id,
      data.pagerduty_user.user5.id,
      data.pagerduty_user.user6.id,
    ]
  }

  teams = [
    pagerduty_team.k8s_platform.id,
    pagerduty_team.k8s_operational.id,
  ]
}
```

**After** — equivalent `pagerduty_schedulev2`:

```hcl
resource "pagerduty_schedulev2" "example_oncall" {
  name        = "Example OnCall Schedule"
  description = "A Example OnCall Schedule."
  time_zone   = "Europe/Amsterdam"

  teams = [
    pagerduty_team.k8s_platform.id,
    pagerduty_team.k8s_operational.id,
  ]

  rotation {
    event {
      name = "OnCall"

      # A 7-day window + weekly recurrence reproduces the legacy
      # rotation_turn_length_seconds = 604800 (one week per turn).
      start_time      = "2025-03-17T07:00:00+02:00"
      end_time        = "2025-03-24T07:00:00+02:00"
      effective_since = "2025-03-17T07:00:00+02:00"
      recurrence      = ["RRULE:FREQ=WEEKLY"]

      # The layer's user list becomes rotating members. shifts_per_member = 1
      # means each member covers one weekly occurrence before handing off.
      assignment_strategy {
        type              = "rotating_member_assignment_strategy"
        shifts_per_member = 1

        member {
          type    = "user_member"
          user_id = data.pagerduty_user.user1.id
        }
        member {
          type    = "user_member"
          user_id = data.pagerduty_user.user2.id
        }
        member {
          type    = "user_member"
          user_id = data.pagerduty_user.user3.id
        }
        member {
          type    = "user_member"
          user_id = data.pagerduty_user.user4.id
        }
        member {
          type    = "user_member"
          user_id = data.pagerduty_user.user5.id
        }
        member {
          type    = "user_member"
          user_id = data.pagerduty_user.user6.id
        }
      }
    }
  }
}
```

Field mapping:

| Legacy (`pagerduty_schedule`)         | Shift-based (`pagerduty_schedulev2`)                                  |
| ------------------------------------- | -------------------------------------------------------------------- |
| `layer`                               | `rotation` (one `rotation` per layer)                                |
| `layer.name`                          | `rotation.event.name`                                                |
| `layer.users`                         | `assignment_strategy.member` (one `member` per user)                 |
| `layer.rotation_turn_length_seconds`  | `event.start_time`/`end_time` window + `recurrence` (e.g. one week → 7-day window + `RRULE:FREQ=WEEKLY`) |
| `layer.rotation_virtual_start`        | `event.start_time` / `effective_since`                               |
| `layer.restriction`                   | a narrower `event` window plus a `BYDAY`/`BYHOUR` `RRULE`            |
| `teams`                               | `teams` (unchanged)                                                  |

~> **Note:** Multiple members rotate when `assignment_strategy.type` is `"rotating_member_assignment_strategy"`; use `"every_member_assignment_strategy"` when everyone should be on call simultaneously. To reproduce a legacy `restriction` (e.g. weekday business hours), narrow the `event` window and encode the days/hours in the `RRULE` — see the *Rotating member assignment* example above.

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
