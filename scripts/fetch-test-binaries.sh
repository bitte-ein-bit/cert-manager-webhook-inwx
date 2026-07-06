#!/bin/sh
# Provision the Kubernetes integration-test binaries (etcd, kube-apiserver)
# required by the cert-manager ACME DNS conformance suite.
#
# The suite locates them via the KUBEBUILDER_ASSETS environment variable, so
# run this script through command substitution:
#
#   export KUBEBUILDER_ASSETS="$(scripts/fetch-test-binaries.sh)"
#
# Progress output goes to stderr; only the assets path is printed to stdout.
set -e

go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest >&2
"$(go env GOPATH)/bin/setup-envtest" use -p path
