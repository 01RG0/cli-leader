# Project Roadmap & Implementation Milestones

## Phase 1: Foundation & Low-Latency Gateway (v0.1.0)
- [ ] Initialize Go module (`github.com/01RG0/cli-leader`).
- [ ] Implement CLI entrypoint with `cobra` (`cli-leader start`, `cli-leader config`).
- [ ] Build the **Low-Latency Anthropic Reverse Proxy Gateway**:
  - [ ] HTTP listener on `:8082` using standard library `http.NewServeMux`.
  - [ ] Implement `POST /v1/messages/count_tokens` locally via `tiktoken-go` (prevents Claude CLI pre-flight aborts).
  - [ ] Implement request transformer: Anthropic `messages` + `tools` to OpenAI/Gemini/Ollama payload format.
  - [ ] Build the **SSE Streaming State Machine** (`content_block_start` $\to$ `content_block_delta` $\to$ `content_block_stop` $\to$ `message_stop`).
  - [ ] Map DeepSeek-R1 `reasoning_content` to Anthropic `thinking` blocks.
  - [ ] Handle tool definitions (`input_schema` $\to$ `parameters`) and sanitize Gemini tool names.
  - [ ] Optimize JSON mutations with `tidwall/gjson` and `tidwall/sjson` (strip `cache_control`).
- [ ] Verification: Launch Claude CLI pointing `ANTHROPIC_BASE_URL=http://localhost:8082` against OpenAI/Ollama models with full tool use.

---

## Phase 2: Worker Subprocess Engine, OTP Supervision & Worktrees (v0.2.0)
- [ ] Subprocess execution engine using `os/exec` and `creack/pty.StartWithSize`.
- [ ] Linux process group isolation: `cmd.SysProcAttr.Setsid = true` and `Pdeathsig = syscall.SIGTERM`.
- [ ] Build **Erlang OTP Supervision Trees** (`thejerf/suture/v4`):
  - [ ] `one_for_one` strategy for isolated speculative candidates.
  - [ ] `one_for_all` strategy for compound worker bundles (PTY + Subprocess + Reader + Buffer).
  - [ ] Restart intensity rate limiting ($M$ crashes in $T$ seconds) bubbling to Claude Brain.
- [ ] Sliding-window regex prompt auto-responder (intercepting `[y/N]` and `Apply changes?` prompts).
- [ ] Git Worktree sandbox manager:
  - [ ] Wrap native `git worktree add`, `git add -A`, `git diff --cached`, `git worktree remove --force`.
  - [ ] Serialize worktree operations with `sync.Mutex` to eliminate `.git/config.lock` collisions.
- [ ] Built-in adapters: `aider`, `gemini-cli`, `ollama`, and generic CLI adapter.

---

## Phase 3: MCP Server, ACP Protocol & Temporal Durable WAL Engine (v0.3.0)
- [ ] Implement Model Context Protocol (MCP) server over `stdio` and local SSE.
- [ ] Implement **Agent Client Protocol (ACP)** JSON-RPC server for Zed and JetBrains IDE integration.
- [ ] Build **Temporal-Grade Durable WAL Engine**:
  - [ ] Append-only event log (`workflow_events` table in SQLite WAL).
  - [ ] Deterministic replay cache: resume interrupted workflows with **zero duplicate LLM tokens**.
  - [ ] Activity heartbeating: terminate hung CLI subprocesses after 15 seconds of silence.
- [ ] Build the **Dual-Ledger State Machine**:
  - [ ] `TaskLedger`: Immutable high-level goals and constraints in SQLite.
  - [ ] `ProgressLedger`: Dynamic sub-step tracking.
  - [ ] Cryptographic state hasher (`SHA256(tool + args + diff)`) and stall counter to prevent repair loops.
- [ ] Boomerang subtask packet squashing (<250 tokens).

---

