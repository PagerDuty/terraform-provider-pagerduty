# Terraform Provider for PagerDuty

- Website: https://registry.terraform.io/providers/PagerDuty/pagerduty/latest
- Documentation: https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs
- Terraform Gitter: [![Terraform Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Terraform Google Groups](http://groups.google.com/group/terraform-tool)

[PagerDuty](https://www.pagerduty.com/) is an alarm aggregation and dispatching service for system administrators and support teams. It collects alerts from your monitoring tools, gives you an overall view of all of your monitoring alarms, and alerts an on duty engineer if there’s a problem. The Terraform Pagerduty provider is a plugin for Terraform that allows for the management of PagerDuty resources using HCL (HashiCorp Configuration Language).

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-pagerduty`

```sh
$ mkdir -p $GOPATH/src/github.com/PagerDuty; cd $GOPATH/src/github.com/PagerDuty
$ git clone git@github.com:PagerDuty/terraform-provider-pagerduty
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/PagerDuty/terraform-provider-pagerduty
$ make build
```

## Using the provider

Please refer to https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs for
examples on how to use the provider and detailed documentation about the
Resources and Data Sources the provider has.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-pagerduty
...
```

### Testing

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
