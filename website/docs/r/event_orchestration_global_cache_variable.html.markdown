---
layout: 'pagerduty'
page_title: 'PagerDuty: pagerduty_event_orchestration_global_cache_variable'
sidebar_current: 'docs-pagerduty-resource-event-orchestration-global-cache-variable'
description: |-
  Creates and manages a Cache Variable for a Global Event Orchestration.
---

# pagerduty_event_orchestration_global_cache_variable

A [Cache Variable][1] can be created on a Global Event Orchestration, in order to temporarily store event data to be referenced later within the Global Event Orchestration

## Example of configuring a Cache Variable for a Global Event Orchestration

This example shows creating a global `Event Orchestration` and a `Cache Variable`. All events that have the `event.source` field will have its `source` value stored in this Cache Variable, and appended as a note for the subsequent incident created by this Event Orchestration.

```hcl
resource "pagerduty_team" "database_team" {
  name = "Database Team"
}

resource "pagerduty_event_orchestration" "event_orchestration" {
  name = "Example Orchestration"
  team = pagerduty_team.database_team.id
}

resource "pagerduty_event_orchestration_global_cache_variable" "cache_var" {
  event_orchestration = pagerduty_event_orchestration.event_orchestration.id
  name = "recent_host"

  condition {
    expression = "event.source exists"
  }

  configuration {
    type = "recent_value"
    source = "event.source"
    regex = ".*"
  }
}

resource "pagerduty_event_orchestration_global" "global" {
  event_orchestration = pagerduty_event_orchestration.event_orchestration.id
  set {
    id = "start"
    rule {
      label = "Always annotate the incident with the event source for all events"
      actions {
        annotate = "Last time, we saw this incident occur on host: {{cache_var.recent_host}}"
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

* `event_orchestration` - (Required) ID of the Global Event Orchestration to which this Cache Variable belongs.
* `name` - (Required) Name of the Cache Variable associated with the Global Event Orchestration.
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

Cache Variables can be imported using colon-separated IDs, which is the combination of the Global Event Orchestration ID followed by the Cache Variable ID, e.g.

```
$ terraform import pagerduty_event_orchestration_global_cache_variable.cache_variable 5e7110bf-0ee7-429e-9724-34ed1fe15ac3:138ed254-3444-44ad-8cc7-701d69def439
```

[1]: https://support.pagerduty.com/docs/event-orchestration-variables
[2]: https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview
[3]: https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview#paths
[4]: https://github.com/google/re2/wiki/Syntax
