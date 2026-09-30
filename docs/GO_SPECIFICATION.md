# Go Technical Specification & Package Layout

This document provides the canonical Go package layout, concrete types, interfaces, method signatures, concurrency semantics, and error handling contracts for `github.com/01RG0/cli-leader`.

> **Documentation Contract**: All specifications adhere strictly to [CONTRACT.md](_review/CONTRACT.md). System design, control flows, and architectural narratives are maintained in [ARCHITECTURE.md](ARCHITECTURE.md), delivery milestones in [ROADMAP.md](ROADMAP.md), and ecosystem research in [ECOSYSTEM_INNOVATIONS.md](ECOSYSTEM_INNOVATIONS.md).

---

## 1. Canonical Package Hierarchy

All packages reside within the module `github.com/01RG0/cli-leader`:

```
github.com/01RG0/cli-leader/
├── cmd/
│   └── cli-leader/                     # [v0.1 Core] CLI application entrypoint & daemon commands
├── internal/
│   ├── config/                         # [v0.1 Core] YAML & ENV configuration loader and validator
│   ├── gateway/                        # [v0.1 Core] FastHTTP reverse proxy translating /v1/messages
│   │   └── providers/                  # [v0.1 Core] Upstream API adapters (OpenAI, Gemini, Ollama)
│   ├── state/                          # [Planned] Dual Ledger, durable replay engine, stall detector
│   ├── mcp/                            # [Planned] Model Context Protocol stdio/SSE server
│   ├── acp/                            # [Planned] Agent Client Protocol JSON-RPC server
│   ├── supervisor/                     # [Planned] Erlang OTP supervisor trees & process lifecycles
│   ├── worker/                         # [Planned] Worker manager, PTY subprocess execution, Git worktrees
│   │   └── adapters/                   # [Planned] Worker CLI adapters (Aider, Gemini, Ollama, Generic)
│   ├── speculative/                    # [Planned] Speculative branch racing dispatcher
│   ├── refinery/                       # [Planned] Serialized merge queue & compensating Saga coordinator
│   ├── verification/                   # [Planned] Fail-to-Pass verification, AST mutation gate, consensus
│   ├── sandbox/                        # [Planned] Bubblewrap and Landlock security containment
│   ├── memory/                         # [Planned] SQLite store, vector index, RepoMap, Thompson Sampling
│   ├── notify/                         # [Planned] Webhook notification dispatcher (Slack, Discord, Telegram)
│   └── tui/                            # [Planned] Charm Bubbletea terminal user interface
│       └── views/                      # [Planned] Dashboard, live streaming logs, memory views
├── docs/                               # System documentation & architectural records
├── go.mod
├── go.sum
└── README.md
```

---

## 2. Package Specifications

---

### 2.1 `cmd/cli-leader` `[v0.1 Core]`

CLI application entrypoint, Cobra subcommands, daemon signal listeners, and graceful shutdown orchestration.

#### Types & Signatures

```go
package main

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type CLIOptions struct {
	ConfigPath string
	Verbose    bool
	ListenAddr string
	DaemonMode bool
}

type Application struct {
	opts       CLIOptions
	rootCmd    *cobra.Command
	shutdownFn func(ctx context.Context) error
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Execute(ctx context.Context) error {
	return nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	return nil
}
```

#### Concurrency Semantics
- Main goroutine listens for OS termination signals (`syscall.SIGINT`, `syscall.SIGTERM`) using `signal.NotifyContext`.
- Context cancellation cascades to all running subsystems (`internal/gateway`, `internal/supervisor`, `internal/mcp`).
- Shutdown grants a 5-second grace period for in-flight requests and PTY child process termination before exiting.

#### Error Handling Contract
- Returns exit code `0` on clean shutdown.
- Returns exit code `1` on unrecoverable initialization error (`ErrConfigInitFailed`, `ErrPortBindingFailed`).
- Returns exit code `130` on interrupted signal cancellation.

---

### 2.2 `internal/config` `[v0.1 Core]`

YAML and environment variable configuration loader and validator. Implements schema keys specified in `CONTRACT.md` Section 3.3.

#### Types & Signatures

