package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	payperprompt "github.com/icreate/payperprompt-x402-sdk-go"
)

func main() {
	baseURL := flag.String("url", "http://127.0.0.1:8084", "PayPerPrompt gateway URL")
	wallet := flag.String("wallet", "", "optional payer address for audit filtering")
	history := flag.Bool("history", false, "show up to 100 historical audit events")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client := payperprompt.New(*baseURL)
	services, err := client.Services(ctx)
	if err != nil {
		log.Fatal(err)
	}
	events, err := client.WorkAudit(ctx, *wallet, *history)
	if err != nil {
		log.Fatal(err)
	}
	payload, _ := json.MarshalIndent(map[string]any{
		"services": services,
		"audit":    events,
	}, "", "  ")
	fmt.Println(string(payload))
}
