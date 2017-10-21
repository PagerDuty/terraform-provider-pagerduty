## 0.1.3 (Unreleased)

IMPROVEMENTS:

* r/pagerduty_service: Add alert_creation attribute [GH-38]
* r/pagerduty_service_integration: Allow for generation of events-api-v2 service integration [GH-40]

BUG FIXES: 
* r/pagerduty_service: Allow disabling service incident timeouts [GH-44]
* r/pagerduty_schedule: Add support for overflow [GH-23]
* r/pagerduty_schedule: Don't read back `start` value for schedule layers [GH-35]
* r/pagerduty_service_integration: Set Limit for /services GET to be at most 100 results when importing a service integration [GH-39]

## 0.1.2 (August 10, 2017)

BUG FIXES: 

* resource/pagerduty_service_integration: Fix panic on nil `integration_key` [#20]

## 0.1.1 (August 04, 2017)

FEATURES:

* **New Resource:** `pagerduty_team_membership` ([#15](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/15))
* **New Resource:** `pagerduty_maintenance_window` ([#17](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/17))

IMPROVEMENTS: 

* r/pagerduty_user: Set time_zone as optional ([#19](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/19))

BUG FIXES:

* resource/pagerduty_service: Fixing updates for `escalation_policy` ([#7](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/7))
* resource/pagerduty_schedule: Fix diff issues related to `start`, `rotation_virtual_start`, `end` ([#4](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/4))
* r/pagerdy_service_integration: Protect against panics on imports ([#16](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/16))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
