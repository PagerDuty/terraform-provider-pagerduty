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
# Assumes the TF_VAR_RUNBOOK_API_KEY variable is defined in the environment

variable "RUNBOOK_API_KEY" {
  type = string
  sensitive = true
}

resource "pagerduty_automation_actions_runner" "example" {
  name = "Runner created via TF"
  description = "Description of the Runner created via TF"
  runner_type = "runbook"
  runbook_base_uri = "rdcat.stg"
  runbook_api_key = var.RUNBOOK_API_KEY
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the runner. Max length is 255 characters.
  * `description` - (Required) The description of the runner. Max length is 1024 characters.
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

-> In the example below the `runbook_api_key` attribute has been omitted to avoid resource replacement after the import.

Runners can be imported using the `id`, e.g.

```
resource "pagerduty_automation_actions_runner" "example" {
  name = "Runner created via TF"
  description = "Description of the Runner created via TF"
  runner_type = "runbook"
  runbook_base_uri = "rdcat.stg"
}
```
```
$ terraform import pagerduty_automation_actions_runner.example 01DER7CUUBF7TH4116K0M4WKPU
```


