---
page_title: "PagerDuty API Support Status"
---

# PagerDuty API Support Status

This document tracks the implementation status of PagerDuty operations REST API across the Resources and Data Sources of the Terraform Provider.
It serves as both a quick reference for available functionality and a development roadmap.
The status is updated as new features are implemented and tested.
Features marked with a “prohibited 🚫” emoji are not relevant for this provider.

- Abilities
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/abilities</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/abilities/{id}</span>

- Add-ons
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/addons</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/addons</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/addons/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/addons/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/addons/{id}</span>

- Alert Grouping Settings
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/alert\_grouping\_settings</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/alert\_grouping\_settings</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/alert\_grouping\_settings/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/alert\_grouping\_settings/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/alert\_grouping\_settings/{id}</span>

- Analytics
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/metrics/incidents/all</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/metrics/incidents/escalation\_policies</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/metrics/incidents/escalation\_policies/all</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/metrics/incidents/services</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/metrics/incidents/services/all</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/metrics/incidents/teams</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/metrics/incidents/teams/all</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/metrics/pd\_advance\_usage/features</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/metrics/responders/all</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/metrics/responders/teams</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/raw/incidents</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/analytics/raw/incidents/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/analytics/raw/incidents/{id}/responses</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/analytics/raw/responders/{responder\_id}/incidents</span>

- Audit
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/audit/records</span>

