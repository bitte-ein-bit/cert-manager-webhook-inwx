# Code Analysis — updates & issues

Snapshot taken 2026-07-06 for the `bitte-ein-bit` fork of
`cert-manager-webhook-inwx` (a maintained fork of the archived
`gitlab.com/smueller18/cert-manager-webhook-inwx`).

The webhook itself is small and functionally sound — the issues below are almost
all about **age**: the dependency tree and toolchain are pinned to 2020-era
versions. Nothing here blocks the container from building and running today, but
the modernization items are worth scheduling.

## Already fixed in this pass (safe changes)

These were applied as part of setting up the fork:

- **Dockerfile Go version.** Was `golang:1.13-alpine` (2019) while GitLab CI used
  `golang:1.21-alpine` — inconsistent and ancient. Now `golang:1.23-alpine`, and
  the build stage uses Buildx `TARGETOS`/`TARGETARCH`/`TARGETVARIANT` so the new
  GitHub Actions workflow can produce `linux/amd64`, `linux/arm64` and
  `linux/arm/v7` from a single definition (replacing the old per-arch matrix +
  manual `docker manifest` dance).
- **`go.mod` Go directive** bumped `1.13` → `1.23` and `go mod tidy` run (added the
  pruned indirect-require block; no dependency versions changed).
- **Format-verb bug** in `main.go` `loadConfig`: `klog.Warningf("... %q", defaultConfig.TTL)`
  used `%q` on an `int`. Changed to `%d`. (`go vet` flags this.)
- **Registry migration** GitLab → GHCR: `values.yaml` `image.repository` and the
  README now point at `ghcr.io/bitte-ein-bit/cert-manager-webhook-inwx`.
- **Missing chart value**: `deployment.yaml` references `.Values.replicaCount`,
  which did not exist in `values.yaml` (rendered an empty `replicas:`). Added
  `replicaCount: 1`.

## Recommended updates (not done — need testing against a live INWX account)

The conformance tests in `main_test.go` hit the real INWX OTE sandbox and need
credentials, so none of these were applied blind.

### 1. cert-manager module rename + major upgrade — highest priority

`main.go` and `go.mod` import `github.com/jetstack/cert-manager v1.0.1`. That
module was **renamed to `github.com/cert-manager/cert-manager`** years ago, and
v1.0.1 is from 2020. Current stable is ~`v1.20.x`.

This also forces an API migration: the code imports
`k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1` (used as
`extapi.JSON`). The `apiextensions.k8s.io/v1beta1` group was **removed in
Kubernetes 1.22** — replace with `.../apiextensions/v1`. The webhook `Solver`
interface and `apis/acme/v1alpha1` types have also moved packages.

Effort: medium. Touch points are small (imports + the `extapi.JSON` type in
`loadConfig`), but it must be validated with the conformance suite.

### 2. Kubernetes client libraries: `v0.19.0` → `v0.3x`

`k8s.io/client-go`, `apimachinery`, `apiextensions-apiserver`, `apiserver`,
`component-base`, `kube-aggregator` are all pinned to `v0.19.0` (K8s 1.19, 2020).
Current stable line is `v0.36.x`. Should be upgraded in lockstep with the
cert-manager bump since they share transitive versions.

### 3. `k8s.io/klog` v1 → `k8s.io/klog/v2`

`main.go` imports `k8s.io/klog` (v1, effectively unmaintained). The v2 module is
already present transitively (`k8s.io/klog/v2` is an indirect dep). Swap the
import; the API is source-compatible for the calls used here (`Errorf`, `Warningf`,
`V(n).Infof`, `Error`).

### 4. `github.com/nrdcg/goinwx`: `v0.6.1` → `v0.12.0`

Five minor releases behind. Newer versions carry API/bugfix changes (record
handling, context support). Review the changelog — method signatures such as
`Nameservers.CreateRecord` / `Info` / `DeleteRecord` may have changed and will
need small edits in `main.go`.

### 5. Retire the deprecated linter

The old GitLab CI used `golang.org/x/lint/golint`, which is **archived and
deprecated**. The new CI workflow drops it in favor of `go vet`. If you want
lint coverage back, adopt `golangci-lint` instead.

## Code-quality issues (low severity, not fixed)

These are cosmetic / defensive and don't affect correctness:

- **`main.go:254` `tryToUnlockWithOTPKey` returns `(error, error)`.** An unusual
  signature where the first return is the raw error and the second a formatted
  one; the first return value is discarded by the only caller. Worth collapsing
  to a single `error`. The 30s `time.Sleep` retry (to satisfy INWX's
  single-OTP-use policy) is intentional and correct, just undocumented in code.
- **`s.ttl` field on `solver` is set but never read** (`newClientFromChallenge`
  assigns `s.ttl = cfg.TTL`). The TTL actually used comes from `cfg.TTL`. Dead
  field — can be removed.
- **`CleanUp` accumulates only the last delete error** (`lastErr`). If multiple
  records fail to delete, earlier errors are logged but not returned. Fine in
  practice (each is logged), but `errors.Join` would be cleaner.
- **`getCredentials` error messages** interpolate the whole `SecretKeySelector`
  struct with `%q` in a couple of "no key" branches (`config.UsernameSecretKeyRef`
  rather than `.Key`), producing noisy messages. Cosmetic.

## Chart / deployment observations (not fixed)

- **`groupName` still `...smueller18.gitlab.com`.** It's baked into `main.go`,
  `values.yaml`, `apiservice.yaml` and the RBAC/`Issuer` examples. Keeping it
  avoids breaking existing installs, but if you want the fork fully de-branded
  it must be changed in all four places at once (breaking change).
- **No `securityContext` hardening** beyond `runAsUser: 65534`. Consider
  `runAsNonRoot: true`, `readOnlyRootFilesystem: true`,
  `allowPrivilegeEscalation: false`, and dropping all capabilities.
- **No resource requests/limits** by default (`resources: {}`).
- **README install instructions** still reference the upstream Helm repo
  `https://smueller18.gitlab.io/helm-charts`, which this fork does not publish.
  Until/unless you host a chart repo (e.g. GitHub Pages or an OCI chart in GHCR),
  install from the local `deploy/` path. Documenting this is a follow-up.

## Suggested follow-up order

1. cert-manager module rename + `apiextensions v1` migration (#1) together with
   the K8s client bump (#2) and `klog/v2` (#3) — one coordinated PR, validated
   against the OTE conformance suite.
2. `goinwx` upgrade (#4).
3. Chart hardening + de-branding decision.
4. Code-quality cleanups.