```go
package config

import (
	"errors"
	"time"
)

var (
	ErrConfigNotFound       = errors.New("configuration file not found")
	ErrInvalidConfiguration = errors.New("invalid configuration parameters")
	ErrMissingAPIKey        = errors.New("required API key environment variable is not set")
)

type Config struct {
	Version       string               `yaml:"version"`
	Gateway       GatewayConfig        `yaml:"gateway"`
	Workers       WorkersConfig        `yaml:"workers"`
	Speculative   SpeculativeConfig    `yaml:"speculative"`
	Refinery      RefineryConfig       `yaml:"refinery"`
	Memory        MemoryConfig         `yaml:"memory"`
	Sandbox       SandboxConfig        `yaml:"sandbox"`
	Notifications NotificationsConfig  `yaml:"notifications"`
}

type GatewayConfig struct {
	ListenAddr      string                    `yaml:"listen_addr"`
	DefaultProvider string                    `yaml:"default_provider"`
	Providers       map[string]ProviderConfig `yaml:"providers"`
}

type ProviderConfig struct {
	BaseURL      string            `yaml:"base_url"`
	APIKeyEnv    string            `yaml:"api_key_env"`
	ModelMapping map[string]string `yaml:"model_mapping"`
}

type WorkersConfig struct {
	Supervisor SupervisorConfig            `yaml:"supervisor"`
	Aider      WorkerToolConfig            `yaml:"aider"`
	Gemini     WorkerToolConfig            `yaml:"gemini"`
	Ollama     WorkerToolConfig            `yaml:"ollama"`
	Generic    map[string]WorkerToolConfig `yaml:"generic,omitempty"`
}

type SupervisorConfig struct {
	MaxRestarts         int `yaml:"max_restarts"`
	PeriodSeconds       int `yaml:"period_seconds"`
	HeartbeatTimeoutSec int `yaml:"heartbeat_timeout_sec"`
}

type WorkerToolConfig struct {
	Binary     string   `yaml:"binary"`
	Flags      []string `yaml:"flags"`
	Categories []string `yaml:"categories"`
}

type SpeculativeConfig struct {
	Enabled               bool `yaml:"enabled"`
	MaxParallelCandidates int  `yaml:"max_parallel_candidates"`
}

type RefineryConfig struct {
	TestCommand     string `yaml:"test_command"`
	MutationTesting bool   `yaml:"mutation_testing"`
	SquashMerges    bool   `yaml:"squash_merges"`
}

type MemoryConfig struct {
	DBPath                string `yaml:"db_path"`
	ConventionsDir        string `yaml:"conventions_dir"`
	SoulFile              string `yaml:"soul_file"`
	MaxRepoMapTokens      int    `yaml:"max_repo_map_tokens"`
	MaxInjectedRuleTokens int    `yaml:"max_injected_rule_tokens"`
	ThompsonSampling      bool   `yaml:"thompson_sampling"`
}

type SandboxConfig struct {
	Tier         string `yaml:"tier"` // "bwrap", "landlock", "none"
	AllowNetwork bool   `yaml:"allow_network"`
}

type NotificationsConfig struct {
	SlackWebhook   string `yaml:"slack_webhook"`
	DiscordWebhook string `yaml:"discord_webhook"`
	TelegramChatID string `yaml:"telegram_chat_id"`
}

func Load(path string) (*Config, error) {
	return nil, nil
}

func (c *Config) Validate() error {
	return nil
}
```

#### Concurrency Semantics
- Configuration is loaded at daemon boot.
- Read operations are safe for concurrent access. Hot-reloading updates an atomic pointer (`atomic.Pointer[Config]`).

#### Error Handling Contract
- `Load` wraps filesystem and YAML parsing errors with `%w`.
- `Validate` aggregates all configuration validation errors into `ErrInvalidConfiguration` with detailed field paths.

---

### 2.3 `internal/gateway` `[v0.1 Core]`

FastHTTP reverse proxy translating Claude CLI Anthropic `/v1/messages` requests to upstream LLM providers. Includes mandatory local BPE token evaluation.

> `TODO(verify): gateway_http_engine` — Confirm whether FastHTTP is canonical or if `net/http` is permitted per ROADMAP.md Phase 1.

#### Types & Signatures

```go
package gateway

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	ErrProviderUnavailable       = errors.New("upstream provider is unavailable")
	ErrModelNotFound             = errors.New("requested model not mapped to provider")
	ErrInvalidTokenCountPayload  = errors.New("invalid payload for token counting")
	ErrStreamAborted             = errors.New("SSE streaming aborted prematurely")
)

type TokenCounter interface {
	CountTokens(ctx context.Context, model string, messages []AnthropicMessage, system string) (int, error)
}

type StreamTranslator interface {
	TranslateStream(ctx context.Context, src io.Reader, dst io.Writer) error
}

type ProxyServer struct {
	listenAddr string
	counter    TokenCounter
	translator StreamTranslator
}

func NewProxyServer(addr string, counter TokenCounter, trans StreamTranslator) *ProxyServer {
	return &ProxyServer{listenAddr: addr, counter: counter, translator: trans}
}

func (s *ProxyServer) Start(ctx context.Context) error {
	return nil
}

func (s *ProxyServer) Stop(ctx context.Context) error {
	return nil
}

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicCountTokensRequest struct {
	Model    string             `json:"model"`
	Messages []AnthropicMessage `json:"messages"`
	System   string             `json:"system,omitempty"`
}

type AnthropicCountTokensResponse struct {
	InputTokens int `json:"input_tokens"`
}
```

#### Concurrency Semantics
- Requests are handled concurrently by FastHTTP worker goroutines.
- SSE stream transformation processes chunks through `gjson` and `sjson` on byte slices without reflection or heap escaping.
- Performance metric: `Target (unmeasured): <100 µs TTFT overhead`.

#### Error Handling Contract
- FastHTTP handler converts upstream connection failures to HTTP 502 with structured JSON error response.
- Missing model mappings trigger HTTP 400 with `ErrModelNotFound`.

---

### 2.4 `internal/gateway/providers` `[v0.1 Core]`

Upstream provider adapters for OpenAI, Gemini, and Ollama. Translates payload schemas, reasoning tokens, and tool definitions.

#### Types & Signatures

```go
package providers

import (
	"context"
	"errors"
	"io"
)

var (
	ErrUnsupportedModel  = errors.New("model mapping not supported by adapter")
	ErrUpstreamTimeout   = errors.New("upstream provider timed out")
	ErrUpstreamAuth      = errors.New("upstream authentication failed")
	ErrRateLimited       = errors.New("upstream rate limit exceeded")
)

type ProviderAdapter interface {
	Name() string
	TranslateRequest(ctx context.Context, anthropicPayload []byte) ([]byte, error)
	Send(ctx context.Context, payload []byte) (io.ReadCloser, error)
	TranslateChunk(chunk []byte) ([]byte, error)
}

type OpenAIAdapter struct {
	baseURL string
	apiKey  string
}

func NewOpenAIAdapter(baseURL, apiKey string) *OpenAIAdapter {
	return &OpenAIAdapter{baseURL: baseURL, apiKey: apiKey}
}

func (a *OpenAIAdapter) Name() string { return "openai" }

func (a *OpenAIAdapter) TranslateRequest(ctx context.Context, payload []byte) ([]byte, error) {
	return nil, nil
}

func (a *OpenAIAdapter) Send(ctx context.Context, payload []byte) (io.ReadCloser, error) {
	return nil, nil
}

func (a *OpenAIAdapter) TranslateChunk(chunk []byte) ([]byte, error) {
	return nil, nil
}

type GeminiAdapter struct {
	baseURL string
	apiKey  string
}

type OllamaAdapter struct {
	baseURL string
}
```

