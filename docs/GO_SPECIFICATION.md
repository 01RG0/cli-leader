# Go Technical Specification & Package Layout

## 1. Standard Project Layout

`cli-leader` adheres to the standard Go project structure:

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
│   │   ├── translator.go           # Anthropic -> OpenAI / Gemini / Ollama translation
│   │   ├── stream.go               # SSE stream adapter & chunk encoder
│   │   └── providers/
│   │       ├── openai.go           # OpenAI payload adapter
│   │       ├── gemini.go           # Gemini payload adapter
│   │       └── ollama.go           # Ollama / DeepSeek adapter
│   ├── mcp/
│   │   ├── server.go               # Model Context Protocol stdio/SSE server
│   │   ├── tools.go                # Tool registration & schemas
│   │   └── handlers.go             # Tool invocation dispatch handlers
│   ├── worker/
│   │   ├── manager.go              # Worker lifecycle, pooling & supervisor
│   │   ├── process.go              # Subprocess execution, PTY & timeouts
│   │   ├── worktree.go             # Git worktree sandbox management
│   │   └── adapters/
│   │       ├── adapter.go          # CLI Adapter interface
│   │       ├── aider.go            # Aider CLI integration
│   │       ├── gemini_cli.go       # Gemini CLI integration
│   │       ├── ollama_cli.go       # Ollama / local CLI integration
│   │       └── generic.go          # Custom generic CLI adapter
│   ├── memory/
│   │   ├── store.go                # SQLite episodic and semantic storage
│   │   ├── matrix.go               # Dynamic CLI capability rating & routing
│   │   ├── rules.go                # Rule extractor & markdown parser
│   │   └── schema.sql              # Database DDL
│   ├── review/
│   │   └── red_blue.go             # Red Team / Blue Team adversarial arbitration
│   └── tui/
│       ├── app.go                  # Bubble Tea main application model
│       ├── views/
│       │   ├── dashboard.go        # Swarm status & task graph
│       │   ├── logs.go             # Real-time worker log streaming
│       │   └── memory_view.go      # Knowledge & rule explorer
│       └── styles/
│           └── theme.go            # Lipgloss styling, colors, and layout
├── pkg/
│   └── protocol/
│       ├── anthropic.go            # Anthropic API message structs
│       ├── openai.go               # OpenAI API message structs
│       └── gemini.go               # Gemini API message structs
├── docs/                           # Architecture, Specs, Roadmap
├── go.mod
├── go.sum
└── README.md
```

---

## 2. Core Go Interfaces & Structs

### 2.1 Worker Adapter Interface
```go
package worker

import (
	"context"
	"io"
)

type CLIAdapter interface {
	// Name returns the identifier of the CLI tool (e.g. "aider", "gemini-cli")
	Name() string
	
	// PrepareCommand constructs the exec command with appropriate flags and environment
	PrepareCommand(ctx context.Context, task Task, worktreeDir string) (*exec.Cmd, error)
	
	// HandlePrompts monitors stdout for interactive confirmation questions and handles auto-reply
	HandlePrompts(in io.Writer, out io.Reader) error
	
	// ParseResult parses stdout/stderr and git diff into a structured TaskResult
	ParseResult(output []byte, diff string) (*TaskResult, error)
}
```

### 2.2 Task & Worker Manager
```go
type TaskStatus string

const (
	StatusQueued    TaskStatus = "queued"
	StatusRunning   TaskStatus = "running"
	StatusReviewing TaskStatus = "reviewing"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
	StatusAborted   TaskStatus = "aborted"
)

