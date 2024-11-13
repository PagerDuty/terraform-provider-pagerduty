---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_automation_actions_action"
sidebar_current: "docs-pagerduty-resource-automation-actions-action"
description: |-
  Creates and manages an Automation Actions action in PagerDuty.
---

# pagerduty\_automation\_actions\_action

An Automation Actions [action](https://developer.pagerduty.com/api-reference/d64584a4371d3-create-an-automation-action) invokes jobs and workflows that are staged in Runbook Automation or Process Automation. It may also execute a command line script run by a Process Automation runner installed in your infrastructure.

## Example Usage

```hcl
resource "pagerduty_automation_actions_action" "pa_action_example" {
  name = "PA Action created via TF"
  description = "Description of the PA Action created via TF"
  action_type = "process_automation"
  action_data_reference {
    process_automation_job_id = "P123456"
  }
}

resource "pagerduty_automation_actions_action" "script_action_example" {
  name = "Script Action created via TF"
  description = "Description of the Script Action created via TF"
  action_type = "script"
  action_data_reference {
    script = "print(\"Hello from a Python script!\")"
    invocation_command = "/usr/local/bin/python3"
  }
}

```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the action. Max length is 255 characters.
  * `description` - (Required) The description of the action. Max length is 1024 characters.
  * `action_type` - (Required) The type of the action. The only allowed values are `process_automation` and `script`. Cannot be changed once set.
  * `action_data_reference` - (Required) Action Data block. Action Data is documented below.
  * `runner_id` - (Optional) The Process Automation Actions runner to associate the action with. Cannot be changed for the `process_automation` action type once set.
  * `action_classification` - (Optional) The category of the action. The only allowed values are `diagnostic` and `remediation`.

Action Data (`action_data_reference`) supports the following:

  * `process_automation_job_id` - (Required for `process_automation` action_type) The ID of the Process Automation job to execute.
  * `process_automation_job_arguments` - (Optional) The arguments to pass to the Process Automation job execution.
  * `process_automation_node_filter` - (Optional) The expression that filters on which nodes a Process Automation Job executes [Learn more](https://docs.rundeck.com/docs/manual/05-nodes.html#node-filtering).
  * `script` - (Required for `script` action_type) Body of the script to be executed on the Runner. Max length is 16777215 characters.
  * `invocation_command` - (Optional) The command to execute the script with.
  * `only_invocable_on_unresolved_incidents` - (Optional) Whether or not the action can be invoked on unresolved incidents.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the action.
* `type` - The type of object. The value returned will be `action`.
* `creation_time` - The time action was created. Represented as an ISO 8601 timestamp.
* `runner_type` - (Optional) The type of the runner associated with the action.
* `modify_time` - (Optional) The last time action has been modified. Represented as an ISO 8601 timestamp.

## Import

Actions can be imported using the `id`, e.g.

```
$ terraform import pagerduty_automation_actions_action.example 01DER7CUUBF7TH4116K0M4WKPU
```
