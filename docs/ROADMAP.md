# Project Roadmap & Implementation Milestones

## Phase 1: Foundation & Gateway (v0.1.0)
- [ ] Initialize Go module (`github.com/01RG0/cli-leader`).
- [ ] Implement CLI entrypoint with `cobra` (`cli-leader start`, `cli-leader config`).
- [ ] Build the **Anthropic-to-OpenAI/Gemini Reverse Proxy Gateway**:
  - [ ] HTTP listener on `:8082`.
  - [ ] Request body parser for `/v1/messages`.
  - [ ] Translation layer to OpenAI Chat Completions API format.
  - [ ] Server-Sent Events (SSE) streaming translator back to Anthropic client format.
- [ ] Verification: Launch Claude CLI pointing `ANTHROPIC_BASE_URL=http://localhost:8082` against OpenAI/Ollama models.

---

## Phase 2: Worker Subprocess Engine & Worktrees (v0.2.0)
- [ ] Subprocess execution engine using `os/exec` and `creack/pty`.
- [ ] Non-interactive prompt handler (intercepting `y/n` prompts).
- [ ] Git Worktree sandbox manager:
  - [ ] Create isolated worktree for tasks.
  - [ ] Capture git diffs upon task completion.
  - [ ] Revert or squash-merge worktree branches.
- [ ] Built-in adapters for initial workers:
  - [ ] `aider` adapter.
  - [ ] `gemini-cli` adapter.
  - [ ] Generic shell/CLI adapter.

---

## Phase 3: MCP Server & Claude Brain Integration (v0.3.0)
- [ ] Implement Model Context Protocol (MCP) server over `stdio` and SSE.
- [ ] Register core tools:
  - [ ] `dispatch_worker`
  - [ ] `get_task_status`
  - [ ] `inspect_diff`
  - [ ] `merge_or_reject`
- [ ] Verification: Claude CLI autonomously invokes and inspects worker CLIs via MCP.

---

## Phase 4: Ever-Learning Cognitive Engine (v0.4.0)
- [ ] Pure-Go SQLite memory database schema (`modernc.org/sqlite`).
- [ ] Episodic logger (task duration, tokens, diffs, exit codes).
- [ ] Semantic rule engine:
  - [ ] Auto-extract bug fixes and project rules into `.brain/conventions.md`.
  - [ ] Inject relevant conventions into worker prompts prior to execution.
- [ ] Dynamic CLI Capability Matrix:
  - [ ] Track success rates and latency per category.
  - [ ] Auto-recommend best worker CLI based on task type.

---

## Phase 5: Adversarial Peer Review & Red Team (v0.5.0)
- [ ] Two-phase execution loop: Builder -> Adversary -> Judge.
- [ ] Adversarial prompt generator (instructing breaker CLI to find regressions and edge cases).
- [ ] Automated rollback if test suite fails or breaker proves defects.

---

## Phase 6: Bubble Tea TUI Mission Control (v0.6.0)
- [ ] Real-time terminal dashboard with Charm's `bubbletea` and `lipgloss`:
  - [ ] Active task pipeline and worker status indicators.
  - [ ] Live scrolling terminal log viewer for running workers.
  - [ ] Real-time cognitive memory feed (recent rules learned and applied).
  - [ ] Cost and token meter.

---

## Phase 7: Dynamic CLI Discovery & Polishing (v1.0.0)
- [ ] Host CLI auto-discovery (`which aider`, `which ollama`, `which gh`, etc.).
- [ ] Zero-shot adapter generator (interrogating `--help` outputs).
- [ ] Cross-platform binary releases (Linux AMD64/ARM64, macOS Apple Silicon, Windows).
- [ ] Complete test suite and integration benchmarks.