type Task struct {
	ID           string            `json:"id"`
	CLIName      string            `json:"cli_name"`
	Prompt       string            `json:"prompt"`
	TargetFiles  []string          `json:"target_files"`
	WorktreePath string            `json:"worktree_path"`
	Branch       string            `json:"branch"`
	Env          map[string]string `json:"env"`
	Status       TaskStatus        `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
}

type TaskResult struct {
	TaskID    string        `json:"task_id"`
	ExitCode  int           `json:"exit_code"`
	Stdout    string        `json:"stdout"`
	Stderr    string        `json:"stderr"`
	GitDiff   string        `json:"git_diff"`
	Duration  time.Duration `json:"duration"`
	Tokens    int64         `json:"tokens_used"`
	Error     error         `json:"error,omitempty"`
}
```

### 2.3 Cognitive Memory Store Interface
```go
package memory

type MemoryStore interface {
	// Episodic
	RecordTaskRun(ctx context.Context, task Task, result TaskResult) error
	GetTaskHistory(ctx context.Context, limit int) ([]TaskResult, error)
	
	// Semantic / Rules
	AddLearnedRule(ctx context.Context, rule Rule) error
	FindRelevantRules(ctx context.Context, contextQuery string) ([]Rule, error)
	
	// Matrix & Routing
	UpdateCLIScore(ctx context.Context, cliName, category string, success bool, latency time.Duration) error
	RecommendCLI(ctx context.Context, category string) (string, error)
}

type Rule struct {
	ID          string    `json:"id"`
	Category    string    `json:"category"`
	Directive   string    `json:"directive"`
	Source      string    `json:"source"` // e.g. "worker:aider:error_fix"
	Confidence  float64   `json:"confidence"`
	CreatedAt   time.Time `json:"created_at"`
}
```

---

## 3. Concurrency & Streaming Architecture

1. **Worker Process Isolation**:
   - Each worker runs in its own goroutine with an attached `creack/pty.Pty`.
   - Output is multiplexed to a thread-safe ring buffer (`ring.Buffer`) for TUI viewing and to an append-only log file.
2. **Context Cancellation & Safety**:
   - When the user cancels a task or Claude CLI calls `abort_worker`, the worker's parent context is canceled.
   - The supervisor issues `SIGTERM` followed by a grace period (3 seconds), escalating to `SIGKILL` on the entire process group.
3. **Event Bus**:
   - A lightweight in-memory Go channel bus (`chan Event`) broadcasts worker state transitions to the TUI and the MCP notification system.

---

## 4. Configuration Schema (`cli-leader.yaml`)

```yaml
version: "1.0"

# Multi-Provider Gateway Settings
gateway:
  listen_addr: "127.0.0.1:8082"
  default_provider: "openai" # Options: anthropic, openai, gemini, ollama
  providers:
    openai:
      base_url: "https://api.openai.com/v1"
      api_key_env: "OPENAI_API_KEY"
      model_mapping:
        "claude-3-7-sonnet-20250219": "o3-mini"
        "claude-3-5-sonnet-20241022": "gpt-4o"
    gemini:
      base_url: "https://generativelanguage.googleapis.com/v1beta"
      api_key_env: "GEMINI_API_KEY"
      model_mapping:
        "claude-3-7-sonnet-20250219": "gemini-2.5-pro"
    ollama:
      base_url: "http://localhost:11434/v1"
      model_mapping:
        "claude-3-5-sonnet-20241022": "deepseek-r1:32b"

# Worker Definitions
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

# Learning & Memory Bank Settings
memory:
  db_path: ".brain/memory.db"
  rules_file: ".brain/conventions.md"
  auto_reflect: true

# Verification & Peer Review
review:
  enable_red_team: true
  adversary_cli: "aider"
  auto_rollback_on_test_failure: true
```

---

## 5. Key Go Dependencies

| Package | Purpose |
|---|---|
| `github.com/charmbracelet/bubbletea` | Elm-architecture Terminal User Interface framework |
| `github.com/charmbracelet/lipgloss` | Terminal styling, borders, and color layouts |
| `github.com/creack/pty` | PTY allocation for interactive CLI worker wrapping |
| `modernc.org/sqlite` | Pure Go embedded SQLite engine (no CGO required) |
| `github.com/spf13/cobra` | CLI command routing (`cli-leader start`, `init`, `status`) |
| `github.com/spf13/viper` | YAML and environment variable configuration management |
| `gopkg.in/yaml.v3` | High-fidelity YAML parsing |
