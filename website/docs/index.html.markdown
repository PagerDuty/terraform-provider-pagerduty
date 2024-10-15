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
      version = ">= 2.2.1"
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

* `token` - (Optional) The v2 authorization token. It can also be sourced from the `PAGERDUTY_TOKEN` environment variable. See [API Documentation](https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTUx-authentication)for more information.
* `user_token` - (Optional) The v2 user level authorization token. It can also be sourced from the `PAGERDUTY_USER_TOKEN` environment variable. See [API Documentation](https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTUx-authentication) for more information.
* `use_app_oauth_scoped_token` - (Optional) Defines the configuration needed for making use of [App Oauth Scoped API token](https://developer.pagerduty.com/docs/e518101fde5f3-obtaining-an-app-o-auth-token) for authenticating API calls.
* `skip_credentials_validation` - (Optional) Skip validation of the token against the PagerDuty API.
* `service_region` - (Optional) The PagerDuty service region to use. Default to empty (uses US region). Supported value: `eu`. This setting also affects configuration of `use_app_oauth_scoped_token` for setting Region of *App Oauth token credentials*. It can also be sourced from the `PAGERDUTY_SERVICE_REGION` environment variable.
* `api_url_override` - (Optional) It can be used to set a custom proxy endpoint as PagerDuty client api url overriding `service_region` setup.
* `insecure_tls` - (Optional) Can be used to disable TLS certificate checking when calling the PagerDuty API. This can be useful if you're behind a corporate proxy.

The `use_app_oauth_scoped_token` block contains the following arguments:

* `pd_client_id` - (Required) An identifier issued when the Scoped OAuth client was added to a PagerDuty App. It can also be sourced from the `PAGERDUTY_CLIENT_ID` environment variable.
* `pd_client_secret` - (Required) A secret issued when the Scoped OAuth client was added to a PagerDuty App. It can also be sourced from the `PAGERDUTY_CLIENT_SECRET` environment variable.
* `pd_subdomain` - (Required) Your PagerDuty account subdomain; i.e: If the *URL* shown by the Browser when you are in your PagerDuty account is some like: https://acme.pagerduty.com, then your PagerDuty subdomain is `acme`. It can also be sourced from the `PAGERDUTY_SUBDOMAIN` environment variable.

## Example using App Oauth scoped token

```hcl
# Configure the PagerDuty provider
terraform {
  required_providers {
    pagerduty = {
      source  = "pagerduty/pagerduty"
      version = ">= 3.0.0" # Mind the supported Provider version
    }
  }
}

provider "pagerduty" {
  # Configure use of App Oauth scoped token
  use_app_oauth_scoped_token {
    pd_client_id     = var.pd_client_id
    pd_client_secret = var.pd_client_secret
    pd_subdomain     = var.pd_subdomain
  }
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

## Debugging Provider Output Using Logs

In addition to the [log levels provided by Terraform](https://developer.hashicorp.com/terraform/internals/debugging), namely `TRACE`, `DEBUG`, `INFO`, `WARN`, and `ERROR` (in descending order of verbosity), the PagerDuty Provider introduces an extra level called `SECURE`. This level offers verbosity similar to Terraform's debug logging level, specifically for the output of API calls and HTTP request/response logs. The key difference is that API keys within the request's Authorization header will be obfuscated, revealing only the last four characters. An example is provided below:

```sh
---[ REQUEST ]---------------------------------------
GET /teams/DER8RFS HTTP/1.1
Accept: application/vnd.pagerduty+json;version=2
Authorization: <OBSCURED>kCjQ
Content-Type: application/json
User-Agent: (darwin arm64) Terraform/1.5.1
```

To enable the `SECURE` log level, you must set two environment variables:

* `TF_LOG=INFO`
* `TF_LOG_PROVIDER_PAGERDUTY=SECURE`