#### Concurrency Semantics
- Provider adapters are stateless and safe for concurrent use across request goroutines.
- HTTP client uses pooled keep-alive TCP connections with idle connection timeouts.

#### Error Handling Contract
- Upstream HTTP 429 errors wrap `ErrRateLimited` with `Retry-After` header extraction.
- HTTP 401/403 errors wrap `ErrUpstreamAuth`.

---

### 2.5 `internal/state` `[Planned]`

Dual Ledger (`TaskLedger`, `ProgressLedger`), SQLite WAL durable replay engine, and cryptographic stall detection.

#### Types & Signatures

```go
package state

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

var (
	ErrWorkflowStalled        = errors.New("workflow stalled: detected circular state loop")
	ErrTaskNotFound           = errors.New("task not found in ledger")
	ErrActivityFailed         = errors.New("activity failed during durable execution")
	ErrDuplicateEventSequence = errors.New("duplicate event sequence number")
)

type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskReviewing TaskStatus = "reviewing"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
)

type Task struct {
	ID          string     `json:"id"`
	Goal        string     `json:"goal"`
	Constraints []string   `json:"constraints"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskStep struct {
	StepID      string     `json:"step_id"`
	TaskID      string     `json:"task_id"`
	SequenceNum int        `json:"sequence_num"`
	Action      string     `json:"action"`
	Status      TaskStatus `json:"status"`
	Output      string     `json:"output,omitempty"`
	DiffHash    string     `json:"diff_hash,omitempty"`
}

type TaskLedger interface {
	CreateTask(ctx context.Context, task *Task) error
	GetTask(ctx context.Context, id string) (*Task, error)
	UpdateTaskStatus(ctx context.Context, id string, status TaskStatus) error
}

type ProgressLedger interface {
	RecordStep(ctx context.Context, step *TaskStep) error
	ListSteps(ctx context.Context, taskID string) ([]TaskStep, error)
}

type DurableEngine struct {
	db *sql.DB
	mu sync.Mutex
}

func NewDurableEngine(db *sql.DB) *DurableEngine {
	return &DurableEngine{db: db}
}

func (e *DurableEngine) ExecuteActivity(
	ctx context.Context,
	workflowID string,
	activityName string,
	input any,
	fn func(ctx context.Context) (any, error),
) (json.RawMessage, error) {
	return nil, nil
}

func (e *DurableEngine) ReplayWorkflow(ctx context.Context, workflowID string) error {
	return nil
}

type StallDetector struct {
	mu           sync.Mutex
	history      map[string][]string // taskID -> list of state hashes
	maxThreshold int
}

func NewStallDetector(maxThreshold int) *StallDetector {
	return &StallDetector{
		history:      make(map[string][]string),
		maxThreshold: maxThreshold,
	}
}

func (sd *StallDetector) RecordStepState(taskID, tool string, args []byte, diffHash string) (bool, error) {
	return false, nil
}

func (sd *StallDetector) Reset(taskID string) {
	sd.mu.Lock()
	defer sd.mu.Unlock()
	delete(sd.history, taskID)
}
```

#### Concurrency Semantics
- SQLite WAL allows concurrent readers with a single writer serialized through `sync.Mutex`.
- `ExecuteActivity` checks replay cache under lock before dispatching work.
- `StallDetector` state map is protected by `sync.Mutex`.

#### Error Handling Contract
- When `StallDetector` detects identical cryptographic state hash repeated beyond threshold (default 3), it returns `ErrWorkflowStalled`.
- Failed activities record `ActivityFailed` events in WAL and return wrapped `ErrActivityFailed`.

---

### 2.6 `internal/mcp` `[Planned]`

Model Context Protocol (MCP) server exposing orchestrator tools to Claude CLI over stdio and local SSE. Incorporates tool definitions from legacy `API_SPECIFICATION.md`.

#### Types & Signatures

```go
package mcp

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrToolNotFound           = errors.New("requested MCP tool not registered")
	ErrInvalidToolParameters  = errors.New("invalid tool invocation arguments")
	ErrToolExecutionTimeout   = errors.New("tool execution timed out")
)

type ToolDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}

type CallToolRequest struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type CallToolResult struct {
	Content []ToolContent `json:"content"`
	IsError bool          `json:"is_error"`
}

type ToolContent struct {
	Type string `json:"type"` // "text", "image", "resource"
	Text string `json:"text,omitempty"`
}

type ToolHandler func(ctx context.Context, req CallToolRequest) (CallToolResult, error)

type Server struct {
	tools map[string]ToolHandler
}

func NewServer() *Server {
	return &Server{tools: make(map[string]ToolHandler)}
}

func (s *Server) RegisterTool(name string, handler ToolHandler) {
	s.tools[name] = handler
}

func (s *Server) HandleCall(ctx context.Context, req CallToolRequest) (CallToolResult, error) {
	return CallToolResult{}, nil
}

// Built-in MCP Tool Parameter Structs
type DispatchWorkerParams struct {
	CLIName     string   `json:"cli_name"`
	TaskPrompt  string   `json:"task_prompt"`
	TargetFiles []string `json:"target_files"`
	Flags       []string `json:"flags,omitempty"`
}

