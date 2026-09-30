# Go Technical Specification & Package Layout

## 1. Standard Project Layout

`cli-leader` is organized as a modular, pure-Go application adhering to standard Go project conventions:

```
cli-leader/
├── cmd/
│   └── cli-leader/
│       └── main.go                 # Entrypoint (CLI flag parsing & daemon launch)
├── internal/
│   ├── config/
│   │   └── config.go               # YAML/ENV configuration loader
│   ├── gateway/
│   │   ├── proxy.go                # HTTP reverse proxy server (:8082)
│   │   ├── stream.go               # SSE stream state machine (OpenAI -> Anthropic)
│   │   ├── tokens.go               # /v1/messages/count_tokens handler (tiktoken-go)
│   │   ├── tools.go                # Tool schema conversion & name sanitization
│   │   └── providers/
│   │       ├── openai.go           # OpenAI payload adapter & reasoning mapper
│   │       ├── gemini.go           # Gemini OpenAI-compatible / REST adapter
│   │       └── ollama.go           # Ollama / local model adapter
│   ├── state/
│   │   ├── dual_ledger.go          # TaskLedger (goals) & ProgressLedger (steps)
│   │   └── stall_detector.go       # Cryptographic state hashing & loop breaker
│   ├── mcp/
│   │   ├── server.go               # Model Context Protocol stdio/SSE server
│   │   ├── tools.go                # Tool registration & schemas
│   │   └── handlers.go             # Tool invocation dispatch handlers
│   ├── worker/
│   │   ├── manager.go              # Worker lifecycle, pooling & supervisor
│   │   ├── process.go              # Subprocess execution, PTY & timeouts
│   │   ├── ringbuffer.go           # High-throughput circular log buffer
│   │   ├── auto_reply.go           # Sliding-window regex prompt auto-responder
│   │   ├── worktree.go             # Native Git worktree wrapper & mutex locks
│   │   └── adapters/
│   │       ├── adapter.go          # CLI Adapter interface
│   │       ├── aider.go            # Aider CLI integration
│   │       ├── gemini_cli.go       # Gemini CLI integration
│   │       ├── ollama_cli.go       # Ollama / local CLI integration
│   │       └── generic.go          # Custom generic CLI adapter
│   ├── speculative/
│   │   └── branch_racing.go        # Parallel candidate dispatch & diff evaluation
│   ├── refinery/
│   │   └── queue.go                # Serialized Bors-style merge queue & test runner
│   ├── verification/
│   │   ├── f2p.go                  # Fail-to-Pass automated TDD engine
│   │   └── red_blue.go             # Adversarial breaker audit loop
│   ├── sandbox/
│   │   ├── bwrap.go                # Bubblewrap unprivileged namespace wrapper
│   │   └── landlock.go             # Linux Landlock LSM re-exec trampoline
│   ├── memory/
│   │   ├── store.go                # SQLite episodic and semantic storage
│   │   ├── vector.go               # Pure-Go chromem-go vector search engine
│   │   ├── reflexion.go            # Verbal reinforcement rule extractor
│   │   ├── rules.go                # Dynamic .cursor/rules/*.mdc glob matcher
│   │   ├── matrix.go               # Thompson Sampling Bayesian capability routing
│   │   └── schema.sql              # Database DDL
│   └── tui/
│       ├── app.go                  # Bubble Tea main application model
│       ├── batcher.go              # 30Hz ticker batcher for smooth 60 FPS logs
│       ├── views/
│       │   ├── dashboard.go        # Swarm status & pipeline graph
│       │   ├── logs.go             # Multi-pane real-time log viewer
│       │   └── memory_view.go      # Knowledge & rule explorer
│       └── styles/
│           └── theme.go            # Lipgloss styling, borders, and colors
├── pkg/
│   └── protocol/
│       ├── anthropic.go            # Anthropic API message structs & SSE types
│       ├── openai.go               # OpenAI API message structs & chunks
│       └── gemini.go               # Gemini API message structs
├── docs/                           # Architecture, Specs, Roadmap
├── go.mod
├── go.sum
└── README.md
```

