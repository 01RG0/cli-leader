package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Drex API Constants
const (
	DefaultBaseURL = "https://drex.nace.ai"
	DefaultModel   = "drex-v1.5"
	DefaultTimeout = 60 * time.Second
)

// Question represents a typed Drex question (noul, choice, or score).
type Question struct {
	Type         string `json:"type"`
	Instructions string `json:"instructions"`
	Criteria     any    `json:"criteria,omitempty"`
}

// SystemOneRequest represents a request to POST /v1/systemone.
type SystemOneRequest struct {
	Model     string              `json:"model"`
	State     string              `json:"state"`
	Questions map[string]Question `json:"questions"`
}

// NoulAnswer represents the result of a noul question.
type NoulAnswer struct {
	Type string  `json:"type"`
	Noul float64 `json:"noul"`
}

// ChoiceAnswer represents the result of a choice question.
type ChoiceAnswer struct {
	Type          string             `json:"type"`
	Choice        string             `json:"choice"`
	Confidence    float64            `json:"confidence"`
	Probabilities map[string]float64 `json:"probabilities"`
}

// AnswerEnvelope handles dynamic parsing of answer objects.
type AnswerEnvelope struct {
	Type          string             `json:"type"`
	Noul          float64            `json:"noul,omitempty"`
	Choice        string             `json:"choice,omitempty"`
	Confidence    float64            `json:"confidence,omitempty"`
	Probabilities map[string]float64 `json:"probabilities,omitempty"`
	Score         float64            `json:"score,omitempty"`
}

// UsageInfo represents token usage.
type UsageInfo struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// SystemOneResponse represents the successful response from POST /v1/systemone.
type SystemOneResponse struct {
	Model            string                    `json:"model"`
	Answers          map[string]AnswerEnvelope `json:"answers"`
	Usage            UsageInfo                 `json:"usage"`
	EvaluationTimeMs float64                   `json:"evaluation_time_ms"`
	RequestID        string                    `json:"request_id"`
}

// Client wraps HTTP requests to Drex with caching, rate limiting, and retries.
type Client struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
	cacheDir   string
	sem        chan struct{} // Concurrency limiter
	cacheMu    sync.Mutex
}

// NewClient initializes a Drex API client.
func NewClient(apiKey, cacheDir string, concurrency int) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("DREX_API_KEY is required")
	}
	if concurrency <= 0 {
		concurrency = 4
	}

	if cacheDir != "" {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache dir: %w", err)
		}
	}

	return &Client{
		apiKey:     apiKey,
		baseURL:    DefaultBaseURL,
		model:      DefaultModel,
		httpClient: &http.Client{Timeout: DefaultTimeout},
		cacheDir:   cacheDir,
		sem:        make(chan struct{}, concurrency),
	}, nil
}

// LoadEnv attempts to load environment variables from a .env file.
func LoadEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

// RedactKey replaces occurrences of the secret API key in errors or text.
func (c *Client) RedactKey(s string) string {
	if c.apiKey != "" && len(c.apiKey) > 8 {
		return strings.ReplaceAll(s, c.apiKey, "[REDACTED_API_KEY]")
	}
	return s
}

