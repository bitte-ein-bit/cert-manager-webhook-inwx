package main

import (
	"context"
	"os"
	"testing"
	"time"

	acmetest "github.com/cert-manager/cert-manager/test/acme"
	"github.com/cert-manager/cert-manager/test/acme/server"
	"github.com/go-logr/logr"
	"gitlab.com/smueller18/cert-manager-webhook-inwx/test"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var (
	zone      = "smueller18.de."
	zoneTwoFA = "smueller18mfa.de."
	fqdn      string
)

func TestRunSuite(t *testing.T) {

	if os.Getenv("TEST_ZONE_NAME") != "" {
		zone = os.Getenv("TEST_ZONE_NAME")
	}
	fqdn = "cert-manager-dns01-tests." + zone

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
			Zones: []string{zone},
		},
	}

	if err := srv.Run(ctx, "udp"); err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer srv.Shutdown()

	d, err := os.ReadFile("testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}

	fixture := acmetest.NewFixture(&solver{},
		acmetest.SetResolvedZone(zone),
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

func TestRunSuiteWithSecret(t *testing.T) {

	if os.Getenv("TEST_ZONE_NAME") != "" {
		zone = os.Getenv("TEST_ZONE_NAME")
	}
	fqdn = "cert-manager-dns01-tests-with-secret." + zone

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
			Zones: []string{zone},
		},
	}

	if err := srv.Run(ctx, "udp"); err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer srv.Shutdown()

	d, err := os.ReadFile("testdata/config.secret.json")
	if err != nil {
		t.Fatal(err)
	}

	fixture := acmetest.NewFixture(&solver{},
		acmetest.SetResolvedZone(zone),
		acmetest.SetResolvedFQDN(fqdn),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetDNSServer(srv.ListenAddr()),
		acmetest.SetManifestPath("testdata/secret-inwx-credentials.yaml"),
		acmetest.SetPropagationLimit(time.Duration(60)*time.Second),
		acmetest.SetUseAuthoritative(false),
		acmetest.SetConfig(&extapi.JSON{
			Raw: d,
		}),
	)

	fixture.RunConformance(t)
}

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
