package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	x402 "github.com/x402-foundation/x402/go/v2"
	x402http "github.com/x402-foundation/x402/go/v2/http"
	nethttpmw "github.com/x402-foundation/x402/go/v2/http/nethttp"
	evm "github.com/x402-foundation/x402/go/v2/mechanisms/evm/exact/server"

	"payperprompt-x402-real/internal/facilitatorpool"
	"payperprompt-x402-real/internal/payperprompt"
)

func main() {
	cfg := config{
		Port:            env("PAYPERPROMPT_OFFICIAL_PORT", "8082"),
		PriceUSD:        env("PAYPERPROMPT_PRICE_USD", "0.01"),
		FacilitatorURLs: facilitatorURLs(),
		Network:         env("X402_NETWORK", "eip155:84532"),
		OllamaURL:       env("OLLAMA_URL", "http://127.0.0.1:11434"),
		OllamaModel:     env("OLLAMA_MODEL", "llama3.1:8b"),
		WebDir:          env("PAYPERPROMPT_OFFICIAL_WEB_DIR", "official-web"),
		ProofPath:       env("OFFICIAL_PROOF_PATH", "proof/official-settlement.json"),
		MerchantAddress: env(
			"MERCHANT_ADDRESS",
			"replace_with_testnet_receiving_wallet",
		),
	}

	if cfg.MerchantAddress == "" || cfg.MerchantAddress == "replace_with_testnet_receiving_wallet" {
		log.Fatal("MERCHANT_ADDRESS must be set to a public EVM testnet receiving address")
	}
	if !common.IsHexAddress(cfg.MerchantAddress) {
		log.Fatal("MERCHANT_ADDRESS is not a valid EVM address")
	}

	facilitatorSpecs := make([]facilitatorpool.ClientSpec, 0, len(cfg.FacilitatorURLs))
	for i, facilitatorURL := range cfg.FacilitatorURLs {
		facilitatorSpecs = append(facilitatorSpecs, facilitatorpool.ClientSpec{
			Name: fmt.Sprintf("facilitator-%d", i+1),
			URL:  facilitatorURL,
			Client: x402http.NewHTTPFacilitatorClient(&x402http.FacilitatorConfig{
				URL: facilitatorURL,
			}),
		})
	}
	facilitatorClient, err := facilitatorpool.New(facilitatorSpecs)
	if err != nil {
		log.Fatalf("invalid facilitator pool: %v", err)
	}
	analyzer := payperprompt.NewAnalyzer(cfg.OllamaURL, cfg.OllamaModel)
	worker := payperprompt.NewWorker(cfg.OllamaURL, cfg.OllamaModel)
	services := officialServices(cfg)

	routes := x402http.RoutesConfig{
		"POST /api/check-prompt": {
			Accepts: x402http.PaymentOptions{
				{
					Scheme:  "exact",
					Price:   "$" + cfg.PriceUSD,
					Network: x402.Network(cfg.Network),
					PayTo:   cfg.MerchantAddress,
				},
			},
			Description: "One standard AI work request with prompt-safety analysis",
			MimeType:    "application/json",
		},
		"POST /api/services/rapid-policy/check-prompt": {
			Accepts: x402http.PaymentOptions{
				{
					Scheme:  "exact",
					Price:   "$0.02",
					Network: x402.Network(cfg.Network),
					PayTo:   cfg.MerchantAddress,
				},
			},
			Description: "One priority AI work request with prompt-safety analysis",
			MimeType:    "application/json",
		},
		"POST /api/services/deep-shield/check-prompt": {
			Accepts: x402http.PaymentOptions{
				{
					Scheme:  "exact",
					Price:   "$0.04",
					Network: x402.Network(cfg.Network),
					PayTo:   cfg.MerchantAddress,
				},
			},
			Description: "One enhanced AI work request with prompt-safety analysis",
			MimeType:    "application/json",
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(cfg.WebDir, "index.html"))
	})
	mux.HandleFunc("GET /static/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" || filepath.Base(name) != name {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(cfg.WebDir, name))
	})
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":           true,
			"service":      "PayPerPrompt official x402 AI lane",
			"mode":         "official-base-sepolia",
			"network":      cfg.Network,
			"facilitators": facilitatorClient.Snapshot(),
			"merchant":     cfg.MerchantAddress,
			"services":     services,
			"ai_model":     cfg.OllamaModel,
			"ollama_ready": analyzer.Health(ctx),
		})
	})
	mux.HandleFunc("GET /api/facilitators/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		probe := r.URL.Query().Get("probe") == "1"
		probeOK := true
		probeError := ""
		if probe {
			ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
			err := facilitatorClient.Probe(ctx)
			cancel()
			if err != nil {
				probeOK = false
				probeError = err.Error()
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"probe_requested": probe,
			"probe_ok":        probeOK,
			"probe_error":     probeError,
			"payment_signed":  false,
			"payment_sent":    false,
			"pool":            facilitatorClient.Snapshot(),
		})
	})
	mux.HandleFunc("GET /api/services", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, http.StatusOK, map[string]any{
			"network":  cfg.Network,
			"merchant": cfg.MerchantAddress,
			"services": services,
		})
	})
	mux.HandleFunc("POST /api/work-preflight", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Prompt   string `json:"prompt"`
			TaskType string `json:"task_type"`
			Quality  string `json:"quality"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&input); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "valid JSON work request is required"})
			return
		}
		input.Prompt = strings.TrimSpace(input.Prompt)
		if input.Prompt == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "work request is required"})
			return
		}
		work, err := worker.Complete(r.Context(), input.TaskType, input.Prompt, input.Quality)
		if err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{
				"work_completed": false,
				"payment_signed": false,
				"payment_sent":   false,
				"error":          err.Error(),
			})
			return
		}
		commitment, err := payperprompt.WorkProductCommitment(work)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		canonicalWork, err := json.Marshal(work)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"work_completed":                true,
			"work":                          work,
			"work_canonical_base64":         base64.StdEncoding.EncodeToString(canonicalWork),
			"deliverable_commitment_sha256": commitment,
			"payment_signed":                false,
			"payment_sent":                  false,
		})
	})
	mux.HandleFunc("GET /api/proof", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		payload, err := os.ReadFile(cfg.ProofPath)
		if err != nil {
			if os.IsNotExist(err) {
				writeJSON(w, http.StatusNotFound, map[string]any{
					"verified": false,
					"reason":   "No official settlement proof has been generated yet.",
				})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if !json.Valid(payload) {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "proof file is not valid JSON"})
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(payload)
	})
	for _, configuredService := range services {
		service := configuredService
		mux.HandleFunc("POST "+service.Path, func(w http.ResponseWriter, r *http.Request) {
			var input struct {
				Prompt                         string                     `json:"prompt"`
				TaskType                       string                     `json:"task_type"`
				PreparedWork                   *payperprompt.WorkProduct  `json:"prepared_work,omitempty"`
				PreparedAnalysis               *payperprompt.Analysis     `json:"prepared_analysis,omitempty"`
				PreparedWorkCommitmentSHA256   string                     `json:"prepared_work_commitment_sha256,omitempty"`
			}
			if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&input); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "valid JSON prompt is required"})
				return
			}

			taskType := payperprompt.NormalizeTaskType(input.TaskType, input.Prompt)
			report := payperprompt.Analysis{}
			aiUsed := false
			var aiErr error
			if input.PreparedAnalysis != nil {
				report = *input.PreparedAnalysis
				report.TaskType = taskType
				report.Strategy = strategyForPaidService(service)
				aiUsed = true
			} else {
				report, aiUsed, aiErr = analyzer.Analyze(r.Context(), input.Prompt)
				report.TaskType = taskType
				report.Strategy = strategyForPaidService(service)
			}
			response := map[string]any{
				"service":       "PayPerPrompt x402",
				"paid_resource": service.Path,
				"provider":      service.Provider,
				"route_id":      service.RouteID,
				"quality":       service.Quality,
				"price_usd":     service.PriceUSD,
				"report":        report,
				"ai_used":       aiUsed,
				"ai_model":      cfg.OllamaModel,
				"settlement":    "handled by official x402 net/http middleware",
				"issued_at":     time.Now().UTC().Format(time.RFC3339Nano),
			}
			if aiErr != nil {
				response["ai_fallback_reason"] = aiErr.Error()
			}
			if !aiUsed {
				response["work_completed"] = false
				response["work_error"] = "Local AI was unavailable after settlement; the safety fallback report is included."
				writeJSON(w, http.StatusServiceUnavailable, response)
				return
			}
			var work payperprompt.WorkProduct
			var workErr error
			if input.PreparedWork != nil {
				work, workErr = payperprompt.ValidatePreparedWork(*input.PreparedWork, taskType, input.Prompt)
				if workErr == nil {
					var commitment string
					commitment, workErr = payperprompt.WorkProductCommitment(work)
					if workErr == nil && !strings.EqualFold(commitment, strings.TrimSpace(input.PreparedWorkCommitmentSHA256)) {
						workErr = fmt.Errorf("prepared work commitment does not match")
					}
					if workErr == nil {
						response["prepared_work_released"] = true
						response["deliverable_commitment_sha256"] = commitment
					}
				}
			} else {
				work, workErr = worker.Complete(r.Context(), taskType, input.Prompt, service.Quality)
			}
			if workErr != nil {
				response["work_completed"] = false
				response["work_error"] = workErr.Error()
				writeJSON(w, http.StatusServiceUnavailable, response)
				return
			}
			response["task_type"] = taskType
			response["work_completed"] = true
			response["work"] = work
			writeJSON(w, http.StatusOK, response)
		})
	}

	handler := nethttpmw.X402Payment(nethttpmw.Config{
		Routes:      routes,
		Facilitator: facilitatorClient,
		Schemes: []nethttpmw.SchemeConfig{
			{Network: x402.Network(cfg.Network), Server: evm.NewExactEvmScheme()},
		},
		Timeout: 150 * time.Second,
	})(mux)

	fmt.Printf("PayPerPrompt official x402 lane listening at http://127.0.0.1:%s\n", cfg.Port)
	fmt.Printf("Network: %s\n", cfg.Network)
	fmt.Printf("Facilitators: %s\n", strings.Join(cfg.FacilitatorURLs, " -> "))
	fmt.Println("Settlement retry: disabled after any ambiguous facilitator response")
	fmt.Printf("PayTo: %s\n", cfg.MerchantAddress)
	fmt.Printf("AI: %s via %s\n", cfg.OllamaModel, cfg.OllamaURL)

	log.Fatal(http.ListenAndServe("127.0.0.1:"+cfg.Port, securityHeaders(handler)))
}

type config struct {
	Port            string
	PriceUSD        string
	FacilitatorURLs []string
	Network         string
	MerchantAddress string
	OllamaURL       string
	OllamaModel     string
	WebDir          string
	ProofPath       string
}

type paidService struct {
	RouteID     string `json:"route_id"`
	Provider    string `json:"provider"`
	Path        string `json:"path"`
	PriceUSD    string `json:"price_usd"`
	LatencyMS   int    `json:"expected_latency_ms"`
	Quality     string `json:"quality"`
	Description string `json:"description"`
}

func officialServices(cfg config) []paidService {
	return []paidService{
		{
			RouteID: "guardrail-economy", Provider: "Local Guard",
			Path: "/api/check-prompt", PriceUSD: cfg.PriceUSD,
			LatencyMS: 180, Quality: "standard",
			Description: "One standard AI work request with prompt-safety analysis",
		},
		{
			RouteID: "guardrail-fast", Provider: "Rapid Policy",
			Path: "/api/services/rapid-policy/check-prompt", PriceUSD: "0.02",
			LatencyMS: 70, Quality: "standard",
			Description: "One priority AI work request with prompt-safety analysis",
		},
		{
			RouteID: "guardrail-deep", Provider: "Deep Shield",
			Path: "/api/services/deep-shield/check-prompt", PriceUSD: "0.04",
			LatencyMS: 420, Quality: "enhanced",
			Description: "One enhanced AI work request with prompt-safety analysis",
		},
	}
}

func strategyForPaidService(service paidService) string {
	switch service.RouteID {
	case "guardrail-deep":
		return "highest-quality"
	case "guardrail-fast":
		return "lowest-latency"
	default:
		return "lowest-cost"
	}
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

func facilitatorURLs() []string {
	raw := env("X402_FACILITATOR_URLS", env("X402_FACILITATOR_URL", "https://x402.org/facilitator"))
	seen := map[string]bool{}
	urls := make([]string, 0, 2)
	for _, value := range strings.Split(raw, ",") {
		value = strings.TrimRight(strings.TrimSpace(value), "/")
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		urls = append(urls, value)
	}
	if len(urls) == 0 {
		return []string{"https://x402.org/facilitator"}
	}
	return urls
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; connect-src 'self'; style-src 'self'; script-src 'self'; base-uri 'none'; frame-ancestors 'none'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}
