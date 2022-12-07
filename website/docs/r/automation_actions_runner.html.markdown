---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_automation_actions_runner"
sidebar_current: "docs-pagerduty-resource-automation-actions-runner"
description: |-
  Creates and manages an Automation Actions runner in PagerDuty.
---

# pagerduty\_automation\_actions\_runner

An Automation Actions [runner](https://developer.pagerduty.com/api-reference/d78999fb7e863-create-an-automation-action-runner) is the method for how actions are executed. This can be done locally using an installed runner agent or as a connection to a PD Runbook Automation instance.

-> Only Runbook Automation (runbook) runners can be created.

## Example Usage

```hcl
resource "pagerduty_automation_actions_runner" "example" {
	name = "Production runner"
	description = "Runner created by the SRE team"
	runner_type = "runbook"
	runbook_base_uri = "prod.cat"
	runbook_api_key = "ABC123456789XYZ"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the runner.
  * `description` - (Required) The description of the runner.
  * `runner_type` - (Required) The type of runner. The only allowed values is `runbook`. 
  * `runbook_base_uri` - (Required) The subdomain for your Runbook Automation Instance. 
  * `runbook_api_key` - (Required) The unique User API Token created in Runbook Automation. 
  
## Attributes Reference

The following attributes are exported:

* `id` - The ID of the runner.
* `type` - The type of object. The value returned will be `runner`.
* `creation_time` - The time runner was created. Represented as an ISO 8601 timestamp.
* `last_seen` - (Optional) The last time runner has been seen. Represented as an ISO 8601 timestamp.

## Import

Runners can be imported using the `id`, e.g.

```
$ terraform import pagerduty_automation_actions_runner.main 01DBJLIGED17S1DQK123
```
