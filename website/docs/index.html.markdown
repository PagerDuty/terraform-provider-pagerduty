---
layout: "pagerduty"
page_title: "Provider: PagerDuty"
sidebar_current: "docs-pagerduty-index"
description: |-
  PagerDuty is an alarm aggregation and dispatching service
---

# PagerDuty Provider

[PagerDuty](https://www.pagerduty.com/) is an incident management platform that provides reliable notifications, automatic escalations, on-call scheduling, and other functionality to help teams detect and address unplanned work in real-time.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the PagerDuty provider
terraform {
  required_providers {
    pagerduty = {
      source  = "pagerduty/pagerduty"
      version = "2.2.1"
    }
  }
}

provider "pagerduty" {
  token = var.pagerduty_token
}

# Create a PagerDuty team
resource "pagerduty_team" "engineering" {
  name        = "Engineering"
  description = "All engineering"
}

# Create a PagerDuty user
resource "pagerduty_user" "earline" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
}

# Create a team membership
resource "pagerduty_team_membership" "earline_engineering" {
  user_id = pagerduty_user.earline.id
  team_id = pagerduty_team.engineering.id
}
```

## Argument Reference

The following arguments are supported:

* `token` - (Required) The v2 authorization token. It can also be sourced from the PAGERDUTY_TOKEN environment variable. See [API Documentation](https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTUx-authentication)for more information.
* `user_token` - (Optional) The v2 user level authorization token. It can also be sourced from the PAGERDUTY_USER_TOKEN environment variable. See [API Documentation](https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTUx-authentication) for more information.
* `skip_credentials_validation` - (Optional) Skip validation of the token against the PagerDuty API.
* `service_region` - (Optional) The PagerDuty service region to use. Default to empty (uses US region). Supported value: `eu`.
* `api_url_override` - (Optional) It can be used to set a custom proxy endpoint as PagerDuty client api url overriding `service_region` setup.

## Additional Resources

- [PagerDuty Knowledge Base](https://support.pagerduty.com/)
- [List of Incident Workflow Actions](https://support.pagerduty.com/docs/incident-workflows#add-actions)
