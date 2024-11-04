---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_jira_cloud_account_mapping"
sidebar_current: "docs-pagerduty-datasource-jira-cloud-account-mapping"
description: |-
  An Account Mapping establishes a connection between a PagerDuty account and a Jira Cloud instance, enabling integration and synchronization between the two platforms.
---

# pagerduty\_jira\_cloud\_account\_mapping

Use this data source to get information about a specific [account mapping][1].

## Example Usage

```hcl
data "pagerduty_jira_cloud_account_mapping" "circular" {
  name = "pdt-circular"
}
```

## Argument Reference

The following arguments are supported:

* `subdomain` - (Required) The service name to use to find a service in the PagerDuty API.

## Attributes Reference

* `id` - The ID of the found account mapping.
* `base_url` - The base URL of the Jira Cloud instance, used for API calls and constructing links.

[1]: https://developer.pagerduty.com/api-reference/8d707b61562b7-get-an-account-mapping
