package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"payperprompt-x402-real/internal/payperprompt"
	"payperprompt-x402-real/internal/x402adapter"
)

func main() {
	cfg := x402adapter.Config{
		FacilitatorURL: env("X402_FACILITATOR_URL", "https://x402.org/facilitator"),
		Network:        env("X402_NETWORK", "eip155:84532"),
		Scheme:         env("X402_SCHEME", "exact"),
		PriceUSD:       env("PAYPERPROMPT_PRICE_USD", "0.01"),
		Asset:          env("X402_ASSET", "USDC_TEST"),
		MerchantAddress: env(
			"MERCHANT_ADDRESS",
			"replace_with_testnet_receiving_wallet",
		),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":      true,
			"service": "PayPerPrompt real x402 Go lane",
			"network": cfg.Network,
		})
	})
	mux.HandleFunc("POST /api/check-prompt", checkPromptHandler(cfg))

	port := env("PAYPERPROMPT_PORT", "8081")
	log.Printf("PayPerPrompt real x402 lane listening at http://127.0.0.1:%s", port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, mux))
}

func checkPromptHandler(cfg x402adapter.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Prompt string `json:"prompt"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&input); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "valid JSON prompt is required"})
			return
		}

		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = x402adapter.NewID()
		}

		paymentPayload := r.Header.Get("PAYMENT-SIGNATURE")
		if paymentPayload == "" {
			challenge := x402adapter.PaymentRequired(cfg, requestID, "/api/check-prompt")
			writePaymentRequired(w, challenge)
			return
		}

		verification, err := x402adapter.Verify(cfg, requestID, paymentPayload)
		if err != nil || !verification.Valid {
			challenge := x402adapter.PaymentRequired(cfg, requestID, "/api/check-prompt")
			writePaymentRequired(w, challenge)
			return
		}

		report := payperprompt.Check(input.Prompt)
		settlement, err := x402adapter.Settle(cfg, requestID, paymentPayload)
		if err != nil || !settlement.Settled {
			writeJSON(w, http.StatusPaymentRequired, map[string]any{
				"error":      "payment settlement failed",
				"request_id": requestID,
			})
			return
		}

		receipt := map[string]any{
			"receipt_id":       x402adapter.NewID(),
			"request_id":       requestID,
			"network":          cfg.Network,
			"asset":            cfg.Asset,
			"amount_usd":       cfg.PriceUSD,
			"transaction_id":   settlement.TransactionID,
			"settled":          true,
			"replay_protected": true,
			"idempotency_key":  requestID + ":" + paymentPayload,
			"facilitator":      cfg.FacilitatorURL,
			"issued_at":        time.Now().UTC().Format(time.RFC3339Nano),
		}

		receiptJSON, _ := json.Marshal(receipt)
		w.Header().Set("PAYMENT-RESPONSE", base64.RawURLEncoding.EncodeToString(receiptJSON))
		writeJSON(w, http.StatusOK, map[string]any{
			"service":       "PayPerPrompt x402",
			"paid_resource": "/api/check-prompt",
			"report":        report,
			"receipt":       receipt,
		})
	}
}

func writePaymentRequired(w http.ResponseWriter, challenge x402adapter.PaymentChallenge) {
	payload, _ := json.Marshal(challenge)
	w.Header().Set("PAYMENT-REQUIRED", base64.RawURLEncoding.EncodeToString(payload))
	writeJSON(w, http.StatusPaymentRequired, challenge)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	payload, _ := json.MarshalIndent(body, "", "  ")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