type CandidateSpec struct {
	CLIName string `json:"cli_name"`
	Prompt  string `json:"prompt"`
}

type DispatchSpeculativeParams struct {
	Candidates  []CandidateSpec `json:"candidates"`
	TargetFiles []string        `json:"target_files"`
}

type VerifyF2PParams struct {
	TaskID               string `json:"task_id"`
	ReproductionTestCmd  string `json:"reproduction_test_cmd"`
	EnableMutationGate   bool   `json:"enable_mutation_gate"`
}

type RefineryEnqueueParams struct {
	TaskID        string `json:"task_id"`
	CommitMessage string `json:"commit_message"`
}

type QueryBrainMemoryParams struct {
	Query    string `json:"query"`
	Category string `json:"category"` // "all", "rules", "history", "symbols", "skills"
}
```

#### Concurrency Semantics
- Tool registrations occur at server initialization before incoming RPCs.
- `HandleCall` dispatches requests in separate goroutines; handlers must be re-entrant and thread-safe.

#### Error Handling Contract
- Unregistered tools return JSON-RPC code `-32601` wrapping `ErrToolNotFound`.
- Schema unmarshal errors return code `-32602` wrapping `ErrInvalidToolParameters`.

---

### 2.7 `internal/acp` `[Planned]`

Agent Client Protocol (ACP) JSON-RPC server enabling editor integrations (Zed, JetBrains).

> `TODO(verify): acp_transport` — Clarify if ACP transport supports UNIX domain sockets in addition to stdio.

#### Types & Signatures

```go
package acp

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

var (
	ErrSessionNotFound         = errors.New("ACP editor session not found")
	ErrProtocolVersionMismatch = errors.New("unsupported ACP protocol version")
	ErrClientDisconnected      = errors.New("editor client connection closed")
)

type InitializeParams struct {
	ProtocolVersion string          `json:"protocolVersion"`
	ClientInfo      ClientInfo      `json:"clientInfo"`
	Capabilities    map[string]bool `json:"capabilities"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
	Capabilities    map[string]interface{} `json:"capabilities"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type AgentSendMessageParams struct {
	ThreadID string       `json:"threadId"`
	Message  AgentMessage `json:"message"`
}

type AgentMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Session struct {
	ID           string
	Client       ClientInfo
	BufferStates map[string]string // filepath -> current buffer content
	mu           sync.RWMutex
}

type Server struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewServer() *Server {
	return &Server{sessions: make(map[string]*Session)}
}

func (s *Server) Serve(ctx context.Context, reader io.Reader, writer io.Writer) error {
	return nil
}
```

#### Concurrency Semantics
- Editor sessions maintain independent read-write locks (`sync.RWMutex`).
- Incoming JSON-RPC requests are decoded sequentially per client stream and dispatched concurrently.

#### Error Handling Contract
- Handshake failures yield JSON-RPC error `-32000` with `ErrProtocolVersionMismatch`.
- Unknown thread IDs return error `-32001` with `ErrSessionNotFound`.

---

### 2.8 `internal/supervisor` `[Planned]`

Erlang OTP-inspired process supervision tree implementing `one_for_one`, `one_for_all`, and `rest_for_one` failure recovery policies.

> `TODO(verify): supervisor_backend_impl` — Confirm whether `Supervisor` directly wraps `thejerf/suture/v4` or implements custom OTP state logic.

#### Types & Signatures

```go
package supervisor

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrRestartIntensityExceeded = errors.New("maximum restart intensity exceeded within period")
	ErrChildAlreadyExists       = errors.New("child with specified ID already registered")
	ErrChildNotFound            = errors.New("child process ID not found")
	ErrSupervisorStopping       = errors.New("supervisor is shutting down")
)

type RestartStrategy string

const (
	OneForOne  RestartStrategy = "one_for_one"  // Isolated worker candidates
	OneForAll  RestartStrategy = "one_for_all"  // Compound worker bundles (PTY + reader)
	RestForOne RestartStrategy = "rest_for_one" // Linear verification pipelines
)

type ProcessStatus string

const (
	StatusStarting ProcessStatus = "starting"
	StatusRunning  ProcessStatus = "running"
	StatusCrashed  ProcessStatus = "crashed"
	StatusStopped  ProcessStatus = "stopped"
)

type ProcessDownMsg struct {
	WorkerID  string
	ExitCode  int
	Err       error
	Timestamp time.Time
}

type WorkerProcess interface {
	ID() string
	PID() int
	Status() ProcessStatus
	Start(ctx context.Context) error
	Stop(timeout time.Duration) error
}

type Supervisor struct {
	mu           sync.Mutex
	strategy     RestartStrategy
	maxIntensity int           // Max failures allowed
	period       time.Duration // Within duration
	failures     []time.Time
	children     map[string]WorkerProcess
	downChan     chan ProcessDownMsg
}

func NewSupervisor(strategy RestartStrategy, maxIntensity int, period time.Duration) *Supervisor {
	return &Supervisor{
		strategy:     strategy,
		maxIntensity: maxIntensity,
		period:       period,
		children:     make(map[string]WorkerProcess),
		downChan:     make(chan ProcessDownMsg, 64),
	}
}

func (s *Supervisor) AddChild(child WorkerProcess) error {
	return nil
}

func (s *Supervisor) RemoveChild(id string) error {
	return nil
}

func (s *Supervisor) Start(ctx context.Context) error {
	return nil
}

