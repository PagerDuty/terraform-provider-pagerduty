# Terraform Provider for PagerDuty

- Website: https://registry.terraform.io/providers/PagerDuty/pagerduty/latest
- Documentation: https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs
- Terraform Gitter: [![Terraform Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Terraform Google Groups](http://groups.google.com/group/terraform-tool)

[PagerDuty](https://www.pagerduty.com/) is an alarm aggregation and dispatching service for system administrators and support teams. It collects alerts from your monitoring tools, gives you an overall view of all of your monitoring alarms, and alerts an on duty engineer if there’s a problem. The Terraform Pagerduty provider is a plugin for Terraform that allows for the management of PagerDuty resources using HCL (HashiCorp Configuration Language).

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 1.1.0
-	[Go](https://golang.org/doc/install) 1.20 (to build the provider plugin)

## Building the Provider

Clone repository to: `$GOPATH/src/github.com/PagerDuty/terraform-provider-pagerduty`

```sh
$ mkdir -p $GOPATH/src/github.com/PagerDuty; cd $GOPATH/src/github.com/PagerDuty
$ git clone git@github.com:PagerDuty/terraform-provider-pagerduty
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/PagerDuty/terraform-provider-pagerduty
$ make build
...
$ $GOPATH/bin/terraform-provider-pagerduty
...
```

This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

## Usage

Please refer to Terraform docs for [PagerDuty Provider](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs)
for examples on how to use the provider and detailed documentation about the Resources and Data Sources the provider has.

## Caching Support

The `go-pagerduty` library relies on various APIs to interact with PagerDuty's resources. However, some of these APIs lack efficient ways to query specific resources by their attributes. When an implementation in the Terraform Provider requires such logic, the library lists all resources of a specific entity and performs a lookup in memory. This can result in inefficient use of the APIs, especially when dealing with a large number of resources, as the repetitive API calls for listing resource definitions can lead to significant time consumption and performance penalties. To address this issue, we have introduced caching to improve the user experience when interacting with Terraform resources that rely on API calls to list all available data for a specific entity, and then perform a lookup by attribute value in memory. With this improvement, the Terraform Provider users can expect better performance and faster response times when working with PagerDuty's resources.

### Resources currently supporting cache of API calls

* `pagerduty_team_membership`
* `pagerduty_user_contact_method`
* `pagerduty_user_notification_rule` 
* `pagerduty_user`

### Caching mechanisms available

* In memory.
* MongoDB.

### To activate caching support

| Environment Variable         | Example Value                                                                      | Description                                                                                                                                  |
|------------------------------|------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| `TF_PAGERDUTY_CACHE`         | memory                                                                             | Activate **In Memory** cache.                                                                                                                |
| `TF_PAGERDUTY_CACHE`         | `mongodb+srv://[mongouser]:[mongopass]@[mongodbname].[mongosubdomain].mongodb.net` | Activate MongoDB cache.                                                                                                                      |
| `TF_PAGERDUTY_CACHE_MAX_AGE` | 30s                                                                                | Only applicable for MongoDB cache. Time in seconds for cached data to become staled. Default value `10s`.                                    |
| `TF_PAGERDUTY_CACHE_PREFILL` | 1                                                                                  | Only applicable for MongoDB cache. Indicates to pre-fill data in cache for *Abilities*, *Users*, *Contact Methods* and *Notification Rules*. |


## Development

### Setup Local Environment

Before developing the provider, ensure that you have go correctly installed.

* Install [Go](http://www.golang.org) on your machine (version 1.11+ is *required*).
* Correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

### Setup Local PagerDuty Provider

Make changes to the PagerDuty provider and post a pull request for review.

1. [Create a fork](https://docs.github.com/en/get-started/quickstart/fork-a-repo) of the **upstream** repository `https://github.com/PagerDuty/terraform-provider-pagerduty`
2. Clone the new **origin** repository to your local go src path: `$GOPATH/src/github.com/<your-github-username>/terraform-provider-pagerduty`
3. optionally make development easier by setting the [**upstream**](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/configuring-a-remote-for-a-fork) repository
   ```
   $ git remote add upstream git@github.com:PagerDuty/terraform-provider-pagerduty.git
   ```
4. Make any changes on your local machine and post a PR to the **upstream** repository

### Run Dev Build with Local Terraform Module

> Note: Development overrides work only in Terraform v0.14 and later. Using a dev_overrides block in your CLI configuration will cause Terraform v0.13 to reject the configuration as invalid.

1. Build the provider with your latest changes. (See [Building the Provider](https://github.com/PagerDuty/terraform-provider-pagerduty#building-the-provider))
2. Override the pagerduty provider with your local build. (See [Development Overrides for Provider Developers](https://www.terraform.io/cli/config/config-file#development-overrides-for-provider-developers))
   * Create the file `$HOME/.terraformrc` and paste the following content into it. Be sure to change the path to wherever your binary is located. It is currently set to the default for go builds.
   ```terraform
   provider_installation {
      dev_overrides {
         "pagerduty/pagerduty" = "/<ABSOLUTE_PATH_TO>/<YOUR_HOME_PATH>/go/bin"
      }
      direct {}
   }
   ```
3. Goto a local terraform module and start running terraform. (See [Using the Provider](https://github.com/PagerDuty/terraform-provider-pagerduty#using-the-provider)). You may need to first install the latest module and provider
   versions allowed within the new configured constraints. Verify with the below warning message.
   ```sh
   $ terraform init -upgrade
   $ terraform plan
   ...
   │ Warning: Provider development overrides are in effect
   ```
4. See `api_url_override` from Terraform docs for [PagerDuty Provider](https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs#argument-reference) to set a custom proxy endpoint as PagerDuty client api url overriding service_region setup.

## Test a specific version of the go-pagerduty API client

Modify the `go.mod` file using a [Go module replacement](https://go.dev/doc/modules/managing-dependencies#external_fork) for `github.com/heimweh/go-pagerduty:
```
$ go mod edit -replace github.com/heimweh/go-pagerduty=/PATH/TO/LOCAL/github.com/<USERNAME>/<REPO>
```

Or update the file directly:
```
replace github.com/heimweh/go-pagerduty => /PATH/TO/LOCAL/go-pagerduty
```

Update vendored dependencies or configure compiler to prefer using downloaded modules based on `go.mod` file:
```
$ export GOFLAGS="-mod=mod"
```
Or:
```
$ go mod vendor
```

### Setup Local Logs

1. See [Debugging Terraform](https://www.terraform.io/internals/debugging). Either add this to your shell's profile
   (example: `~/.bashrc`), or just execute these commands:
   ```
   export TF_LOG=trace
   export TF_LOG_PATH="/PATH/TO/YOUR/LOG_FILE.log"
   ```
2. stream logs
   ```
   $ tail -f /PATH/TO/YOUR/LOG_FILE.log
   ```

### Secure Logs Level

In addition to the [log levels provided by Terraform](https://developer.hashicorp.com/terraform/internals/debugging), namely `TRACE`, `DEBUG`, `INFO`, `WARN`, and `ERROR` (in descending order of verbosity), the PagerDuty Provider introduces an extra level called `SECURE`. This level offers verbosity similar to Terraform's debug logging level, specifically for the output of API calls and HTTP request/response logs. The key difference is that API keys within the request's Authorization header will be obfuscated, revealing only the last four characters. An example is provided below:

```sh
---[ REQUEST ]---------------------------------------
GET /teams/DER8RFS HTTP/1.1
Accept: application/vnd.pagerduty+json;version=2
Authorization: <OBSCURED>kCjQ
Content-Type: application/json
User-Agent: (darwin arm64) Terraform/1.5.1
```

To enable the `SECURE` log level, you must set two environment variables:

* `TF_LOG=INFO`
* `TF_LOG_PROVIDER_PAGERDUTY=SECURE`

## Testing

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`. Running the acceptance tests requires
that the `PAGERDUTY_TOKEN` environment variable be set to a valid API Token and that the
`PAGERDUTY_USER_TOKEN` environment variable be set to a valid API User Token. Many tests also
require that the [Email Domain Restriction](https://support.pagerduty.com/docs/account-settings#email-domain-restriction) feature
either be disabled *or* be configured to include `foo.test` as an allowed domain.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

*Additional Note:* In order for the tests on the Slack Connection resources to pass you will need valid Slack workspace and channel IDs from a [Slack workspace connected to your PagerDuty account](https://support.pagerduty.com/docs/slack-integration-guide#integration-walkthrough).

Run a specific subset of tests by name use the `TESTARGS="-run TestName"` option which will run all test functions with "TestName" in their name.

```sh
$ make testacc TESTARGS="-run TestAccPagerDutyTeam"
```

Some tests require additional environment variables to be set to enable them due to account restrictions on certain
features. Similarly to [`TF_ACC`](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests#environment-variables),
the value of the environment variable is not relevant.

For example:
```sh
PAGERDUTY_ACC_INCIDENT_WORKFLOWS=1 make testacc TESTARGS="-run PagerDutyIncidentWorkflow"
PAGERDUTY_ACC_SERVICE_INTEGRATION_GENERIC_EMAIL_NO_FILTERS="user@<your_domain>.pagerduty.com" make testacc TESTARGS="-run PagerDutyServiceIntegration_GenericEmailNoFilters"
PAGERDUTY_ACC_INCIDENT_CUSTOM_FIELDS=1 make testacc TESTARGS="-run PagerDutyIncidentCustomField"
PAGERDUTY_ACC_LICENSE_NAME="Full User" make testacc TESTARGS="-run DataSourcePagerDutyLicense_Basic"
PAGERDUTY_ACC_SCHEDULE_USED_BY_EP_W_1_LAYER=1 make testacc TESTARGS="-run PagerDutyScheduleWithTeams_EscalationPolicyDependantWithOneLayer"
```

| Variable Name                                                | Feature Set         |
|--------------------------------------------------------------|---------------------|
| `PAGERDUTY_ACC_INCIDENT_WORKFLOWS`                           | Incident Workflows  |
| `PAGERDUTY_ACC_SERVICE_INTEGRATION_GENERIC_EMAIL_NO_FILTERS` | Service Integration |
| `PAGERDUTY_ACC_INCIDENT_CUSTOM_FIELDS`                       | Custom Fields       |
| `PAGERDUTY_ACC_LICENSE_NAME`                                 | Licenses            |
| `PAGERDUTY_ACC_SCHEDULE_USED_BY_EP_W_1_LAYER`                | Schedule            |
| `PAGERDUTY_ACC_JIRA_ACCOUNT_MAPPING_ID`                      | Set Jira account-mapping ID to use during acceptance tests |
| `PAGERDUTY_ACC_EXTERNAL_PROVIDER_VERSION`                    | Modifies the version used to compare plans between sdkv2 and framework implementations. Default `~> 3.6`. |
