package main

import (
	"os"
	"reflect"
	"testing"
)

func TestOfficialServicesExposeThreePricePoints(t *testing.T) {
	services := officialServices(config{PriceUSD: "0.01"})
	if len(services) != 3 {
		t.Fatalf("expected three official services, got %d", len(services))
	}
	want := []struct {
		routeID string
		price   string
	}{
		{"guardrail-economy", "0.01"},
		{"guardrail-fast", "0.02"},
		{"guardrail-deep", "0.04"},
	}
	for i, expected := range want {
		if services[i].RouteID != expected.routeID || services[i].PriceUSD != expected.price {
			t.Fatalf("service %d mismatch: %+v", i, services[i])
		}
	}
}

func TestPaidReportStrategyMatchesSettledService(t *testing.T) {
	services := officialServices(config{PriceUSD: "0.01"})
	want := []string{"lowest-cost", "lowest-latency", "highest-quality"}
	for index, service := range services {
		if got := strategyForPaidService(service); got != want[index] {
			t.Fatalf("route %s mapped to %s, want %s", service.RouteID, got, want[index])
		}
	}
}

func TestFacilitatorURLsAreOrderedAndDeduplicated(t *testing.T) {
	const key = "X402_FACILITATOR_URLS"
	previous, hadPrevious := os.LookupEnv(key)
	t.Cleanup(func() {
		if hadPrevious {
			_ = os.Setenv(key, previous)
		} else {
			_ = os.Unsetenv(key)
		}
	})
	_ = os.Setenv(key, " https://primary.example/, https://backup.example,https://primary.example ")
	got := facilitatorURLs()
	want := []string{"https://primary.example", "https://backup.example"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected ordered facilitator URLs: got=%v want=%v", got, want)
	}
}
