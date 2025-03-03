---
layout: 'pagerduty'
page_title: 'PagerDuty: pagerduty_event_orchestration_global_cache_variable'
sidebar_current: 'docs-pagerduty-datasource-event-orchestration-global-cache-variable'
description: |-
  Get information about a Cache Variable for a Global Event Orchestration that you have created.
---

# pagerduty_event_orchestration_global_cache_variable

Use this data source to get information about a specific [Cache Variable][1] for a Global Event Orchestration.

## Example Usage

```hcl

resource "pagerduty_event_orchestration" "event_orchestration" {
  name = "Test Event Orchestration"
}

data "pagerduty_event_orchestration_global_cache_variable" "cache_variable" {
  event_orchestration = pagerduty_event_orchestration.event_orchestration.id
  name = "example_cache_variable"
}

```

## Argument Reference

The following arguments are supported:

* `event_orchestration` - (Required) ID of the Global Event Orchestration to which this Cache Variable belongs.
* `id` - (Optional) ID of the Cache Variable associated with the Global Event Orchestration. Specify either `id` or `name`. If both are specified `id` takes precedence.
* `name` - (Optional) Name of the Cache Variable associated with the Global Event Orchestration. Specify either `id` or `name`. If both are specified `id` takes precedence.

## Attributes Reference

* `disabled` - Indicates whether the Cache Variable is disabled and would therefore not be evaluated.
* `condition` - Conditions to be evaluated in order to determine whether or not to update the Cache Variable's stored value. This attribute can only be used when `configuration.0.type` is `recent_value` or `trigger_event_count`.
  * `expression`- A [PCL condition][2] string.
* `configuration` - A configuration object to define what and how values will be stored in the Cache Variable.
  * `type` - The [type of value][1] to store into the Cache Variable. Can be one of: `recent_value`, `trigger_event_count` or `external_data`.
  * `source` - The path to the event field where the `regex` will be applied to extract a value. You can use any valid [PCL path][3]. This field is only used when `type` is `recent_value`
  * `regex` - A [RE2 regular expression][4] that will be matched against the field specified via the `source` argument. This field is only used when `type` is `recent_value`
  * `ttl_seconds` - The number of seconds indicating how long to count incoming trigger events for. This field is only used when `type` is `trigger_event_count` or `external_data`
  * `data_type` - The type of data that will eventually be set for the Cache Variable via an API request. This field is only used when type is `external_data`


[1]: https://support.pagerduty.com/docs/event-orchestration-variables
[2]: https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview
[3]: https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview#paths
[4]: https://github.com/google/re2/wiki/Syntax
