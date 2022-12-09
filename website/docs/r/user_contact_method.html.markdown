---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_user_contact_method"
sidebar_current: "docs-pagerduty-resource-user-contact-method"
description: |-
  Creates and manages contact methods for a user in PagerDuty.
---

# pagerduty_user_contact_method

-> This resource behaves a little differently than may be expected. If the defined contact method already exists for the user in PagerDuty this resource will import the values of the existing contact method into your Terraform state.

A [contact method](https://developer.pagerduty.com/api-reference/b3A6Mjc0ODI0MA-create-a-user-contact-method) is a contact method for a PagerDuty user (email, phone or SMS).


## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
  teams = [pagerduty_team.example.id]
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
```

## Argument Reference

The following arguments are supported:

  * `user_id` - (Required) The ID of the user.
  * `type` - (Required) The contact method type. May be (`email_contact_method`, `phone_contact_method`, `sms_contact_method`, `push_notification_contact_method`).
  * `send_short_email` - (Optional) Send an abbreviated email message instead of the standard email output.
  * `country_code` - (Optional) The 1-to-3 digit country calling code. Required when using `phone_contact_method` or `sms_contact_method`.
  * `label` - (Required) The label (e.g., "Work", "Mobile", etc.).
  * `address` - (Required) The "address" to deliver to: `email`, `phone number`, etc., depending on the type.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the contact method.
  * `blacklisted` - If true, this phone has been blacklisted by PagerDuty and no messages will be sent to it.
  * `enabled` - If true, this phone is capable of receiving SMS messages.

## Import

Contact methods can be imported using the `user_id` and the `id`, e.g.

```
$ terraform import pagerduty_user_contact_method.main PLBP09X:PLBP09X
```
