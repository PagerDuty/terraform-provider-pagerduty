## 1.1.0 (Unreleased)

FEATURES:

* **New Data Source:** `pagerduty_team` [GH-65]

IMPROVEMENTS:

* resource/pagerduty_service: Don't re-create services if support hours or scheduled actions change [GH-68]


## 1.0.0 (March 08, 2018)

FEATURES:

IMPROVEMENTS:

* **New Resource:** `pagerduty_extension` ([#69](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/69))
* **New Data Source:** `pagerduty_extension_schema` ([#69](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/69))

BUG FIXES: 
* r/service_integration: Add `html_url` as a computed value ([#59](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/59))
* r/pagerduty_service: allow incident_urgency_rule to be computed ([#63](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/63))
* d/pagerduty_vendor: Match whole string and fallback to partial matching ([#55](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/55))

## 0.1.3 (January 16, 2018)

FEATURES:

* **New Resource:** `pagerduty_user_contact_method` ([#29](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/29))

IMPROVEMENTS:

* r/pagerduty_service: Add alert_creation attribute ([#38](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/38))
* r/pagerduty_service_integration: Allow for generation of events-api-v2 service integration ([#40](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/40))

BUG FIXES: 
* r/pagerduty_service: Allow disabling service incident timeouts ([#44](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/44)] [[#52](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/52))
* r/pagerduty_schedule: Add support for overflow ([#23](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/23))
* r/pagerduty_schedule: Don't read back `start` value for schedule layers ([#35](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/35))
* r/pagerduty_service_integration: Fix import issue when more than 100 services exist ([#39](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/39)] [[#47](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/47))

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
