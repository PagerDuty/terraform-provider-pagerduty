---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_user_contact_method"
sidebar_current: "docs-pagerduty-datasource-user-contact-method"
description: |-
  Get information about a contact method of a PagerDuty user (email, phone, SMS or push notification).
---

# pagerduty\_user\_contact\_method

Use this data source to get information about a specific [contact method][1] of a PagerDuty [user][2] that you can use for other PagerDuty resources.

## Example Usage

```hcl
data "pagerduty_user" "me" {
  email = "me@example.com"
}

data "pagerduty_user_contact_method" "phone_push" {
  user_id = data.pagerduty_user.me.id
  type    = "push_notification_contact_method"
  label   = "iPhone (John)"
}

resource "pagerduty_user_notification_rule" "low_urgency_sms" {
  user_id                = data.pagerduty_user.me.id
  start_delay_in_minutes = 5
  urgency                = "high"

  contact_method = {
    type = "push_notification_contact_method"
    id   = data.pagerduty_user_contact_method.phone_push.id
  }
}
```

## Argument Reference

The following arguments are supported:

  * `user_id` - (Required) The ID of the user.
  * `type` - (Required) The contact method type. May be (`email_contact_method`, `phone_contact_method`, `sms_contact_method`, `push_notification_contact_method`).
  * `label` - (Required) The label (e.g., "Work", "Mobile", "Ashley's iPhone", etc.).

## Attributes Reference
  * `id` - The ID of the found user.
  * `type` - The type of the found contact method. May be (`email_contact_method`, `phone_contact_method`, `sms_contact_method`, `push_notification_contact_method`).
  * `send_short_email` - Send an abbreviated email message instead of the standard email output. (Email contact method only.)
  * `country_code` - The 1-to-3 digit country calling code. (Phone and SMS contact methods only.)
  * `label` - The label (e.g., "Work", "Mobile", "Ashley's iPhone", etc.).
  * `address` - The "address" to deliver to: `email`, `phone number`, etc., depending on the type.
  * `blacklisted` - If true, this phone has been blacklisted by PagerDuty and no messages will be sent to it. (Phone and SMS contact methods only.)
  * `enabled` - If true, this phone is capable of receiving SMS messages. (Phone and SMS contact methods only.)
  * `device_type` - Either `ios` or `android`, depending on the type of the device receiving notifications. (Push notification contact method only.)

[1]: https://developer.pagerduty.com/api-reference/reference/REST/openapiv3.json/paths/~1users~1%7Bid%7D~1contact_methods~1%7Bcontact_method_id%7D/get
[2]: https://developer.pagerduty.com/api-reference/reference/REST/openapiv3.json/paths/~1users~1%7Bid%7D/get
