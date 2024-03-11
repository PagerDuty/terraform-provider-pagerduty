---
layout: 'pagerduty'
page_title: 'PagerDuty: pagerduty_event_orchestration_service_cache_variable'
sidebar_current: 'docs-pagerduty-datasource-event-orchestration-service-cache-variable'
description: |-
  Get information about a Cache Variable for a Service Event Orchestration that you have created.
---

# pagerduty_event_orchestration_service_cache_variable

Use this data source to get information about a specific [Cache Variable][1] for a Service Event Orchestration.

## Example Usage

```hcl

resource "pagerduty_team" "engineering" {
  name = "Engineering"
}

resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
  teams = [pagerduty_team.engineering.id]
}

resource "pagerduty_escalation_policy" "example" {
  name      = "Engineering Escalation Policy"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user"
      id   = pagerduty_user.example.id
    }
  }
}

resource "pagerduty_service" "service" {
  name                    = "My Web App"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.example.id
  alert_creation          = "create_alerts_and_incidents"
}

data "pagerduty_event_orchestration_service_cache_variable" "cache_variable" {
  service = pagerduty_service.service.id
  name = "example_cache_variable"
}

```

## Argument Reference

The following arguments are supported:

* `service` - (Required) ID of the Service Event Orchestration to which this Cache Variable belongs.
* `id` - (Optional) ID of the Cache Variable associated with the Service Event Orchestration. Specify either `id` or `name`. If both are specified `id` takes precedence.
* `name` - (Optional) Name of the Cache Variable associated with the Service Event Orchestration. Specify either `id` or `name`. If both are specified `id` takes precedence.

## Attributes Reference

* `disabled` - Indicates whether the Cache Variable is disabled and would therefore not be evaluated.
* `condition` - Conditions to be evaluated in order to determine whether or not to update the Cache Variable's stored value.
  * `expression`- A [PCL condition][2] string.
* `configuration` - A configuration object to define what and how values will be stored in the Cache Variable.
  * `type` - The [type of value][1] to store into the Cache Variable. Can be one of: `recent_value` or `trigger_event_count`.
  * `source` - The path to the event field where the `regex` will be applied to extract a value. You can use any valid [PCL path][3]. This field is only used when `type` is `recent_value`
  * `regex` - A [RE2 regular expression][4] that will be matched against the field specified via the `source` argument. This field is only used when `type` is `recent_value`
  * `ttl_seconds` - The number of seconds indicating how long to count incoming trigger events for. This field is only used when `type` is `trigger_event_count`


[1]: https://support.pagerduty.com/docs/event-orchestration-variables
[2]: https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview
[3]: https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview#paths
[4]: https://github.com/google/re2/wiki/Syntax
