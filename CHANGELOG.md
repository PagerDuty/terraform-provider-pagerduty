## 1.6.2 (Unreleased)
BUG FIXES:
* resource_pagerduty_service_integration: Fix panic when reading [GH-214] 

## 1.6.1 (April 09, 2020)
BUG FIXES:
* Added links to `pagerduty_ruleset` and `pagerduty_ruleset_rule` to side nav([#198](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/198))

* Fixed importing on `pagerduty_ruleset` and `pagerduty_ruleset_rule` also added import testing. ([#199](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/199))

## 1.6.0 (April 07, 2020)

FEATURES:
* **New Resources:** `pagerduty_ruleset` and `pagerduty_ruleset_rule` ([#195](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/195))

BUG FIXES:
* resource/resource_pagerduty_team_membership: Docs: Fixed Team membership role defaults to manager ([#194](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/194))

IMPROVEMENTS:
* `resource/resource_pagerduty_service` and `resource/resource_pagerduty_service_integration`:mplement retry logic on read([#191](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/191))

## 1.5.1 (March 19, 2020)

FEATURES:
* **New Resource:** `pagerduty_user_notification_rule` ([#180](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/180))

BUG FIXES:
* resource/resource_pagerduty_service: Fixed alert_grouping_timeout to be able to accept values of 0 ([#190](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/190))

IMPROVEMENTS:
* resource/pagerduty_team_membership: add `role` to the resource ([#151](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/151))

* Update terraform plugin SDK from 1.0.0 to 1.7.0 ([#188](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/188))
## 1.5.0 (March 05, 2020)

BUG FIXES:

* data_source_pagerduty_user: Docs: update `team_responder` role thath as been renamed to `observer` ([#179](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/179))
* data_source_pagerduty_vendor: Update vendor id to fix test failures ([#178](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/178))

IMPROVEMENTS:

* resource/resource_pagerduty_user: Remove deprecated teams field from example in doc ([#185](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/185))
* Add retry to team read and schedule read ([#186](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/186))
* Update description to match other official providers ([#183](https://github.com/terraform-providers/terraform-provider-pagerduty/pull/183))

## 1.4.2 (January 30, 2020)

BUG FIXES:

* resource/resource_pagerduty_service: Fix service to populate the `alert_grouping` and `alert_grouping_timeout` fields when reading resource ([#177](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/177))
* resource/resource_pagerduty_event_rule: Changing pagerduty_event_rule.catch_all field to Computed ([#169](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/169))
* data-source/pagerduty_vendor: Fix the exact matching of vendor name when it contains special chars ([#166](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/166))

IMPROVEMENTS:
* resource/resource_pagerduty_service: improve formatting in document to better highlight `intelligent`([#172](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/172))
* resource/resource_pagerduty_extension: clarified `endpoint_url` with a note that sometimes it is required ([#164](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/164))

## 1.4.1 (October 24, 2019)

BUG FIXES:

* resource/pagerduty_team_membership: Handle missing user referenced by team membership ([#153](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/153))

* resource/pagerduty_event_rule: Fix perpetual diff issue with advanced conditions ([#157](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/157)). Labeled advanced condition field as optional in documentation  ([#160](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/160))

* resource/pagerduty_user: Documentation fixed list of valid colors ([#154](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/154))

IMPROVEMENTS:
* Switch to standalone Terraform Plugin SDK: ([#158](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/158))

* Add html_url read-only attribute to resource_pagerduty_service, resource_pagerduty_extension, resource_pagerduty_team ([#162](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/162))  

* resource/pagerduty_event_rule: Documentation for `depends_on` field ([#152](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/152)).

## 1.4.0 (August 23, 2019)

NOTES:

* resource/pagerduty_user: The `teams` attribute has been deprecated in favor of the `pagerduty_team_membership` resource ([#146](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/146))

FEATURES:

* **New Data Source:** `pagerduty_service` ([#141](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/141))
* **New Resource:** `pagerduty_event_rule` ([#150](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/150))

BUG FIXES:

* resource/pagerduty_maintenance_window: Allow services to be unordered ([#142](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/142))

IMPROVEMENTS:

* resource/pagerduty_service: Add support for alert_grouping and alert_grouping_timeout ([#143](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/143))

## 1.3.1 (July 29, 2019)

BUG FIXES:

* resource/pagerduty_user: Remove invalid role types ([#135](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/135))
* resource/pagerduty_service: Remove status from payload ([#133](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/133))

## 1.3.0 (May 29, 2019)

BUG FIXES:

* data-source/pagerduty_team: Fix team search issue [[#110](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/110)] 
* resource/pagerduty_maintenance_window: Suppress spurious diff in `start_time` & `end_time` ([#116](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/116))
* resource/pagerduty_service: Set invitation_sent [[#127](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/127)] 
* resource/pagerduty_escalation_policy: Correctly set teams ([#129](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/129))

IMPROVEMENTS:

* Switch to Terraform 0.12 SDK which is required for Terraform 0.12 support. This is the first release to use the 0.12 SDK required for Terraform 0.12 support. Some provider behaviour may have changed as a result of changes made by the new SDK version ([#126](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/126))

## 1.2.1 (November 21, 2018)

BUG FIXES:

* resource/pagerduty_service: Fix `scheduled_actions` bug ([#99](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/99))

## 1.2.0 (August 16, 2018)

IMPROVEMENTS:

* resource/pagerduty_extension: `endpoint_url` is now an optional field ([#83](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/83))
* resource/pagerduty_extension: Manage extension configuration as plain JSON ([#84](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/84))
* resource/pagerduty_service_integration: Documentation regarding integration url for events ([#91](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/91))
* resource/pagerduty_user_contact_method: Add support for push_notification_contact_method ([#93](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/93))

## 1.1.1 (May 30, 2018)

BUG FIXES:

* Fix `Unable to locate any extension schema` bug ([#79](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/79))

## 1.1.0 (April 12, 2018)

FEATURES:

* **New Data Source:** `pagerduty_team` ([#65](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/65))

IMPROVEMENTS:

* resource/pagerduty_service: Don't re-create services if support hours or scheduled actions change ([#68](https://github.com/terraform-providers/terraform-provider-pagerduty/issues/68))


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
