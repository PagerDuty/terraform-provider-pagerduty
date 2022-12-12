---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_automation_actions_runner"
sidebar_current: "docs-pagerduty-datasource-automation-actions-runner"
description: |-
  Get information about an Automation Actions runner that you have created.
---

# pagerduty\_automation\_actions\_runner

Use this data source to get information about a specific [automation actions runner][1].

## Example Usage

```hcl
data "pagerduty_automation_actions_runner" "example" {
  id = "01DBJLIGED17S1DQKQC2AV8XYZ" 
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The id of the automation actions runner in the PagerDuty API.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the found runner.
* `name` - The name of the found runner.
* `type` - The type of object. The value returned will be `runner`.
* `runner_type` - The type of runner. Allowed values are `sidecar` and `runbook`.
* `creation_time` - The time runner was created. Represented as an ISO 8601 timestamp.
* `description` - (Optional) The description of the runner.
* `last_seen` - (Optional) The last time runner has been seen. Represented as an ISO 8601 timestamp.
* `runbook_base_uri` - (Optional) The base URI of the Runbook server to connect to. Applicable to `runbook` type runners only.

[1]: https://developer.pagerduty.com/api-reference/aace61f84cbd0-get-an-automation-action-runner