---

## 2. Core Go Interfaces & Structs

### 2.1 Dual-Ledger State Machine
```go
package state

import (
	"context"
	"time"
)

type StepStatus string

const (
	StepPending   StepStatus = "pending"
	StepRunning   StepStatus = "running"
	StepVerifying StepStatus = "verifying"
	StepCompleted StepStatus = "completed"
	StepFailed    StepStatus = "failed"
	StepAborted   StepStatus = "aborted"
)

// TaskLedger stores the high-level immutable user goal and constraints
type TaskLedger struct {
	ID                 string    `json:"id"`
	Goal               string    `json:"goal"`
	AcceptanceCriteria []string  `json:"acceptance_criteria"`
	BaseBranch         string    `json:"base_branch"`
	CreatedAt          time.Time `json:"created_at"`
}

// ProgressLedger tracks dynamic execution steps and breaks thrashing loops
type ProgressStep struct {
	ID             string     `json:"id"`
	TaskID         string     `json:"task_id"`
	StepIndex      int        `json:"step_index"`
	AssignedWorker string     `json:"assigned_worker"`
	ActionPrompt   string     `json:"action_prompt"`
	WorktreePath   string     `json:"worktree_path"`
	Branch         string     `json:"branch"`
	Status         StepStatus `json:"status"`
	StateHash      string     `json:"state_hash"` // SHA256(tool_name + args + diff)
	RetryCount     int        `json:"retry_count"`
	StallCount     int        `json:"stall_count"`
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

type LedgerEngine interface {
	CreateTask(ctx context.Context, task TaskLedger) error
	RecordStep(ctx context.Context, step ProgressStep) error
	CheckLoopOrStall(ctx context.Context, taskID, stateHash string) (isLoop bool, stallCount int, err error)
	MarkStepComplete(ctx context.Context, stepID string) error
}
```

### 2.2 Streaming SSE State Machine
```go
package gateway

import (
	"net/http"
)

type BlockType int

const (
	BlockNone BlockType = iota
	BlockThinking
	BlockText
	BlockToolUse
)

// SSEStreamTranslator converts upstream OpenAI chunk deltas into framed Anthropic SSE envelopes
type SSEStreamTranslator struct {
	w                  http.ResponseWriter
	flusher            http.Flusher
	currentBlockIndex  int
	currentBlockType   BlockType
	openAIToBlockIndex map[int]int // Maps OpenAI tool_call index -> Anthropic block index
	messageID          string
	model              string
}

func NewSSEStreamTranslator(w http.ResponseWriter, flusher http.Flusher, msgID, model string) *SSEStreamTranslator {
	return &SSEStreamTranslator{
		w:                  w,
		flusher:            flusher,
		openAIToBlockIndex: make(map[int]int),
		messageID:          msgID,
		model:              model,
	}
}

func (t *SSEStreamTranslator) HandleOpenAIChunk(chunk *OpenAIStreamChunk) error {
	// Translates reasoning_content -> thinking blocks, content -> text blocks,
	// and tool_calls -> input_json_delta blocks. Flushes immediately.
	return nil
}
```

### 2.3 Thompson Sampling Capability Matrix
```go
package memory

import (
	"context"
	"math/rand"
	"time"
)

type WorkerScore struct {
	CLIName     string  `json:"cli_name"`
	Category    string  `json:"category"`
	AlphaScore  float64 `json:"alpha_score"` // 1 + successes (Beta distribution)
	BetaScore   float64 `json:"beta_score"`  // 1 + failures (Beta distribution)
	AvgDuration time.Duration `json:"avg_duration"`
	TotalCost   float64 `json:"total_cost_usd"`
}

type CapabilityRouter interface {
	// SampleWorker draws from Beta(alpha, beta) for each candidate and picks the highest sample
	SampleWorker(ctx context.Context, category string, candidates []string) (string, error)
	RecordSuccess(ctx context.Context, cliName, category string, duration time.Duration, cost float64) error
	RecordFailure(ctx context.Context, cliName, category string) error
}
```

