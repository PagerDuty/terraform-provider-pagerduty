## v3.19.3 (January 27, 2025)

IMPROVEMENTS:

* Expand upon documentation for Event Orchestrations ([971](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/971))
* `pagerduty_incident_workflow_trigger` remove value check on conditional ([964](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/964))

## v3.19.2 (January 15, 2025)

FIXES:

* Address: Unable to import service dependencies ([970](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/970))

## v3.19.1 (January 13, 2025)

FIXES:

* Add support for Incident Types (retry failed release)

## v3.19.0 (January 13, 2025)

FEATURES:

* Add support for **Incident Types** ([962](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/962))
  - `data/pagerduty_incident_type`
  - `data/pagerduty_incident_type_custom_field`
  - `resource/pagerduty_incident_type`
  - `resource/pagerduty_incident_type_custom_field`

## v3.18.3 (December 24, 2024)

IMPROVEMENTS:

* Fix service update conflict on alert grouping rules ([959](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/959))

## v3.18.2 (December 19, 2024)

IMPROVEMENTS:

* Install a fresh terraform binary if we're unable to locate a matching version in PATH. ([941](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/941))

## v3.18.1 (November 20, 2024)

BUG FIXES:

* Avoid panic at null workflow in IW trigger, allowing it to refresh ([955](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/955))

## v3.18.0 (November 19, 2024)

FEATURES:

* Add `only_invocable_on_unresolved_incidents` to Automation Action's schema ([945](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/945))

IMPROVEMENTS:

* Add description to `alert_grouping_setting` docs ([953](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/953))

## v3.17.2 (November 11, 2024)

BUG FIXES:

* Address: Using `pagerduty_alert_grouping_setting` causes error HTTP response failed with status code 404 and no JSON error object was present ([952](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/952))

## v3.17.1 (November 8, 2024)

BUG FIXES:

* Fix cursor for datasource alert grouping setting ([950](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/950))

## v3.17.0 (November 4, 2024)

FEATURES:

* `resource/pagerduty_jira_cloud_account_mapping_rule`: Add support for Jira Cloud integration ([942](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/942))
* `data/pagerduty_jira_cloud_account_mapping`: Add support for Jira Cloud integration ([942](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/942))

IMPROVEMENTS:

* `resource/pagerduty_service`: Add link to migration docs in `alert_grouping_parameters` deprecation ([947](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/947))
* `PagerDuty/pagerduty`: Fix `pagerudty.com` typo in docs ([939](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/939))

## v3.16.0 (October 14, 2024)

FEATURES:

* `resource/pagerduty_alert_grouping_setting`: Add resource and data source for Alert Grouping Setting ([935](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/935))

## v3.15.7 (October 9, 2024)

IMPROVEMENTS:

