---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_tag_assignment"
sidebar_current: "docs-pagerduty-resource-tag-assignment"
description: |-
  Creates and manages a tag assignment in PagerDuty.
---

# pagerduty\_tag\_assignment

A [tag](https://developer.pagerduty.com/api-reference/b3A6Mjc0ODEwMA-assign-tags) is applied to Escalation Policies, Teams or Users and can be used to filter them.

## Example Usage

```hcl
resource "pagerduty_tag" "example" {
  label = "API"
}
resource "pagerduty_team" "engteam" {
  name = "Engineering"
}
resource "pagerduty_tag_assignment" "example" {
  tag_id      = pagerduty_tag.example.id
  entity_type = "teams"
  entity_id   = pagerduty_team.engteam.id
}
```

## Argument Reference

The following arguments are supported:

  * `tag_id` - (Required) The ID of the tag.
  * `entity_type` - (Required) Type of entity in the tag assignment. Possible values can be `users`, `teams`, and `escalation_policies`.
  * `entity_id` - (Required) The ID of the entity.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the tag assignment.

## Import

Tag assignments can be imported using the `id` which is constructed by taking the `entity` Type, `entity` ID and the `tag` ID separated by a dot, e.g.

```
$ terraform import pagerduty_tag_assignment.main users.P7HHMVK.PYC7IQQ
```
