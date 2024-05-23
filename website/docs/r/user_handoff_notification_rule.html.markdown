---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_user_handoff_notification_rule"
sidebar_current: "docs-pagerduty-resource-user-handoff-notification-rule"
description: |-
  Creates and manages an user handoff notification rule in PagerDuty.
---

# pagerduty\_user_handoff_notification_rule

An [user handoff notification rule](https://developer.pagerduty.com/api-reference/f2ab7a3c1418a-create-a-user-handoff-notification-rule) is a rule that specifies how a user should be notified when they are handed off an incident.

## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@foo.test"
}

resource "pagerduty_user_contact_method" "phone" {
  user_id      = pagerduty_user.example.id
  type         = "phone_contact_method"
  country_code = "+1"
  address      = "2025550199"
  label        = "Work"
}

resource "pagerduty_user_handoff_notification_rule" "example-oncall-offcall" {
  user_id                   = pagerduty_user.example.id
  handoff_type              = "both"
  notify_advance_in_minutes = 180
  contact_method {
    id   = pagerduty_user_contact_method.phone.id
    type = pagerduty_user_contact_method.phone.type
  }
}
```

## Argument Reference

The following arguments are supported:

  * `user_id` - (Required) The ID of the user.
  * `handoff_type` - (Optional) The type of handoff to notify the user about. Possible values are `oncall`, `offcall`, `both`.
  * `notify_advance_in_minutes` - (Required) The number of minutes before the handoff that the user should be notified. Must be a positive integer greater than or equal to 0.
  * `contact_method` - (Required) The contact method to notify the user. Contact method documented below.

Contact method supports the following:

  * `id` - (Required) The ID of the contact method.
  * `type` - (Required) The type of the contact method. May be (`email_contact_method`, `email_contact_method_reference`, `phone_contact_method`, `phone_contact_method_reference`, `push_notification_contact_method`, `push_notification_contact_method_reference`, `sms_contact_method`, `sms_contact_method_reference`).

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the user handoff notification rule.

## Import

User handoff notification rules can be imported using the `user_id` and `id` separated by a dot, e.g.

```
$ terraform import pagerduty_user_handoff_notification_rule.main PX4IAP4.PULREBP
```
