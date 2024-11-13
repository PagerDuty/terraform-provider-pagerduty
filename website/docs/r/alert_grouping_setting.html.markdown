---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_alert_grouping_setting"
sidebar_current: "docs-pagerduty-resource-alert-grouping-setting"
description: |-
  Creates and manages an alert grouping setting in PagerDuty.
---

# pagerduty\_alert\_grouping\_setting

An [alert grouping setting](https://developer.pagerduty.com/api-reference/create-an-alert-grouping-setting)
stores and centralize the configuration used during grouping of the alerts.

## Example Usage

```hcl
data "pagerduty_escalation_policy" "default" {
	name = "Default"
}

resource "pagerduty_service" "basic" {
	name = "Example"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_alert_grouping_setting" "%[1]s" {
  name = "Configuration for type-1 devices"
  type = "content_based"
  services = [pagerduty_service.basic.id]
  config {
    time_window = 300
    aggregate = "all"
    fields = ["fields"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the alert groupig settings.
* `description` - A human-friendly text to describe and identify this alert grouping setting.
* `type` - (Required) The type of alert grouping; one of `intelligent`, `time`, `content_based` or  `content_based_intelligent`.
* `services` - (Required)  [Updating can cause a resource replacement] The list IDs of services associated to this setting.
* `config` - (Required) The set of values used for configuration.

The `config` block contains the following arguments:

* `timeout` - (Optional) The duration in minutes within which to automatically group incoming alerts. This setting is only required and applies when `type` is set to `time`. To continue grouping alerts until the incident is resolved leave this value unset or set it to `null`.
* `aggregate` - (Optional) One of `any` or `all`. This setting is only required and applies when `type` is set to `content_based` or `content_based_intelligent`. Group alerts based on one or all of `fields` value(s).
* `fields` - (Optional) Alerts will be grouped together if the content of these fields match. This setting is only required and applies when `type` is set to `content_based` or `content_based_intelligent`.
* `time_window` - (Optional) The maximum amount of time allowed between Alerts. This setting applies only when `type` is set to `intelligent`, `content_based`, `content_based_intelligent`. Value must be between `300` and `3600` or exactly `86400` (86400 is supported only for `content_based` alert grouping). Any Alerts arriving greater than `time_window` seconds apart will not be grouped together. This is a rolling time window and is counted from the most recently grouped alert. The window is extended every time a new alert is added to the group, up to 24 hours. To use the recommended time window leave this value unset or set it to `null`.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the alert grouping setting.

## Migration from `alert_grouping_parameters`

To migrate from using the field `alert_grouping_parameters` of a
[service](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/resources/service)
to a `pagerduty_alert_grouping_setting` resource, you can cut-and-paste the
contents of an `alert_grouping_parameters` field from a `pagerduty_service`
resource into the new resource, but you also need to add at least one value in
the field `services` to create the alert grouping setting with a service
associated to it.

If you are using `timeout = 0` or `time_window = 0` in order to use the values
recommended by PagerDuty you also need to set its value to null or delete it,
since a value of `0` is no longer accepted.

Since the `alert_grouping_parameters` field creates an Alert Grouping Setting
behind the scenes, it is necessary to import them if you want to keep your
configuration the same as it is right now.

**Example**:

Before:
```
data "pagerduty_escalation_policy" "default" {
    name = "Default"
}

resource "pagerduty_service" "foo" {
    name              = "Foo"
    escalation_policy = data.pagerduty_escalation_policy.default.id
    alert_grouping_parameters {
        type = "time"
        config {
            timeout = 0
        }
    }
}
```

After:
```
data "pagerduty_escalation_policy" "default" {
    name = "Default"
}

resource "pagerduty_service" "foo" {
    name              = "Foo"
    escalation_policy = data.pagerduty_escalation_policy.default.id
}

data "pagerduty_alert_grouping_setting" "foo_alert" {
    name = "Foo"
}

import {
  id = data.pagerduty_alert_grouping_setting.foo_alert.id
  to = pagerduty_alert_grouping_setting.foo_alert
}

resource "pagerduty_alert_grouping_setting" "foo_alert" {
    name = "Alert Grouping for Foo-like services"
    type = "time"
    config {
        time = null
    }
    services = [pagerduty_service.foo.id]
}
```

But if you prefer to have a clean restart, you can do it in two steps: delete
the current `alert_grouping_parameters` and later create a new
`alert_grouping_setting` associated to your resources now free to be associated
with this Alert Grouping Setting.

**Example**:

Before:
```
data "pagerduty_escalation_policy" "default" {
    name = "Default"
}

resource "pagerduty_service" "foo" {
    name              = "Foo"
    escalation_policy = data.pagerduty_escalation_policy.default.id
    alert_grouping_parameters {
        type = "content_based"
        config {
            time_window = 300
            aggregate = "all"
            fields = ["summary"]
        }
    }
}

resource "pagerduty_service" "bar" {
    name              = "Bar"
    escalation_policy = data.pagerduty_escalation_policy.default.id
    alert_grouping_parameters {
        type = "content_based"
        config {
            time_window = 300
            aggregate = "all"
            fields = ["summary"]
        }
    }
}
```

Step 1:
```
data "pagerduty_escalation_policy" "default" {
    name = "Default"
}

resource "pagerduty_service" "foo" {
    name              = "Foo"
    escalation_policy = data.pagerduty_escalation_policy.default.id
    alert_grouping_parameters {}
}

resource "pagerduty_service" "bar" {
    name              = "Bar"
    escalation_policy = data.pagerduty_escalation_policy.default.id
    alert_grouping_parameters {}
}
```

Step 2:
```
data "pagerduty_escalation_policy" "default" {
    name = "Default"
}

resource "pagerduty_service" "foo" {
    name              = "Foo"
    escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_service" "bar" {
    name              = "Bar"
    escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_alert_grouping_setting" "type_a" {
    name = "Type A"
    description = "Configuration used for all services of type A"
    type = "content_based"
    config {
        time_window = 300
        aggregate = "all"
        fields = ["summary"]
    }
    services = [
        pagerduty_service.foo.id,
        pagerduty_service.bar.id,
    ]
}
```

## Import

Alert grouping settings can be imported using its `id`, e.g.

```
$ terraform import pagerduty_alert_grouping_setting.example P3DH5M6
```
