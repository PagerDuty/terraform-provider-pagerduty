# Documentation Index

Entry point for this repo's documentation. Add new docs here as
`YYYYMMDD-HHMM_<topic>.md` files.

## Repository layout

- `pagerduty/` — provider implementation on the Terraform Plugin SDKv2. Most
  existing resources and data sources live here.
- `pagerdutyplugin/` — provider implementation on the newer Terraform Plugin
  Framework. New resources are added here going forward.
- `website/docs/` — Terraform Registry documentation for resources and data
  sources (rendered at registry.terraform.io).
- `vendor/` — checked-in Go module dependencies.
- `CHANGELOG.md` — release notes, one entry per PR/feature.

## Building and testing

Defined in `GNUmakefile`:

- `make build` — `go install` the provider binary.
- `make test` — unit tests (excludes `vendor/`).
- `make testacc` — acceptance tests; requires `TF_ACC=1` and real PagerDuty
  API credentials (`PAGERDUTY_TOKEN`), runs against live PagerDuty resources.
- `make vet` / `make fmt` / `make fmtcheck` — static checks and formatting.
- `make update-go-pagerduty` — bump the `go-pagerduty` client and re-vendor.

## Agent memory

See `AGENTS.md` at the repository root for conventions agents should follow
when making changes here.
