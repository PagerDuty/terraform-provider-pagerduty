---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_jira_cloud_account_mapping_rule"
sidebar_current: "docs-pagerduty-resource-jira-cloud-account-mapping-rule"
description: |-
  Creates and manages a Jira Cloud account mapping Rule to integrate with PagerDuty.
---

# pagerduty\_jira\_cloud\_account\_mapping\_rule

An Jira Cloud's account mapping [rule](https://developer.pagerduty.com/api-reference/85dc30ba966a6-create-a-rule)
configures the bidirectional synchronization between Jira issues and PagerDuty
incidents.

## Example Usage

```hcl
data "pagerduty_escalation_policy" "default" {
  name = "Default"
}

data "pagerduty_priority" "p1" {
  name = "P1"
}

data "pagerduty_priority" "p2" {
  name = "P2"
}

data "pagerduty_priority" "p3" {
  name = "P3"
}

resource "pagerduty_service" "foo" {
  name              = "My Web App"
  escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_user" "foo" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
}

resource "pagerduty_jira_cloud_account_mapping_rule" "foo" {
  name = "Integration with My Web App"
  account_mapping = "PLBP09X"
  config {
    service = pagerduty_service.foo.id
    jira {
      autocreate_jql = "priority = Highest"
      create_issue_on_incident_trigger = true
      custom_fields {
        source_incident_field = "incident_description"
        target_issue_field = "description"
        target_issue_field_name = "Description"
        type = "attribute"
      }
      custom_fields {
        target_issue_field = "security"
        target_issue_field_name = "Security Level"
        type = "jira_value"
        value = jsonencode({
          displayName = "Sec Level 1"
          id = "10000"
        })
      }
      issue_type {
        id = "10001"
        name = "Incident"
      }
      priorities {
        jira_id = "1"
        pagerduty_id = data.pagerduty_priority.p1.id
      }
      priorities {
        jira_id = "2"
        pagerduty_id = data.pagerduty_priority.p2.id
      }
      priorities {
        jira_id = "3"
        pagerduty_id = data.pagerduty_priority.p3.id
      }
      project {
        id = "10100"
        key = "ITS"
        name = "IT Support"
      }
      status_mapping {
        acknowledged {
          id = "2"
          name = "In Progress"
        }
        resolved {
          id = "3"
          name = "Resolved"
        }
        triggered {
          id = "1"
          name = "Open"
        }
      }
      sync_notes_user = pagerduty_user.foo.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the rule.
* `account_mapping` - (Required) [Updating can cause a resource replacement] The account mapping this rule belongs to. 
* `config` - (Required) Configuration for bidirectional synchronization between Jira issues and PagerDuty incidents.

The `config` block contains the following arguments:

* `service` - (Required) [Updating can cause a resource replacement] The ID of the linked PagerDuty service.
* `jira` - (Required) Synchronization settings.

The `jira` block contains the following arguments:

* `autocreate_jql` - JQL query to automatically create PagerDuty incidents when matching Jira issues are created. Leave empty to disable this feature.
* `create_issue_on_incident_trigger` - When enabled, automatically creates a Jira issue whenever a PagerDuty incident is triggered.
* `custom_fields` - Defines how Jira fields are populated when a Jira Issue is created from a PagerDuty Incident.
* `issue_type` - (Required) Specifies the Jira issue type to be created or synchronized with PagerDuty incidents.
* `priorities` - Maps PagerDuty incident priorities to Jira issue priorities for synchronization.
* `project` - (Required) [Updating can cause a resource replacement] Defines the Jira project where issues will be created or synchronized.
* `status_mapping` - (Required) Maps PagerDuty incident statuses to corresponding Jira issue statuses for synchronization.
* `sync_notes_user` - ID of the PagerDuty user for syncing notes and comments between Jira issues and PagerDuty incidents. If not provided, note synchronization is disabled.

A `custom_fields` block contains the following arguments:

* `type` - (Required) The type of the value that will be set; one of `attribute`, `const` or `jira_value`.
* `source_incident_field` - The PagerDuty incident field from which the value will be extracted (only applicable if `type` is `attribute`); one of `incident_number`, `incident_title`, `incident_description`, `incident_status`, `incident_created_at`, `incident_service`, `incident_escalation_policy`, `incident_impacted_services`, `incident_html_url`, `incident_assignees`, `incident_acknowledgers`, `incident_last_status_change_at`, `incident_last_status_change_by`, `incident_urgency` or `incident_priority`.
* `target_issue_field` - (Required) The unique identifier key of the Jira field that will be set.
* `target_issue_field_name` - (Required) The human-readable name of the Jira field.
* `value` - The value to be set for the Jira field (only applicable if `type` is `const` or `jira_value`). It must be set as a JSON string.

The `issue_type` block contains the following arguments:

* `id` - (Required) Unique identifier for the Jira issue type.
* `name` - (Required) The name of the Jira issue type.

A `priorities` block contains the following arguments:

* `jira_id` - (Required) The ID of the Jira priority.
* `pagerduty_id` - (Required) The ID of the PagerDuty priority.

The `project` block contains the following arguments:

* `id` - (Required) Unique identifier for the Jira project.
* `key` - (Required) The short key name of the Jira project.
* `name` - (Required) The name of the Jira project.

The `status_mapping` block contains the following arguments:

* `acknowledged` - Jira status that maps to the PagerDuty `acknowledged` status.
* `resolved` - Jira status that maps to the PagerDuty `resolved` status.
* `triggered` - (Required) Jira status that maps to the PagerDuty `triggered` status.

The `acknowledged`, `resolved` and `triggered` blocks contains the following arguments:

* `id` - Unique identifier for the Jira status.
* `name` - Name of the Jira status.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service.
* `autocreate_jql_disabled_reason` - If auto-creation using JQL is disabled, this field provides the reason for the disablement.
* `autocreate_jql_disabled_until` - The timestamp until which the auto-creation using JQL feature is disabled.

## Import

Jira Cloud account mapping rules can be imported using the `account_mapping_id` and `rule_id`, e.g.

```
$ terraform import pagerduty_jira_cloud_account_mapping_rule.main PLBP09X:PLB09Z
```
