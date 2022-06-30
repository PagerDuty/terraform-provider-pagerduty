---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_response_play"
sidebar_current: "docs-pagerduty-datasource-response-play"
description: |-
  Get information about a response play in PagerDuty.
---

# pagerduty_response_play

Use this data source to infomation about a specific response play that you can use for other PagerDuty resources.

## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
}

data "pagerduty_response_play" "example" {
  name = "My Response Play"
  from = pagerduty_user.example.email
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the response play.
  * `from` - (Required) The email of the user attributed to the request. Needs to be a valid email address of a user in the PagerDuty account.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the response play.
