---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_standards_resources_scores"
sidebar_current: "docs-pagerduty-datasource-standards-resources-scores"
description: |-
  Get information about the scores received for the standards associated to
  multiple resources.
---

# pagerduty\_standards\_resource\_scores

Use this data source to get information about the [scores for the standards for
many resources][1].

## Example Usage

```hcl
data "pagerduty_service" "foo" {
  name = "foo"
}

data "pagerduty_service" "bar" {
  name = "bar"
}

data "pagerduty_service" "baz" {
  name = "baz"
}

data "pagerduty_standards_resources_scores" "scores" {
  resource_type = "technical_services"

  ids = [
    data.pagerduty_service.foo.id,
    data.pagerduty_service.bar.id,
    data.pagerduty_service.baz.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - Type of the object the standards are associated to. Allowed values are `technical_services`.
* `ids` - List of identifiers of the resources to query.

## Attributes Reference

* `resources` - List of score results for each queried resource.

### Resource (`resources`) is a map with following attributes:

* `resource_type` - Type of the object the standards are associated to.
* `resource_id` - Unique Identifier.
* `score` - Summary of the scores for standards associated with this resource.
  * `passing` - Number of standards this resource successfully complies to.
  * `total` - Number of standards associated to this resource.
* `standards` - The list of standards evaluated against.
  * `id` - A unique identifier for the standard.
  * `name` - The human-readable name of the standard.
  * `active` - Indicates whether the standard is currently active and applicable to the resource.
  * `description` - Provides a textual description of the standard.
  * `type` - The type of the standard.
  * `resource_type` - Specifies the type of resource to which the standard applies.
  * `pass` - Indicates whether the resource complies to this standard.

[1]: https://developer.pagerduty.com/api-reference/2e832500ae129-list-resources-standards-scores
