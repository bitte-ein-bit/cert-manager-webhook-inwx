# Changelog

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
