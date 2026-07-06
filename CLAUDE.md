# CLAUDE.md

Guidance for Claude Code (and humans) working in this repository.

## What this is

A [cert-manager](https://cert-manager.io/) ACME DNS-01 webhook for the
[INWX](https://inwx.de/) domain registrar, plus a Helm chart to deploy it.
When cert-manager needs to solve a DNS-01 challenge, it calls this webhook,
which creates/deletes the required `_acme-challenge` TXT record via the INWX
XML-RPC API (using [`nrdcg/goinwx`](https://github.com/nrdcg/goinwx)).

This is a maintained fork of the archived
`gitlab.com/smueller18/cert-manager-webhook-inwx`.

## Layout

- `main.go` — the whole webhook implementation. Implements cert-manager's
  `webhook.Solver` interface: `Name`, `Initialize`, `Present`, `CleanUp`.
  Credentials come from the solver config inline or from a Kubernetes Secret
  (`usernameSecretKeyRef` / `passwordSecretKeyRef` / `otpKeySecretKeyRef`).
  2FA/OTP is supported via `pquerna/otp` (TOTP).
- `main_test.go` + `test/server.go` — conformance tests against the real INWX
  OTE (sandbox) API. They need real test credentials and cannot run in CI.
- `deploy/cert-manager-webhook-inwx/` — the Helm chart.
- `Dockerfile` — multi-stage build, static binary on `scratch`.
- `.github/workflows/` — CI (build/vet/helm lint) and container publish to
  `ghcr.io`.

## Common commands

```bash
go build ./...        # compile
go vet ./...          # static checks
helm lint deploy/cert-manager-webhook-inwx

# Full conformance tests (require real INWX OTE credentials + test binaries):
scripts/fetch-test-binaries.sh
TEST_ZONE_NAME="example.com." go test -v -cover .
```

The conformance tests hit a live sandbox account and are slow (the OTP path
sleeps 30s to satisfy INWX's single-OTP-use policy). Prefer `go build`/`go vet`
for quick feedback; don't expect `go test` to pass without credentials.

## Conventions

- The API `groupName` is `cert-manager-webhook-inwx.smueller18.gitlab.com`.
  It is baked into `main.go`, the chart's `values.yaml`, and `Issuer` examples.
  Changing it is a breaking change for existing deployments — keep them in sync
  if you ever do.
- The container image is published to
  `ghcr.io/bitte-ein-bit/cert-manager-webhook-inwx`. Keep `values.yaml`
  `image.repository`/`image.tag`, `Chart.yaml` versions, and the README table
  in sync on a release.
- INWX enforces a minimum DNS TTL of 300s; `loadConfig` clamps lower values.

## Known modernization backlog

See `ANALYSIS.md` for the full list. Highlights: the module still imports the
renamed `github.com/jetstack/cert-manager` (now `github.com/cert-manager/cert-manager`),
pins k8s libraries at v0.19.0 and `goinwx` at v0.6.1, and uses the removed
`apiextensions-apiserver/.../v1beta1` API. These upgrades are intentionally not
done yet because they require testing against a live INWX account.
