package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

func main() {
	apiKey := os.Getenv("PPLX_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: Set PPLX_API_KEY=pplx-... (full key)")
		os.Exit(1)
	}

	reqBody := Request{
		Model: "sonar",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "What are the top AI trends in 2026?"},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("JSON error: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.perplexity.ai/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-Subscription-Token", apiKey) // Often required alongside Bearer
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("❌ API Error %d: %s\n", resp.StatusCode, string(body))
		fmt.Println("\n🔍 Debug tips:")
		fmt.Println("  - Key format? Must start with 'pplx-' (check Settings > API)")
		fmt.Println("  - Credits added? (Settings > API > Add credits)")
		fmt.Println("  - Pro active? API needs Pro + credits")
		os.Exit(1)
	}

	fmt.Println("✅ Success!")
	fmt.Printf("Response: %s\n", string(body))
}