- Automation Actions
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/actions</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/automation\_actions/actions</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/actions/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/automation\_actions/actions/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/automation\_actions/actions/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/automation\_actions/actions/{id}/invocations</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/actions/{id}/services</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/automation\_actions/actions/{id}/services</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/actions/{id}/services/{service\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/automation\_actions/actions/{id}/services/{service\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/actions/{id}/teams</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/automation\_actions/actions/{id}/teams</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/actions/{id}/teams/{team\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/automation\_actions/actions/{id}/teams/{team\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/invocations</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/invocations/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/runners</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/automation\_actions/runners</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/runners/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/automation\_actions/runners/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/automation\_actions/runners/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/runners/{id}/teams</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/automation\_actions/runners/{id}/teams</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/automation\_actions/runners/{id}/teams/{team\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/automation\_actions/runners/{id}/teams/{team\_id}</span>

- Business Services
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/business\_services</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/business\_services</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/business\_services/impactors</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/business\_services/impacts</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/business\_services/priority\_thresholds</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/business\_services/priority\_thresholds</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/business\_services/priority\_thresholds</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/business\_services/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/business\_services/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/business\_services/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/business\_services/{id}/account\_subscription</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/business\_services/{id}/account\_subscription</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/business\_services/{id}/subscribers</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/business\_services/{id}/subscribers</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/business\_services/{id}/supporting\_services/impacts</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/business\_services/{id}/unsubscribe</span>

- Change Events
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/change\_events</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/change\_events</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/change\_events/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/change\_events/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/related\_change\_events</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/services/{id}/change\_events</span>

- Escalation Policies
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/escalation\_policies</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/escalation\_policies</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/escalation\_policies/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/escalation\_policies/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/escalation\_policies/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/escalation\_policies/{id}/audit/records</span>

- Event Orchestrations
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/event\_orchestrations</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/services/{service\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/event\_orchestrations/services/{service\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/services/{service\_id}/active</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/event\_orchestrations/services/{service\_id}/active</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/services/{service\_id}/cache\_variables</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/event\_orchestrations/services/{service\_id}/cache\_variables</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/services/{service\_id}/cache\_variables/{cache\_variable\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/event\_orchestrations/services/{service\_id}/cache\_variables/{cache\_variable\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/event\_orchestrations/services/{service\_id}/cache\_variables/{cache\_variable\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/event\_orchestrations/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/event\_orchestrations/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/{id}/cache\_variables</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/event\_orchestrations/{id}/cache\_variables</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/{id}/cache\_variables/{cache\_variable\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/event\_orchestrations/{id}/cache\_variables/{cache\_variable\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/event\_orchestrations/{id}/cache\_variables/{cache\_variable\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/{id}/global</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/event\_orchestrations/{id}/global</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/{id}/integrations</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/event\_orchestrations/{id}/integrations</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/event\_orchestrations/{id}/integrations/migration</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/{id}/integrations/{integration\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/event\_orchestrations/{id}/integrations/{integration\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/event\_orchestrations/{id}/integrations/{integration\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/{id}/router</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/event\_orchestrations/{id}/router</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/event\_orchestrations/{id}/unrouted</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/event\_orchestrations/{id}/unrouted</span>

- Extension Schemas
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/extension\_schemas</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/extension\_schemas/{id}</span>

- Extensions
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/extensions</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/extensions</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/extensions/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/extensions/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/extensions/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/extensions/{id}/enable</span>

- Incident Custom Fields
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/custom\_fields</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/custom\_fields</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/custom\_fields/{field\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/custom\_fields/{field\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/incidents/custom\_fields/{field\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/custom\_fields/{field\_id}/field\_options</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/custom\_fields/{field\_id}/field\_options</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/custom\_fields/{field\_id}/field\_options/{field\_option\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/incidents/custom\_fields/{field\_id}/field\_options/{field\_option\_id}</span>

- Incident Types
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/types</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/types</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/types/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/types/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/types/{id}/custom\_fields</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/types/{id}/custom\_fields</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/types/{id}/custom\_fields/{field\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/types/{id}/custom\_fields/{field\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/incidents/types/{id}/custom\_fields/{field\_id}</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/types/{id}/custom\_fields/{field\_id}/field\_options</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/types/{id}/custom\_fields/{field\_id}/field\_options</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/types/{id}/custom\_fields/{field\_id}/field\_options/{field\_option\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/types/{id}/custom\_fields/{field\_id}/field\_options/{field\_option\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/incidents/types/{id}/custom\_fields/{field\_id}/field\_options/{field\_option\_id}</span>

- Incident Workflows
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incident\_workflows</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incident\_workflows</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incident\_workflows/actions</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incident\_workflows/actions/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incident\_workflows/triggers</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incident\_workflows/triggers</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incident\_workflows/triggers/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incident\_workflows/triggers/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/incident\_workflows/triggers/{id}</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incident\_workflows/triggers/{id}/services</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/incident\_workflows/triggers/{trigger\_id}/services/{service\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incident\_workflows/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incident\_workflows/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/incident\_workflows/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incident\_workflows/{id}/instances</span>

- Incidents
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/alerts</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/{id}/alerts</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/alerts/{alert\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/{id}/alerts/{alert\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/business\_services/impacts</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/{id}/business\_services/{business\_service\_id}/impacts</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/custom\_fields/values</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/{id}/custom\_fields/values</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/log\_entries</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/{id}/merge</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/notes</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/{id}/notes</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/outlier\_incident</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/past\_incidents</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/related\_incidents</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/{id}/responder\_requests</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/{id}/snooze</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/{id}/status\_updates</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{id}/status\_updates/subscribers</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/{id}/status\_updates/subscribers</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/{id}/status\_updates/unsubscribe</span>

- Licenses
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/license\_allocations</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/licenses</span>

- Log Entries
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/log\_entries</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/log\_entries/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/log\_entries/{id}/channel</span>

- Maintenance Windows
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/maintenance\_windows</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/maintenance\_windows</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/maintenance\_windows/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/maintenance\_windows/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/maintenance\_windows/{id}</span>

- Notifications
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/notifications</span>

- OAuth Delegations
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/oauth\_delegations</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/oauth\_delegations/revocation\_requests/status</span>

- On-Calls
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/oncalls</span>

- Paused Incident Reports
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/paused\_incident\_reports/alerts</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/paused\_incident\_reports/counts</span>

- Priorities
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/priorities</span>

- Response Plays
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/response\_plays</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/response\_plays</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/response\_plays/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/response\_plays/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/response\_plays/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/response\_plays/{response\_play\_id}/run</span>

- Rulesets
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/rulesets</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/rulesets</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/rulesets/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/rulesets/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/rulesets/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/rulesets/{id}/rules</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/rulesets/{id}/rules</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/rulesets/{id}/rules/{rule\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/rulesets/{id}/rules/{rule\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/rulesets/{id}/rules/{rule\_id}</span>

- Schedules
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/schedules</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/schedules</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/schedules/preview</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/schedules/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/schedules/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/schedules/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/schedules/{id}/audit/records</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/schedules/{id}/overrides</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/schedules/{id}/overrides</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/schedules/{id}/overrides/{override\_id}</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/schedules/{id}/users</span>

- Service Dependencies
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/service\_dependencies/associate</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/service\_dependencies/business\_services/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/service\_dependencies/disassociate</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/service\_dependencies/technical\_services/{id}</span>

- Services
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/services</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/services</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/services/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/services/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/services/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/services/{id}/audit/records</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/services/{id}/integrations</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/services/{id}/integrations/{integration\_id}</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/services/{id}/integrations/{integration\_id}</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/services/{id}/rules</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/services/{id}/rules</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/services/{id}/rules/convert</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/services/{id}/rules/{rule\_id}</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/services/{id}/rules/{rule\_id}</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/services/{id}/rules/{rule\_id}</span>

- Standards
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/standards</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/standards/scores/{resource\_type}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/standards/scores/{resource\_type}/{id}</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/standards/{id}</span>

- Status Dashboards
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_dashboards</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_dashboards/url\_slugs/{url\_slug}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_dashboards/url\_slugs/{url\_slug}/service\_impacts</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_dashboards/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_dashboards/{id}/service\_impacts</span>

- Status Pages
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/impacts</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/impacts/{impact\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/posts</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/status\_pages/{id}/posts</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}/post\_updates</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}/post\_updates</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}/post\_updates/{post\_update\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}/post\_updates/{post\_update\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}/post\_updates/{post\_update\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}/postmortem</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}/postmortem</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}/postmortem</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/status\_pages/{id}/posts/{post\_id}/postmortem</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/services</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/services/{service\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/severities</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/severities/{severity\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/statuses</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/statuses/{status\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/subscriptions</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/status\_pages/{id}/subscriptions</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/status\_pages/{id}/subscriptions/{subscription\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/status\_pages/{id}/subscriptions/{subscription\_id}</span>

- Tags
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/tags</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/tags</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/tags/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/tags/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/tags/{id}/{entity\_type}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/{entity\_type}/{id}/change\_tags</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/{entity\_type}/{id}/tags</span>

- Teams
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/teams</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/teams</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/teams/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/teams/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/teams/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/teams/{id}/audit/records</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/teams/{id}/escalation\_policies/{escalation\_policy\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/teams/{id}/escalation\_policies/{escalation\_policy\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/teams/{id}/members</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/teams/{id}/notification\_subscriptions</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/teams/{id}/notification\_subscriptions</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/teams/{id}/notification\_subscriptions/unsubscribe</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/teams/{id}/users/{user\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/teams/{id}/users/{user\_id}</span>

- Templates
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/templates</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/templates</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/templates/fields</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/templates/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/templates/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/templates/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/templates/{id}/render</span>

- Users
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/users</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/me</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/users/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/users/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/audit/records</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/contact\_methods</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/users/{id}/contact\_methods</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/contact\_methods/{contact\_method\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/users/{id}/contact\_methods/{contact\_method\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/users/{id}/contact\_methods/{contact\_method\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/license</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/notification\_rules</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/users/{id}/notification\_rules</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/notification\_rules/{notification\_rule\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/users/{id}/notification\_rules/{notification\_rule\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/users/{id}/notification\_rules/{notification\_rule\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/notification\_subscriptions</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/users/{id}/notification\_subscriptions</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/users/{id}/notification\_subscriptions/unsubscribe</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/oncall\_handoff\_notification\_rules</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/users/{id}/oncall\_handoff\_notification\_rules</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/oncall\_handoff\_notification\_rules/{oncall\_handoff\_notification\_rule\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/users/{id}/oncall\_handoff\_notification\_rules/{oncall\_handoff\_notification\_rule\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/users/{id}/oncall\_handoff\_notification\_rules/{oncall\_handoff\_notification\_rule\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/sessions</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/users/{id}/sessions</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/sessions/{type}/{session\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/users/{id}/sessions/{type}/{session\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/status\_update\_notification\_rules</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/users/{id}/status\_update\_notification\_rules</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/users/{id}/status\_update\_notification\_rules/{status\_update\_notification\_rule\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/users/{id}/status\_update\_notification\_rules/{status\_update\_notification\_rule\_id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/users/{id}/status\_update\_notification\_rules/{status\_update\_notification\_rule\_id}</span>

- Vendors
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/vendors</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/vendors/{id}</span>

- Webhooks
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/webhook\_subscriptions</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/webhook\_subscriptions</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/webhook\_subscriptions/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/webhook\_subscriptions/{id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/webhook\_subscriptions/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/webhook\_subscriptions/{id}/enable</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/webhook\_subscriptions/{id}/ping</span>

- Workflow Integrations
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/workflows/integrations</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/workflows/integrations/connections</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/workflows/integrations/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/workflows/integrations/{integration\_id}/connections</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/workflows/integrations/{integration\_id}/connections</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/workflows/integrations/{integration\_id}/connections/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/workflows/integrations/{integration\_id}/connections/{id}</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PATCH</span> <span style="font-size:.9em">/workflows/integrations/{integration\_id}/connections/{id}</span>

- Jira Cloud Account Mappings
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/integration-jira-cloud/accounts\_mappings</span>
    - [ ] 🚫 <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/integration-jira-cloud/accounts\_mappings/{id}</span>

- Jira Cloud Rules
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/integration-jira-cloud/accounts\_mappings/{id}/rules</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/integration-jira-cloud/accounts\_mappings/{id}/rules</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/integration-jira-cloud/accounts\_mappings/{id}/rules/{rule\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/integration-jira-cloud/accounts\_mappings/{id}/rules/{rule\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/integration-jira-cloud/accounts\_mappings/{id}/rules/{rule\_id}</span>

- Slack Connections
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/workspaces/{slack\_team\_id}/connections</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/workspaces/{slack\_team\_id}/connections</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/workspaces/{slack\_team\_id}/connections/{connection\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/workspaces/{slack\_team\_id}/connections/{connection\_id}</span>
    - [x] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/workspaces/{slack\_team\_id}/connections/{connection\_id}</span>

- Slack Dedicated Channels
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{incident\_id}/dedicated\_channel</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/{incident\_id}/dedicated\_channel</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/incidents/{incident\_id}/dedicated\_channel</span>

- Slack Notification Channels
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#05b870;}">GET</span> <span style="font-size:.9em">/incidents/{incident\_id}/notification\_channels</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#19abff;}">POST</span> <span style="font-size:.9em">/incidents/{incident\_id}/notification\_channels</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f46d2a;}">PUT</span> <span style="font-size:.9em">/incidents/{incident\_id}/notification\_channels</span>
    - [ ] <span style="position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:#f05151;}">DELETE</span> <span style="font-size:.9em">/incidents/{incident\_id}/notification\_channels/{channel\_id}</span>