func (s *Supervisor) Stop() error {
	return nil
}
```

#### Concurrency Semantics
- Supervisor loop runs on a dedicated actor goroutine reading from `downChan`.
- Mutex serializes child additions, removals, and intensity threshold evaluation.
- When `OneForAll` detects a crash, it terminates all peer children concurrently via `sync.WaitGroup`.

#### Error Handling Contract
- When crash count exceeds `maxIntensity` within `period`, supervisor transitions to `StatusCrashed` and returns `ErrRestartIntensityExceeded` to trigger Brain escalation.

---

### 2.9 `internal/worker` `[Planned]`

Worker Manager, PTY subprocess execution (`creack/pty`), thread-safe circular log buffer, and native Git worktree lifecycle.

> `TODO(verify): pty_subreaper_support` — Confirm if `prctl(PR_SET_CHILD_SUBREAPER, 1)` is invoked alongside `Pdeathsig`.

#### Types & Signatures

```go
package worker

import (
	"context"
	"errors"
	"os/exec"
	"sync"
	"time"

	"github.com/creack/pty"
)

var (
	ErrWorkerTimeout        = errors.New("worker execution timed out")
	ErrWorktreeLockFailed   = errors.New("failed to acquire git worktree mutex lock")
	ErrProcessKilled        = errors.New("process terminated via signal")
	ErrBufferOverflow       = errors.New("log ring buffer write overflow")
)

type WorkerRequest struct {
	ID          string
	ToolName    string
	Prompt      string
	TargetFiles []string
	WorktreeDir string
	Timeout     time.Duration
}

type WorkerResult struct {
	WorkerID string
	ExitCode int
	DiffHash string
	Logs     []byte
	Duration time.Duration
}

type RingLogBuffer struct {
	mu       sync.RWMutex
	capacity int
	entries  [][]byte
	cursor   int
}

func NewRingLogBuffer(capacity int) *RingLogBuffer {
	return &RingLogBuffer{
		capacity: capacity,
		entries:  make([][]byte, capacity),
	}
}

func (r *RingLogBuffer) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (r *RingLogBuffer) ReadTail(lines int) [][]byte {
	return nil
}

type WorktreeManager struct {
	repoRoot string
	mu       sync.Mutex // Serializes git worktree additions to prevent .git/config.lock collisions
}

func NewWorktreeManager(repoRoot string) *WorktreeManager {
	return &WorktreeManager{repoRoot: repoRoot}
}

func (w *WorktreeManager) CreateWorktree(ctx context.Context, branchName, targetDir string) error {
	return nil
}

func (w *WorktreeManager) RemoveWorktree(ctx context.Context, targetDir string) error {
	return nil
}

type PTYRunner struct{}

func (p *PTYRunner) StartPTY(ctx context.Context, cmd *exec.Cmd, sz *pty.Winsize) (*exec.Cmd, error) {
	return nil, nil
}
```

#### Concurrency Semantics
- PTY streams are read on background goroutines and written to `RingLogBuffer`.
- `WorktreeManager` guards Git mutations with a package-level mutex to eliminate `.git/config.lock` contention.
- Subprocesses run in isolated process groups (`Setsid: true`) with `Pdeathsig = syscall.SIGTERM`.

#### Error Handling Contract
- Execution exceeding `Timeout` triggers `SIGTERM`, waits a 3-second grace period, then issues `SIGKILL` on `-pgid`, returning `ErrWorkerTimeout`.

---

### 2.10 `internal/worker/adapters` `[Planned]`

CLI tool adapters for Aider, Gemini CLI, Ollama CLI, and generic terminal commands.

#### Types & Signatures

```go
package adapters

import (
	"context"
	"errors"
	"os/exec"
)

var (
	ErrUnknownCLITool      = errors.New("unsupported worker CLI adapter")
	ErrAdapterParseFailure = errors.New("failed to parse worker output stream")
	ErrUnsupportedFlag     = errors.New("unsupported adapter flag")
)

type AutoReplyRule struct {
	Pattern  string
	Response string
}

type CLIAdapter interface {
	Name() string
	BuildCommand(ctx context.Context, worktreeDir, prompt string, targetFiles []string) (*exec.Cmd, error)
	AutoReplyRules() []AutoReplyRule
	ExtractDiff(rawOutput []byte) ([]byte, error)
}

type AiderAdapter struct {
	binaryPath string
}

func (a *AiderAdapter) Name() string { return "aider" }

func (a *AiderAdapter) BuildCommand(ctx context.Context, worktreeDir, prompt string, targetFiles []string) (*exec.Cmd, error) {
	return nil, nil
}

func (a *AiderAdapter) AutoReplyRules() []AutoReplyRule {
	return []AutoReplyRule{
		{Pattern: `(?i)Apply changes\?`, Response: "y\n"},
		{Pattern: `(?i)Run test\?`, Response: "y\n"},
	}
}

func (a *AiderAdapter) ExtractDiff(raw []byte) ([]byte, error) {
	return nil, nil
}

type GeminiCLIAdapter struct{}
type OllamaCLIAdapter struct{}
type GenericCLIAdapter struct{}

// Dynamic Discovery & Zero-Shot Synthesis [Research]
type ZeroShotAdapterSynthesizer interface {
	SynthesizeAdapter(ctx context.Context, binaryPath string) (CLIAdapter, error)
}
```

#### Concurrency Semantics
- Adapters are stateless; command generation is thread-safe and re-entrant.

#### Error Handling Contract
- Malformed outputs wrap `ErrAdapterParseFailure` with line offset context.

---

### 2.11 `internal/speculative` `[Planned]`

Speculative branch racing dispatcher. Concurrently executes candidate workers across isolated Git worktrees.

#### Types & Signatures

```go
package speculative

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrAllCandidatesFailed   = errors.New("all speculative candidates failed verification")
	ErrBranchRaceTimeout     = errors.New("speculative branch race timed out")
	ErrSpeculativeCanceled   = errors.New("speculative execution canceled by winner")
)

