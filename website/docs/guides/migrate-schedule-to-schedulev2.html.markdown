---
layout: "pagerduty"
page_title: "Migrating from pagerduty_schedule to pagerduty_schedulev2"
sidebar_current: "docs-pagerduty-guides-migrate-schedule-to-schedulev2"
description: |-
  Step-by-step guide for migrating on-call schedules from the legacy pagerduty_schedule resource to the new pagerduty_schedulev2 resource backed by the PagerDuty v3 Schedules API.
---

# Migrating from `pagerduty_schedule` to `pagerduty_schedulev2`

`pagerduty_schedule` is deprecated and will be removed in a future provider release. `pagerduty_schedulev2` is its replacement, backed by the PagerDuty v3 Schedules API.

---

## Feasibility Assessment

| Pattern | Migratable? |
|---|---|
| Single-user, 24/7 coverage | ✅ Yes |
| Single-user with daily time restrictions (e.g. 09:00–17:00 every day) | ✅ Yes |
| Single-user with weekly restrictions (e.g. MON–FRI business hours) | ✅ Yes |
| Multi-user layers (`users` list with more than one entry) | ❌ Not supported — v2 uses static per-event assignment |
| Multiple overlapping layers (`final_schedule` merging) | ⚠️ Redesign needed — v2 rotations are independent |
| `teams` association | ❌ Not yet supported in v2 |
| `overflow` parameter | ❌ No equivalent in v3 API |

~> **Before starting:** audit each `layer.users` list. Any layer with more than one user requires architectural redesign and cannot be migrated directly.