### 2.4 Decoupled High-Throughput TUI Ring Buffer
```go
package worker

import (
	"sync"
)

type RingLogBuffer struct {
	mu       sync.Mutex
	lines    []string
	maxLines int
	dirty    bool
}

func NewRingLogBuffer(maxLines int) *RingLogBuffer {
	return &RingLogBuffer{
		lines:    make([]string, 0, maxLines),
		maxLines: maxLines,
	}
}

func (b *RingLogBuffer) Append(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.lines) >= b.maxLines {
		b.lines = b.lines[1:]
	}
	b.lines = append(b.lines, line)
	b.dirty = true
}

func (b *RingLogBuffer) DrainNew(fromIndex int) ([]string, int, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.dirty || fromIndex >= len(b.lines) {
		return nil, len(b.lines), false
	}
	newLines := make([]string, len(b.lines)-fromIndex)
	copy(newLines, b.lines[fromIndex:])
	return newLines, len(b.lines), true
}
```

---

## 3. Configuration Schema (`cli-leader.yaml`)

```yaml
version: "1.0"

# Multi-Provider Gateway Settings
gateway:
  listen_addr: "127.0.0.1:8082"
  default_provider: "openai"
  providers:
    openai:
      base_url: "https://api.openai.com/v1"
      api_key_env: "OPENAI_API_KEY"
      model_mapping:
        "claude-3-7-sonnet-20250219": "o3-mini"
        "claude-3-5-sonnet-20241022": "gpt-4o"
    gemini:
      base_url: "https://generativelanguage.googleapis.com/v1beta/openai"
      api_key_env: "GEMINI_API_KEY"
      model_mapping:
        "claude-3-7-sonnet-20250219": "gemini-2.5-pro"
    ollama:
      base_url: "http://localhost:11434/v1"
      model_mapping:
        "claude-3-5-sonnet-20241022": "deepseek-r1:32b"

# Worker Swarm Topology
workers:
  aider:
    binary: "aider"
    flags: ["--no-auto-commits", "--yes-always"]
    categories: ["refactoring", "code_edit"]
  gemini:
    binary: "gemini"
    categories: ["indexing", "search"]
  ollama:
    binary: "ollama"
    categories: ["test_generation", "boilerplate"]

# Speculative Execution & Merge Refinery
speculative:
  enabled: true
  max_parallel_candidates: 2

refinery:
  test_command: "go test -v ./..."
  squash_merges: true

# Cognitive Memory & Reflexion
memory:
  db_path: ".brain/memory.db"
  conventions_dir: ".brain/conventions"
  max_injected_rule_tokens: 800
  thompson_sampling: true

# Verification & Sandboxing
sandbox:
  tier: "auto" # options: bwrap, landlock, none
  allow_network: false
```

---

## 4. Key Go Dependencies

| Package | Purpose |
|---|---|
| `modernc.org/sqlite` | Pure Go embedded SQLite engine (zero CGO required) |
| `github.com/philippgille/chromem-go` | Pure Go in-memory vector store with persistence (<2ms search) |
| `github.com/charmbracelet/bubbletea` | Elm-architecture Terminal User Interface framework |
| `github.com/charmbracelet/lipgloss` | Terminal styling, borders, and layouts |
| `github.com/creack/pty` | PTY allocation for interactive CLI worker wrapping |
| `github.com/pkoukk/tiktoken-go` | Pure Go BPE tokenization for `/v1/messages/count_tokens` |
| `github.com/tidwall/gjson` & `sjson` | Zero-allocation JSON manipulation for low-latency proxying |
| `github.com/bmatcuk/doublestar/v4` | High-performance glob matching for `.cursor/rules/*.mdc` rules |
| `github.com/landlock-lsm/go-landlock` | Unprivileged Linux Landlock LSM sandboxing |
| `github.com/spf13/cobra` | CLI command routing (`cli-leader start`, `status`) |
