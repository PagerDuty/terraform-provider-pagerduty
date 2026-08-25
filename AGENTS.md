# Agent Memory

Terraform provider for PagerDuty. See `docs/README.md` for the full
documentation index.

## Provider packages

- `pagerduty/` — legacy provider on the Terraform Plugin SDKv2. Existing
  resources/data sources mostly live here; only touch this package to fix or
  extend something already implemented on the SDKv2.
- `pagerdutyplugin/` — provider on the Terraform Plugin Framework. Add new
  resources and data sources here unless there is a specific reason to use
  the SDKv2 package instead.

Each package has its own `config.go` and `provider.go`; resources follow the
naming convention `resource_pagerduty_<name>.go` /
`data_source_pagerduty_<name>.go`, each paired with a `_test.go` file.

## Build and test

- `make build` — build/install the provider binary.
- `make test` — unit tests.
- `make testacc` — acceptance tests. These call the real PagerDuty API and
  require `TF_ACC=1` plus a valid `PAGERDUTY_TOKEN` (and related env vars in
  `pagerduty/config.go` / `pagerdutyplugin/config.go`). Do not run these
  without explicit credentials and user awareness that they hit live
  infrastructure.
- `make vet`, `make fmt`, `make fmtcheck` — run before submitting changes.

## Dependencies

`vendor/` is checked into the repo. After changing `go.mod`/`go.sum`, run
`go mod vendor` (or `make update-go-pagerduty` for the go-pagerduty client
specifically) and commit the resulting vendor changes.

## Changelog

Add an entry to `CHANGELOG.md` for user-facing changes (new resource,
behavior change, bug fix), following the existing format in that file.

## Documentation

Terraform Registry docs for resources/data sources live in `website/docs/`
and must be updated alongside any schema change.
