# Project Roadmap & Implementation Milestones

## Phase 1: Foundation & Multi-Provider Gateway (v0.1.0)
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

## Phase 2: Worker Subprocess Engine & Worktrees (v0.2.0)
- [ ] Subprocess execution engine using `os/exec` and `creack/pty.StartWithSize`.
- [ ] Linux process group isolation: `cmd.SysProcAttr.Setsid = true` and `Pdeathsig = syscall.SIGTERM`.
- [ ] Graceful teardown escalation: `SIGTERM` on `-pgid`, 3-second grace period, `SIGKILL` on `-pgid`.
- [ ] Sliding-window regex prompt auto-responder (intercepting `[y/N]` and `Apply changes?` prompts).
- [ ] Git Worktree sandbox manager:
  - [ ] Wrap native `git worktree add`, `git add -A`, `git diff --cached`, `git worktree remove --force`.
  - [ ] Serialize worktree operations with `sync.Mutex` to eliminate `.git/config.lock` collisions.
- [ ] Built-in adapters: `aider`, `gemini-cli`, `ollama`, and generic CLI adapter.

---

## Phase 3: MCP Server & Dual-Ledger State Machine (v0.3.0)
- [ ] Implement Model Context Protocol (MCP) server over `stdio` and local SSE.
- [ ] Build the **Dual-Ledger State Machine**:
  - [ ] `TaskLedger`: Immutable high-level goals and constraints in SQLite.
  - [ ] `ProgressLedger`: Dynamic sub-step tracking.
  - [ ] Cryptographic state hasher (`SHA256(tool + args + diff)`) and stall counter to prevent repair loops.
- [ ] Register core MCP tools:
  - [ ] `dispatch_worker`
  - [ ] `get_task_status`
  - [ ] `inspect_diff`
  - [ ] `abort_worker`
- [ ] Verification: Claude CLI autonomously orchestrates worker CLIs without conversational context bloat.

---

## Phase 4: Ever-Learning Cognitive Engine (v0.4.0)
- [ ] Pure-Go SQLite memory schema using `modernc.org/sqlite` (zero CGO).
- [ ] Pure-Go embedded vector index using `github.com/philippgille/chromem-go`.
- [ ] Build the **Reflexion Verbal Learning Loop**:
  - [ ] Parse compiler errors, panics, and test failures into verbal directives (`[WHEN-DO-BECAUSE]`).
  - [ ] Persist rules in `.brain/conventions/*.mdc` with YAML frontmatter (`globs`, `description`).
  - [ ] Dynamic rule injector: Glob matching (`bmatcuk/doublestar`) + vector similarity, capped at `<800` tokens.
- [ ] Implement **Thompson Sampling Capability Matrix**:
  - [ ] Bayesian Beta distributions ($\alpha, \beta$) per worker and task category.
  - [ ] Dynamic routing balancing exploration of local models with exploitation of frontier models.

---

## Phase 5: Speculative Execution, F2P Verification & Refinery (v0.5.0)
- [ ] **Speculative Multi-Worktree Execution ("Branch Racing")**:
  - [ ] MCP tool: `dispatch_speculative_candidates`.
  - [ ] Concurrently spawn 2 workers in separate worktrees; auto-evaluate diffs and test passes.
- [ ] **Fail-to-Pass (F2P) Automated TDD Protocol**:
  - [ ] Generate reproduction test; assert `exit_code != 0` on unmodified HEAD.
  - [ ] Apply fix; assert reproduction test passes (`exit_code == 0`) and existing test suite passes.
- [ ] **Bors-Style Refinery Merge Queue**:
  - [ ] Serialized merge queue rebasing completed worktrees on latest HEAD before merging.
- [ ] Tiered sandboxing: Bubblewrap (`bwrap`) with fallback to Linux Landlock LSM (`go-landlock`).

---

## Phase 6: Bubble Tea TUI Mission Control (v0.6.0)
- [ ] Real-time terminal dashboard with Charm's `bubbletea` and `lipgloss`:
  - [ ] Thread-safe circular ring buffer (`RingLogBuffer`) decoupling high-rate PTY outputs.
  - [ ] 30Hz ticker batcher for smooth 60 FPS rendering without terminal freezing.
  - [ ] Multi-pane layout: Swarm Status, Live Worker PTY Logs, Git Diff Reviewer, Cognitive Memory Stream.
  - [ ] Keyboard navigation (`Tab`, `j/k`, `f` auto-follow, `Enter`).

---

## Phase 7: Dynamic CLI Discovery & Polishing (v1.0.0)
- [ ] Host CLI auto-discovery (`which aider`, `which ollama`, `which gh`, `which cursor`, etc.).
- [ ] Zero-shot adapter generator (interrogating `--help` outputs).
- [ ] Cross-platform static binary releases (Linux AMD64/ARM64, macOS Apple Silicon, Windows).
- [ ] End-to-end integration test suite and benchmarking.
