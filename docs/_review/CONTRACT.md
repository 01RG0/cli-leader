# Documentation Overhaul Contract

> **MANDATORY CONTRACT**: All workers, sub-agents, and documentation contributors MUST strictly adhere to this document. No deviations are permitted without an approved amendment to this contract.

---

## 1. Document Ownership Map

Each documentation topic has exactly **ONE** canonical owner file. Information must reside in its owner file; cross-references should be relative links rather than duplicated text.

| Topic | Canonical Owner File | Permitted Scope & Purpose |
| :--- | :--- | :--- |
| **Vision, Quickstart & Overview** | `README.md` | High-level mission, core value proposition, quickstart guide, technology summary table, top-level architecture overview, links to deep documentation. No deep implementation code or low-level interface definitions. |
| **System Design & Architecture** | `docs/ARCHITECTURE.md` | Deep design decisions, architectural component diagrams, control flows, Erlang OTP supervisor trees, state engines (Dual Ledger, SQLite WAL replay), sandbox models, and protocol interactions (MCP & ACP). |
| **Technical Specification & Go Packages**| `docs/GO_SPECIFICATION.md` | Canonical Go package layouts, concrete Go types, structs, interfaces, method signatures, concurrency semantics, and error handling contracts. |
| **Phases, Milestones & Status** | `docs/ROADMAP.md` | Delivery phases (Phase 0 through Phase 4+), milestone deliverables, feature status matrices, timelines, and inter-component dependencies. |
| **Ecosystem Research & Innovations** | `docs/ECOSYSTEM_INNOVATIONS.md` | Comparative analysis, surveys of external projects (Aider, Claude Code, OpenHands, Hermes Agent, LiteLLM, RouteLLM, Zellij, etc.), theoretical inspirations, and algorithmic notes. |
| **Contribution & Developer Guidelines** | `CONTRIBUTING.md` | Contribution workflow, git branch naming, pull request rules, development environment setup, testing requirements, and code review criteria. |

*(Note: `docs/API_SPECIFICATION.md` is considered a legacy/supplementary document. Its internal Go interfaces belong in `docs/GO_SPECIFICATION.md`, and its external MCP/HTTP protocol contracts belong in `docs/ARCHITECTURE.md`.)*

---

## 2. Feature Status Tags

Every feature or capability mentioned anywhere in documentation must have **exactly one** status tag attached to its heading or specification block.

1. **`[v0.1 Core]`**: Fully implemented or part of the immediate minimal viable orchestrator release.
2. **`[Planned]`**: Actively architected and scheduled for upcoming milestone phases.
3. **`[Research]`**: Speculative exploration, theoretical design, or survey of emerging techniques.

> **RULE**: Do NOT use arbitrary tags like `[In Progress]`, `[Beta]`, `[Experimental]`, or `[v0.2]`. Map all features to one of the three canonical tags above.

---

## 3. Canonical Glossary & Term Meanings

To prevent loose or vague usage of borrowed terms, all documents must use the exact definitions below.

### 3.1 Borrowed Systems & AI Terms

- **Saga**: A sequence of local transactions with compensating rollback actions executed in LIFO order to undo non-git side effects (such as created directories, external API calls, or modified system state) when a worker task fails.
- **Temporal / Durable Engine**: A SQLite WAL-backed event replay mechanism that persists step outcomes and allows resumption of long-running workflows across process crashes and restarts without re-executing completed operations.
- **Byzantine Quorum / Consensus**: A multi-agent voting mechanism where independent worker reviews are aggregated with Thompson-sampling-weighted confidence to prevent hallucinated or flawed code merges.
- **OTP Supervision Tree**: An Erlang/Elixir-inspired hierarchical process supervision model implementing `one_for_one` and `one_for_all` strategies with intensity rate limiting and automatic escalation to the Brain.
- **Refinery**: A serialized Bors-style merge queue that verifies candidate worktree branches with automated test suites and mutation testing gates before squash-merging them into the active branch.
- **Branch Racing / Speculative Execution**: Dispatching identical or alternative task prompts to multiple worker worktrees in parallel and selecting the first candidate that passes the Refinery test gate.
- **F2P (Fail-to-Pass)**: A Test-Driven Development verification protocol where an agent first generates a test that fails against the current codebase, and then validates that the proposed patch makes that exact test pass.
- **Thompson Sampling**: A Bayesian multi-armed bandit algorithm that dynamically estimates and updates capability scores (success probabilities) for each worker CLI to route tasks optimally.
- **Dual Ledger**: Separation of execution state into a `TaskLedger` (high-level user goals and invariants) and a `ProgressLedger` (atomic step actions and outputs).
- **Stall Detector**: Cryptographic hashing of step state and outputs to identify repetitive loops, thrashing, or lack of forward progress and trigger emergency intervention.
- **RepoMap**: A Tree-Sitter AST symbol dependency graph ranked by Personalized PageRank to select relevant context snippets within a token budget.
- **Reflexion**: A post-execution verbal reinforcement mechanism that extracts reusable rules and failure diagnoses from failed tasks into `.brain/conventions/`.
- **Soul File**: `.brain/SOUL.md`, a persistent personality, mission, and hard boundary definition injected into the cognitive brain.
- **Landlock / Bubblewrap**: Linux kernel security mechanisms (unprivileged namespaces and LSM filesystem sandboxes) isolating worker subprocess execution.
- **MCP (Model Context Protocol)**: Stdio/SSE protocol used by Claude CLI to invoke cli-leader's orchestrator tools.
- **ACP (Agent Client Protocol)**: JSON-RPC protocol enabling external IDEs (Zed, JetBrains) to interact directly with cli-leader.

