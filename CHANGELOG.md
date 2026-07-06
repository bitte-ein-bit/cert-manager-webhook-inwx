# Changelog

## v0.7.0

- **Breaking:** renamed the API `groupName` to `cert-manager-webhook-inwx.bitte-ein-bit.github.com`. Update the `groupName` in your `Issuer`/`ClusterIssuer` `webhook` config to match when upgrading.
- Made the 2FA/OTP unlock robust: retry across TOTP windows so rapid or concurrent logins no longer fail on INWX's single-use OTP policy
- Renamed the Go module to `github.com/bitte-ein-bit/cert-manager-webhook-inwx`

## v0.6.0

- Modernized dependencies: cert-manager `v1.20.3` (module rename `jetstack` → `cert-manager`), Kubernetes libraries `v0.35.2`, `klog/v2`, goinwx `v0.12.0`, Go `1.25`
- Migrated apiextensions API `v1beta1` → `v1`
- Added an end-to-end conformance workflow running against the INWX OTE sandbox on every PR
- Dropped support for non-2FA INWX accounts (2FA is now required)

## v0.5.1

- Maintained fork: container images now published to `ghcr.io/bitte-ein-bit/cert-manager-webhook-inwx`
- Helm chart published via GitHub Pages at <https://bitte-ein-bit.github.io/cert-manager-webhook-inwx>
- Image tag now defaults to the chart `appVersion`
- Add missing `replicaCount` chart value; fix log format verb

## v0.5.0

- Support for multiple credentialsSecretRefs [#7](https://gitlab.com/smueller18/cert-manager-webhook-inwx/-/issues/7)

## v0.4.1

- Add CA certificates to Docker image

## v0.4.0

- Add multi arch container images
- Support INWX accounts protected by multi factor authentication
