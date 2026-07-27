package payperprompt

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientDiscoversServicesAndCurrentAudit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/services":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"services": []Service{{RouteID: "guardrail-economy", PriceUSD: 0.01}},
			})
		case "/api/official/work-audit":
			if r.URL.Query().Get("history") != "false" {
				t.Fatalf("SDK did not request current-only audit")
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"events": []WorkAuditEvent{{Stage: "settlement", Status: "settled"}},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := New(server.URL)
	services, err := client.Services(context.Background())
	if err != nil || len(services) != 1 || services[0].PriceUSD != 0.01 {
		t.Fatalf("unexpected services: %+v err=%v", services, err)
	}
	events, err := client.WorkAudit(context.Background(), "0x123", false)
	if err != nil || len(events) != 1 || events[0].Status != "settled" {
		t.Fatalf("unexpected audit: %+v err=%v", events, err)
	}
}

func TestClientSurfacesPolicyDenialWithoutSigning(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"reason": "Daily spending limit would be exceeded.",
		})
	}))
	defer server.Close()

	_, err := New(server.URL).PrepareWork(context.Background(), PlanRequest{
		Prompt:        "Review this code.",
		ExpectedPayer: "0x826154a3d58aea3fbd2aa64aad424594ade927ef",
	})
	if err == nil {
		t.Fatal("policy denial was not returned")
	}
}