// Evaluate sends a SystemOneRequest, utilizing cache if available, with retries and concurrency control.
func (c *Client) Evaluate(ctx context.Context, req SystemOneRequest) (*SystemOneResponse, bool, error) {
	if req.Model == "" {
		req.Model = c.model
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, false, err
	}

	// 1. Check cache
	cacheHash := sha256.Sum256(reqBody)
	cacheKey := hex.EncodeToString(cacheHash[:])
	var cacheFile string
	if c.cacheDir != "" {
		cacheFile = filepath.Join(c.cacheDir, cacheKey+".json")
		c.cacheMu.Lock()
		cachedData, err := os.ReadFile(cacheFile)
		c.cacheMu.Unlock()
		if err == nil {
			var cachedResp SystemOneResponse
			if err := json.Unmarshal(cachedData, &cachedResp); err == nil {
				return &cachedResp, true, nil // Cache hit
			}
		}
	}

	// 2. Acquire semaphore token
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return nil, false, ctx.Err()
	}

	// 3. Execute with retries & exponential backoff
	maxRetries := 4
	for attempt := 0; attempt <= maxRetries; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/systemone", bytes.NewReader(reqBody))
		if err != nil {
			return nil, false, err
		}
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			if attempt >= maxRetries || ctx.Err() != nil {
				return nil, false, errors.New(c.RedactKey(err.Error()))
			}
			time.Sleep(c.backoffDuration(attempt, nil))
			continue
		}

		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			if attempt >= maxRetries {
				return nil, false, errors.New(c.RedactKey(readErr.Error()))
			}
			time.Sleep(c.backoffDuration(attempt, nil))
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var sysResp SystemOneResponse
			if err := json.Unmarshal(respBody, &sysResp); err != nil {
				return nil, false, fmt.Errorf("failed to parse Drex response: %w", err)
			}

			// Store in cache
			if c.cacheDir != "" {
				c.cacheMu.Lock()
				_ = os.WriteFile(cacheFile, respBody, 0644)
				c.cacheMu.Unlock()
			}

			return &sysResp, false, nil
		}

		// Non-retryable errors
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusPaymentRequired || resp.StatusCode == 422 {
			return nil, false, fmt.Errorf("Drex %d (%s): %s", resp.StatusCode, resp.Header.Get("x-request-id"), c.RedactKey(string(respBody)))
		}

		// Retryable: 429, 500, 529
		if attempt >= maxRetries {
			return nil, false, fmt.Errorf("Drex %d (%s) exhausted retries: %s", resp.StatusCode, resp.Header.Get("x-request-id"), c.RedactKey(string(respBody)))
		}

		time.Sleep(c.backoffDuration(attempt, resp.Header))
	}

	return nil, false, errors.New("request failed after retries")
}

func (c *Client) backoffDuration(attempt int, headers http.Header) time.Duration {
	if headers != nil {
		if msStr := headers.Get("retry-after-ms"); msStr != "" {
			if ms, err := strconv.Atoi(msStr); err == nil && ms > 0 {
				return time.Duration(ms) * time.Millisecond
			}
		}
		if sStr := headers.Get("retry-after"); sStr != "" {
			if s, err := strconv.Atoi(sStr); err == nil && s > 0 {
				return time.Duration(s) * time.Second
			}
		}
	}

	// Exponential backoff: min 500ms, cap at 10s with 50% random jitter
	base := 500 * (1 << attempt)
	if base > 10000 {
		base = 10000
	}
	jitter := rand.Intn(base / 2)
	return time.Duration(base/2+jitter) * time.Millisecond
}

// Standard Section Questions
func SectionQuestions() map[string]Question {
	return map[string]Question{
		"makes_unmeasured_performance_claim": {
			Type:         "noul",
			Instructions: "Does this text make a performance, speed, latency, or throughput claim without citing a specific measurement, benchmark, or proof?",
		},
		"states_feature_without_implementation_detail": {
			Type:         "noul",
			Instructions: "Does this text state or promise a technical feature without describing any implementation detail, mechanism, or algorithm?",
		},
		"uses_borrowed_term_loosely": {
			Type:         "noul",
			Instructions: "Does this text use an external technical framework or domain term loosely without grounding it in the actual system architecture?",
		},
		"contains_marketing_language": {
			Type:         "noul",
			Instructions: "Does this text contain promotional, exaggerated, or marketing buzzwords rather than objective technical description?",
		},
		"section_role": {
			Type:         "choice",
			Instructions: "What is the primary role of this documentation section?",
			Criteria: map[string]any{
				"vision":         "High-level mission, philosophy, or future ambitions",
				"design":         "System architecture, internal mechanisms, or component interactions",
				"interface_spec": "API contracts, schemas, structs, or protocol specifications",
				"roadmap":        "Milestones, phases, delivery plans, or timelines",
				"research":       "Survey of external projects, comparative analysis, or background theory",
				"other":          "Miscellaneous information, contributing guides, or general notes",
			},
		},
	}
}

// Standard Cross-File Pair Questions
func PairQuestions() map[string]Question {
	return map[string]Question{
		"sections_contradict": {
			Type:         "noul",
			Instructions: "Do Section A and Section B make contradictory, conflicting, or incompatible claims with each other?",
		},
		"sections_duplicate_each_other": {
			Type:         "noul",
			Instructions: "Do Section A and Section B substantially duplicate the same information, explanations, or specifications?",
		},
	}
}
