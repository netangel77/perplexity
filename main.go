package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

//go:embed index.html
var indexHTML string

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Key      string    `json:"key"`
	Messages []Message `json:"messages"`
	Model    string    `json:"model,omitempty"`
}

type pplxResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Citations []string `json:"citations,omitempty"`
	Usage     any      `json:"usage,omitempty"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, indexHTML)
	})

	mux.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		req.Key = strings.TrimSpace(req.Key)
		if req.Key == "" {
			http.Error(w, "missing api key", http.StatusBadRequest)
			return
		}

		model := strings.TrimSpace(req.Model)
		if model == "" {
			model = "sonar-small-online"
		}

		answer, citations, usage, err := callPerplexity(r.Context(), req.Key, model, req.Messages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		resp := map[string]any{
			"answer":    answer,
			"citations": citations,
			"usage":     usage,
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(resp)
	})

	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// Bind to a free local port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	addr := "http://" + ln.Addr().String()

	log.Printf("Perplexity GUI: %s\n", addr)
	_ = openBrowser(addr) // if it fails, just open the printed URL manually

	log.Fatal(srv.Serve(ln))
}

func callPerplexity(ctx context.Context, key, model string, messages []Message) (string, []string, any, error) {
	apiReq := map[string]any{
		"model":    model,
		"messages": messages,
		"stream":   false,
	}

	body, _ := json.Marshal(apiReq)

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.perplexity.ai/chat/completions", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("X-Subscription-Token", key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		// Return body for easier debugging.
		return "", nil, nil, fmt.Errorf("api error %d: %s", resp.StatusCode, string(respBody))
	}

	var pr pplxResponse
	if err := json.Unmarshal(respBody, &pr); err != nil {
		return "", nil, nil, fmt.Errorf("bad api json: %w", err)
	}
	if len(pr.Choices) == 0 {
		return "", pr.Citations, pr.Usage, fmt.Errorf("empty response")
	}
	return pr.Choices[0].Message.Content, pr.Citations, pr.Usage, nil
}

func openBrowser(url string) error {
	// Common approach: Windows uses "rundll32 url.dll,FileProtocolHandler <url>". [web:89][web:99]
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
