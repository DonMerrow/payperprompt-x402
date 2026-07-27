package payperprompt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

type Service struct {
	RouteID    string  `json:"route_id"`
	Provider   string  `json:"provider"`
	Capability string  `json:"capability"`
	Resource   string  `json:"resource"`
	PriceUSD   float64 `json:"price_usd"`
	LatencyMS  int     `json:"latency_ms"`
	Quality    string  `json:"quality"`
}

type PlanRequest struct {
	Prompt        string `json:"prompt"`
	TaskType      string `json:"task_type"`
	ExpectedPayer string `json:"expected_payer"`
}

type PlanResponse struct {
	Planned      bool            `json:"planned"`
	Model        string          `json:"model"`
	AIUsed       bool            `json:"ai_used"`
	Reason       string          `json:"reason,omitempty"`
	WorkOrder    json.RawMessage `json:"work_order"`
	Service      Service         `json:"selected_service"`
	Policy       json.RawMessage `json:"policy_decision"`
	Challenge    json.RawMessage `json:"challenge"`
	PreparedWork json.RawMessage `json:"prepared_work"`
}

type WorkAuditEvent struct {
	Stage         string `json:"stage"`
	Status        string `json:"status"`
	TaskType      string `json:"task_type,omitempty"`
	Title         string `json:"title,omitempty"`
	RouteID       string `json:"route_id,omitempty"`
	Provider      string `json:"provider,omitempty"`
	AmountUSD     string `json:"amount_usd,omitempty"`
	Wallet        string `json:"wallet,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
	Reason        string `json:"reason,omitempty"`
	RecordedAt    string `json:"recorded_at"`
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{Timeout: 220 * time.Second},
	}
}

func (c *Client) Services(ctx context.Context) ([]Service, error) {
	var response struct {
		Services []Service `json:"services"`
	}
	if err := c.doJSON(ctx, http.MethodGet, "/api/services", nil, &response); err != nil {
		return nil, err
	}
	return response.Services, nil
}

func (c *Client) PrepareWork(ctx context.Context, request PlanRequest) (PlanResponse, error) {
	var response PlanResponse
	if strings.TrimSpace(request.Prompt) == "" ||
		strings.TrimSpace(request.ExpectedPayer) == "" {
		return response, errors.New("prompt and expected payer are required")
	}
	err := c.doJSON(ctx, http.MethodPost, "/api/official/plan", request, &response)
	return response, err
}

func (c *Client) WorkAudit(ctx context.Context, wallet string, history bool) ([]WorkAuditEvent, error) {
	query := url.Values{}
	query.Set("wallet", wallet)
	query.Set("history", fmt.Sprintf("%t", history))
	var response struct {
		Events []WorkAuditEvent `json:"events"`
	}
	if err := c.doJSON(ctx, http.MethodGet, "/api/official/work-audit?"+query.Encode(), nil, &response); err != nil {
		return nil, err
	}
	return response.Events, nil
}

func (c *Client) VerifyEvidence(ctx context.Context, prompt string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.doJSON(ctx, http.MethodPost, "/api/official/evidence", map[string]string{
		"prompt": prompt,
	}, &response)
	return response, err
}

func (c *Client) doJSON(ctx context.Context, method, path string, input, output any) error {
	var body io.Reader
	if input != nil {
		payload, err := json.Marshal(input)
		if err != nil {
			return err
		}
		body = bytes.NewReader(payload)
	}
	request, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
	if err != nil {
		return err
	}
	if input != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response, err := c.HTTP.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(response.Body, 8<<20))
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var failure struct {
			Reason string `json:"reason"`
			Error  string `json:"error"`
		}
		_ = json.Unmarshal(payload, &failure)
		message := failure.Reason
		if message == "" {
			message = failure.Error
		}
		if message == "" {
			message = strings.TrimSpace(string(payload))
		}
		return fmt.Errorf("PayPerPrompt returned HTTP %d: %s", response.StatusCode, message)
	}
	if err := json.Unmarshal(payload, output); err != nil {
		return fmt.Errorf("decode PayPerPrompt response: %w", err)
	}
	return nil
}
