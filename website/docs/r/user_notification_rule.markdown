---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_user_notification_rule"
sidebar_current: "docs-pagerduty-resource-user-notification-rule"
description: |-
  Creates and manages notification rules for a user in PagerDuty.
---

#pagerduty_user_notification_rule

A [notification rule](https://v2.developer.pagerduty.com/v2/page/api-reference#!/Users/get_users_id_notification_rules_notification_rule_id) is a notification rule for a PagerDuty user.


## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
  teams = ["${pagerduty_team.example.id}"]
}

resource "pagerduty_user_contact_method" "email" {
  user_id = "${pagerduty_user.example.id}"
  type    = "email_contact_method"
  address = "foo@bar.com"
  label   = "Work"
}

resource "pagerduty_user_notification_rule" "phone" {
  user_id                = "${pagerduty_user.example.id}"
  contact_method_id      = "${pagerduty_user_contact_method.email.id}"
  contact_method_type    = "${pagerduty_user_contact_method.email.type}"
  start_delay_in_minutes = 0
  urgency                = "high"
}
```

## Argument Reference

The following arguments are supported:

  * `user_id` - (Required) The ID of the user.
  * `contact_method_id` - (Required) The ID of the contact method.
  * `contact_method_type` - (Required) The type of the contact method.
  * `start_delay_in_minutes` - (Required) The delay before firing the rule, in minutes.
  * `urgency` - (Required) Which incident urgency this rule is used for.
  Account must have the `urgencies` ability to have a low urgency notification rule. Can be `high` or `low`.

## Attributes Reference

The following attributes are exported:
  * `id` the ID of the notification rule.
  * `blacklisted` - If true, this phone has been blacklisted by PagerDuty and no messages will be sent to it.
  * `enabled` - If true, this phone is capable of receiving SMS messages.

## Import

Notification rules can be imported using the `user_id` and the `id`, e.g.

```
$ terraform import pagerduty_user_notification_rule.main PLBP09X:PLBP09X
```
