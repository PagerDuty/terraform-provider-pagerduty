---
layout: 'pagerduty'
page_title: 'PagerDuty: pagerduty_event_orchestration_service_cache_variable'
sidebar_current: 'docs-pagerduty-resource-event-orchestration-service-cache-variable'
description: |-
  Creates and manages a Cache Variable for a Service Event Orchestration.
---

# pagerduty_event_orchestration_service_cache_variable

A [Cache Variable][1] can be created on a Service Event Orchestration, in order to temporarily store event data to be referenced later within the Service Event Orchestration

## Example of configuring a Cache Variable for a Service Event Orchestration

This example shows creating a service `Event Orchestration` and a `Cache Variable`. This Cache Variable will count and store the number of trigger events with 'database' in its title. Then all alerts sent to this Event Orchestration will have its severity upped to 'critical' if the count has reached at least 5 triggers within the last 1 minute.

```hcl
resource "pagerduty_team" "database_team" {
  name = "Database Team"
}

resource "pagerduty_user" "user_1" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
  teams = [pagerduty_team.database_team.id]
}

resource "pagerduty_escalation_policy" "db_ep" {
  name      = "Database Escalation Policy"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user"
      id   = pagerduty_user.user_1.id
    }
  }
}

resource "pagerduty_service" "svc" {
  name                    = "My Database Service"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.db_ep.id
  alert_creation          = "create_alerts_and_incidents"
}

resource "pagerduty_event_orchestration_service_cache_variable" "num_db_triggers" {
  service = pagerduty_service.svc.id
  name = "num_db_triggers"

  condition {
    expression = "event.summary matches part 'database'"
  }

  configuration {
    type = "trigger_event_count"
    ttl_seconds = 60
  }
}

resource "pagerduty_event_orchestration_service" "event_orchestration" {
  service = pagerduty_service.svc.id
  enable_event_orchestration_for_service = true

  set {
    id = "start"
    rule {
      label = "Set severity to critical if we see at least 5 triggers on the DB within the last 1 minute"
      condition {
        expression = "cache_var.num_db_triggers >= 5"
      }
      actions {
        severity = "critical"
      }
    }
  }

  catch_all {
    actions { }
  }
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required) ID of the Service Event Orchestration to which this Cache Variable belongs.
* `name` - (Required) Name of the Cache Variable associated with the Service Event Orchestration.
* `disabled` - (Optional) Indicates whether the Cache Variable is disabled and would therefore not be evaluated.
* `condition` - Conditions to be evaluated in order to determine whether or not to update the Cache Variable's stored value.
  * `expression`- A [PCL condition][2] string.
* `configuration` - A configuration object to define what and how values will be stored in the Cache Variable.
  * `type` - The [type of value][1] to store into the Cache Variable. Can be one of: `recent_value` or `trigger_event_count`.
  * `source` - The path to the event field where the `regex` will be applied to extract a value. You can use any valid [PCL path][3]. This field is only used when `type` is `recent_value`
  * `regex` - A [RE2 regular expression][4] that will be matched against the field specified via the `source` argument. This field is only used when `type` is `recent_value`
  * `ttl_seconds` - The number of seconds indicating how long to count incoming trigger events for. This field is only used when `type` is `trigger_event_count`

## Attributes Reference

The following attributes are exported:

- `id` - ID of this Cache Variable.

## Import

Cache Variables can be imported using colon-separated IDs, which is the combination of the Service Event Orchestration ID followed by the Cache Variable ID, e.g.

```
$ terraform import pagerduty_event_orchestration_service_cache_variable.cache_variable PLBP09X:138ed254-3444-44ad-8cc7-701d69def439
```

[1]: https://support.pagerduty.com/docs/event-orchestration-variables
[2]: https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview
[3]: https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview#paths
[4]: https://github.com/google/re2/wiki/Syntax