type CandidateRequest struct {
	CLIName        string
	Prompt         string
	WorktreeBranch string
}

type CandidateOutcome struct {
	CandidateID string
	ExitCode    int
	DiffHash    string
	PassedF2P   bool
	Duration    time.Duration
	Err         error
}

type SpeculativeDispatcher struct {
	maxParallel int
}

func NewSpeculativeDispatcher(maxParallel int) *SpeculativeDispatcher {
	return &SpeculativeDispatcher{maxParallel: maxParallel}
}

func (d *SpeculativeDispatcher) DispatchRace(
	ctx context.Context,
	taskID string,
	candidates []CandidateRequest,
	evaluator func(ctx context.Context, outcome CandidateOutcome) (bool, error),
) (*CandidateOutcome, error) {
	return nil, nil
}
```

#### Concurrency Semantics
- Spawns parallel worker goroutines managed by `sync.WaitGroup`.
- Context cancellation propagates immediately to losing candidates once a candidate passes evaluation.

#### Error Handling Contract
- If all parallel workers crash or fail verification, returns wrapped `ErrAllCandidatesFailed`.

---

### 2.12 `internal/refinery` `[Planned]`

Serialized Bors-style merge queue and LIFO compensating Saga coordinator.

#### Types & Signatures

```go
package refinery

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrMergeConflict          = errors.New("refinery merge conflict on rebase")
	ErrRefineryTestFailed     = errors.New("test suite failed in refinery staging branch")
	ErrSagaCompensationFailed = errors.New("critical failure during saga compensation rollback")
	ErrQueueFull              = errors.New("refinery merge queue is full")
)

type MergeRequest struct {
	TaskID        string
	WorktreeDir   string
	CommitMessage string
	ResultChan    chan MergeResult
}

type MergeResult struct {
	MergedCommitSHA string
	Err             error
	Duration        time.Duration
}

type MergeQueue struct {
	mu       sync.Mutex
	requests chan MergeRequest
	running  bool
}

func NewMergeQueue(bufferSize int) *MergeQueue {
	return &MergeQueue{
		requests: make(chan MergeRequest, bufferSize),
	}
}

func (q *MergeQueue) Enqueue(ctx context.Context, req MergeRequest) error {
	return nil
}

func (q *MergeQueue) ProcessLoop(ctx context.Context) error {
	return nil
}

type SagaAction struct {
	Name       string
	Execute    func(ctx context.Context) error
	Compensate func(ctx context.Context) error // LIFO rollback action
}

type SagaCoordinator struct {
	mu           sync.Mutex
	compensations []func(ctx context.Context) error
}

func NewSagaCoordinator() *SagaCoordinator {
	return &SagaCoordinator{}
}

func (s *SagaCoordinator) ExecuteStep(ctx context.Context, step SagaAction) error {
	return nil
}

func (s *SagaCoordinator) CompensateAll(ctx context.Context) error {
	return nil
}
```

#### Concurrency Semantics
- `MergeQueue` processes items sequentially in a single dedicated worker goroutine to guarantee serial git rebases.
- `SagaCoordinator` maintains thread-safe LIFO slice of compensations executed sequentially on failure.

#### Error Handling Contract
- If staging rebase fails, returns `ErrMergeConflict`.
- If automated test command exits non-zero, staging branch is discarded and `ErrRefineryTestFailed` is returned.
- If compensation action fails during rollback, error is recorded and wrapped in `ErrSagaCompensationFailed`.

---

### 2.13 `internal/verification` `[Planned]`

Fail-to-Pass (F2P) automated verification, AST mutation testing gates, and Byzantine quorum consensus.

> `TODO(verify): mutation_testing_toolchain` — Clarify if AST mutation runs pure-Go `go/ast` rewriting or calls an external binary.

#### Types & Signatures

```go
package verification

import (
	"context"
	"errors"
)

var (
	ErrF2PReproductionFailed = errors.New("F2P reproduction test unexpectedly passed on unmodified HEAD")
	ErrF2PPostFixFailed      = errors.New("F2P test failed after patch application")
	ErrMutantSurvived        = errors.New("mutation testing gate failed: survived mutants detected")
	ErrConsensusRejected     = errors.New("byzantine quorum consensus rejected patch")
)

type F2PReceipt struct {
	ReproExitCodePre  int    // MUST be != 0
	ReproExitCodePost int    // MUST be == 0
	FullSuiteExitCode int    // MUST be == 0
	ReproOutput       string
}

type F2PEngine interface {
	VerifyF2P(ctx context.Context, worktreeDir, reproTestCmd, fullSuiteCmd string) (*F2PReceipt, error)
}

type MutationReceipt struct {
	TotalMutants    int
	KilledMutants   int
	SurvivedMutants int
	MutantDiffs     []string
}

type MutationTester interface {
	RunMutationGate(ctx context.Context, worktreeDir string, modifiedFiles []string) (*MutationReceipt, error)
}

type AgentVote struct {
	AgentID     string
	Category    string
	Approve     bool
	DiffHash    string
	HasF2PProof bool
	LintClean   bool
}

type ConsensusArbiter struct {
	Threshold float64 // Net weighted threshold e.g. +0.5
}

func (ca *ConsensusArbiter) EvaluateConsensus(
	votes []AgentVote,
	getThompsonWeights func(agentID, cat string) (alpha, beta float64),
) (bool, error) {
	return false, nil
}

