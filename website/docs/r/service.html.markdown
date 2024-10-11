---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_service"
sidebar_current: "docs-pagerduty-resource-service"
description: |-
  Creates and manages a service in PagerDuty.
---

# pagerduty\_service

A [service](https://developer.pagerduty.com/api-reference/b3A6Mjc0ODE5Nw-create-a-service) represents something you monitor (like a web service, email service, or database service). It is a container for related incidents that associates them with escalation policies.


## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "Engineering Escalation Policy"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user_reference"
      id   = pagerduty_user.example.id
    }
  }
}

resource "pagerduty_service" "example" {
  name                    = "My Web App"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.foo.id
  alert_creation          = "create_alerts_and_incidents"

  auto_pause_notifications_parameters {
    enabled = true
    timeout = 300
  }
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the service.
  * `description` - (Optional) A human-friendly description of the service.
    If not set, a placeholder of "Managed by Terraform" will be set.
  * `auto_resolve_timeout` - (Optional) Time in seconds that an incident is automatically resolved if left open for that long. Disabled if set to the `"null"` string.
  * `acknowledgement_timeout` - (Optional) Time in seconds that an incident changes to the Triggered State after being Acknowledged. Disabled if set to the `"null"` string.  If not passed in, will default to '"1800"'.
  * `escalation_policy` - (Required) The escalation policy used by this service.
  * `response_play` - (Optional) (Deprecated) The response play used by this service.
  * `alert_creation` - (Optional) (Deprecated) This attribute has been deprecated as all services will be migrated to use alerts and incidents. The incident only service setting will be no longer available and this attribute will be removed in an upcoming version. See knowledge base for details https://support.pagerduty.com/docs/alerts#enable-and-disable-alerts-on-a-service. 
  * `alert_grouping` - (Optional) (Deprecated) Defines how alerts on this service will be automatically grouped into incidents. Note that the alert grouping features are available only on certain plans. If not set, each alert will create a separate incident; If value is set to `time`: All alerts within a specified duration will be grouped into the same incident. This duration is set in the `alert_grouping_timeout` setting (described below). Available on Standard, Enterprise, and Event Intelligence plans; If value is set to `intelligent` - Alerts will be intelligently grouped based on a machine learning model that looks at the alert summary, timing, and the history of grouped alerts. Available on Enterprise and Event Intelligence plan. This field is deprecated, use `alert_grouping_parameters.type` instead,
  * `alert_grouping_timeout` - (Optional) (Deprecated) The duration in minutes within which to automatically group incoming alerts. This setting applies only when `alert_grouping` is set to `time`. To continue grouping alerts until the incident is resolved, set this value to `0`. This field is deprecated, use `alert_grouping_parameters.config.timeout` instead,
  * `alert_grouping_parameters` - (Optional) (Deprecated) Defines how alerts on this service will be automatically grouped into incidents. Note that the alert grouping features are available only on certain plans. If not set, each alert will create a separate incident. Instructions on how to migrate this configuration to `pagerduty_alert_grouping_setting` resource can be found [here](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/resources/alert_grouping_setting#migration-from-alert_grouping_parameters).
  * `auto_pause_notifications_parameters` - (Optional) Defines how alerts on this service are automatically suspended for a period of time before triggering, when identified as likely being transient. Note that automatically pausing notifications is only available on certain plans as mentioned [here](https://support.pagerduty.com/docs/auto-pause-incident-notifications).

The `alert_grouping_parameters` block contains the following arguments:

* `type` (Optional) - The type of alert grouping; one of `intelligent`, `time` or `content_based`.
* `config` (Optional) - Alert grouping parameters dependent on `type`. If `type` is set to `intelligent` or empty then `config` can be empty.
    * `timeout` - (Optional) The duration in minutes within which to automatically group incoming alerts. This setting applies only when `type` is set to `time`. To continue grouping alerts until the incident is resolved, set this value to `0`.
    * `aggregate` - (Optional) One of `any` or `all`. This setting applies only when `type` is set to `content_based`. Group alerts based on one or all of `fields` value(s).
    * `fields` - (Optional) Alerts will be grouped together if the content of these fields match. This setting applies only when `type` is set to `content_based`.
    * `time_window` - (Optional) The maximum amount of time allowed between Alerts. This setting applies only when `type` is set to `intelligent` or `content_based`. Value must be between `300` and `3600` or exactly `86400` (86400 is supported only for `content_based` alert grouping). Any Alerts arriving greater than `time_window` seconds apart will not be grouped together. This is a rolling time window and is counted from the most recently grouped alert. The window is extended every time a new alert is added to the group, up to 24 hours.

The `auto_pause_notifications_parameters` block contains the following arguments:

* `enabled` (Optional) - Indicates whether alerts should be automatically suspended when identified as transient.  If not passed in, will default to 'false'.
* `timeout` (Optional) - Indicates in seconds how long alerts should be suspended before triggering. Allowed values: `120`, `180`, `300`, `600`, `900` if `enabled` is `true`. Must be omitted or set to `null` if `enabled` is `false`.


You may specify one optional `incident_urgency_rule` block configuring what urgencies to use.
Your PagerDuty account must have the `urgencies` ability to assign an incident urgency rule.
The block contains the following arguments:

  * `type` - The type of incident urgency: `constant` or `use_support_hours` (when depending on specific support hours; see `support_hours`).
  * `urgency` - The urgency: `low` Notify responders (does not escalate), `high` (follows escalation rules) or `severity_based` Set's the urgency of the incident based on the severity set by the triggering monitoring tool.
  * `during_support_hours` - (Optional) Incidents' urgency during support hours.
  * `outside_support_hours` - (Optional) Incidents' urgency outside support hours.

When using `type = "use_support_hours"` in `incident_urgency_rule` you must specify exactly one (otherwise optional) `support_hours` block.
Your PagerDuty account must have the `service_support_hours` ability to assign support hours.
The block contains the following arguments:

  * `type` - The type of support hours. Can be `fixed_time_per_day`.
  * `time_zone` - The time zone for the support hours.
  * `days_of_week` - Array of days of week as integers. `1` to `7`, `1` being
    Monday and `7` being Sunday.
  * `start_time` - The support hours' starting time of day.
  * `end_time` - The support hours' ending time of day.

A `scheduled_actions` block is required when using `type = "use_support_hours"` in `incident_urgency_rule`.

The block contains the following arguments:

  * `type` - The type of scheduled action. Currently, this must be set to `urgency_change`.
  * `to_urgency` - The urgency to change to: `low` (does not escalate), or `high` (follows escalation rules).
  * `at` - A block representing when the scheduled action will occur.

The `at` block contains the following arguments:
  * `type` - The type of time specification. Currently, this must be set to `named_time`.
  * `name` - Designates either the start or the end of the scheduled action. Can be `support_hours_start` or `support_hours_end`.

Note that it is currently only possible to define the scheduled action when urgency is set to `high` for `during_support_hours` and to `low`  for `outside_support_hours` in `incident_urgency_rule`.

Below is an example for a `pagerduty_service` resource with `incident_urgency_rules` with `type = "use_support_hours"`, `support_hours` and a default `scheduled_action` as well.

```hcl
resource "pagerduty_service" "foo" {
  name                    = "bar"
  description             = "bar bar bar"
  auto_resolve_timeout    = 3600
  acknowledgement_timeout = 3600
  escalation_policy       = pagerduty_escalation_policy.foo.id

  incident_urgency_rule {
    type = "use_support_hours"

    during_support_hours {
      type    = "constant"
      urgency = "high"
    }

    outside_support_hours {
      type    = "constant"
      urgency = "low"
    }
  }

  support_hours {
    type         = "fixed_time_per_day"
    time_zone    = "America/Lima"
    start_time   = "09:00:00"
    end_time     = "17:00:00"
    days_of_week = [1, 2, 3, 4, 5]
  }

  scheduled_actions {
    type       = "urgency_change"
    to_urgency = "high"

    at {
      type = "named_time"
      name = "support_hours_start"
    }
  }
}
```

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the service.
  * `last_incident_timestamp`- Last incident timestamp of the service.
  * `created_at`- Creation timestamp of the service.
  * `status`- The status of the service.
  * `html_url`- URL at which the entity is uniquely displayed in the Web app.
  * `type` - The type of object. The value returned will be `service`. Can be used for passing to a service dependency.

## Import

Services can be imported using the `id`, e.g.

```
$ terraform import pagerduty_service.main PLBP09X
```
