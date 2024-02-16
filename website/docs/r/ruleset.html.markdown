---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_ruleset"
sidebar_current: "docs-pagerduty-resource-ruleset"
description: |-
  Creates and manages an ruleset in PagerDuty.
---

# pagerduty\_ruleset

[Rulesets](https://support.pagerduty.com/docs/rulesets) allow you to route events to an endpoint and create collections of event rules, which define sets of actions to take based on event content.

<div role="alert" class="alert alert-warning">
  <div class="alert-title"><i class="fa fa-warning"></i>End-of-Life</div>
  <p>
    Rulesets and Event Rules will end-of-life soon. We highly recommend that you
    <a
      href="https://support.pagerduty.com/docs/migrate-to-event-orchestration"
      rel="noopener noreferrer"
      target="_blank"
      >migrate to Event Orchestration</a>
    as soon as possible so you can take advantage of the new functionality, such
    as improved UI, rule creation, REST APIs and Terraform support, advanced
    conditions, and rule nesting.
  </p>
</div>

## Example Usage

```hcl
resource "pagerduty_team" "foo" {
	name = "Engineering (Seattle)"
}

resource "pagerduty_ruleset" "foo" {
	name = "Primary Ruleset"
	team {
		id = pagerduty_team.foo.id
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the ruleset.
* `team` - (Optional) Reference to the team that owns the ruleset. If none is specified, only admins have access.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the ruleset.
* `routing_keys` - Routing keys routed to this ruleset.
* `type` - Type of ruleset. Currently, only sets to `global`.

## Import

Rulesets can be imported using the `id`, e.g.

```
$ terraform import pagerduty_ruleset.main 19acac92-027a-4ea0-b06c-bbf516519601
```