## Phase 4: Ever-Learning Cognitive Engine, Tree-Sitter RepoMap & Skills (v0.4.0)
- [ ] Pure-Go SQLite memory schema using `modernc.org/sqlite` (zero CGO).
- [ ] Pure-Go embedded vector index using `github.com/philippgille/chromem-go`.
- [ ] Build **Tree-Sitter Repo Map Engine** (`smacker/go-tree-sitter`):
  - [ ] Extract definitions and references across languages.
  - [ ] Run Personalized PageRank to rank symbol importance.
  - [ ] Render ultra-compact AST context (<1024 tokens) for Claude Brain.
- [ ] Build the **Reflexion Verbal Learning Loop**:
  - [ ] Parse compiler errors, panics, and test failures into verbal directives (`[WHEN-DO-BECAUSE]`).
  - [ ] Persist rules in `.brain/conventions/*.mdc` with YAML frontmatter (`globs`, `description`).
  - [ ] Dynamic rule injector: Glob matching (`bmatcuk/doublestar`) + vector similarity, capped at `<800` tokens.
- [ ] Hermes-style **Self-Evolving Skills** (`.brain/skills/*.yaml`) and OpenClaw declarative `SOUL.md`.
- [ ] Two-Stage Router: Sub-5ms semantic centroid pre-routing + Thompson Sampling Beta capability matrix.

---

## Phase 5: Speculative Branch Racing, Sagas & Mutation Testing Gate (v0.5.0)
- [ ] **Speculative Multi-Worktree Execution ("Branch Racing")**:
  - [ ] MCP tool: `dispatch_speculative_candidates`.
  - [ ] Concurrently spawn 2 workers in separate worktrees; auto-evaluate diffs and test passes.
- [ ] **Saga Non-Git Side-Effect Coordinator**:
  - [ ] LIFO compensation stack for package installs (`npm`, `pip`), Docker daemons, and database migrations.
  - [ ] Backward compensation executes automatically when candidate branches are abandoned.
- [ ] **Fail-to-Pass (F2P) Automated TDD Protocol**:
  - [ ] Assert reproduction test fails (`exit_code != 0`) on unmodified HEAD.
  - [ ] Assert reproduction test passes (`exit_code == 0`) and regression suite passes on fix.
- [ ] **Mutation Testing Verification Gate (`cargo-mutants` style)**:
  - [ ] Introduce synthetic AST mutations into modified code.
  - [ ] Verify test suite kills all mutants; trigger test synthesis on survived mutants.
- [ ] **Byzantine Quorum Consensus**:
  - [ ] Require machine receipts (exit code 0, Git diff hash, F2P proof receipts, linter zero-errors).
  - [ ] Thompson-weighted agent review voting.
- [ ] **Bors-Style Refinery Merge Queue**:
  - [ ] Serialized merge queue rebasing completed worktrees on latest HEAD before merging.
- [ ] Tiered sandboxing: Bubblewrap (`bwrap`) with fallback to Linux Landlock LSM (`go-landlock`).

---

## Phase 6: Bubble Tea TUI Mission Control & Notifications (v0.6.0)
- [ ] Real-time terminal dashboard with Charm's `bubbletea` and `lipgloss`:
  - [ ] Thread-safe circular ring buffer (`RingLogBuffer`) decoupling high-rate PTY outputs.
  - [ ] 30Hz ticker batcher for smooth 60 FPS rendering without terminal freezing.
  - [ ] Multi-pane layout: Swarm Status, Live Worker PTY Logs, Git Diff Reviewer, Cognitive Memory Stream.
  - [ ] Keyboard navigation (`Tab`, `j/k`, `f` auto-follow, `Enter`).
- [ ] OpenClaw-style async webhook notifications (Slack, Discord, Telegram) for long-running branch races.

---

## Phase 7: Dynamic CLI Discovery & Polishing (v1.0.0)
- [ ] Host CLI auto-discovery (`which aider`, `which ollama`, `which gh`, `which cursor`, etc.).
- [ ] Zero-shot adapter generator (interrogating `--help` outputs).
- [ ] Cross-platform static binary releases (Linux AMD64/ARM64, macOS Apple Silicon, Windows).
- [ ] Complete test suite and integration benchmarks.
