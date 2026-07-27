package x402adapter

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type Config struct {
	FacilitatorURL  string
	Network         string
	Scheme          string
	PriceUSD        string
	Asset           string
	MerchantAddress string
}

type PaymentChallenge struct {
	X402Version     string              `json:"x402_version"`
	Mode            string              `json:"mode"`
	RequestID       string              `json:"request_id"`
	Reason          string              `json:"reason"`
	Accepts         []PaymentAcceptance `json:"accepts"`
	RetryWithHeader string              `json:"retry_with_header"`
	ProductionNote  string              `json:"production_note"`
}

type PaymentAcceptance struct {
	Scheme      string `json:"scheme"`
	Network     string `json:"network"`
	Price       string `json:"price"`
	PayTo       string `json:"pay_to"`
	Resource    string `json:"resource"`
	Description string `json:"description"`
}

type Verification struct {
	Valid bool
}

type Settlement struct {
	Settled       bool
	TransactionID string
}

func PaymentRequired(cfg Config, requestID string, resource string) PaymentChallenge {
	return PaymentChallenge{
		X402Version: "integration-target",
		Mode:        "real-x402-go-skeleton",
		RequestID:   requestID,
		Reason:      "Payment required before the prompt safety report is generated.",
		Accepts: []PaymentAcceptance{
			{
				Scheme:      cfg.Scheme,
				Network:     cfg.Network,
				Price:       "$" + cfg.PriceUSD,
				PayTo:       cfg.MerchantAddress,
				Resource:    resource,
				Description: "One prompt guardrail safety check",
			},
		},
		RetryWithHeader: "PAYMENT-SIGNATURE",
		ProductionNote:  "Wire this adapter to github.com/x402-foundation/x402/go/v2 facilitator verify and settle.",
	}
}

func Verify(_ Config, _ string, paymentPayload string) (Verification, error) {
	if paymentPayload == "" {
		return Verification{Valid: false}, errors.New("missing payment payload")
	}

	// Replace this with official x402 facilitator /verify using the Go SDK.
	return Verification{Valid: true}, nil
}

func Settle(_ Config, _ string, paymentPayload string) (Settlement, error) {
	if paymentPayload == "" {
		return Settlement{Settled: false}, errors.New("missing payment payload")
	}

	// Replace this with official x402 facilitator /settle using the Go SDK.
	return Settlement{
		Settled:       true,
		TransactionID: "replace-with-facilitator-settlement-transaction",
	}, nil
}

func NewID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "id-unavailable"
	}
	return hex.EncodeToString(bytes[:])
}
