---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_response_play"
sidebar_current: "docs-pagerduty-resource-response-play"
description: |-
  Creates and manages a response play in PagerDuty.
---

# pagerduty_response_play

A [response play](https://developer.pagerduty.com/api-reference/b3A6Mjc0ODE2Ng-create-a-response-play) allows you to create packages of Incident Actions that can be applied during an Incident's life cycle.

<div role="alert" class="alert alert-warning">
  <div class="alert-title"><i class="fa fa-warning"></i>End-of-Life</div>
  <p>
    Response Play will end-of-life soon. We highly recommend that you
    <a
      href="https://support.pagerduty.com/docs/upgrade-response-plays-to-incident-workflows"
      rel="noopener noreferrer"
      target="_blank"
      >migrate to Incident Workflows</a>
    as soon as possible so you can take advantage of the new functionality.
    With Incident Workflows, customers are able to define if-this-then-that
    logic to effortlessly trigger a sequence of common incident actions, advanced conditions, REST APIs
    and Terraform support.
  </p>
</div>

## Example Usage

```hcl
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
  teams = [pagerduty_team.example.id]
}

resource "pagerduty_escalation_policy" "example" {
  name      = "Engineering Escalation Policy"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user"
      id   = pagerduty_user.example.id
    }
  }
}

resource "pagerduty_response_play" "example" {
  name = "My Response Play"
  from = pagerduty_user.example.email

  responder {
    type = "escalation_policy_reference"
    id   = pagerduty_escalation_policy.example.id
  }

  subscriber {
    type = "user_reference"
    id   = pagerduty_user.example.id
  }

  runnability = "services"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the response play.
  * `from` - (Required) The email of the user attributed to the request. Needs to be a valid email address of a user in the PagerDuty account.
  * `description` - (Optional) A human-friendly description of the response play.
    If not set, a placeholder of "Managed by Terraform" will be set.
  * `type` - (Optional)  A string that determines the schema of the object. If not set, the default value is "response_play".
  * `team` - (Optional) The ID of the team associated with the response play.
  * `subscriber` - (Required) A user and/or team to be added as a subscriber to any incident on which this response play is run. There can be multiple subscribers defined on a single response play.
  * `subscribers_message` - (Optional) The content of the notification that will be sent to all incident subscribers upon the running of this response play. Note that this includes any users who may have already been subscribed to the incident prior to the running of this response play. If empty, no notifications will be sent.
  * `responder` - (Required) A user and/or escalation policy to be requested as a responder to any incident on which this response play is run. There can be multiple responders defined on a single response play.
  * `responders_message` - (Optional) The message body of the notification that will be sent to this response play's set of responders. If empty, a default response request notification will be sent.
  * `runnability` - (Optional) String representing how this response play is allowed to be run. Valid options are:

    * `services`: This response play cannot be manually run by any users. It will run automatically for new incidents triggered on any services that are configured with this response play.
    * `teams`: This response play can be run manually on an incident only by members of its configured team. This option can only be selected when the team property for this response play is not empty.
    * `responders`: This response play can be run manually on an incident by any responders in this account.

* `conference_number` - (Optional) The telephone number that will be set as the conference number for any incident on which this response play is run.
* `conference_url` - (Optional) The URL that will be set as the conference URL for any incident on which this response play is run.

### Responders (`responder`) can have two different objects and supports the following:

**User Responders**
* `id` - ID of the user defined as the responder
* `type` - Should be set as `user_reference` for user responders. `escalation_policy`

**Escalation Policy Responders**
* `id` - ID of the user defined as the responder
* `type` - Should be set as `escalation_policy` for escalation policy responders.
* `name` - Name of the escalation policy
* `description` - Description of escalation policy
* `num_loops` - The number of times the escalation policy will repeat after reaching the end of its escalation.
* `on_call_handoff_notifications` - Determines how on call handoff notifications will be sent for users on the escalation policy. Defaults to "if_has_services". Could be "if_has_services", "always
* `escalation_rule` - The escalation rules
  * `escalation_delay_in_minutes` - The number of minutes before an unacknowledged incident escalates away from this rule.
  * `target` - The targets an incident should be assigned to upon reaching this rule.
    * `type` - Type of object of the target. Supported types are `user_reference`, `schedule_reference`.
* `service` - There can be multiple services associated with a policy.
* `team` - (Optional) Teams associated with the policy. Account must have the `teams` ability to use this parameter. There can be multiple teams associated with a policy.


## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the response play.

## Import

Response Plays can be imported using the `id.from(email)`, e.g.

```
$ terraform import pagerduty_response_play.main 16208303-022b-f745-f2f5-560e537a2a74.user@email.com
```
