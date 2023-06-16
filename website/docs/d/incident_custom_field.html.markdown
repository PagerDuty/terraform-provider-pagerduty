---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_incident_custom_field"
sidebar_current: "docs-pagerduty-datasource-incident-custom-field"
description: |-
  Get information about an Incident Custom Field in PagerDuty.
---

# pagerduty\_incident\_custom\_field

Use this data source to get information about a specific [Incident Custom Field](https://support.pagerduty.com/docs/custom-fields-on-incidents).

## Example Usage

```hcl
data "pagerduty_incident_custom_field" "environment" {
  name      = "environment"
}

resource "pagerduty_incident_custom_field_option" "dev_environment" {
  field    = data.pagerduty_incident_custom_field.environment.id
  datatype = "string"
  value    = "dev"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the field.

## Attributes Reference

* `id` - The ID of the found field.
