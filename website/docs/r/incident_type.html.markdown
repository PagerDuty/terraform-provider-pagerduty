---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_incident_type"
sidebar_current: "docs-pagerduty-resource-incident_type"
description: |-
  Creates and manages a incident_type in PagerDuty.
---

# pagerduty\_incident\_type

An [incident\_type](https://developer.pagerduty.com/api-reference/1981087c1914c-create-an-incident-type)
is a feature which allows customers to categorize incidents, such as a security
incident, a major incident, or a fraud incident.

<div role="alert" class="alert alert-warning">
  <div class="alert-title"><i class="fa fa-warning"></i>Resource limitation</div>
  <p>Incident Types cannot be deleted, only disabled</p>
  <p>If you want terraform to stop tracking this resource please use <code>terraform state rm</code>.</p>
</div>


## Example Usage

```hcl
data "pagerduty_incident_type" "base" {
    display_name = "Base Incident"
}

resource "pagerduty_incident_type" "example" {
  name = "backoffice"
  display_name = "Backoffice Incident"
  parent_type = data.pagerduty_incident_type.base.id
  description = "Internal incidents not facing customer"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Incident Type. Usage of the suffix `_default` is prohibited. This cannot be changed once the incident type has been created.
* `display_name` - (Required) The display name of the Incident Type. Usage of the prefixes PD, PagerDuty, or the suffixes Default, or (Default) is prohibited.
* `parent_type` - (Required) The parent type of the Incident Type. Either name or id of the parent type can be used.
* `description` - A succinct description of the Incident Type.
* `enabled`  - State of this Incident Type object. Defaults to true if not provided.

## Attributes Reference

* `id`  - The unique identifier of the incident type.
* `type`  - A string that determines the schema of the object.

## Import

Services can be imported using the `id`, e.g.

```
$ terraform import pagerduty_incident_type.main P12345
```
