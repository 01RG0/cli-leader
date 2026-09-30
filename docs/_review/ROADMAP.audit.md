# ROADMAP.md Audit Report

This document audits `docs/ROADMAP.md` against `docs/_review/CONTRACT.md` as part of the `cli-leader` documentation overhaul.

---

## 1. Audit Summary

| Metric | Count |
| :--- | :--- |
| **Total Audit Findings** | 18 |
| **Non-Canonical Status Tags** | 7 |
| **Missing Status Tags & Packages** | 5 |
| **Unmeasured Performance / Metric Claims** | 6 |
| **Architectural / Code Leakage** | 6 |
| **Missing Roadmap Structural Requirements** | 3 |

---

## 2. Detailed Findings

### ROAD-01: Non-Canonical Version & Milestone Status Tags
- **Location**: `docs/ROADMAP.md` (Lines 3, 18, 33, 48, 64, 86, 96)
- **Exact Quote**:
  ```markdown
  ## Phase 1: Foundation & Low-Latency Gateway (v0.1.0)
  ## Phase 2: Worker Subprocess Engine, OTP Supervision & Worktrees (v0.2.0)
  ## Phase 3: MCP Server, ACP Protocol & Temporal Durable WAL Engine (v0.3.0)
  ## Phase 4: Ever-Learning Cognitive Engine, Tree-Sitter RepoMap & Skills (v0.4.0)
  ## Phase 5: Speculative Branch Racing, Sagas & Mutation Testing Gate (v0.5.0)
  ## Phase 6: Bubble Tea TUI Mission Control & Notifications (v0.6.0)
  ## Phase 7: Dynamic CLI Discovery & Polishing (v1.0.0)
  ```
- **Violation**: CONTRACT Section 2 states:
  > *"Every feature or capability mentioned anywhere in documentation must have exactly one status tag attached to its heading or specification block: `[v0.1 Core]`, `[Planned]`, or `[Research]`. RULE: Do NOT use arbitrary tags like `[In Progress]`, `[Beta]`, `[Experimental]`, or `[v0.2]`."*
  Arbitrary semver tags (`(v0.1.0)`, `(v0.2.0)`, etc.) were used in place of canonical status tags.
- **Resolution**: Reorganize into canonical phases (Phase 0 through Phase 4) and tag every milestone feature with exactly one of `[v0.1 Core]`, `[Planned]`, or `[Research]`.

---

### ROAD-02: Phase Structure Mismatch
- **Location**: `docs/ROADMAP.md` (Lines 3–100)
- **Exact Quote**:
  ```markdown
  ## Phase 1: Foundation & Low-Latency Gateway (v0.1.0)
  ...
  ## Phase 7: Dynamic CLI Discovery & Polishing (v1.0.0)
  ```
- **Violation**: Worker Roadmap mandate and CONTRACT Section 1 specify delivery phases structured by:
  - Phase 0: Foundation / v0.1 Core
  - Phase 1: Robust Swarm & Gateway
  - Phase 2: Refinery & Speculative Racing
  - Phase 3: Cognitive Memory & Ecosystem Integrations
  - Phase 4: Production Hardening & Operational Control
  The legacy document used 7 ad-hoc phases that fragmented core foundation and gateway milestones.
- **Resolution**: Restructure roadmap sections strictly into Phase 0 through Phase 4.

---

### ROAD-03: Missing Canonical Go Package Mappings
- **Location**: `docs/ROADMAP.md` (All phases, lines 4–100)
- **Exact Quote**:
  ```markdown
  - [ ] Initialize Go module (`github.com/01RG0/cli-leader`).
  - [ ] Implement CLI entrypoint with `cobra` (`cli-leader start`, `cli-leader config`).
  - [ ] Subprocess execution engine using `os/exec` and `creack/pty.StartWithSize`.
  - [ ] Erlang OTP Supervision Trees (`thejerf/suture/v4`)
  ```
- **Violation**: Worker Roadmap mandate and CONTRACT Section 3.2 require every feature to explicitly map to its canonical Go package under `github.com/01RG0/cli-leader` (e.g. `cmd/cli-leader`, `internal/gateway`, `internal/supervisor`, `internal/refinery`).
- **Resolution**: Map every milestone deliverable to its canonical Go package path defined in `docs/GO_SPECIFICATION.md`.

---

### ROAD-04: Contradictory Gateway Implementation Details
- **Location**: `docs/ROADMAP.md` (Line 7)
- **Exact Quote**:
  ```markdown
  - [ ] HTTP listener on `:8082` using standard library `http.NewServeMux`.
  ```
- **Violation**: CONTRACT Section 3.2, `docs/GO_SPECIFICATION.md`, and `docs/ARCHITECTURE.md` establish that `internal/gateway` uses FastHTTP (`proxy.go # HTTP reverse proxy server (:8082, FastHTTP)`). Using `http.NewServeMux` contradicts the specification and leaks low-level implementation details into the roadmap.
- **Resolution**: Remove conflicting `http.NewServeMux` reference. Align with canonical FastHTTP gateway package `internal/gateway`.

