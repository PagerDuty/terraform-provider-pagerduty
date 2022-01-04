## 2.2.1 (December 10, 2021)
BUG FIX:
* `resource/pagerduty_user`: Fix in go library for user object ([ref](https://github.com/heimweh/go-pagerduty/pull/74))

## 2.2.0 (December 6, 2021)
FEATURES:
* `resource/pagerduty_webhook_subscriptions`: Added resource to support Webhooks v3 ([#420](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/420))
* `resource/pagerduty_business_service_subscriber`: Added resource to support Business Service Subscribers ([#414](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/414))

IMPROVEMENTS:
* `pagerduty/provider`: Added `api_url_override` to support routing PagerDuty API traffic through a proxy ([#366](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/366))
* `pagerduty/provider`: Adding various fetch functions to fix race conditions with PagerDuty API ([#380](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/380))

BUG FIXES:
* `data_source/pagerduty_vendor`: Correct example in docs ([#417](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/417))
* Documentation fixes ([#418](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/418))

## 2.1.1 (October 29, 2021)
BUG FIXES:
* `resource/pagerduty_user`: Fixed caching ([#413](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/413))

IMPROVEMENTS:
* `resource/pagerduty_extension`: Added plumbing to enable paging in a future release  ([#413](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/413))
* `resource/pagerduty_service`: Added plumbing to enable email filtering in a future release  ([#413](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/413))

## 2.1.0 (October 25, 2021)

FEATURES:
* `pagerduty/provider`: Added support for `service_region` and `user_token` tags on the provider block ([#384](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/384))
* Added GitHub Actions workflow to run tests ([#406](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/406))

IMPROVEMENTS:
* Docs: `resource/pagerduty_slack_connection`: Added `notification_type` to sample code ([#408](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/408))
* `resource/pagerduty_escalation_policy`: Added validation for teams maxitems 1 ([#412](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/412))

BUG FIXES:
* Docs: `resource/pagerduty_tag`: Corrected sample code ([#409](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/409))
* Docs: `index`: Removed deprecated `teams` field from sample code ([#411](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/411))


## 2.0.0 (October 11, 2021)
  FEATURES:
* `resource/pagerduty_tag`: Added resource to manage Tags ([#402](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/402))

* `resource/pagerduty_tag_assignment`: Added resource to manage Tag Assignments ([#402](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/402))

* `data_source_pagerduty_tag`: Added data source for Tags ([#402](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/402))

* Update Terraform Plugin SDK to v2 ([#375](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/375))


IMPROVEMENTS:
* `resource/pagerduty_service_integration`: Added validation that ensures an email address is specified for email integrations ([#382](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/382))

* `resource/pagerduty_schedule`: Added validation for `start_time_of_day` format ([#383](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/383))

* `resource/pagerduty_schedule`: Added validation for `start_day_of_week` format ([#385](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/385))

* `resource/pagerduty_schedule`: Added validation that `start_day_of_week` is only set when `weekly_restriction` is set as `restriction.type` ([#386](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/386))


* `resource/pagerduty_service`: CustomizeDiff to ensure general urgency rule is not set for an urgency rule of type support hours ([#387](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/387))

* Docs: `resource/pagerduty_rulset_rule`: Update severity order to reflect criticality ([#392](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/392))

* `resource/pagerduty_escalation_policy`: Added validation to ensure `num_loops` stays between `0` and `9` ([#395](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/395))

* `resource/pagerduty_escalation_policy`: Added validation to ensure `escalation_delay_in_minutes` is an int and greater than `0` ([#401](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/401))

* `resource/pagerduty_schedule`: Added validation to `rotation_turn_length_seconds` ([#405](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/405))



BUG FIXES:
* Docs: `resource/pagerduty_service`: Corrected `acknowledgement_timeout` to add statement about default being set to 1800.  ([#393](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/393))

* `resource/pagerduty_service` and `resource/pagerduty_service_dependency`: Fix `alert_grouping` and `alert_grouping_timeout` conflicting with `alert_grouping_parameters` ([#377](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/377))

* `resource/pagerduty_service_dependency`: Fix sporadic panic when processing service dependencies ([#376](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/376))

## 1.11.0 (September 9, 2021)
FEATURES:
* `resource/pagerduty_slack_connection`: Added resource to manage Slack Connections (aka Slack v2) ([#381](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/381))

BUG FIXES:
* `resource/pagerduty_service`: Clarified docs surrounding the `scheduled_actions` block on the service resource ([#360](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/360))

## 1.10.1 (August 16, 2021)
FEATURES:
* `data_source_pagerduty_service_integration`: Add new data source for Service Integrations ([#363](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/363))
* `resource/pagerduty_schedule`: Add `teams` field to schedule resource ([#368](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/368))

BUG FIXES:
* `resource/pagerduty_extension_servicenow`: Fixed blank `snow_password` when changes ignored ([#371](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/371))
* `data_source_pagerduty_business_service`: Fixed bug where any business service can now be a data source ([#372](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/372))
* Various Resources: Updated links to API reference throughout documentation ([#369](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/369))

## 1.10.0 (July 29, 2021)
FEATURES:
* `resource/pagerduty_extension_servicenow`: Added `pagerduty_extension_servicenow` resource to account for specific values that need to be set for ServiceNow extension ([#348](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/348))
* `resource/pagerduty_schedule`: Added `alert_grouping_parameters`field to `resource_pagerduty_service` ([#342](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/342))


BUG FIXES:
* `resource/pagerduty_user`: Fixed broken tests on user objects because of new API constraints ([#362](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/362))


IMPROVEMENTS:
* `resource/pagerduty_schedule`: Set `MinItem: 1` parameter for the schedule layers field ([#350](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/350))
* `resource/pagerduty_user`: Added option for caching user objects ([#362](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/362))

## 1.9.9 (June 23, 2021)
IMPROVEMENT:
* `resource/pagerduty_ruleset_rule`: Update sample code in documentation to reference a Terraform resource rather than hard coded service ID ([#349](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/349))


## 1.9.8 (June 17, 2021)
BUG FIXES:
* `resource/pagerduty_schedule`: Fixed the spurious diff in the `start` field when it's set to the past ([#343](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/343))

## 1.9.7 (May 27, 2021)
BUG FIXES:
* `resource/pagerduty_escalation_policy`: Fixed `num_loops` field so that it could be unset to match API behavior ([#324](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/324))
* `resource/pagerduty_user_contact_method`: Corrected error message on exceptions during import ([#327](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/327))


IMPROVEMENTS:
* `resource/pagerduty_team`,`resource/pagerduty_ruleset_rule`,`resource/pagerduty_schedule`,`resource/pagerduty_service_event_rule`: Documentation clarifications, fixes and improvements ([#322](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/322))
* `resource/pagerduty_team_membership`,`resource/pagerduty_user`: Update user role documentation to add clarity over user vs team roles  ([#325](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/325))
* General provider improvement: Update dependencies to use go1.16 ([#326](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/326))
* General provider improvement: Update release task to use go1.16 ([#333](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/333))
* `resource/pagerduty_ruleset_rule`: Clarified documentation on the `start_time` of scheduled event rules ([#340](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/340))


## 1.9.6 (April 6, 2021)
BUG FIXES:
* `resource/pagerduty_response_play`: Fixed `responder` field to be optional to match API behavior ([#316](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/316))
* `resource/pagerduty_schedule`: Suppressing the equal diff on schedule layer timestamps ([#321](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/321))

IMPROVEMENTS:
* `resource/ruleset_rule`: Added clarification to documentation on event rule actions ([#317](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/317))

FEATURES:
* `resource/pagerduty_team`: Added `parent` field to resource to support team hierarchy ([#319](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/319))

## 1.9.5 (March 11, 2021)
BUG FIXES:
* `data_source_pagerduty_ruleset`: Fixed bug by adding `routing_keys` to data source schema ([#312](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/312))

IMPROVEMENTS:
* `resource/pagerduty_escalation_policy`: Added retry logic to escalation policy delete ([#309](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/309))
* `resource/pagerduty_user`: Trimmed leading and trailing spaces on the value for `name` field to match behavior of PagerDuty API ([#312](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/312))

FEATURES:
* `resource/pagerduty_ruleset_rule` and `resource/pagerduty_service_event_rule`: Add `template` field to the rule object([#314](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/314))

## 1.9.4 (February 26, 2021)
BUG FIXES:
* `resource/pagerduty_team_membership`: Fixed issue with importing team members to teams with more than 100 users ([#305](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/305))

IMPROVEMENTS:
* `data_source_pagerduty_ruleset`: Added `routing_keys` field to the `ruleset` object ([#305](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/305))

## 1.9.3 (February 11, 2021)
BUG FIXES: 
* `resource/pagerduty_service_event_rule`,`resource/pagerduty_ruleset_rule`: Fixed Bug with Event Rule Suppress Action ([#302](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/302))

## 1.9.2 (February 10, 2021)
BUG FIXES: 
* `resource/pagerduty_service_event_rule`,`resource/pagerduty_ruleset_rule`: Fixed Bug with Event Rule Positioning ([#301](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/301))

## 1.9.1 (February 8, 2021)
FEATURES:
* `resource/pagerduty_service_event_rule`: Add service event rule resource ([#296](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/296))

IMPROVEMENTS:
* `resource/resource_pagerduty_user`: Add retry logic to user update error ([#286](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/286))
* `data_source_pagerduty_busines_service`, `data_source_pagerduty_escalation_policy`, `data_source_pagerduty_extension_schema`, `data_source_pagerduty_priority`, `data_source_pagerduty_ruleset`, `data_source_pagerduty_schedule`, `data_source_pagerduty_service`, `data_source_pagerduty_team`, `data_source_pagerduty_user`, `data_source_pagerduty_vendor`: Add retry logic on error ([#287](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/287))
* `resource/pagerduty_ruleset_rule`: Add `suspend` and `variables` fields ([#296](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/296))

BUG FIXES:
* `resource/pagerduty_ruleset_rule`: Fixed bug where `position` wasn't setting properly. Add retry to allow for `position` to be correctly set before apply is complete ([#296](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/296))

## 1.8.0 (November 30, 2020)
FEATURES:
* `resource/resource_pagerduty_response_play` Add response play resource ([#278](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/278))

IMPROVEMENTS:
* `resource/resource_pagerduty_extension` Make endpoint_url sensitive ([#281](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/281))

BUG FIXES:
* `resource/resource_ruleset_rule` Fix issue with disabling ruleset rule ([#282](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/282))
* `resource/resource_service_integration` Fix documentation to have `integration_email` show fully-qualified email address ([#284](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/284))

## 1.7.11 (November 3, 2020)
IMPROVEMENTS:
* `resource/resource_pagerduty_service` Added validation to Service Name field. ([#275](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/275))

BUG FIXES:
* `resource/resource_pagerduty_escalation_policy` Fixed retry logic on Escalation Policy Read function. ([#275](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/275))

* README -- fixed typos ([#270](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/270))
* `resource/resource_pagerduty_event_rule` Docs -- fixed typo in NOTE ([#267](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/267))

## 1.7.10 (August 27, 2020)
IMPROVEMENTS:
* `resource/resource_pagerduty_schedule`, `resource/resource_pagerduty_escalation_policy` added retry logic to updating schedule and escalation_policy ([#264](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/264))

## 1.7.4 (July 27, 2020)

FEATURES:
* `resource/resource_pagerduty_business_service` Add team to business_service resource ([#246](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/246))

BUG FIXES:
* `resource/resource_pagerduty_user` Docs -- fixed typo in user_roles info ([#248](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/248))
* `resource/resource_pagerduty_user` Docs -- added time_zone field ([#256](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/256))



IMPROVEMENTS:
* `resource/resource_pagerduty_ruleset` Docs -- added example of Default Global Ruleset ([#239](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/239))
* `resource/resource_pagerduty_escalation_policy` Docs -- extending example of multiple targets([#242](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/242))
* `resource/resource_pagerduty_service` Docs -- extending example of alert creation on service([#243](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/243))
* `resource/resource_pagerduty_service_integration` Docs -- extending example of events v2 and email integrations on service_integration([#244](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/244))
* `resource/resource_pagerduty_schedule` Allowing "Removal of Schedule Layers ([#257](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/257))
* `resource/resource_pagerduty_schedule`,`resource/resource_pagerduty_team_membership`, `resource/resource_pagerduty_team_user` adding retry logic to deleting schedule,team_membership and user ([#258](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/258))

## 1.7.3 (June 12, 2020)
FEATURES

* Update service_dependency to support technical service dependencies ([#238](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/238))
* Implement retry logic on all reads ([#208](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/208))
* Bump golang to v1.14.1 ([#193](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/193))

BUG FIXES: 
* data_source_ruleset: add example of Default Global Ruleset in Docs ([#239](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/239))

## 1.7.2 (June 01, 2020)
FEATURES
* **New Data Source:** `pagerduty_ruleset`  ([#237](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/237))
* Update docs/tests to TF 0.12 syntax ([#223](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/223))

BUG FIXES: 
* testing: update sweepers ([#220](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/220))
* data_source_priority: adding doc to sidebar nav ([#221](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/221))

## 1.7.1 (April 29, 2020)
FEATURES:
* **New Data Source:** `pagerduty_priority`  ([#219](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/219))

BUG FIXES:
* resource_pagerduty_service: Fix panic  ([#218](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/218))

## 1.7.0 (April 20, 2020)
FEATURES:
* **New Resources:** `pagerduty_business_service` and `pagerduty_service_dependency` ([#213](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/213))

BUG FIXES:
* resource_pagerduty_service_integration: Fix panic when reading ([#214](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/214))
* resource_pagerduty_ruleset_rule: Fix Import of catch_all rules ([#205](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/205))
* resource_pagerduty_ruleset_rule: Fixing multi-rule creation bug and suppress rule panic ([#211](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/211))

## 1.6.1 (April 09, 2020)
BUG FIXES:
* Added links to `pagerduty_ruleset` and `pagerduty_ruleset_rule` to side nav([#198](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/198))

* Fixed importing on `pagerduty_ruleset` and `pagerduty_ruleset_rule` also added import testing. ([#199](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/199))

## 1.6.0 (April 07, 2020)

FEATURES:
* **New Resources:** `pagerduty_ruleset` and `pagerduty_ruleset_rule` ([#195](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/195))

BUG FIXES:
* resource/resource_pagerduty_team_membership: Docs: Fixed Team membership role defaults to manager ([#194](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/194))

IMPROVEMENTS:
* `resource/resource_pagerduty_service` and `resource/resource_pagerduty_service_integration`: implement retry logic on read([#191](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/191))

## 1.5.1 (March 19, 2020)

FEATURES:
* **New Resource:** `pagerduty_user_notification_rule` ([#180](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/180))

BUG FIXES:
* resource/resource_pagerduty_service: Fixed alert_grouping_timeout to be able to accept values of 0 ([#190](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/190))

IMPROVEMENTS:
* resource/pagerduty_team_membership: add `role` to the resource ([#151](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/151))

* Update terraform plugin SDK from 1.0.0 to 1.7.0 ([#188](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/188))
## 1.5.0 (March 05, 2020)

BUG FIXES:

* data_source_pagerduty_user: Docs: update `team_responder` role that as been renamed to `observer` ([#179](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/179))
* data_source_pagerduty_vendor: Update vendor id to fix test failures ([#178](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/178))

IMPROVEMENTS:

* resource/resource_pagerduty_user: Remove deprecated teams field from example in doc ([#185](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/185))
* Add retry to team read and schedule read ([#186](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/186))
* Update description to match other official providers ([#183](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/183))

## 1.4.2 (January 30, 2020)

BUG FIXES:

* resource/resource_pagerduty_service: Fix service to populate the `alert_grouping` and `alert_grouping_timeout` fields when reading resource ([#177](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/177))
* resource/resource_pagerduty_event_rule: Changing pagerduty_event_rule.catch_all field to Computed ([#169](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/169))
* data-source/pagerduty_vendor: Fix the exact matching of vendor name when it contains special chars ([#166](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/166))

IMPROVEMENTS:
* resource/resource_pagerduty_service: improve formatting in document to better highlight `intelligent`([#172](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/172))
* resource/resource_pagerduty_extension: clarified `endpoint_url` with a note that sometimes it is required ([#164](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/164))

## 1.4.1 (October 24, 2019)

BUG FIXES:

* resource/pagerduty_team_membership: Handle missing user referenced by team membership ([#153](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/153))

* resource/pagerduty_event_rule: Fix perpetual diff issue with advanced conditions ([#157](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/157)). Labeled advanced condition field as optional in documentation  ([#160](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/160))

* resource/pagerduty_user: Documentation fixed list of valid colors ([#154](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/154))

IMPROVEMENTS:
* Switch to standalone Terraform Plugin SDK: ([#158](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/158))

* Add html_url read-only attribute to resource_pagerduty_service, resource_pagerduty_extension, resource_pagerduty_team ([#162](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/162))  

* resource/pagerduty_event_rule: Documentation for `depends_on` field ([#152](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/152)).

## 1.4.0 (August 23, 2019)

NOTES:

* resource/pagerduty_user: The `teams` attribute has been deprecated in favor of the `pagerduty_team_membership` resource ([#146](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/146))

FEATURES:

* **New Data Source:** `pagerduty_service` ([#141](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/141))
* **New Resource:** `pagerduty_event_rule` ([#150](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/150))

BUG FIXES:

* resource/pagerduty_maintenance_window: Allow services to be unordered ([#142](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/142))

IMPROVEMENTS:

* resource/pagerduty_service: Add support for alert_grouping and alert_grouping_timeout ([#143](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/143))

## 1.3.1 (July 29, 2019)

BUG FIXES:

* resource/pagerduty_user: Remove invalid role types ([#135](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/135))
* resource/pagerduty_service: Remove status from payload ([#133](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/133))

## 1.3.0 (May 29, 2019)

BUG FIXES:

* data-source/pagerduty_team: Fix team search issue [[#110](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/110)] 
* resource/pagerduty_maintenance_window: Suppress spurious diff in `start_time` & `end_time` ([#116](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/116))
* resource/pagerduty_service: Set invitation_sent [[#127](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/127)] 
* resource/pagerduty_escalation_policy: Correctly set teams ([#129](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/129))

IMPROVEMENTS:

* Switch to Terraform 0.12 SDK which is required for Terraform 0.12 support. This is the first release to use the 0.12 SDK required for Terraform 0.12 support. Some provider behaviour may have changed as a result of changes made by the new SDK version ([#126](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/126))

## 1.2.1 (November 21, 2018)

BUG FIXES:

* resource/pagerduty_service: Fix `scheduled_actions` bug ([#99](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/99))

## 1.2.0 (August 16, 2018)

IMPROVEMENTS:

* resource/pagerduty_extension: `endpoint_url` is now an optional field ([#83](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/83))
* resource/pagerduty_extension: Manage extension configuration as plain JSON ([#84](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/84))
* resource/pagerduty_service_integration: Documentation regarding integration url for events ([#91](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/91))
* resource/pagerduty_user_contact_method: Add support for push_notification_contact_method ([#93](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/93))

## 1.1.1 (May 30, 2018)

BUG FIXES:

* Fix `Unable to locate any extension schema` bug ([#79](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/79))

## 1.1.0 (April 12, 2018)

FEATURES:

* **New Data Source:** `pagerduty_team` ([#65](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/65))

IMPROVEMENTS:

* resource/pagerduty_service: Don't re-create services if support hours or scheduled actions change ([#68](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/68))


## 1.0.0 (March 08, 2018)

FEATURES:

IMPROVEMENTS:

* **New Resource:** `pagerduty_extension` ([#69](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/69))
* **New Data Source:** `pagerduty_extension_schema` ([#69](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/69))

BUG FIXES:

* r/service_integration: Add `html_url` as a computed value ([#59](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/59))
* r/pagerduty_service: allow incident_urgency_rule to be computed ([#63](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/63))
* d/pagerduty_vendor: Match whole string and fallback to partial matching ([#55](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/55))

## 0.1.3 (January 16, 2018)

FEATURES:

* **New Resource:** `pagerduty_user_contact_method` ([#29](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/29))

IMPROVEMENTS:

* r/pagerduty_service: Add alert_creation attribute ([#38](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/38))
* r/pagerduty_service_integration: Allow for generation of events-api-v2 service integration ([#40](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/40))

BUG FIXES:

* r/pagerduty_service: Allow disabling service incident timeouts ([#44](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/44)] [[#52](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/52))
* r/pagerduty_schedule: Add support for overflow ([#23](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/23))
* r/pagerduty_schedule: Don't read back `start` value for schedule layers ([#35](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/35))
* r/pagerduty_service_integration: Fix import issue when more than 100 services exist ([#39](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/39)] [[#47](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/47))

## 0.1.2 (August 10, 2017)

BUG FIXES:

* resource/pagerduty_service_integration: Fix panic on nil `integration_key` [#20]

## 0.1.1 (August 04, 2017)

FEATURES:

* **New Resource:** `pagerduty_team_membership` ([#15](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/15))
* **New Resource:** `pagerduty_maintenance_window` ([#17](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/17))

IMPROVEMENTS:

* r/pagerduty_user: Set time_zone as optional ([#19](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/19))

BUG FIXES:

* resource/pagerduty_service: Fixing updates for `escalation_policy` ([#7](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/7))
* resource/pagerduty_schedule: Fix diff issues related to `start`, `rotation_virtual_start`, `end` ([#4](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/4))
* r/pagerduty_service_integration: Protect against panics on imports ([#16](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/16))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
