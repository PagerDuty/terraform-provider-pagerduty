---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_standards_resource_scores"
sidebar_current: "docs-pagerduty-datasource-standards-resource-scores"
description: |-
  Get information about the scores received for the standards associated to a
  single resource.
---

# pagerduty\_standards\_resource\_scores

Use this data source to get information about the [scores for the standards of a
resource][1].

## Example Usage

```hcl
data "pagerduty_service" "example" {
  name = "My Service"
}

data "pagerduty_standards_resource_scores" "scores" {
    resource_type = "technical_services"
    id = data.pagerduty_service.example.id
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - Type of the object the standards are associated to. Allowed values are `technical_services`.
* `id` - Identifier of said resource.

## Attributes Reference

* `score` - Summary of the scores for standards associated with this resource.
* `standards` - The list of standards evaluated against.

### Score (`score`) is a map with following attributes:

* `passing` - Number of standards this resource successfully complies to.
* `total` - Number of standards associated to this resource.

### Standards (`standards`) is a list of objects that support the following:

* `id` - A unique identifier for the standard.
* `name` - The human-readable name of the standard.
* `active` - Indicates whether the standard is currently active and applicable to the resource.
* `description` - Provides a textual description of the standard.
* `type` - The type of the standard.
* `resource_type` - Specifies the type of resource to which the standard applies.
* `pass` - Indicates whether the resource complies to this standard.

[1]: https://developer.pagerduty.com/api-reference/f339354b607d5-list-a-resource-s-standards-scores
