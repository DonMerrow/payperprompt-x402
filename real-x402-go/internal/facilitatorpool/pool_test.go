package facilitatorpool

import (
	"context"
	"errors"
	"testing"

	x402 "github.com/x402-foundation/x402/go/v2"
)

type fakeClient struct {
	supported      x402.SupportedResponse
	supportedErr   error
	verify         *x402.VerifyResponse
	verifyErr      error
	settle         *x402.SettleResponse
	settleErr      error
	supportedCalls int
	verifyCalls    int
	settleCalls    int
}

func (f *fakeClient) GetSupported(context.Context) (x402.SupportedResponse, error) {
	f.supportedCalls++
	return f.supported, f.supportedErr
}

func (f *fakeClient) Verify(context.Context, []byte, []byte) (*x402.VerifyResponse, error) {
	f.verifyCalls++
	return f.verify, f.verifyErr
}

func (f *fakeClient) Settle(context.Context, []byte, []byte) (*x402.SettleResponse, error) {
	f.settleCalls++
	return f.settle, f.settleErr
}

func supported() x402.SupportedResponse {
	return x402.SupportedResponse{Kinds: []x402.SupportedKind{
		{X402Version: 2, Scheme: "exact", Network: "eip155:84532"},
	}}
}

func TestSupportedFailsOverToHealthyEndpoint(t *testing.T) {
	primary := &fakeClient{supportedErr: errors.New("primary offline")}
	secondary := &fakeClient{supported: supported()}
	pool, err := New([]ClientSpec{
		{Name: "primary", URL: "https://primary.example", Client: primary},
		{Name: "secondary", URL: "https://secondary.example", Client: secondary},
	})
	if err != nil {
		t.Fatal(err)
	}
	response, err := pool.GetSupported(context.Background())
	if err != nil || len(response.Kinds) != 1 {
		t.Fatalf("expected secondary supported response, response=%+v err=%v", response, err)
	}
	snapshot := pool.Snapshot()
	if snapshot.ActiveFacilitator != "secondary" ||
		snapshot.Endpoints[0].State != "degraded" ||
		snapshot.Endpoints[1].State != "healthy" {
		t.Fatalf("unexpected status: %+v", snapshot)
	}
}

func TestProbeChecksEveryEndpointAndRestoresConfiguredPriority(t *testing.T) {
	primary := &fakeClient{supportedErr: errors.New("offline")}
	secondary := &fakeClient{supported: supported()}
	pool, _ := New([]ClientSpec{
		{Name: "primary", Client: primary},
		{Name: "secondary", Client: secondary},
	})
	if _, err := pool.GetSupported(context.Background()); err != nil {
		t.Fatal(err)
	}
	primary.supportedErr = nil
	primary.supported = supported()
	if err := pool.Probe(context.Background()); err != nil {
		t.Fatal(err)
	}
	snapshot := pool.Snapshot()
	if snapshot.ActiveFacilitator != "primary" {
		t.Fatalf("probe did not restore configured priority: %+v", snapshot)
	}
	if primary.supportedCalls != 2 || secondary.supportedCalls != 2 {
		t.Fatalf("probe did not check every endpoint: primary=%d secondary=%d", primary.supportedCalls, secondary.supportedCalls)
	}
}

func TestVerifyFailoverPinsSettlementToVerifier(t *testing.T) {
	primary := &fakeClient{verifyErr: errors.New("verify timeout")}
	secondary := &fakeClient{
		verify: &x402.VerifyResponse{IsValid: true},
		settle: &x402.SettleResponse{
			Success: true, Transaction: "0xabc", Network: x402.Network("eip155:84532"),
		},
	}
	pool, _ := New([]ClientSpec{
		{Name: "primary", Client: primary},
		{Name: "secondary", Client: secondary},
	})
	payload, requirements := []byte(`{"payload":1}`), []byte(`{"amount":"10000"}`)
	verified, err := pool.Verify(context.Background(), payload, requirements)
	if err != nil || !verified.IsValid {
		t.Fatalf("verification should fail over: response=%+v err=%v", verified, err)
	}
	settled, err := pool.Settle(context.Background(), payload, requirements)
	if err != nil || !settled.Success {
		t.Fatalf("settlement should use verifier: response=%+v err=%v", settled, err)
	}
	if primary.settleCalls != 0 || secondary.settleCalls != 1 {
		t.Fatalf("wrong settlement target: primary=%d secondary=%d", primary.settleCalls, secondary.settleCalls)
	}
}

func TestAmbiguousSettlementNeverRetriesAnotherFacilitator(t *testing.T) {
	primary := &fakeClient{
		verify:    &x402.VerifyResponse{IsValid: true},
		settleErr: errors.New("connection reset after request body"),
	}
	secondary := &fakeClient{
		verify: &x402.VerifyResponse{IsValid: true},
		settle: &x402.SettleResponse{Success: true, Transaction: "0xduplicate"},
	}
	pool, _ := New([]ClientSpec{
		{Name: "primary", Client: primary},
		{Name: "secondary", Client: secondary},
	})
	payload, requirements := []byte(`{"payload":2}`), []byte(`{"amount":"20000"}`)
	if _, err := pool.Verify(context.Background(), payload, requirements); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Settle(context.Background(), payload, requirements); err == nil {
		t.Fatal("expected ambiguous settlement error")
	}
	if primary.settleCalls != 1 || secondary.settleCalls != 0 {
		t.Fatalf("settlement was retried: primary=%d secondary=%d", primary.settleCalls, secondary.settleCalls)
	}
	status := pool.Snapshot().LastSettlementState
	if status.State != "unknown_requires_reconciliation" || status.AutomaticSettlementRetry {
		t.Fatalf("unsafe settlement state: %+v", status)
	}
}

func TestDefinitiveRejectionIsNotMarkedAmbiguous(t *testing.T) {
	client := &fakeClient{
		verify: &x402.VerifyResponse{IsValid: true},
		settle: &x402.SettleResponse{Success: false, ErrorReason: "invalid_payment"},
	}
	pool, _ := New([]ClientSpec{{Name: "only", Client: client}})
	payload, requirements := []byte(`{"payload":3}`), []byte(`{"amount":"10000"}`)
	_, _ = pool.Verify(context.Background(), payload, requirements)
	response, err := pool.Settle(context.Background(), payload, requirements)
	if err != nil || response.Success {
		t.Fatalf("expected definitive rejection response=%+v err=%v", response, err)
	}
	if state := pool.Snapshot().LastSettlementState.State; state != "rejected" {
		t.Fatalf("expected rejected state, got %s", state)
	}
}
