package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bitte-ein-bit/cert-manager-webhook-inwx/test"
	acmetest "github.com/cert-manager/cert-manager/test/acme"
	"github.com/cert-manager/cert-manager/test/acme/server"
	"github.com/go-logr/logr"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// This project only supports INWX accounts protected by two-factor
// authentication, so the suite exercises the OTP code paths exclusively.
var (
	zoneTwoFA = "example.de."
	fqdn      string
)

func TestRunSuiteWithTwoFA(t *testing.T) {

	if os.Getenv("TEST_ZONE_NAME_WITH_TWO_FA") != "" {
		zoneTwoFA = os.Getenv("TEST_ZONE_NAME_WITH_TWO_FA")
	}

	fqdn = "cert-manager-dns01-tests." + zoneTwoFA

	ctx := context.Background()

	srv := &server.BasicServer{
		Handler: &test.Handler{
			Log: logr.Discard(),
			TxtRecords: map[string][][]string{
				fqdn: {
					{},
					{},
					{"123d=="},
					{"123d=="},
				},
			},
			Zones: []string{zoneTwoFA},
		},
	}

	if err := srv.Run(ctx, "udp"); err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer srv.Shutdown()

	d, err := os.ReadFile("testdata/config-otp.json")
	if err != nil {
		t.Fatal(err)
	}

	fixture := acmetest.NewFixture(&solver{},
		acmetest.SetResolvedZone(zoneTwoFA),
		acmetest.SetResolvedFQDN(fqdn),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetDNSServer(srv.ListenAddr()),
		acmetest.SetPropagationLimit(time.Duration(60)*time.Second),
		acmetest.SetUseAuthoritative(false),
		// Set to false because INWX implementation deletes all records
		acmetest.SetStrict(false),
		acmetest.SetConfig(&extapi.JSON{
			Raw: d,
		}),
	)

	fixture.RunConformance(t)
}

func TestRunSuiteWithSecretAndTwoFA(t *testing.T) {

	if os.Getenv("TEST_ZONE_NAME_WITH_TWO_FA") != "" {
		zoneTwoFA = os.Getenv("TEST_ZONE_NAME_WITH_TWO_FA")
	}
	fqdn = "cert-manager-dns01-tests-with-secret." + zoneTwoFA

	ctx := context.Background()

	srv := &server.BasicServer{
		Handler: &test.Handler{
			Log: logr.Discard(),
			TxtRecords: map[string][][]string{
				fqdn: {
					{},
					{},
					{"123d=="},
					{"123d=="},
				},
			},
			Zones: []string{zoneTwoFA},
		},
	}

	if err := srv.Run(ctx, "udp"); err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer srv.Shutdown()

	d, err := os.ReadFile("testdata/config-otp.secret.json")
	if err != nil {
		t.Fatal(err)
	}

	fixture := acmetest.NewFixture(&solver{},
		acmetest.SetResolvedZone(zoneTwoFA),
		acmetest.SetResolvedFQDN(fqdn),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetDNSServer(srv.ListenAddr()),
		acmetest.SetManifestPath("testdata/secret-inwx-credentials-otp.yaml"),
		acmetest.SetPropagationLimit(time.Duration(60)*time.Second),
		acmetest.SetUseAuthoritative(false),
		acmetest.SetConfig(&extapi.JSON{
			Raw: d,
		}),
	)

	fixture.RunConformance(t)
}