### 3.2 Canonical Package Names

All Go code specifications must use these exact package paths under `github.com/01RG0/cli-leader`:

- `cmd/cli-leader`: CLI application entrypoint.
- `internal/config`: YAML and environment variable configuration loader.
- `internal/gateway`: FastHTTP reverse proxy translating `/v1/messages` to upstream providers.
- `internal/gateway/providers`: Upstream API adapters (`openai`, `gemini`, `ollama`).
- `internal/state`: Dual Ledger (`TaskLedger`, `ProgressLedger`), durable replay, and stall detection.
- `internal/mcp`: Model Context Protocol server exposing orchestrator tools to Claude CLI.
- `internal/acp`: Agent Client Protocol server for IDE integration.
- `internal/supervisor`: Erlang OTP supervision tree and process lifecycles.
- `internal/worker`: Worker manager, PTY subprocess execution, circular log buffer, and Git worktrees.
- `internal/worker/adapters`: Worker CLI tool adapters (`aider`, `gemini_cli`, `ollama_cli`, `generic`).
- `internal/speculative`: Speculative branch racing dispatcher.
- `internal/refinery`: Serialized merge queue and compensating Saga coordinator.
- `internal/verification`: Fail-to-Pass verification, AST mutation testing, and consensus review.
- `internal/sandbox`: Bubblewrap and Landlock security containment.
- `internal/memory`: SQLite store, vector embeddings, Tree-Sitter RepoMap, and Thompson Sampling.
- `internal/notify`: Webhook notification dispatcher (Slack, Discord, Telegram).
- `internal/tui`: Charm Bubbletea terminal user interface.
- `internal/tui/views`: Dashboard, live streaming logs, and memory views.

### 3.3 Canonical Configuration Keys

Workers documenting configuration schemas must use the keys defined in `cli-leader.example.yaml`:

- `version`: Configuration schema version string.
- `gateway`: `listen_addr`, `default_provider`, `providers` (`anthropic`, `openai`, `gemini`, `ollama`).
- `workers`: `supervisor` (`max_restarts`, `period_seconds`, `heartbeat_timeout_sec`), and tool registries (`aider`, `gemini`, `ollama`).
- `speculative`: `enabled`, `max_parallel_candidates`.
- `refinery`: `test_command`, `mutation_testing`, `squash_merges`.
- `memory`: `db_path`, `conventions_dir`, `soul_file`, `max_repo_map_tokens`, `max_injected_rule_tokens`, `thompson_sampling`.
- `sandbox`: `tier` (`bwrap`, `landlock`, `none`), `allow_network`.
- `notifications`: `slack_webhook`, `discord_webhook`, `telegram_chat_id`.

### 3.4 Canonical Component Names

- **Brain**: The primary cognitive supervisor (Claude CLI) responsible for planning and delegating.
- **Gateway**: The internal Anthropic-compatible reverse proxy forwarding requests to LLM providers.
- **Supervisor**: The OTP process supervisor managing child worker subprocesses.
- **Worker Manager**: The subsystem managing PTY execution, worker pools, and isolated Git worktrees.
- **Refinery**: The gated CI and mutation verification pipeline controlling branch merges.
- **Durable Engine**: The SQLite WAL state replay engine ensuring fault-tolerant task execution.
- **Memory Store**: The multi-tiered semantic and episodic knowledge base.

---

## 4. Global Documentation Rules

All documentation workers and sub-agents must strictly enforce these global rules:

1. **Docs Only**: Modify only markdown documentation. Never edit Go source code or configuration templates during the documentation overhaul.
2. **Never Delete an Idea**: If an idea, feature, or section belongs to another owner file, move it to `docs/_review/<FILE>.ideas.md` with columns `idea | old location | new location` so the orchestrator can route it. Never discard ideas.
3. **Never Invent Facts**: If a technical specification, parameter, or name is not established in this CONTRACT, do NOT invent it. Output `TODO(verify): <what>` and record the question in `docs/_review/<FILE>.questions.md`.
4. **Target (unmeasured)**: Any claim of performance, throughput, or speed (e.g., "60 FPS", "sub-10ms latency", "zero-overhead") that lacks an empirical benchmark citation must be written as `Target (unmeasured): <claim>` or revised to an architectural description.
5. **Strict File Boundaries**: A worker assigned to `<FILE>` may edit ONLY its assigned file and its corresponding ledger files (`docs/_review/<FILE>.audit.md`, `docs/_review/<FILE>.ideas.md`, `docs/_review/<FILE>.questions.md`). All other files must be treated as READ-ONLY.
6. **Integrity & Verification**: Ensure all Mermaid blocks parse without syntax errors, relative markdown links resolve to existing files, and tables use consistent GitHub-flavored markdown alignment.