---

### ROAD-05: Protocol Envelopes Leaked into Roadmap
- **Location**: `docs/ROADMAP.md` (Line 10)
- **Exact Quote**:
  ```markdown
  - [ ] Build the **SSE Streaming State Machine** (`content_block_start` $\to$ `content_block_delta` $\to$ `content_block_stop` $\to$ `message_stop`).
  ```
- **Violation**: CONTRACT Section 1 (Document Ownership Map): Protocol interaction schemas and wire format deep dives belong in `docs/ARCHITECTURE.md` and `docs/API_SPECIFICATION.md`. The roadmap should describe milestone deliverables and acceptance criteria rather than protocol sequence specs.
- **Resolution**: Describe the milestone capability (`internal/gateway: SSE Streaming State Machine`) and link to `docs/ARCHITECTURE.md#subsystem-2-go-multi-provider-gateway-anthropic-protocol-emulator` for protocol event breakdowns.

---

### ROAD-06: Code Snippets & Syscall Leaks in Process Engine Deliverables
- **Location**: `docs/ROADMAP.md` (Line 20)
- **Exact Quote**:
  ```markdown
  - [ ] Linux process group isolation: `cmd.SysProcAttr.Setsid = true` and `Pdeathsig = syscall.SIGTERM`.
  ```
- **Violation**: CONTRACT Section 1: Concrete Go code, structs, and syscall configurations belong exclusively in `docs/GO_SPECIFICATION.md` (`internal/worker/pty.go`).
- **Resolution**: Rephrase as a milestone deliverable (`Linux Process Group Isolation and Lifecycle Teardown`) under canonical package `internal/worker`, linking to `docs/GO_SPECIFICATION.md`.

---

### ROAD-07: Unmeasured Claim: Zero Duplicate LLM Tokens
- **Location**: `docs/ROADMAP.md` (Line 38)
- **Exact Quote**:
  ```markdown
  - [ ] Deterministic replay cache: resume interrupted workflows with **zero duplicate LLM tokens**.
  ```
- **Violation**: CONTRACT Section 4, Rule 4 states:
  > *"Any claim of performance, throughput, or speed (e.g., '60 FPS', 'sub-10ms latency', 'zero-overhead') that lacks an empirical benchmark citation must be written as `Target (unmeasured): <claim>` or revised to an architectural description."*
  Claiming "zero duplicate LLM tokens" is an unmeasured promise without empirical validation.
- **Resolution**: Revise to `Target (unmeasured): zero duplicate LLM tokens during workflow resumption via SQLite WAL replay cache`.

---

### ROAD-08: Unmeasured Claim: Boomerang Packet Token Bound
- **Location**: `docs/ROADMAP.md` (Line 44)
- **Exact Quote**:
  ```markdown
  - [ ] Boomerang subtask packet squashing (<250 tokens).
  ```
- **Violation**: CONTRACT Section 4, Rule 4: Claiming `<250 tokens` without empirical benchmark verification violates the rule against unmeasured performance promises.
- **Resolution**: Rewrite as `Target (unmeasured): Boomerang subtask packet squashing to <250 tokens`.

---

### ROAD-09: Unmeasured Claim: Tree-Sitter RepoMap Context Budget
- **Location**: `docs/ROADMAP.md` (Line 54)
- **Exact Quote**:
  ```markdown
  - [ ] Render ultra-compact AST context (<1024 tokens) for Claude Brain.
  ```
- **Violation**: CONTRACT Section 4, Rule 4: Claiming `<1024 tokens` without qualification is an unmeasured target.
- **Resolution**: Rephrase as `Target (unmeasured): AST context serialization capped at <1024 tokens`.

---

### ROAD-10: Unmeasured Claim: Dynamic Rule Injection Budget
- **Location**: `docs/ROADMAP.md` (Line 58)
- **Exact Quote**:
  ```markdown
  - [ ] Dynamic rule injector: Glob matching (`bmatcuk/doublestar`) + vector similarity, capped at `<800` tokens.
  ```
- **Violation**: CONTRACT Section 4, Rule 4: Token limit claims require the unmeasured target qualifier.
- **Resolution**: Rephrase as `Target (unmeasured): dynamic rule prompt injection capped at <800 tokens`.

---

### ROAD-11: Unmeasured Latency Claim: Sub-5ms Semantic Centroid Pre-Routing
- **Location**: `docs/ROADMAP.md` (Line 60)
- **Exact Quote**:
  ```markdown
  - [ ] Two-Stage Router: Sub-5ms semantic centroid pre-routing + Thompson Sampling Beta capability matrix.
  ```