-> **Every `pagerduty_schedule` can be migrated.** Simple single-user configurations map directly using the field reference below. For complex configurations (multi-user layers, overlapping restrictions), use the [PagerDuty REST API `GET /schedules/{id}`](https://developer.pagerduty.com/api-reference/3f03afb2c84a4-get-a-schedule) to retrieve the computed `final_schedule`. The `final_schedule` contains the resolved, non-overlapping on-call windows that `pagerduty_schedulev2` events map to directly, making the conversion straightforward to automate with a script.

---

## Concept Mapping

The two resources use different abstractions for on-call coverage.

| Concept | `pagerduty_schedule` (v1) | `pagerduty_schedulev2` (v2) |
|---|---|---|
| Core unit | `layer` — users rotate through turns | `rotation` + `event` — explicit time window with fixed assignment |
| Recurrence | Implied by `rotation_turn_length_seconds` | Explicit RFC 5545 RRULE string |
| Time restriction | `restriction` block limits when a turn is active | `start_time`/`end_time` + `recurrence` define the window directly |
| Layer lifecycle | Layers can only be ended, never deleted | Rotations and events are fully mutable |

---

## Field Reference

### Schedule-level

| `pagerduty_schedule` | `pagerduty_schedulev2` | Notes |
|---|---|---|
| `name` | `name` | v1 is Optional; v2 is Required |
| `time_zone` | `time_zone` | v1 rejects `"UTC"` — use `"Etc/UTC"`. v2 accepts `"UTC"` directly. Prefer the full IANA name (e.g. `"America/New_York"`) in both. |
| `description` | `description` | v1 defaults to `"Managed by Terraform"`; v2 has no default |
| `overflow` | — | No equivalent |
| `teams` | — | Not yet supported in v2 |
| `final_schedule` (computed) | — | Not exposed in v2 |

### Layer → Rotation + Event

| `layer` attribute | `event` attribute | Notes |
|---|---|---|
| `name` | `name` | Direct mapping |
| `start` | `effective_since` | When the layer/event begins producing shifts |
| `end` | `effective_until` | Omit for indefinite coverage |
| `rotation_virtual_start` | `start_time` | First occurrence start time |
| `rotation_turn_length_seconds` | `recurrence` | Replaced by RRULE (e.g. weekly turn → `RRULE:FREQ=WEEKLY`) |
| `users[0]` | `assignment_strategy.member.user_id` | Single user maps to one `member` block |

### Restriction → Recurrence

`end_time = start_time + duration_seconds`

| v1 restriction | v2 `recurrence` RRULE |
|---|---|
| `daily_restriction` at 09:00, 8 h | `RRULE:FREQ=DAILY` with `start_time=T09:00`, `end_time=T17:00` |
| `weekly_restriction` Mon–Fri (5 blocks) | `RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR` (single event) |
| `weekly_restriction` Sat + Sun (or one block, `duration=172740s`) | `RRULE:FREQ=WEEKLY;BYDAY=SA,SU` with `end_time` on Sun 23:59 |

---

## Migration Procedure

**Step 1 — Translate the config.** For each `layer`, create one `rotation` with one `event` per distinct time window using the field mapping above.

**Step 2 — Add the new resource without removing the old one.**

```bash
terraform plan   # 1 to add, 0 to change, 0 to destroy
terraform apply
```

**Step 3 — Verify in the PagerDuty UI.** Confirm that both schedules produce identical on-call windows (start time, recurrence, user).

**Step 4 — Remove the old schedule** once you have confirmed the new one is correct.

```bash
terraform plan   # 0 to add, 0 to change, 1 to destroy
terraform apply
```

---

## Worked Example

Business hours (MON–FRI 09:00–17:00 UTC) + weekend on-call (SAT–SUN 00:00–23:59 UTC).

### Before — `pagerduty_schedule`

```hcl
resource "pagerduty_schedule" "business_hours" {
  name      = "Business Hours On-Call"
  time_zone = "Etc/UTC"

  layer {
    name                         = "Weekday Business Hours"
    start                        = "2026-02-21T09:00:00Z"
    rotation_virtual_start       = "2026-02-21T09:00:00Z"
    rotation_turn_length_seconds = 604800
    users                        = [pagerduty_user.oncall.id]

    restriction {
      type = "weekly_restriction"; start_day_of_week = 1
      start_time_of_day = "09:00:00"; duration_seconds = 28800
    }
    restriction {
      type = "weekly_restriction"; start_day_of_week = 2
      start_time_of_day = "09:00:00"; duration_seconds = 28800
    }
    restriction {
      type = "weekly_restriction"; start_day_of_week = 3
      start_time_of_day = "09:00:00"; duration_seconds = 28800
    }
    restriction {
      type = "weekly_restriction"; start_day_of_week = 4
      start_time_of_day = "09:00:00"; duration_seconds = 28800
    }
    restriction {
      type = "weekly_restriction"; start_day_of_week = 5
      start_time_of_day = "09:00:00"; duration_seconds = 28800
    }
  }

  layer {
    name                         = "Weekend On-Call"
    start                        = "2026-02-22T00:00:00Z"
    rotation_virtual_start       = "2026-02-22T00:00:00Z"
    rotation_turn_length_seconds = 604800
    users                        = [pagerduty_user.oncall.id]

    restriction {
      type              = "weekly_restriction"
      start_day_of_week = 6
      start_time_of_day = "00:00:00"
      duration_seconds  = 172740   # 47h59m
    }
  }
}
```

### After — `pagerduty_schedulev2`

```hcl
resource "pagerduty_schedulev2" "business_hours" {
  name      = "Business Hours On-Call"
  time_zone = "UTC"   # v2 accepts "UTC" directly

  rotation {
    event {
      name            = "Weekday Business Hours"
      start_time      = "2026-02-21T09:00:00Z"
      end_time        = "2026-02-21T17:00:00Z"
      effective_since = "2026-02-21T09:00:00Z"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"]

      assignment_strategy {
        type = "user_assignment_strategy"
        member { type = "user_member"; user_id = pagerduty_user.oncall.id }
      }
    }
  }

  rotation {
    event {
      name            = "Weekend On-Call"
      start_time      = "2026-02-22T00:00:00Z"
      end_time        = "2026-02-23T23:59:00Z"
      effective_since = "2026-02-22T00:00:00Z"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=SA,SU"]

      assignment_strategy {
        type = "user_assignment_strategy"
        member { type = "user_member"; user_id = pagerduty_user.oncall.id }
      }
    }
  }
}
```

The five `weekly_restriction` blocks collapse into one RRULE `BYDAY` clause. The weekend `duration_seconds = 172740` becomes an explicit `end_time` on Sunday 23:59.

---

## Known Limitations

**`effective_since` adjusted by the API.** The v3 API moves a past `effective_since` to the current time silently. Always set it to a future date to avoid state drift.

**`start_time`/`end_time` UTC normalization.** The v3 API normalizes these to UTC in responses. The provider preserves the original config value when both represent the same instant, so offset values like `"2026-06-01T09:00:00-05:00"` will not cause perpetual diffs.

**`daily_restriction` cannot use `start_day_of_week`.** v1 validation rejects this combination. For weekday-only daily windows in v2, use `RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR` instead of `RRULE:FREQ=DAILY`.

**Importing an existing v3 schedule.**

```bash
terraform import pagerduty_schedulev2.example P1234AB
```

After import, `start_time`, `end_time`, and `effective_since` may differ from config due to UTC normalization and API time adjustment.