// Distributed Byzantine Swarm Consensus [Research]
type DistributedConsensusArbiter interface {
	ProposeCommit(ctx context.Context, diffHash string, receipts []F2PReceipt) (bool, error)
}
```

#### Concurrency Semantics
- Test commands run in isolated worktree child processes.
- Consensus evaluation is purely functional and compute-bound; thread-safe and re-entrant.

#### Error Handling Contract
- If reproduction test returns code 0 before code modification, `ErrF2PReproductionFailed` is returned.
- If mutation testing detects `SurvivedMutants > 0`, `ErrMutantSurvived` is returned.

---

### 2.14 `internal/sandbox` `[Planned]`

Bubblewrap namespace isolation and Linux Landlock LSM sandboxing.

#### Types & Signatures

```go
package sandbox

import (
	"context"
	"errors"
	"os/exec"
)

var (
	ErrSandboxUnsupported      = errors.New("requested sandbox tier is not supported on host OS")
	ErrLandlockApplyFailed     = errors.New("failed to apply Landlock filesystem rules")
	ErrBubblewrapExecFailed    = errors.New("bubblewrap execution failed")
)

type Tier string

const (
	TierBwrap    Tier = "bwrap"
	TierLandlock Tier = "landlock"
	TierNone     Tier = "none"
)

type SandboxPolicy struct {
	Tier         Tier
	WorktreeDir  string
	AllowNetwork bool
	ReadOnlyDirs []string
}

type SandboxEngine interface {
	WrapCommand(ctx context.Context, cmd *exec.Cmd, policy SandboxPolicy) (*exec.Cmd, error)
}

type BwrapEngine struct{}

func (b *BwrapEngine) WrapCommand(ctx context.Context, cmd *exec.Cmd, policy SandboxPolicy) (*exec.Cmd, error) {
	return nil, nil
}

type LandlockEngine struct{}

func (l *LandlockEngine) WrapCommand(ctx context.Context, cmd *exec.Cmd, policy SandboxPolicy) (*exec.Cmd, error) {
	return nil, nil
}
```

#### Concurrency Semantics
- Sandbox wrappers modify command invocation arguments prior to process execution. Completely stateless and thread-safe.

#### Error Handling Contract
- If host kernel lacks Landlock LSM or Bubblewrap binary, engine falls back or returns `ErrSandboxUnsupported`.

---

### 2.15 `internal/memory` `[Planned]`

SQLite episodic store, embedded vector embeddings, Tree-Sitter RepoMap, and Thompson Sampling capability routing.

> `TODO(verify): vector_db_persistence` — Confirm if `chromem-go` vectors persist to SQLite BLOB tables or a dedicated directory.

#### Types & Signatures

```go
package memory

import (
	"context"
	"database/sql"
	"errors"
	"sync"
)

var (
	ErrMemoryNotFound       = errors.New("memory record not found")
	ErrVectorIndexCorrupt    = errors.New("vector index file is corrupted")
	ErrTokenBudgetExceeded   = errors.New("rendered context exceeds token budget")
	ErrRuleParseFailed       = errors.New("failed to parse cursor rule mdc frontmatter")
)

type SymbolTag struct {
	File      string `json:"file"`
	Symbol    string `json:"symbol"`
	Kind      string `json:"kind"` // "def" or "ref"
	Line      int    `json:"line"`
	Signature string `json:"signature"`
}

type RepoMapEngine interface {
	ExtractTags(ctx context.Context, repoPath string) ([]SymbolTag, error)
	ComputePageRank(ctx context.Context, activeFiles []string, promptText string) ([]string, error)
	RenderContext(ctx context.Context, maxTokens int) (string, error)
}

type VectorIndex interface {
	Add(ctx context.Context, id, text string, metadata map[string]string) error
	Search(ctx context.Context, query string, topK int) ([]string, error)
}

type ThompsonSamplingMatrix struct {
	mu     sync.RWMutex
	params map[string]map[string]*BetaDistribution // worker -> category -> Beta(alpha, beta)
}

type BetaDistribution struct {
	Alpha float64
	Beta  float64
}

func (m *ThompsonSamplingMatrix) Sample(category string, workers []string) string {
	return ""
}

func (m *ThompsonSamplingMatrix) Update(worker, category string, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
}
```

#### Concurrency Semantics
- Vector search latency: `Target (unmeasured): <2ms vector search latency`.
- Thompson sampling distribution parameters are protected by `sync.RWMutex`.
- SQLite database operates in WAL mode (`_journal_mode=WAL`).

#### Error Handling Contract
- When PageRank rendered tags exceed `maxTokens`, output is truncated deterministically and returned without failing.

---

### 2.16 `internal/notify` `[Planned]`

Webhook notification dispatcher (Slack, Discord, Telegram).

#### Types & Signatures

```go
package notify

import (
	"context"
	"errors"
	"time"
)

var (
	ErrWebhookFailed         = errors.New("failed to dispatch webhook notification")
	ErrRateLimitedWebhook    = errors.New("webhook rate limit exceeded by destination")
	ErrInvalidWebhookURL     = errors.New("malformed webhook endpoint URL")
)

type NotificationEvent struct {
	TaskID    string        `json:"task_id"`
	Status    string        `json:"status"`
	Summary   string        `json:"summary"`
	DiffStats string        `json:"diff_stats"`
	Duration  time.Duration `json:"duration"`
}

type Notifier interface {
	Notify(ctx context.Context, event NotificationEvent) error
}

type Dispatcher struct {
	queue chan NotificationEvent
}

func NewDispatcher(bufferSize int) *Dispatcher {
	return &Dispatcher{queue: make(chan NotificationEvent, bufferSize)}
}

