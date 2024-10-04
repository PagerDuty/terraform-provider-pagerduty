---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_alert_grouping_setting"
sidebar_current: "docs-pagerduty-datasource-alert-grouping-setting"
description: |-
  Get information about an alert grouping setting that you have created.
---

# pagerduty\_alert\_grouping\_setting

Use this data source to get information about a specific [alert grouping setting][1].

## Example Usage

```hcl
data "pagerduty_alert_grouping_setting" "example" {
  name = "My example setting"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use to find an alert grouping setting in the PagerDuty API.

## Attributes Reference

* `id` - The ID of the found alert grouping setting.
* `name` - The short name of the found alert grouping setting.
* `description` - A description of this alert grouping setting.
* `type` - The type of object. The value returned will be one of `content_based`, `content_based_intelligent`, `intelligent` or `time`.
* `config` - The values for the configuration setup for this setting.
* `services` - A list of string containing the IDs of services associated with this setting.

The `config` block contains the following arguments:

* `timeout` - The duration in minutes within which to automatically group incoming alerts. This setting is only required and applies when `type` is set to `time`. To continue grouping alerts until the incident is resolved leave this value unset or set it to `null`.
* `aggregate` - One of `any` or `all`. This setting is only required and applies when `type` is set to `content_based` or `content_based_intelligent`. Group alerts based on one or all of `fields` value(s).
* `fields` - Alerts will be grouped together if the content of these fields match. This setting is only required and applies when `type` is set to `content_based` or `content_based_intelligent`.
* `time_window` - The maximum amount of time allowed between Alerts. This setting applies only when `type` is set to `intelligent`, `content_based`, `content_based_intelligent`. Value must be between `300` and `3600` or exactly `86400` (86400 is supported only for `content_based` alert grouping). Any Alerts arriving greater than `time_window` seconds apart will not be grouped together. This is a rolling time window and is counted from the most recently grouped alert. The window is extended every time a new alert is added to the group, up to 24 hours. To use the recommended time window leave this value unset or set it to `null`.

[1]: https://developer.pagerduty.com/api-reference/9b5a6c8d7379b-get-an-alert-grouping-setting
