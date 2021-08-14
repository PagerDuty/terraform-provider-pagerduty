---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_user_notification_rule"
sidebar_current: "docs-pagerduty-resource-user-notification-rule"
description: |-
  Creates and manages notification rules for a user in PagerDuty.
---

# pagerduty_user_notification_rule

A [notification rule](https://developer.pagerduty.com/api-reference/reference/REST/openapiv3.json/paths/~1users~1%7Bid%7D~1notification_rules~1%7Bnotification_rule_id%7D/get) configures where and when a PagerDuty user is notified when a triggered incident is assigned to them. Unique notification rules can be created for both high and low-urgency incidents.

## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
}

resource "pagerduty_user_contact_method" "email" {
  user_id = pagerduty_user.example.id
  type    = "email_contact_method"
  address = "foo@bar.com"
  label   = "Work"
}

resource "pagerduty_user_contact_method" "phone" {
  user_id      = pagerduty_user.example.id
  type         = "phone_contact_method"
  country_code = "+1"
  address      = "2025550199"
  label        = "Work"
}

resource "pagerduty_user_contact_method" "sms" {
  user_id      = pagerduty_user.example.id
  type         = "sms_contact_method"
  country_code = "+1"
  address      = "2025550199"
  label        = "Work"
}

resource "pagerduty_user_notification_rule" "high_urgency_phone" {
  user_id                = pagerduty_user.example.id
  start_delay_in_minutes = 1
  urgency                = "high"

  contact_method = {
    type = "phone_contact_method"
    id   = pagerduty_user_contact_method.phone.id
  }
}

resource "pagerduty_user_notification_rule" "low_urgency_email" {
  user_id                = pagerduty_user.example.id
  start_delay_in_minutes = 1
  urgency                = "low"

  contact_method = {
    type = "email_contact_method"
    id   = pagerduty_user_contact_method.email.id
  }
}

resource "pagerduty_user_notification_rule" "low_urgency_sms" {
  user_id                = pagerduty_user.example.id
  start_delay_in_minutes = 10
  urgency                = "low"

  contact_method = {
    type = "sms_contact_method"
    id   = pagerduty_user_contact_method.sms.id
  }
}
```

## Argument Reference

The following arguments are supported:

  * `user_id` - (Required) The ID of the user.
  * `start_delay_in_minutes` - (Required) The delay before firing the rule, in minutes.
  * `urgency` - (Required) Which incident urgency this rule is used for. Account must have the `urgencies` ability to have a low urgency notification rule. Can be `high` or `low`.
  * `contact_method` - (Required) A contact method block, configured as a block described below.

Contact methods (`contact_method`) supports the following:

  * `id` - (Required) The id of the referenced contact method.
  * `type` - (Required) The type of contact method. Can be `email_contact_method`, `phone_contact_method`, `push_notification_contact_method` or `sms_contact_method`.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the user notification rule.

## Import

User notification rules can be imported using the `user_id` and the `id`, e.g.

```
$ terraform import pagerduty_user_notification_rule.main PXPGF42:PPSCXAN
```