func (d *Dispatcher) Enqueue(event NotificationEvent) bool {
	select {
	case d.queue <- event:
		return true
	default:
		return false
	}
}

func (d *Dispatcher) Start(ctx context.Context) error {
	return nil
}
```

#### Concurrency Semantics
- Notifications are non-blocking; callers push to buffered channel (`capacity: 256`).
- Dedicated worker goroutine drains events and issues HTTP POSTs.

#### Error Handling Contract
- HTTP transport retries with exponential backoff on HTTP 5xx responses up to 3 attempts.

---

### 2.17 `internal/tui` `[Planned]`

Bubble Tea terminal user interface main model, batching ticker, and mission control loop.

#### Types & Signatures

```go
package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TUIOptions struct {
	RefreshRate time.Duration
	EnableMouse bool
}

type AppModel struct {
	options TUIOptions
	active  int
	width   int
	height  int
}

func NewAppModel(opts TUIOptions) *AppModel {
	return &AppModel{options: opts}
}

func (m *AppModel) Init() tea.Cmd {
	return nil
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *AppModel) View() string {
	return ""
}

type TickerBatcher struct {
	interval time.Duration
	incoming chan []byte
}

func NewTickerBatcher(interval time.Duration) *TickerBatcher {
	return &TickerBatcher{
		interval: interval,
		incoming: make(chan []byte, 1024),
	}
}

func (tb *TickerBatcher) Push(data []byte) {
	select {
	case tb.incoming <- data:
	default:
	}
}
```

#### Concurrency Semantics
- Rendering performance target: `Target (unmeasured): smooth 60 FPS log rendering at 30Hz ticker batching`.
- Log line streams are ingested asynchronously into `TickerBatcher` and flushed on each tick to prevent event queue flooding.

#### Error Handling Contract
- Terminal window resize messages (`tea.WindowSizeMsg`) adjust pane viewports with boundary safety to prevent panic on terminal collapse.

---

### 2.18 `internal/tui/views` `[Planned]`

Dashboard view, multi-pane live streaming log view, and memory/rule explorer views.

#### Types & Signatures

```go
package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ViewComponent interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (ViewComponent, tea.Cmd)
	View(width, height int) string
}

type DashboardView struct {
	ActiveWorkers int
	PipelineStage string
}

func (v *DashboardView) Init() tea.Cmd { return nil }
func (v *DashboardView) Update(msg tea.Msg) (ViewComponent, tea.Cmd) { return v, nil }
func (v *DashboardView) View(width, height int) string { return "" }

type LogsView struct {
	AutoFollow bool
	LogLines   [][]byte
}

func (v *LogsView) Init() tea.Cmd { return nil }
func (v *LogsView) Update(msg tea.Msg) (ViewComponent, tea.Cmd) { return v, nil }
func (v *LogsView) View(width, height int) string { return "" }

type MemoryView struct {
	SelectedCategory string
	Items            []string
}

func (v *MemoryView) Init() tea.Cmd { return nil }
func (v *MemoryView) Update(msg tea.Msg) (ViewComponent, tea.Cmd) { return v, nil }
func (v *MemoryView) View(width, height int) string { return "" }
```

#### Concurrency Semantics
- View state is mutated exclusively within Bubble Tea's single-threaded event loop.
- Thread-safe buffer snapshots are acquired before passing data into view models.

#### Error Handling Contract
- Malformed ANSI sequences in log lines are sanitized to prevent terminal garbling.

---

## 3. Key Go Dependencies

| Package | Status Tag | Purpose | Benchmark / Target |
| :--- | :--- | :--- | :--- |
| `modernc.org/sqlite` | `[v0.1 Core]` | Pure-Go embedded SQLite engine with WAL support (zero CGO required) | Deterministic WAL event storage |
| `github.com/philippgille/chromem-go` | `[Planned]` | Pure-Go in-memory vector store with persistence | `Target (unmeasured): <2ms vector search latency` |
| `github.com/smacker/go-tree-sitter` | `[Planned]` | Tree-Sitter AST parser for definitions and references | AST parsing across multiple languages |
| `thejerf/suture/v4` | `[Planned]` | Erlang-style supervisor trees with rate-limited backoff | OTP supervision tree lifecycle management |
| `github.com/charmbracelet/bubbletea` | `[Planned]` | Elm-architecture Terminal User Interface framework | `Target (unmeasured): smooth 60 FPS log rendering at 30Hz ticker batching` |
| `github.com/charmbracelet/lipgloss` | `[Planned]` | Terminal styling, borders, and layouts | Pure functional styling |
| `github.com/creack/pty` | `[Planned]` | PTY allocation for interactive CLI worker execution | Linux/macOS pseudo-terminal allocation |
| `github.com/pkoukk/tiktoken-go` | `[v0.1 Core]` | Pure-Go BPE tokenization for `/v1/messages/count_tokens` | Local zero-network token evaluation |
| `github.com/tidwall/gjson` & `sjson` | `[v0.1 Core]` | Byte-slice JSON parsing and modification for reverse proxy | `Target (unmeasured): <100 µs TTFT overhead` |
| `github.com/bmatcuk/doublestar/v4` | `[Planned]` | High-performance glob matching for `.cursor/rules/*.mdc` rules | Glob pattern matching |
| `github.com/landlock-lsm/go-landlock` | `[Planned]` | Unprivileged Linux Landlock LSM sandboxing | Kernel-level unprivileged sandbox |
| `github.com/spf13/cobra` | `[v0.1 Core]` | CLI command routing (`cli-leader start`, `status`, `config`) | POSIX-compliant CLI flag parsing |
