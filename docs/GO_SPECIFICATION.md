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
│   │   ├── proxy.go                # HTTP reverse proxy server (:8082, FastHTTP)
│   │   ├── stream.go               # SSE stream state machine (OpenAI -> Anthropic)
│   │   ├── tokens.go               # /v1/messages/count_tokens handler (tiktoken-go)
│   │   ├── tools.go                # Tool schema conversion & name sanitization
│   │   └── providers/
│   │       ├── openai.go           # OpenAI payload adapter & reasoning mapper
│   │       ├── gemini.go           # Gemini OpenAI-compatible / REST adapter
│   │       └── ollama.go           # Ollama / local model adapter
│   ├── state/
│   │   ├── dual_ledger.go          # TaskLedger (goals) & ProgressLedger (steps)
│   │   ├── durable_engine.go       # Temporal-grade SQLite WAL replay cache
│   │   └── stall_detector.go       # Cryptographic state hashing & loop breaker
│   ├── mcp/
│   │   ├── server.go               # Model Context Protocol stdio/SSE server
│   │   ├── tools.go                # Tool registration & schemas
│   │   └── handlers.go             # Tool invocation dispatch handlers
│   ├── acp/
│   │   ├── server.go               # Agent Client Protocol (Zed / JetBrains JSON-RPC)
│   │   └── session.go              # Editor buffer and thread state management
│   ├── supervisor/
│   │   ├── tree.go                 # Erlang OTP supervisor (one_for_one, one_for_all)
│   │   ├── process.go              # Fate-sharing worker bundle lifecycle
│   │   └── intensity.go            # Restart rate limiter & Brain escalation
│   ├── worker/
│   │   ├── manager.go              # Worker pooling & task assignment
│   │   ├── pty.go                  # Subprocess execution, creack/pty & Pdeathsig
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
│   │   ├── queue.go                # Serialized Bors-style merge queue & CI runner
│   │   └── saga.go                 # LIFO compensating non-git rollback coordinator
│   ├── verification/
│   │   ├── f2p.go                  # Fail-to-Pass automated TDD engine
│   │   ├── mutation.go             # AST mutation testing gate (cargo-mutants style)
│   │   └── consensus.go            # Byzantine quorum consensus (Thompson-weighted)
│   ├── sandbox/
│   │   ├── bwrap.go                # Bubblewrap unprivileged namespace wrapper
│   │   └── landlock.go             # Linux Landlock LSM re-exec trampoline
│   ├── memory/
│   │   ├── store.go                # SQLite episodic and semantic storage
│   │   ├── vector.go               # Pure-Go chromem-go vector search engine
│   │   ├── repomap.go              # Tree-Sitter Personalized PageRank symbol map
│   │   ├── reflexion.go            # Verbal reinforcement rule extractor
│   │   ├── rules.go                # Dynamic .cursor/rules/*.mdc glob matcher
│   │   ├── skills.go               # Hermes self-evolving skill catalog
│   │   ├── matrix.go               # Thompson Sampling Bayesian capability routing
│   │   └── schema.sql              # Database DDL
│   ├── notify/
│   │   └── webhook.go              # OpenClaw-style notifications (Slack/Discord)
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
│       └── acp.go                  # Agent Client Protocol message envelopes
├── docs/                           # Architecture, Specs, Innovations, Roadmap
├── go.mod
├── go.sum
└── README.md
```

---

## 2. Core Go Interfaces & Structs

### 2.1 Erlang OTP Supervisor Tree
```go
package supervisor

import (
	"context"
	"sync"
	"syscall"
	"time"
)

type RestartStrategy string

const (
	OneForOne  RestartStrategy = "one_for_one"  // Candidate workers isolated
	OneForAll  RestartStrategy = "one_for_all"  // Compound worker bundles (PTY + reader)
	RestForOne RestartStrategy = "rest_for_one" // Linear verification pipelines
)

type ProcessDownMsg struct {
	WorkerID  string
	ExitCode  int
	Err       error
	Timestamp time.Time
}

type Supervisor struct {
	mu           sync.Mutex
	strategy     RestartStrategy
	maxIntensity int           // Max failures allowed
	period       time.Duration // Within duration
	failures     []time.Time
	children     map[string]*WorkerProcess
	downChan     chan ProcessDownMsg
}
```

### 2.2 Temporal-Grade Durable Engine (SQLite WAL)
```go
package state

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"
)

type DurableEngine struct {
	db *sql.DB
	mu sync.Mutex
}