- **Violation**: CONTRACT Section 4, Rule 4: "Sub-5ms" latency claim is unmeasured.
- **Resolution**: Rewrite as `Target (unmeasured): Sub-5ms semantic centroid pre-routing combined with Bayesian Thompson Sampling`.

---

### ROAD-12: Unmeasured UI Rendering Claim: 60 FPS Terminal Output
- **Location**: `docs/ROADMAP.md` (Line 89)
- **Exact Quote**:
  ```markdown
  - [ ] 30Hz ticker batcher for smooth 60 FPS rendering without terminal freezing.
  ```
- **Violation**: CONTRACT Section 4, Rule 4: Claiming "smooth 60 FPS rendering without terminal freezing" lacks empirical benchmark citations.
- **Resolution**: Rewrite as `Target (unmeasured): 30Hz ticker batcher targeting 60 FPS terminal rendering under high PTY throughput`.

---

### ROAD-13: Inconsistent Tool and Component Naming
- **Location**: `docs/ROADMAP.md` (Line 29, Line 66)
- **Exact Quote**:
  ```markdown
  - [ ] Built-in adapters: `aider`, `gemini-cli`, `ollama`, and generic CLI adapter.
  - [ ] MCP tool: `dispatch_speculative_candidates`.
  ```
- **Violation**: CONTRACT Section 3.2 lists `internal/worker/adapters/gemini_cli.go` and `ollama_cli.go`. `docs/API_SPECIFICATION.md` Section 2.2 defines the speculative tool name as `dispatch_speculative`, not `dispatch_speculative_candidates`.
- **Resolution**: Align names with CONTRACT and `API_SPECIFICATION.md`: `gemini_cli`, `ollama_cli`, and tool name `dispatch_speculative`.

---

### ROAD-14: Redundant Acronym in ACP Mention
- **Location**: `docs/ROADMAP.md` (Line 35)
- **Exact Quote**:
  ```markdown
  - [ ] Implement **Agent Client Protocol (ACP)** JSON-RPC server for Zed and JetBrains IDE integration.
  ```
- **Violation**: Minor naming redundancy ("ACP Protocol" in heading on line 33). ACP expands to Agent Client Protocol.
- **Resolution**: Use the canonical name `Agent Client Protocol (ACP)` without redundant "Protocol" suffix in headings.

---

### ROAD-15: Missing Milestone Dependency Graph
- **Location**: `docs/ROADMAP.md` (Entire document)
- **Exact Quote**: *N/A (Missing component)*
- **Violation**: Mandate Step 5 requires:
  > *"Include a clear milestone dependency graph (Mermaid flowchart) and deliverable acceptance criteria."*
  The legacy roadmap contained no dependency graph or Mermaid visualization.
- **Resolution**: Add an explicit Mermaid flowchart showing the dependency relationships between milestone deliverables across Phase 0 through Phase 4.

---

### ROAD-16: Missing Deliverable Acceptance Criteria
- **Location**: `docs/ROADMAP.md` (Entire document)
- **Exact Quote**: *N/A (Missing component)*
- **Violation**: Mandate Step 5 requires measurable deliverable acceptance criteria for each phase. The legacy document had only bulleted task lists with arbitrary checkmarks.
- **Resolution**: Define concrete, testable acceptance criteria for every phase (e.g. CLI test commands, verification receipts, error code validations).

---

### ROAD-17: Low-Level Git & Mutex Details in Worktree Deliverables
- **Location**: `docs/ROADMAP.md` (Lines 27–28)
- **Exact Quote**:
  ```markdown
  - [ ] Wrap native `git worktree add`, `git add -A`, `git diff --cached`, `git worktree remove --force`.
  - [ ] Serialize worktree operations with `sync.Mutex` to eliminate `.git/config.lock` collisions.
  ```
- **Violation**: CONTRACT Section 1: Detailed Go mutex locks and git command invocations belong in `docs/GO_SPECIFICATION.md` under `internal/worker/worktree.go`.
- **Resolution**: Summarize as `Native Git Worktree Isolation & Mutex Concurrency Control` under `internal/worker` and link to `docs/GO_SPECIFICATION.md`.

---

### ROAD-18: Missing Verification of Release Timelines and Schedules
- **Location**: `docs/ROADMAP.md` (All phases)
- **Exact Quote**: *N/A (No dates or delivery calendar)*
- **Violation**: CONTRACT Section 4, Rule 3 states:
  > *"If a technical specification, parameter, or name is not established in this CONTRACT, do NOT invent it. Output `TODO(verify): <what>` and record the question in `docs/_review/<FILE>.questions.md`."*
  Dates and milestone target schedules were omitted or unverified.
- **Resolution**: Mark schedule timelines with `TODO(verify): target delivery timeline for Phase <X>` and log open questions in `docs/_review/ROADMAP.questions.md`.
