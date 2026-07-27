package facilitatorpool

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	x402 "github.com/x402-foundation/x402/go/v2"
)

// ClientSpec gives a facilitator a stable operator-facing identity.
type ClientSpec struct {
	Name   string
	URL    string
	Client x402.FacilitatorClient
}

type EndpointStatus struct {
	Name                 string `json:"name"`
	URL                  string `json:"url"`
	Order                int    `json:"order"`
	State                string `json:"state"`
	LastOperation        string `json:"last_operation,omitempty"`
	LastCheckedAt        string `json:"last_checked_at,omitempty"`
	LastError            string `json:"last_error,omitempty"`
	SuccessfulOperations uint64 `json:"successful_operations"`
	FailedOperations     uint64 `json:"failed_operations"`
}

type SettlementStatus struct {
	State                    string `json:"state"`
	Facilitator              string `json:"facilitator,omitempty"`
	Transaction              string `json:"transaction,omitempty"`
	LastError                string `json:"last_error,omitempty"`
	UpdatedAt                string `json:"updated_at,omitempty"`
	AutomaticSettlementRetry bool   `json:"automatic_settlement_retry"`
	RecoveryAction           string `json:"recovery_action"`
}

type Snapshot struct {
	Mode                string           `json:"mode"`
	SelectionPolicy     string           `json:"selection_policy"`
	VerifyFailover      bool             `json:"verify_failover"`
	SettleFailover      bool             `json:"settle_failover"`
	DuplicatePayGuard   bool             `json:"duplicate_payment_guard"`
	ActiveFacilitator   string           `json:"active_facilitator,omitempty"`
	Endpoints           []EndpointStatus `json:"endpoints"`
	LastSettlementState SettlementStatus `json:"last_settlement_state"`
}

type endpoint struct {
	spec   ClientSpec
	status EndpointStatus
}

// Pool is an x402 FacilitatorClient that safely provides ordered redundancy.
//
// GetSupported and Verify may try the next endpoint because those operations do
// not transfer funds. Settle is deliberately attempted exactly once. A transport
// error after Settle begins is ambiguous: another facilitator is never called,
// and the operator must reconcile the chain/proof state before another payment.
type Pool struct {
	mu                 sync.Mutex
	endpoints          []endpoint
	active             int
	verifiedSelections map[[32]byte]int
	lastSettlement     SettlementStatus
}

func New(specs []ClientSpec) (*Pool, error) {
	if len(specs) == 0 {
		return nil, errors.New("at least one facilitator is required")
	}
	p := &Pool{
		active:             -1,
		verifiedSelections: make(map[[32]byte]int),
		lastSettlement: SettlementStatus{
			State:                    "idle",
			AutomaticSettlementRetry: false,
			RecoveryAction:           "none",
		},
	}
	for i, spec := range specs {
		if spec.Client == nil {
			return nil, fmt.Errorf("facilitator %d has no client", i+1)
		}
		name := strings.TrimSpace(spec.Name)
		if name == "" {
			name = fmt.Sprintf("facilitator-%d", i+1)
		}
		p.endpoints = append(p.endpoints, endpoint{
			spec: ClientSpec{Name: name, URL: strings.TrimRight(spec.URL, "/"), Client: spec.Client},
			status: EndpointStatus{
				Name: name, URL: strings.TrimRight(spec.URL, "/"), Order: i + 1, State: "unknown",
			},
		})
	}
	return p, nil
}

func (p *Pool) GetSupported(ctx context.Context) (x402.SupportedResponse, error) {
	var failures []error
	for index := range p.endpoints {
		client := p.client(index)
		response, err := client.GetSupported(ctx)
		if err != nil {
			p.record(index, "supported", err)
			failures = append(failures, fmt.Errorf("%s: %w", p.name(index), err))
			continue
		}
		p.record(index, "supported", nil)
		p.setActive(index)
		return response, nil
	}
	return x402.SupportedResponse{}, fmt.Errorf("all facilitators failed supported check: %w", errors.Join(failures...))
}

// Probe checks every configured endpoint without signing or settling. It
// selects the first healthy endpoint in configured priority order.
func (p *Pool) Probe(ctx context.Context) error {
	var failures []error
	firstHealthy := -1
	for index := range p.endpoints {
		endpointCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
		response, err := p.client(index).GetSupported(endpointCtx)
		cancel()
		if err != nil {
			p.record(index, "supported", err)
			failures = append(failures, fmt.Errorf("%s: %w", p.name(index), err))
			continue
		}
		if len(response.Kinds) == 0 {
			err = errors.New("supported response contains no payment kinds")
			p.record(index, "supported", err)
			failures = append(failures, fmt.Errorf("%s: %w", p.name(index), err))
			continue
		}
		p.record(index, "supported", nil)
		if firstHealthy < 0 {
			firstHealthy = index
		}
	}
	if firstHealthy < 0 {
		return fmt.Errorf("all facilitators failed health probe: %w", errors.Join(failures...))
	}
	p.setActive(firstHealthy)
	return nil
}