* `resource/pagerduty_escalation_policy`: Ensure escalation rule targets are added in same order as the plan ([#937](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/937))

## v3.15.6 (August 23, 2024)

IMPROVEMENTS:

* `resource/pagerduty_team_membership`: Return remediation measures for team membership deletion with EP dependencies ([#863](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/863))

## v3.15.5 (August 22, 2024)

IMPROVEMENTS:

* `PagerDuty/pagerduty`: Second attempt to address: API calls cancelled by client because of Timeout ([#928](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/928))

## v3.15.4 (August 16, 2024)

IMPROVEMENTS:

* `PagerDuty/pagerduty`: Address: API calls cancelled by client because of Timeout ([#926](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/926))

## v3.15.3 (August 14, 2024)

IMPROVEMENTS:

* `resource/pagerduty_business_service`: Address: Errors when trying to apply a business service that was manually deleted from PD ([#925](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/925))

## v3.15.2 (August 12, 2024)

IMPROVEMENTS:

* `resource/pagerduty_team_membership`: fix: improve team_membership deletion logic ([#918](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/918))

## v3.15.1 (July 31, 2024)

IMPROVEMENTS:

* `resource/pagerduty_event_orchestration_router`: Update code samples in the pagerduty_event_orchestration_router resource docs ([#917](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/917))

BUG FIXES:

* `resource/pagerduty_service`: Update flattening of auto pause notif params in service ([#919](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/919))
* `resource/pagerduty_service_integration`: Address: Service email integrations - unable to configure text after and text before value extractors ([#920](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/920))
* `PagerDuty/pagerduty`: Allow appending to UserAgent with ldflag ([#921](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/921))

## v3.15.0 (July 22, 2024)

FEATURES:

* `resource/pagerduty_event_orchestration_router`, `resource/pagerduty_event_orchestration_global`, `resource/pagerduty_event_orchestration_service`: Event Orchestration: add support for Dynamic Routing and Dynamic Escalation Policy Assignment ([#885](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/885))


## v3.14.5 (July 4, 2024)

IMPROVEMENTS:

* Migrate pack #3 of datasources and resources to terraform plugin framework ([#896](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/896))
  - `data/pagerduty_license`
  - `data/pagerduty_licenses`
  - `data/pagerduty_priority`
  - `resource/pagerduty_team`
* `PagerDuty/pagerduty`: Address: API calls canceled by client because of Timeout ([#900](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/900))

## v3.14.4 (June 26, 2024)

IMPROVEMENTS:

* Use all pages to search for `data.pagerduty_escalation_policy` ([#894](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/894))

BUG FIXES:

* Prevent null-pointer panic at escalation policy ([#894](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/894))

## v3.14.3 (June 20, 2024)

BUG FIXES:

* Handle 403 at EO path service refresh as orphan ([#890](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/890))

## v3.14.2 (June 19, 2024)

BUG FIXES:

* Remove invalid `alert_grouping_parameters` fields after strict check in API ([#888](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/888))

## v3.14.1 (June 19, 2024)

BUG FIXES:

* Fix datasource `pagerduty_service_integration` not finding service when they're many ([#886](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/886))

## 3.14.0 (June 14, 2024)

FEATURES:

* `PagerDuty/pagerduty`: Add the option to ignore TLS certificate errors when calling the PD API. ([#881](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/881))

IMPROVEMENTS:

* Migrate pack #2 of datasources and resources to terraform plugin framework ([#866](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/866))
  - `resource/pagerduty_service_dependency`
  - `data/pagerduty_service_integration`
  - `data/pagerduty_service`
  - `resource/pagerduty_addon`
* `resource/pagerduty_user`: fix invalid target type in docs ([#856](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/856))
* `resource/pagerduty_service`: feat: improve PD service time window validation to accept 86400 as a valid value ([#876](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/876))
* `resource/pagerduty_event_orchestration_router`,  `resource/pagerduty_event_orchestration_service`: Ignore import state verification for EO routing rules ids ([#883](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/883))

## 3.13.1 (June 12, 2024)

FEATURES:

* `resource/pagerduty_user_handoff_notification_rule`: Add support for user handoff notification rules ([#875](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/875))

BUG FIXES:

* `resource/pagerduty_extension`: Allow resource extension's endpoint_url to be imported as null ([#880](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/880))

IMPROVEMENTS:

* `PagerDuty/pagerduty`: migrate goreleaser deprecated config attributes ([#882](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/882))

## 3.13.0 (June 12, 2024) - Not released because of failed release process

FEATURES:

* `resource/pagerduty_user_handoff_notification_rule`: Add support for user handoff notification rules ([#875](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/875))

BUG FIXES:

* `resource/pagerduty_extension`: Allow resource extension's endpoint_url to be imported as null ([#880](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/880))

## 3.12.2 (May 20, 2024)

IMPROVEMENTS:

* `data/pagerduty_service`: Address: Service DS not locating queried name when is not in first page of results ([#874](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/874))

## 3.12.1 (May 14, 2024)

BUG FIXES:

* `resource/pagerduty_service`: Prevent timeout for alert_grouping_parameters to be explicitly set to zero when they are of type "content_based" ([#871](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/871))

## 3.12.0 (May 9, 2024)

FEATURES:

* `resource/pagerduty_incident_workflow_trigger`: Add support for Incident Workflow triggers team restrictions ([#861](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/861))

IMPROVEMENTS:

* Migrate pack of datasources and resources to terraform plugin framework ([#816](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/816))
  - `data/pagerduty_extension_schema`
  - `data/pagerduty_tag`
  - `resource/pagerduty_extension_servicenow`
  - `resource/pagerduty_extension`
  - `resource/pagerduty_tag`
  - `resource/pagerduty_tag_assignment`

BUG FIXES:

* `resource/pagerduty_service`: Address: Support hours should be required in the plan if the incident_urgency_rule type is "use_support_hours" ([#868](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/868))

## 3.11.4 (Apr 16, 2024)

BUG FIXES:

* `resource/pagerduty_escalation_policy`: Handle malformed 403 Forbidden errors on EP updates ([#858](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/858))

## 3.11.3 (Apr 12, 2024)

BUG FIXES:

* `resource/pagerduty_service`: Handle nil pointer conversions across `pagerduty_service` implementation ([#854](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/854))

## 3.11.2 (Apr 10, 2024)

BUG FIXES:

* `resource/pagerduty_escalation_policy`: Re-address: 403 Forbidden error when updating existing pagerduty_escalation_policy's target ([#851](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/851))
* `resource/pagerduty_service`: Address: Pagerduty Provider Plugin Crashed ([#852](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/852))

## 3.11.1 (Apr 5, 2024)

BUG FIXES:

* `resource/pagerduty_escalation_policy`: Address: 403 Forbidden error when updating existing pagerduty_escalation_policy's target ([#846](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/846))

## 3.11.0 (Apr 3, 2024)

FEATURES:

* Add support for Event Orchestration Cache Variables ([#822](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/822))
  - `data/pagerduty_event_orchestration_global_cache_variable`
  - `data/pagerduty_event_orchestration_service_cache_variable`
  - `resource/pagerduty_event_orchestration_global_cache_variable`
  - `resource/pagerduty_event_orchestration_service_cache_variable`

IMPROVEMENTS:

* `resource/*`, `data/*`: Add retry config for `PagerDuty/go-pagerduty` ([#841](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/841))

## 3.10.1 (Mar 26, 2024)

BUG FIXES:

* `resource/pagerduty_team_membership`: fix(team membership): Allow 404 to propagate to avoid registering it in disassociated eps list ([#838](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/838))
* `resource/pagerduty_service_dependency`: Fix #832 service_dependency breaks if dependent services are deleted externally ([#834](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/834))

## 3.10.0 (Mar 15, 2024)

IMPROVEMENTS:

* `resource/resource_pagerduty_service`: Remove default value and enable diff suppression to account for planned end-of-life of create_incidents option.

## 3.9.0 (Feb 26, 2024)

FEATURES:

* `resource/pagerduty_business_sevice`: Migrate Resource pagerduty_business_service to TF Plugin Framework ([#808](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/808))
* `resource/pagerduty_event_orchestration_global`, `resource/pagerduty_event_orchestration_service`: Support for incident custom fields for Event Orchestration ([#749](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/749))

## 3.8.1 (Feb 20, 2024)

IMPROVEMENTS:

* `pagerduty/pagerduty`: Fix go module name ([#616](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/616))

## 3.8.0 (Feb 16, 2024)

FEATURES:

* `data/pagerduty_standards_resource_scores`,  `data/pagerduty_standards_resources_scores`: Add datasource standards resource scores and standards resources scores ([#812](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/812))

BUG FIXES:

* `resource/pagerduty_tag_assignment`: Address: pagerduty_tag_assignment teams change_tag 404 not found ([#818](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/818))
* `resource/pagerduty_incident_workflow`: Fix import for incident workflow ([#820](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/820))

IMPROVEMENTS:

* `resource/pagerduty_response_play`: Response play documentation EOL heads up ([#821](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/821))

## 3.7.1 (Feb 6, 2024)

BUG FIXES:

* `resource/pagerduty_event_orchestration_router`: Address: Fail to enable rules in pagerduty_event_orchestration_router ([#814](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/814))

## 3.7.0 (Jan 31, 2024)

FEATURES:

* `data/pagerduty_standards`: Add provider and a Standards datasource using Terraform Plugin Framework ([#787](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/787))

## 3.6.0 (Jan 26, 2024)

FEATURES:

* `resource/*`, `data/*`: feat: support for sourcing the provider service region from env var (take 2) ([#805](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/805))

IMPROVEMENTS:

* `resource/pagerduty_user`: Address permadiff when user name has repeated whitespaces ([#810](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/810))

## 3.5.2 (Jan 25, 2024)

IMPROVEMENTS:

* `resource/*`, `data/*`: Address: timeout while waiting for state to become 'success' ([#807](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/807))

## 3.5.1 (Jan 24, 2024)

BUG FIXES:

* `resource/*`, `data/*`: Revert "Support for sourcing the provider service region from env var" ([#804](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/804))

## 3.5.0 (Jan 23, 2024)

FEATURES:

* `resource/*`, `data/*`: feat: support for sourcing the provider service region from an env var ([#797](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/797))

IMPROVEMENTS:

* `resource/*`, `data/*`: Add API Client timeouts configuration for http transport ([#802](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/802))

## 3.4.1 (Jan 18, 2024)

BUG FIXES:

* `resource/pagerduty_service_integration`: Fix attribute deprecation warning ([#801](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/801))

IMPROVEMENTS:

* `resource/*`, `data/*`: Update go version to 1.20 ([#794](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/794))
* `resource/pagerduty_service`: Allow time_window to be set for content-based grouping (issue 788) ([#795](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/795))

## 3.4.0 (Dec 21, 2023)

FEATURES:

* `data/pagerduty_team_members`: add `pagerduty_team_members` data source ([#717](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/717))

IMPROVEMENTS:

* `resource/pagerduty_contact_method`: Update `pagerduty_user_contact_method` address validation logic ([#792](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/792))
* `resource/*`,  `data/*`: Bump golang.org/x/crypto from 0.11.0 to 0.17.0 ([#790](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/790))

## 3.3.1 (Dec 12, 2023)

IMPROVEMENTS:

* `resource/*`,  `data/*`: Add jitter correction to ratelimit headers handling ([#784](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/784))

## 3.3.0 (Dec 7, 2023)

FEATURES:

* `resource/pagerduty_escalation_policy`: Add Round Robin support to `pagerduty_escalation_policy` ([#781](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/781))

BUG FIXES:

* More gracefully reject colonCompoundIDs that aren't compound IDs ([#762](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/762))
  - `resource/pagerduty_automation-actions_action_service_association`
  - `resource/pagerduty_automation-actions_action_team_association`
  - `resource/pagerduty_automation-actions_runner_team_association`
  - `resource/pagerduty_event_orchestration_integration`
  - `resource/pagerduty_team_membership`

## 3.2.2 (Dec 4, 2023)

BUG FIXES:

* `resource/pagerduty_service`: Hotfix - Alert grouping parameters input validation broken ([#779](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/779))

## 3.2.1 (Dec 1, 2023)

IMPROVEMENTS:

* `resource/*`,  `data/*`: Support rate limiting throttling ([#777](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/777))

## 3.2.0 (Dec 1, 2023)

FEATURES:

* `resource/pagerduty_incident_workflow`: Feat/add iw inline steps inputs support ([#768](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/768))
* `resource/pagerduty_service`: Support for Intelligent Time Window to Alert Grouping Parameters ([#773](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/773))

IMPROVEMENTS:

* `resource/pagerduty_service_integration`: Deprecate integration_key attribute mutation for Service Integrations ([#775](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/775))

BUG FIXES:

* Revert add remaining delays for retries ([#776](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/776))
  - `data/source_pagerduty_event_orchestration`
  - `data/source_pagerduty_event_orchestrations`
  - `resource/pagerduty_automation_actions_action_service_association`
  - `resource/pagerduty_automation_actions_action_team_association`
  - `resource/pagerduty_automation_actions_action`
  - `resource/pagerduty_automation_actions_runner_team_association`
  - `resource/pagerduty_automation_actions_runner`
  - `resource/pagerduty_business_service_subscriber`
  - `resource/pagerduty_business_service`
  - `resource/pagerduty_escalation_policy`
  - `resource/pagerduty_event_orchestration_integration`
  - `resource/pagerduty_event_orchestration_path_global`
  - `resource/pagerduty_event_orchestration_path_router`
  - `resource/pagerduty_event_orchestration_path_service`
  - `resource/pagerduty_event_orchestration_path_unrouted`
  - `resource/pagerduty_event_orchestration`
  - `resource/pagerduty_schedule`
  - `resource/pagerduty_service_dependency`
  - `resource/pagerduty_service_event_rule`
  - `resource/pagerduty_service_integration`
  - `resource/pagerduty_service`
  - `resource/pagerduty_slack_connection`
  - `resource/pagerduty_tag_assignment`
  - `resource/pagerduty_tag`
  - `resource/pagerduty_team_membership`
  - `resource/pagerduty_team`
  - `resource/pagerduty_user`
  - `resource/pagerduty_webhook_subscription`

## 3.1.2 (Nov 17, 2023)

IMPROVEMENTS:

* `resource/pagerduty_schedule`: Address `pagerduty_schedule` validation error on weekly restriction wihtout Start of week day ([#764](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/764))
* `resource/pagerduty_event_orchestration_service`: Patch - Service Orchestration enable status doesn't opt in ([#769](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/769))
* `resource/pagerduty_incident_custom_field`,  `resource/pagerduty_incident_custom_field_option`: add missing properties to incident_custom_field resource documentation ([#751](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/751))
* `resource/pagerduty_service`: Address Service Alert Grouping parameters config permadiff ([#771](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/771))

## 3.1.1 (Nov 2, 2023)

IMPROVEMENTS:

* `resource/pagerduty_user_contact_method`: Specify the default country code for contact methods ([#723](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/723))
* `resource/pagerduty_service`: Update retry delay for Technical Service state read ([#763](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/763))

## 3.1.0 (Oct 25, 2023)

FEATURES:

* `resource/pagerduty_team`, `data/pagerduty_team`: teams: add support for private teams ([#612](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/612))

## 3.0.3 (Oct 24, 2023)

IMPROVEMENTS:

* Add delays to API calls retries lacking of them on various TF Objects ([#758](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/758))
  - `data/source_pagerduty_event_orchestration`
  - `data/source_pagerduty_event_orchestrations`
  - `resource/pagerduty_automation_actions_action`
  - `resource/pagerduty_automation_actions_action_service_association`
  - `resource/pagerduty_automation_actions_action_team_association`
  - `resource/pagerduty_automation_actions_runner`
  - `resource/pagerduty_automation_actions_runner_team_association`
  - `resource/pagerduty_business_service`
  - `resource/pagerduty_business_service_subscriber`
  - `resource/pagerduty_escalation_policy`
  - `resource/pagerduty_event_orchestration`
  - `resource/pagerduty_event_orchestration_integration`
  - `resource/pagerduty_event_orchestration_path_global`
  - `resource/pagerduty_event_orchestration_path_router`
  - `resource/pagerduty_event_orchestration_path_service`
  - `resource/pagerduty_event_orchestration_path_unrouted`
  - `resource/pagerduty_schedule`
  - `resource/pagerduty_service_dependency`
  - `resource/pagerduty_service_event_rule`
  - `resource/pagerduty_service_integration`
  - `resource/pagerduty_slack_connection`
  - `resource/pagerduty_tag`
  - `resource/pagerduty_tag_assignment`
  - `resource/pagerduty_team`
  - `resource/pagerduty_team_membership`
  - `resource/pagerduty_user`
  - `resource/pagerduty_webhook_subscription`

## 3.0.2 (Oct 6, 2023)

BUG FIXES:
* `resource/pagerduty_service`, `resource/pagerduty_escalation_policy`: Update name validation logic for technical services ([#752](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/752))

## 3.0.1 (Sep 20, 2023)

BUG FIXES:
* `resource/pagerduty_schedule`: Fix - Provider crashing on Terraform state snapshot validation ([#747](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/747))

## 3.0.0 (Sep 8, 2023)

BREAKING CHANGES:

* `pagerduty/pagerduty`: Support for PagerDuty Apps Oauth scoped token ([#708](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/708))
* The following data sources and resources has been deprecated ([#744](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/744))
  - `data/pagerduty_custom_field`
  - `data/pagerduty_custom_field_schema`
  - `resource/pagerduty_custom_field`
  - `resource/pagerduty_custom_field_option`
  - `resource/pagerduty_custom_field_schema`
  - `resource/pagerduty_custom_field_schema_field_configuration`
  - `resource/pagerduty_custom_field_schema_assignment`

NOTES:

* Provider configuration attribute `token` has become `optional`.

## 2.16.2 (Aug 30, 2023)

BUG FIXES:
* `resource/pagerduty_service_dependency`: Address service dependency drift for external removals ([#741](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/741))

## 2.16.1 (Aug 24, 2023)

IMPROVEMENTS:
* `pagerduty/pagerduty`: bump terraform-plugin-sdk/v2 to solve #732 ([#733](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/733))
* Docs: `README.md`: update README to ref SECURE logs level ([#734](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/734))

## 2.16.0 (Aug 21, 2023)

FEATURES:

* `pagerduty/pagerduty`: Add support for SECURE logging level ([#730](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/730))

IMPROVEMENTS:
* `resource/pagerduty_escalation_policy`, `resource/pagerduty_service`: Update name validation func to not accept white spaces at the end ([#731](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/731))

## 2.15.3 (Aug 14, 2023)

IMPROVEMENTS:
* `resource/pagerduty_service_dependency`: [TFPROVDEV-27] Avoid Concurrent calls for Service Dependencies creation ([#724](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/724))
* `resource/pagerduty_user`: [TFPROVDEV-30] Address update user role failed as license is not re-computed correctly ([#725](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/725))

## 2.15.2 (Jul 21, 2023)

IMPROVEMENTS:
* `resource/pagerduty_schedule`: Stop retrying on Schedule deletion when open incidents are untraceable ([#714](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/714))
* Address: Too long and unneeded timeouts for call retries with 400 http errors ([#713](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/713))
    - `data/pagerduty_automation_actions_action`
    - `data/pagerduty_automation_actions_runner`
    - `data/pagerduty_business_service`
    - `data/pagerduty_escalation_policy`
    - `data/pagerduty_event_orchestration`
    - `data/pagerduty_event_orchestration_integration`
    - `data/pagerduty_event_orchestrations`
    - `data/pagerduty_extension_schema`
    - `data/pagerduty_incident_custom_field`
    - `data/pagerduty_incident_workflow`
    - `data/pagerduty_license`
    - `data/pagerduty_licenses`
    - `data/pagerduty_priority`
    - `data/pagerduty_ruleset`
    - `data/pagerduty_schedule`
    - `data/pagerduty_service`
    - `data/pagerduty_service_integration`
    - `data/pagerduty_tag`
    - `data/pagerduty_team`
    - `data/pagerduty_user`
    - `data/pagerduty_user_contact_method`
    - `data/pagerduty_users`
    - `data/pagerduty_vendor`
    - `resource/pagerduty_addon`
    - `resource/pagerduty_automation_actions_action`
    - `resource/pagerduty_automation_actions_action_service_association`
    - `resource/pagerduty_automation_actions_action_team_association`
    - `resource/pagerduty_automation_actions_runner`
    - `resource/pagerduty_automation_actions_runner_team_association`
    - `resource/pagerduty_business_service`
    - `resource/pagerduty_business_service_subscriber`
    - `resource/pagerduty_escalation_policy`
    - `resource/pagerduty_event_orchestration`
    - `resource/pagerduty_event_orchestration_integration`
    - `resource/pagerduty_event_orchestration_path_global`
    - `resource/pagerduty_event_orchestration_path_router`
    - `resource/pagerduty_event_orchestration_path_service`
    - `resource/pagerduty_event_orchestration_path_unrouted`
    - `resource/pagerduty_event_rule`
    - `resource/pagerduty_extension`
    - `resource/pagerduty_extension_servicenow`
    - `resource/pagerduty_incident_custom_field`
    - `resource/pagerduty_incident_custom_field_option`
    - `resource/pagerduty_incident_workflow`
    - `resource/pagerduty_incident_workflow_trigger`
    - `resource/pagerduty_maintenance_window`
    - `resource/pagerduty_response_play`
    - `resource/pagerduty_ruleset`
    - `resource/pagerduty_ruleset_rule`
    - `resource/pagerduty_schedule`
    - `resource/pagerduty_service`
    - `resource/pagerduty_service_dependency`
    - `resource/pagerduty_service_event_rule`
    - `resource/pagerduty_service_integration`
    - `resource/pagerduty_slack_connection`
    - `resource/pagerduty_tag`
    - `resource/pagerduty_tag_assignment`
    - `resource/pagerduty_team`
    - `resource/pagerduty_team_membership`
    - `resource/pagerduty_user`
    - `resource/pagerduty_user_contact_method`
    - `resource/pagerduty_user_notification_rule`
    - `resource/pagerduty_webhook_subscription`

## 2.15.1 (Jul 12, 2023)

IMPROVEMENTS:
* `resource/pagerduty_escalation_policy`, `resource/pagerduty_service`: Address name format validation on Escalation Policies ([#712](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/712))
* `dependency/google.golang.org/grpc`: bump google.golang.org/grpc from 1.33.2 to 1.53.0 ([#711](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/711))
* `resource/pagerduty_slack_connection`: fix slack_connection doc ([#587](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/587))
* Custom Fields - remove early access marker from incident custom fields pages ([#701](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/701))
    - `data/pagerduty_incident_custom_field`
    - `resource/pagerduty_incident_custom_field`
    - `resource/pagerduty_incident_custom_field_option`


## 2.15.0 (May 30, 2023)

BREAKING CHANGES:

* The following data sources and resources has been deprecated ([#684](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/684))
  - `data/pagerduty_custom_field`
  - `data/pagerduty_custom_field_schema`
  - `resource/pagerduty_custom_field`
  - `resource/pagerduty_custom_field_option`
  - `resource/pagerduty_custom_field_schema`
  - `resource/pagerduty_custom_field_schema_field_configuration`
  - `resource/pagerduty_custom_field_schema_assignment`

NOTES:

* The following data sources and resources shall be removed in the next **Major** release
  - `data/pagerduty_custom_field`
  - `data/pagerduty_custom_field_schema`
  - `resource/pagerduty_custom_field`
  - `resource/pagerduty_custom_field_option`
  - `resource/pagerduty_custom_field_schema`
  - `resource/pagerduty_custom_field_schema_field_configuration`
  - `resource/pagerduty_custom_field_schema_assignment`

FEATURES:

* adapt Terraform provider to use reflect simplified Custom Fields API. The following data sources and resources were added ([#684](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/684))
  - `data/pagerduty_incident_custom_field`
  - `resource/pagerduty_incident_custom_field`
  - `resource/pagerduty_incident_custom_field_option`

## 2.14.6 (May 29, 2023)

IMPROVEMENTS:
*  `resource/schedule`: Update handling of format errors on `pagerduty_schedule.start` ([#691](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/691))
*  `chore`: Update go module directive go version to `1.17` ([#694](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/694))

BUG FIXES:
*  `resource/schedule`: Address Schedule can't be deleted when used by EP with one layer configured ([#693](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/693))

## 2.14.5 (May 15, 2023)

IMPROVEMENTS:
* Plan recreation of tag assignments and teams on external to Terraform deletion ([#686](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/686))
*  `resource/pagerduty_tag_assignment`
*  `resource/pagerduty_team`

## 2.14.4 (May 2, 2023)

IMPROVEMENTS:
*  `resource/schedule`: Improve `resource/pagerduty_schedule` open incidents handling on deletion ([#681](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/681))

## 2.14.3 (April 20, 2023)

IMPROVEMENTS:
*  `data/vendor`: Doc update for `data.pagerduty_vendor` regarding PagerDuty AIOps feature gate ([#678](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/678))
*  `resource/schedule`: Add schedule's users as query param when listing open incidents associated to EP snashot ([#679](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/679))
*  `resource/service`: Service response play no op update ([#680](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/680))

## 2.14.2 (April 17, 2023)

IMPROVEMENTS:
*  `resource/pagerduty_escalation_policy`: Handle retries and state drift clean up for Escalation Policy ([#677](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/677))

## 2.14.1 (April 14, 2023)

IMPROVEMENTS:
* Support for deleting remote configuration of Event Orchestration Paths ([#676](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/676))
  - `resource/pagerduty_event_orchestration_service`
  - `resource/pagerduty_event_orchestration_router`
  - `resource/pagerduty_event_orchestration_unrouted`
  - `resource/pagerduty_event_orchestration_global`

## 2.14.0 (April 13, 2023)

FEATURES:
* Adds license resource for user management ([#657](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/657))
  - `data/pagerduty_license`
  - `data/pagerduty_licenses`
  - `resource/pagerduty_user`

## 2.13.0 (April 12, 2023)

FEATURES:
* Event Orchestration Updates: Orchestration Warnings, Global Orchestrations, Orchestration Integrations ([#618](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/618))
  - `data/pagerduty_event_orchestration_integration`
  - `resource/pagerduty_event_orchestration_global`
  - `resource/pagerduty_event_orchestration_integration`

## 2.12.2 (April 11, 2023)

IMPROVEMENTS:
* [ORCA-3999] Add EOL banner to Ruleset and Service Rules. ([#672](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/672))
  - `data/pagerduty_ruleset`
  - `resource/pagerduty_ruleset`
  - `resource/pagerduty_ruleset_rule`
  - `resource/pagerduty_service_event_rule`

## 2.12.1 (April 6, 2023)

BUG FIXES:
*  `resource/pagerduty_service`: Address: Service Read lifecycle wasn't detecting drift for auto pause notif and alert grouping params ([#673](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/673))

## 2.12.0 (April 6, 2023)

FEATURES:
* `data/pagerduty_service`: service: Compute additional fields already included in API response ([#660](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/660))

IMPROVEMENTS:
* `resource/pagerduty_schedule`: Handle a schedule being deleted in the UI ([#661](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/661))
* `resource/pagerduty_tag`: Address #655 Tags are not cleaned up from State after removed externally ([#670](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/670))
* `resource/pagerduty_tag`: Address #655 Tags are not cleaned up from State after removed externally ([#670](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/670))
* `resource/pagerduty_user_contact_method`: Print number and error for failing contact method validation ([#568](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/568))

## 2.11.3 (April 5, 2023)

IMPROVEMENTS:
*  `resource/pagerduty_team_memberhsip`: Upgrade `go-pagerduty` to support caching for `pagerduty_team_membership` ([#666](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/666))

## 2.11.2 (March 10, 2023)

IMPROVEMENTS:
*  `resource/pagerduty_custom_field_schema_field_configuration`: Addressing name typo on custom fields schema field config docs and test ([#654](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/654))
*  `resource/pagerduty_service_dependency`: Addressing GOAWAY error on `pagerduty_service_dependency` ([#653](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/653))
* Add support for `process_automation_node_filter` on Automation Actions ([#647](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/647))
  - `data_source/pagerduty_automation_actions_action`
  - `resource/pagerduty_automation_actions_action`

## 2.11.1 (March 9, 2023)

IMPROVEMENTS:
*  `resource/pagerduty_event_orchestration_service`: Enable event orchestration active status for service ([#649](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/649))
* Remove early access header for incident workflows. ([#645](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/645))
  - `data_source/pagerduty_incident_workflow`
  - `resource/pagerduty_incident_workflow`
  - `resource/pagerduty_incident_workflow_trigger`

## 2.11.0 (February 15, 2023)

FEATURES:

* Support **Custom Fields** via several new resources. ([#623](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/623))
  - `resource/pagerduty_custom_field`
  - `resource/pagerduty_custom_field_option`
  - `resource/pagerduty_custom_field_schema`
  - `resource/pagerduty_custom_field_schema_assignment`
  - `resource/pagerduty_custom_field_schema_field_configuration`
  - `data_source/pagerduty_custom_field`
  - `data_source/pagerduty_custom_field_schema`

## 2.10.2 (February 9, 2023)

BUG FIXES:
*  `resource/pagerduty_schedule`: Cannot destroy pagerduty_schedule. You must first resolve the following incidents related with Escalation Policies ([#619](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/619))

## 2.10.1 (February 9, 2023)

BUG FIXES:
*  `resource/pagerduty_service_integration`: Address: Service integration perm diff with Generic email and empty/omitted `email_filter` ([#625](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/625))

## 2.10.0 (February 9, 2023)

FEATURES:

* `data/pagerduty_event_orchestrations`: feat: add `pagerduty_event_orchestrations` datasource ([#581](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/581))

## 2.9.3 (January 26, 2023)

BUG FIXES:
*  `resource/pagerduty_team_membership`: EF-3964 Address `team_membership` inconsistency after `create` and `update` ([#621](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/621))

## 2.9.2 (January 17, 2023)

BUG FIXES:
*  `dependency/go-pagerduty`: update `github.com/heimweh/go-pagerduty` to fix offset pagination ([#615](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/615))

## 2.9.1 (January 9, 2023)

IMPROVEMENTS:

* `resource/pagerduty_incident_workflow`: add team support for incident workflows ([#609](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/609))
* replace `ValidateValueFunc` with `ValidateValueDiagFunc` due to deprecation on the following resources. ([#605](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/605))
	- `resource/pagerduty_event_orchestration`
	- `resource/pagerduty_automation_actions_action`
	- `resource/pagerduty_automation_actions_runner`
	- `resource/pagerduty_business_service`
	- `resource/pagerduty_business_service_subscriber`
	- `resource/pagerduty_escalation_policy`
	- `resource/pagerduty_event_orchestration_service`
	- `resource/pagerduty_event_orchestration_unrouted`
	- `resource/pagerduty_incident_workflow_trigger`
	- `resource/pagerduty_ruleset_rule`
	- `resource/pagerduty_schedule`
	- `resource/pagerduty_service`
	- `resource/pagerduty_service_dependency`
	- `resource/pagerduty_service_event_rule`
	- `resource/pagerduty_service_integration`
	- `resource/pagerduty_slack_connection`
	- `resource/pagerduty_tag_assignment`
	- `resource/pagerduty_team_membership`
	- `resource/pagerduty_user`
	- `resource/pagerduty_user_contact_method`
	- `resource/pagerduty_user_notification_rule`
	- `resource/pagerduty_webhook_subscription`

## 2.9.0 (January 9, 2023)

FEATURES:

* `data/pagerduty_automation_actions_action`: Add support for data.pagerduty_automation_actions_action ([#601](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/601))
* `resource/pagerduty_automation_actions_runner_team_association`: Add support for Automation Actions' Runner association with a Team ([#607](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/607))
* `resource/pagerduty_automation_actions_action_service_association`: Add support for Automation Actions' Action association to a Service ([#608](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/608))


IMPROVEMENTS:
* `resource/pagerduty_automation_actions_action`, `resource/pagerduty_automation_actions_runner`: Add support for the update operation on pagerduty_automation_actions_action and pagerduty_automation_actions_runner ([#603](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/603))

## 2.8.1 (December 23, 2022)

IMPROVEMENTS:
* `resource/pagerduty_schedule`: Address: Schedules can't be deleted when they have open incidents ([#602](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/602))

## 2.8.0 (December 23, 2022)

FEATURES:
* Support **Incident Workflows** via several new resources. ([#596](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/596))
  - `resource/pagerduty_incident_workflow_trigger`
  - `resource/pagerduty_incident_workflow`
  - `data_source/pagerduty_incident_workflow`
* `resource/pagerduty_automation_actions_action_team_association`: Add support for Automation Actions' Action association to a Team ([#600](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/600))
* `resource/pagerduty_automation_actions_action`: Add support for pagerduty_automation_actions_action ([#599](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/599))

## 2.7.0 (December 12, 2022)

FEATURES:
* `resource/pagerduty_automation_actions_runner`, `data_source/pagerduty_automation_actions_runner`: Add support for `pagerduty_automation_actions_runner` and `data.pagerduty_automation_actions_runner` ([#595](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/595))

## 2.6.5 (December 12, 2022)

BUG FIXES:
* `resource/pagerduty_user_contact_method`: Address: Unique contact method error not being captured ([#586](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/586))

## 2.6.4 (November 3, 2022)

BUG FIXES:
* `resource/pagerduty_service`: Test and handle time-based alert grouping parameters ([#582](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/582))

## 2.6.3 (October 11, 2022)

BUG FIXES:
* `resource/pagerduty_service`: Address: `pagerduty_service.alert_grouping_parameters.config` block parsing ([#570](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/570))
* `resource/pagerduty_service`: resource_pagerduty_service: skip response_play with "null" value ([#573](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/573))

IMPROVEMENTS:
* Docs: `resource/tag_assignment`: Adds proper documentation for pagerduty_tag resource in tag_assignmen... ([#563](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/563))
* Docs: `README.md`: Updates developer docs with more helpful setup information ([#571](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/571))

## 2.6.2 (September 8, 2022)

BUG FIXES:
* `resource/pagerduty_schedule`: Test rule removal rather than update for Escalation Policy Dependent Schedule ([#564](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/564))

IMPROVEMENTS:
* Docs: `changelog.md`: Fixing PR Links in Changelog ([#566](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/566))

## 2.6.1 (August 25, 2022)

BUG FIXES:
* `resource/pagerduty_schedule`: Add support for gracefully destroy `pagerduty_schedule` ([#561](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/561))


## 2.6.0 (August 17, 2022)

FEATURES:
* `data_source/pagerduty_users`: Add `pagerduty_users` data source ([#545](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/545))
* `resource/pagerduty_service`: Add support for service level `auto_pause_notifications_parameters` ([#525](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/525))

IMPROVEMENTS:
* `resource/pagerduty_service`: Add `response_play` field to service ([#515](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/515)))
* Docs: `resource/pagerduty_service`: Fixed Docs Bugs in pagerduty_service Ref: Issue #522 ([#554](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/554))

BUG FIXES:
* `resource/pagerduty_team_membership`: Add support for gracefully destroy `pagerduty_team_membership` ([#558](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/558))


## 2.5.2 (July 12, 2022)
IMPROVEMENTS:
* `goreleaser`: update to go-version 1.17 ([#543](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/543))
* `resource/pagerduty_schedule`: Addressing output not showing rendered_coverage_percentage ([#528](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/528))
* Docs: `resource/pagerduty_event_orchestration_router`: Fix typos in the event orchestration router docs ([#536](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/536))
* Docs: `resource/pagerduty_event_orchestration_router`, `resource/pagerduty_event_orchestration_service`, `resource/pagerduty_event_orchestration_unrouted`: Fix docs for event_orchestration resources import ([#529](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/529))
* Docs: `resource/pagerduty_extension_servicenow`: Wrong Extension Schema ([#487](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/487))

BUG FIXES:
* `resource/pagerduty_service`: remove `expectNonEmptyPlanFromTest` from Service test responding to feedback left in PR#527 ([#542](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/542))
* `resource/pagerduty_service_dependency`: add input validation and drift detection during deletion for `service_dependency` ([#530](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/530))
* `data/pagerduty_business_services`, `data/extension_schema`, `data/priority`, `data/service`, `data/integration`, `data/tag`, `data/team`: Remove 429 check on remaining data sources ([#537](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/537))
* `resource/pagerduty_service`: Address unable to switch off alert grouping on a service ([#527](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/527))


## 2.5.1 (June 9, 2022)
FEATURES:
* `resource/pagerduty_service`: Address: Unable to switch off alert grouping on a service ([#455](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/455))

BUG FIXES:
* `resource/pagerduty_slack_connection`: Addressing pagerduty_slack_connection unable to set "No Priority" vs "Any Priority"([#519](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/519))
* `data_source/pagerduty_user`,`data_source/pagerduty_schedule`,`data_source/pagerduty_vendor`: Changed logic to retry on all errors returned by PD API. Remedies GOAWAY error. ([#521](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/521))


## 2.5.0 (June 1, 2022)

FEATURES:
* Support for Event Orchestration via several new resources. ([#512](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/512))
    * `resource/pagerduty_event_orchestration`
    * `resource/pagerduty_event_orchestration_router`
    * `resource/pagerduty_event_orchestration_unrouted`
    * `resource/pagerduty_event_orchestration_service`
    * `data_source/pagerduty_event_orchestration`

IMPROVEMENTS:
* `data_source/pagerduty_user`: Add support for pagination. Gets all users. ([#511](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/511))

BUG FIXES: 
* `resource/pagerduty_service_integration`: Fix permadiff in email_parser with type regex & minor docs update ([#479](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/479))

## 2.4.2 (May 20, 2022)

IMPROVEMENTS:
* Acceptance Test Improvements: Use "@foo.test" email addresses in tests. ([#491](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/491))
* Acceptance Test Improvements: Adding better notes to README on running ACC ([#503](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/503)) 
* `resource/pagerduty_ruleset_rule`: Introduce support for `catch_all` rules. ([#481](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/481))
* Docs: `resource/pagerduty_slack_connection`: Improved notes on resource supporting Slack V2 Next Generation ([#496](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/496))

BUG FIXES:
* Documentation: Fixed all broken links to the PagerDuty API documentation ([#464](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/464))
* Docs: `resource/pagerduty_escalation_policy`: Fixed `user` -> `user_reference` in samples ([#497](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/497))
* Build Process: Include `timezdata` build tag in goreleaser config ([#488](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/488))
* `data_source/pagerduty_escalation_policy`,`data_source/pagerduty_ruleset`: Changed logic to retry on all errors returned by PD API. Remedies GOAWAY error. ([#507](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/507))

## 2.4.1 (April 22, 2022)
IMPROVEMENTS:
* `resource/pagerduty_user_notification`: Create user/notification rule: allow using existing ones ([#482](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/482))
* `resource/pagerduty_schedule`: Enforce 0 second time restriction on schedules ([#483](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/483))

BUG FIXES:
* Embed time zone data into the provider's binary([#478](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/478))
* Documentation fix: update all broken links pointing to PagerDuty API documentation ([#464](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/464))


## 2.4.0 (April 1, 2022)
FEATURES:
* `resource/pagerduty_service_integration`: Add Email Filters to Service Integrations ([#468](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/468))

IMPROVEMENTS:
* `resource/ruleset_rule`,`resource/schedule`,`resource/service`,`resource/service_event_rule`,`resource/user`: Validate time zones ([#473](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/473))
* `resource/user_contact_method`: Validate phone numbers starting with 0 are not supported ([#475](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/475))
* `resource/pagerduty_schedule`: Send `"end": null` when layer end removed([#460](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/460))
* `resource/maintenance_window`: Ignore error code 405 on delete ([#466](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/466))


## 2.3.0 (February 10, 2022)
IMPROVEMENTS:
* Updated TF SDK to v2.10.1 and added `depends_on` to eventrule tests ([#446](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/446))
* `resource/pagerduty_schedule`: Added validation to `duration_seconds` ([#433](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/433))
* Documentation fix: update code sample on index ([#436](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/436)
* `resource/pagerduty_escalation_policy`: Validate user and schedule reference in escalation policy targets ([#435](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/435))
* `resource/pagerduty_service`,`resource/pagerduty_business_service`,`data_source/pagerduty_service`,`data_source/pagerduty_business_service`: Adding computed `type` field to be used in `service_dependencies` ([#364](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/364))
* `resource/pagerduty_service_integration`: Support emails that are only known after apply ([#425](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/425))
* Safter HTTP client initialization and usage ([#458](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/458))
* Increase Retry Time on Data Sources ([#454](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/454))

BUG FIXES:
* Documentation fix: update broken links to auth docs in PagerDuty ([#449](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/449))
* Documentation fix: update description of PagerDuty ([#441](https://github.com/PagerDuty/terraform-provider-pagerduty/pull/441))

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

* data-source/pagerduty_team: Fix team search issue ([#110](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/110))
* resource/pagerduty_maintenance_window: Suppress spurious diff in `start_time` & `end_time` ([#116](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/116))
* resource/pagerduty_service: Set invitation_sent ([#127](https://github.com/PagerDuty/terraform-provider-pagerduty/issues/127))
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