// ExecuteActivity returns cached result if step was already completed, avoiding duplicate tokens
func (e *DurableEngine) ExecuteActivity(
	ctx context.Context,
	workflowID string,
	activityName string,
	input any,
	fn func(ctx context.Context) (any, error),
) (json.RawMessage, error) {
	// 1. Check SQLite WAL for existing ActivityCompleted event (Replay Cache)
	// 2. Cache hit: Return cached output immediately.
	// 3. Cache miss: Execute fn(), record ActivityCompleted, and return.
	return nil, nil
}
```

### 2.3 Saga Non-Git Side-Effect Coordinator
```go
package refinery

import (
	"context"
	"sync"
)

type SagaAction struct {
	Name       string
	Execute    func(ctx context.Context) error
	Compensate func(ctx context.Context) error // LIFO rollback action
}

type SagaCoordinator struct {
	mu           sync.Mutex
	executedComp []func(ctx context.Context) error
}

func (s *SagaCoordinator) ExecuteStep(ctx context.Context, step SagaAction) error {
	// Executes step. If error, calls all compensations in reverse LIFO order.
	return nil
}
```

### 2.4 Byzantine Quorum Consensus Arbiter
```go
package verification

type AgentVote struct {
	AgentID     string
	Category    string
	Approve     bool
	DiffHash    string
	HasF2PProof bool // Machine verified F2P test exit code 0
	LintClean   bool // Zero static analysis errors
}

type ConsensusArbiter struct {
	Threshold float64 // Net weighted threshold e.g. +0.5
}

func (ca *ConsensusArbiter) EvaluateConsensus(
	votes []AgentVote,
	getThompsonWeights func(agentID, cat string) (alpha, beta float64),
) (bool, error) {
	// Discards unverified votes; weights remaining votes by Thompson Sampling Beta distribution.
	return true, nil
}
```

### 2.5 Tree-Sitter PageRank Repo Map Engine
```go
package memory

import (
	"context"
)

type SymbolTag struct {
	File      string `json:"file"`
	Symbol    string `json:"symbol"`
	Kind      string `json:"kind"` // def or ref
	Line      int    `json:"line"`
	Signature string `json:"signature"`
}

type RepoMapEngine interface {
	ExtractTags(ctx context.Context, repoPath string) ([]SymbolTag, error)
	ComputePageRank(ctx context.Context, activeFiles []string, promptText string) ([]string, error)
	RenderContext(ctx context.Context, maxTokens int) (string, error)
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

# Worker Swarm & OTP Policies
workers:
  supervisor:
    max_restarts: 3
    period_seconds: 30
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

# Speculative Branch Racing & Merge Refinery
speculative:
  enabled: true
  max_parallel_candidates: 2

refinery:
  test_command: "go test -v ./..."
  mutation_testing: true
  squash_merges: true

# Cognitive Memory, Repo Map & Rules
memory:
  db_path: ".brain/memory.db"
  conventions_dir: ".brain/conventions"
  soul_file: ".brain/SOUL.md"
  max_repo_map_tokens: 1024
  max_injected_rule_tokens: 800
  thompson_sampling: true

# Verification & Sandboxing
sandbox:
  tier: "auto" # options: bwrap, landlock, none
  allow_network: false

# Notification Webhooks (OpenClaw-style)
notifications:
  slack_webhook: ""
  discord_webhook: ""
  telegram_chat_id: ""
```

---

## 4. Key Go Dependencies

| Package | Purpose |
|---|---|
| `modernc.org/sqlite` | Pure Go embedded SQLite engine with WAL support (zero CGO required) |
| `github.com/philippgille/chromem-go` | Pure Go in-memory vector store with persistence (<2ms search) |
| `github.com/smacker/go-tree-sitter` | Tree-Sitter AST parser for definitions, references, and syntax check |
| `thejerf/suture/v4` | Erlang-style supervisor trees in Go with rate-limited exponential backoff |
| `github.com/charmbracelet/bubbletea` | Elm-architecture Terminal User Interface framework |
| `github.com/charmbracelet/lipgloss` | Terminal styling, borders, and layouts |
| `github.com/creack/pty` | PTY allocation for interactive CLI worker wrapping |
| `github.com/pkoukk/tiktoken-go` | Pure Go BPE tokenization for `/v1/messages/count_tokens` |
| `github.com/tidwall/gjson` & `sjson` | Zero-allocation JSON manipulation for low-latency proxying |
| `github.com/bmatcuk/doublestar/v4` | High-performance glob matching for `.cursor/rules/*.mdc` rules |
| `github.com/landlock-lsm/go-landlock` | Unprivileged Linux Landlock LSM sandboxing |
| `github.com/spf13/cobra` | CLI command routing (`cli-leader start`, `status`) |