func (p *Pool) Verify(ctx context.Context, payloadBytes, requirementsBytes []byte) (*x402.VerifyResponse, error) {
	var failures []error
	for _, index := range p.orderedIndexes() {
		client := p.client(index)
		response, err := client.Verify(ctx, payloadBytes, requirementsBytes)
		if err != nil {
			p.record(index, "verify", err)
			failures = append(failures, fmt.Errorf("%s: %w", p.name(index), err))
			continue
		}
		p.record(index, "verify", nil)
		p.setActive(index)
		if response != nil && response.IsValid {
			p.rememberSelection(fingerprint(payloadBytes, requirementsBytes), index)
		}
		// A reachable facilitator's invalid-payment response is definitive.
		return response, nil
	}
	return nil, fmt.Errorf("all facilitators failed payment verification: %w", errors.Join(failures...))
}

func (p *Pool) Settle(ctx context.Context, payloadBytes, requirementsBytes []byte) (*x402.SettleResponse, error) {
	key := fingerprint(payloadBytes, requirementsBytes)
	index := p.selection(key)
	if index < 0 {
		index = p.activeIndex()
	}
	if index < 0 {
		return nil, errors.New("no healthy facilitator selected; run GetSupported or Verify before settlement")
	}

	// Remove the selection before the external call. Even if the call times out,
	// this exact payload cannot silently reuse a cached route for another attempt.
	p.forgetSelection(key)
	client := p.client(index)
	response, err := client.Settle(ctx, payloadBytes, requirementsBytes)
	if err != nil {
		p.record(index, "settle", err)
		p.setSettlement(SettlementStatus{
			State:                    "unknown_requires_reconciliation",
			Facilitator:              p.name(index),
			LastError:                err.Error(),
			UpdatedAt:                time.Now().UTC().Format(time.RFC3339Nano),
			AutomaticSettlementRetry: false,
			RecoveryAction:           "inspect Base Sepolia and reconcile durable authorization before any retry",
		})
		return nil, fmt.Errorf(
			"settlement outcome is unknown after calling %s; automatic failover disabled to prevent duplicate payment: %w",
			p.name(index), err,
		)
	}
	p.record(index, "settle", nil)
	state := "rejected"
	recovery := "review facilitator rejection before creating a new authorization"
	transaction := ""
	if response != nil && response.Success {
		state = "settled"
		recovery = "none"
		transaction = response.Transaction
	}
	p.setSettlement(SettlementStatus{
		State:                    state,
		Facilitator:              p.name(index),
		Transaction:              transaction,
		UpdatedAt:                time.Now().UTC().Format(time.RFC3339Nano),
		AutomaticSettlementRetry: false,
		RecoveryAction:           recovery,
	})
	return response, nil
}

func (p *Pool) Snapshot() Snapshot {
	p.mu.Lock()
	defer p.mu.Unlock()
	statuses := make([]EndpointStatus, len(p.endpoints))
	for i := range p.endpoints {
		statuses[i] = p.endpoints[i].status
	}
	active := ""
	if p.active >= 0 && p.active < len(p.endpoints) {
		active = p.endpoints[p.active].spec.Name
	}
	return Snapshot{
		Mode:                "ordered-health-aware",
		SelectionPolicy:     "configured-order health probe; pin successful verifier to exact settlement payload",
		VerifyFailover:      true,
		SettleFailover:      false,
		DuplicatePayGuard:   true,
		ActiveFacilitator:   active,
		Endpoints:           statuses,
		LastSettlementState: p.lastSettlement,
	}
}

func (p *Pool) orderedIndexes() []int {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := make([]int, 0, len(p.endpoints))
	if p.active >= 0 && p.active < len(p.endpoints) {
		result = append(result, p.active)
	}
	for i := range p.endpoints {
		if i != p.active {
			result = append(result, i)
		}
	}
	return result
}

func (p *Pool) client(index int) x402.FacilitatorClient {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.endpoints[index].spec.Client
}

func (p *Pool) name(index int) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.endpoints[index].spec.Name
}

func (p *Pool) record(index int, operation string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	status := &p.endpoints[index].status
	status.LastOperation = operation
	status.LastCheckedAt = time.Now().UTC().Format(time.RFC3339Nano)
	if err != nil {
		status.State = "degraded"
		status.LastError = err.Error()
		status.FailedOperations++
		return
	}
	status.State = "healthy"
	status.LastError = ""
	status.SuccessfulOperations++
}

func (p *Pool) setActive(index int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.active = index
}

func (p *Pool) activeIndex() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.active
}

func (p *Pool) rememberSelection(key [32]byte, index int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.verifiedSelections[key] = index
}

func (p *Pool) selection(key [32]byte) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	index, ok := p.verifiedSelections[key]
	if !ok {
		return -1
	}
	return index
}

func (p *Pool) forgetSelection(key [32]byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.verifiedSelections, key)
}

func (p *Pool) setSettlement(status SettlementStatus) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lastSettlement = status
}

func fingerprint(payloadBytes, requirementsBytes []byte) [32]byte {
	combined := make([]byte, 0, len(payloadBytes)+1+len(requirementsBytes))
	combined = append(combined, payloadBytes...)
	combined = append(combined, 0)
	combined = append(combined, requirementsBytes...)
	return sha256.Sum256(combined)
}
