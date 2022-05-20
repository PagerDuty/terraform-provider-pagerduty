---
layout: "pagerduty"
page_title: "PagerDuty: event_orchestration_separate_resources_per_type"
sidebar_current: "docs-pagerduty-resource-event-orchestration-router"
description: |-
  Creates and manages an orchestration path (router|unrouted|service) associated with a Global Event Orchestration in PagerDuty.
---

# Separate resources for each "type" of Orchestration path & no rule resource

In this example, we have separate resources for each "type" of Orchestration:

* `pagerduty_event_orchestration_router`
* `pagerduty_event_orchestration_unrouted`
* `pagerduty_event_orchestration_service`

Rules are defined as attributes of the orchestration path resource


## Example of configuring a Global Event Orchestration

```hcl
resource "pagerduty_team" "engineering" {
  name = "Engineering"
}

resource "pagerduty_event_orchestration" "my_monitor" {
  name = "My Monitoring Orchestration"
  description = "Send events to a pair of services"
  team = pagerduty_team.foo.id
}
```

## Example of configuring Router rules for an Orchestration

```hcl
resource "pagerduty_event_orchestration_router" "my_monitor" {
  event_orchestration = pagerduty_event_orchestration.my_monitor
  # Note: the router currently only supports a single set of rules.
  # If we change routers to support multiple sets in the future we'd need to introde a `sets` argument.
  rules = [
    {
      name = "Events relating to our relational database"
      conditions = [
        { expression = "event.summary matches part 'database'" },
        { expression = "event.source matches regex 'db[0-9]+-server'" }
      ]
      route_to = pagerduty_service.database.id
    },
    {
      name = "Events relating to our www app server"
      conditions = [
        { expression = "event.summary matches part 'www'" }
      ]
      route_to = pagerduty_service.www.id
    }
  ]
  catch_all = {
    route_to = "unrouted"
  }
}
```

**Argument Reference:** supports same rule and catch-all arguments as the [Router Orchestration endpoint](https://developer.pagerduty.com/api-reference/b3A6MzU3MDU0Mzk-update-the-router-for-a-global-event-orchestration) of the Public API

## Example of configuring a Unrouted Rules for an Orchestration

```hcl
resource "pagerduty_event_orchestration_unrouted" "my_monitor" {
  event_orchestration = pagerduty_event_orchestration.my_monitor
  sets = [
    {
      id = "start"
      rules = [
        {
          name = "Update the summary of un-matched Critical alerts so they're easier to spot"
          conditions = [
            { expression = "event.severity matches 'critical'" }
          ]
          actions {
            extractions = [
              {
                target = "event.summary"
                template = "[Critical Unrouted] {{event.summary}}"
              }
            ]
          }
        },
        {
          name = "Reduce the severity of all other unrouted events"
          conditions = []
          actions = {
            severity = "info"
          }
        }
      ]
    }
  ]
  # In this example the user has defined their own "all other unrouted events" rule with no conditions
  # so they aren't defining global catch_all behavior for the pagerduty_event_orchestration_unrouted resource
  # but using the catch_all would be a totally legit alternative.
}
```

**Argument Reference:** supports same rule, set, and catch-all arguments as the [Unrouted Orchestration endpoint](https://developer.pagerduty.com/api-reference/b3A6MzU3MDU0NDE-update-the-unrouted-orchestration-for-a-global-event-orchestration) of the Public API

## Example of configuring a Service Orchestration

```hcl
resource "pagerduty_event_orchestration_service" "www" {
  service = route_to = pagerduty_service.www.id
  sets = [
    {
      id: "start"
      rules = [
        {
          name = "Always apply some consistent event transformations to all events"
          conditions = []
          actions {
            variables {
              name = "hostname"
              path = "event.component"
              value = "hostname: (.*)"
              type = "regex"
            }
            extractions = [
              {
                # Demonstrating a template-style extraction
                template = "{{variables.hostname}}"
                target = "event.custom_details.hostname"
              },
              {
                # Demonstrating a regex-style extractions
                source = "event.source"
                regex = "www (.*) service"
                target = "event.source"
              }
            ]
            route_to = "step-two"                 
          }
        }
      ]
    },
    {
      id: "step-two"
      rules = [
        {
          name = "All critical alerts should be treated as P1 incidents"
          conditions = [
            { expression = "event.severity matches 'critical'" }
          ]
          actions {
            annotate = "Please use our P1 runbook: https://docs.test/p1-runbook"
            priority = "P0IN2KQ"
            suppress = false
          }
        },
        {
          name = "If there's something wrong on the canary let the team know about it in our deployments Slack channel"
          conditions = [
            { expression = "event.custom_details.hostname matches part 'canary'" }
          ]
          actions {
            automation_actions = [
              {
                name = "Canary Slack Notification"
                url = "https://our-slack-listerner.test/canary-notification"
                auto_send = true
                parameters = []
              }
            ]
          }
        },
        {
          name = "Never bother the on-call for info-level events outside of work hours"
          conditions = [
            { expression = "event.severity matches 'info' and not (now in Mon,Tue,Wed,Thu,Fri 09:00:00 to 17:00:00 America/Los_Angeles)" }
          ]
          actions {
            suppress = true
          }
        }
      ]
    }
  ]
  catch_all = {
    suppress = true
  }
}
```
**Argument Reference:** supports same rule, set, and catch-all arguments as the [Service Orchestration endpoint](https://developer.pagerduty.com/api-reference/b3A6MzU3MDU0NDM-update-the-service-orchestration-for-a-service) of the Public API