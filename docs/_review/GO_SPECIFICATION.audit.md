# Audit: GO_SPECIFICATION.md against CONTRACT.md

This document audits `docs/GO_SPECIFICATION.md` against the mandates in `docs/_review/CONTRACT.md`.

---

## Audit Findings

### GOSPEC-01: Missing Status Tags Across All Sections and Packages
- **Contract Reference**: Section 2 ("Every feature or capability mentioned anywhere in documentation must have exactly one status tag attached to its heading or specification block: `[v0.1 Core]`, `[Planned]`, or `[Research]`.")
- **Exact Quotes from `docs/GO_SPECIFICATION.md`**:
  - Line 1: `# Go Technical Specification & Package Layout`
  - Line 3: `## 1. Standard Project Layout`
  - Line 96: `## 2. Core Go Interfaces & Structs`
  - Line 98: `### 2.1 Erlang OTP Supervisor Tree`
  - Line 135: `### 2.2 Temporal-Grade Durable Engine (SQLite WAL)`
  - Line 166: `### 2.3 Saga Non-Git Side-Effect Coordinator`
  - Line 192: `### 2.4 Byzantine Quorum Consensus Arbiter`
  - Line 218: `### 2.5 Tree-Sitter PageRank Repo Map Engine`
  - Line 243: `## 3. Configuration Schema (cli-leader.yaml)`
  - Line 318: `## 4. Key Go Dependencies`
- **Violation**: Zero canonical status tags exist in the file.
- **Resolution**: Apply `[v0.1 Core]`, `[Planned]`, or `[Research]` to every section, package specification, interface, and capability block based on `docs/ROADMAP.md` and `docs/_review/CONTRACT.md`.

---

### GOSPEC-02: Non-Canonical Package Layout (`pkg/protocol`)
- **Contract Reference**: Section 3.2 ("All Go code specifications must use these exact package paths under `github.com/01RG0/cli-leader`...")
- **Exact Quote from `docs/GO_SPECIFICATION.md`**:
  - Lines 83–87:
    ```
    ├── pkg/
    │   └── protocol/
    │       ├── anthropic.go            # Anthropic API message structs & SSE types
    │       ├── openai.go               # OpenAI API message structs & chunks
    │       └── acp.go                  # Agent Client Protocol message envelopes
    ```
- **Violation**: `pkg/protocol` is not in CONTRACT.md Section 3.2. All types must reside in canonical packages (`internal/gateway`, `internal/mcp`, `internal/acp`).
- **Resolution**: Remove `pkg/protocol` from package tree. Relocate types to their canonical internal packages (`internal/gateway` for Anthropic/OpenAI wire envelopes, `internal/acp` for ACP envelopes, and `internal/mcp` for MCP types).

---

### GOSPEC-03: Non-Canonical Subpackage (`internal/tui/styles`)
- **Contract Reference**: Section 3.2 ("`internal/tui`: Charm Bubbletea terminal user interface", "`internal/tui/views`: Dashboard, live streaming logs, and memory views")
- **Exact Quote from `docs/GO_SPECIFICATION.md`**:
  - Lines 81–82:
    ```
    │       └── styles/
    │           └── theme.go            # Lipgloss styling, borders, and colors
    ```
- **Violation**: `internal/tui/styles` is not defined as an independent canonical package in Section 3.2.
- **Resolution**: Consolidate styling definitions directly into `internal/tui` or `internal/tui/views`.

---

### GOSPEC-04: Unmeasured Performance Claim in TUI Batcher
- **Contract Reference**: Section 4 Rule 4 ("Any claim of performance, throughput, or speed (e.g., '60 FPS', 'sub-10ms latency', 'zero-overhead') that lacks an empirical benchmark citation must be written as `Target (unmeasured): <claim>` or revised to an architectural description.")
- **Exact Quote from `docs/GO_SPECIFICATION.md`**:
  - Line 76: `├── batcher.go # 30Hz ticker batcher for smooth 60 FPS logs`
- **Violation**: "60 FPS" is an unmeasured throughput claim without benchmark citation.
- **Resolution**: Rephrase to `Target (unmeasured): smooth 60 FPS log rendering at 30Hz ticker batching`.

---

### GOSPEC-05: Unmeasured Performance Claim in Dependencies Table
- **Contract Reference**: Section 4 Rule 4 ("Any claim of performance, throughput, or speed... must be written as `Target (unmeasured): <claim>`...")
- **Exact Quote from `docs/GO_SPECIFICATION.md`**:
  - Line 323: `| github.com/philippgille/chromem-go | Pure Go in-memory vector store with persistence (<2ms search) |`
- **Violation**: "(<2ms search)" is an unmeasured latency assertion.
- **Resolution**: Convert to `Target (unmeasured): <2ms vector search latency`.

---

### GOSPEC-06: Narrative Architectural Scope Belonging in ARCHITECTURE.md
- **Contract Reference**: Section 1 Document Ownership Map (`docs/ARCHITECTURE.md` owns system design, control flows, state engine design; `docs/GO_SPECIFICATION.md` owns Go package layouts, types, structs, interfaces, method signatures, concurrency semantics, and error handling contracts).
- **Exact Quote from `docs/GO_SPECIFICATION.md`**:
  - Lines 243–314: Full YAML configuration example file block `cli-leader.yaml` embedded directly in the technical specification without corresponding Go configuration structs, validation methods, or environment variable parsing contracts.
- **Violation**: Raw YAML configuration text belongs in `cli-leader.example.yaml` and `docs/ARCHITECTURE.md`. `docs/GO_SPECIFICATION.md` must specify the Go struct `internal/config.Config` and unmarshaling contracts.
- **Resolution**: Log redirection in `GO_SPECIFICATION.ideas.md`. Replace raw YAML in `GO_SPECIFICATION.md` with concrete Go `Config` structs, tags, validation signatures, and concurrency guarantees.

---

### GOSPEC-07: Incomplete Package Coverage (Missing 12 of 17 Canonical Packages)
- **Contract Reference**: Section 1 & Section 3.2 (Defining all packages under `github.com/01RG0/cli-leader`).
- **Exact Quote from `docs/GO_SPECIFICATION.md`**:
  - Section 2 only defines 5 packages (`supervisor`, `state`, `refinery`, `verification`, `memory`).
  - Missing entirely from Section 2:
    1. `cmd/cli-leader`
    2. `internal/config`
    3. `internal/gateway`
    4. `internal/gateway/providers`
    5. `internal/mcp`
    6. `internal/acp`
    7. `internal/worker`
    8. `internal/worker/adapters`
    9. `internal/speculative`
    10. `internal/sandbox`
    11. `internal/notify`
    12. `internal/tui` & `internal/tui/views`
- **Violation**: Over 70% of the project's canonical packages have no Go type definitions, interface contracts, or method signatures.
- **Resolution**: Add comprehensive Go specifications with concrete structs, interfaces, and signatures for all 17 canonical packages.

---

### GOSPEC-08: Incomplete Type Definitions in Defined Packages
- **Contract Reference**: Section 1 ("concrete Go types, structs, interfaces, method signatures...").
- **Exact Quote from `docs/GO_SPECIFICATION.md`**:
  - Lines 135–164: Only `DurableEngine` is defined; `TaskLedger`, `ProgressLedger`, and `StallDetector` are missing.
  - Lines 98–133: `WorkerProcess` struct is referenced (`children map[string]*WorkerProcess`) but never defined; process lifecycle methods (`Start`, `Stop`, `Restart`) are omitted.
  - Lines 166–190: Only `SagaCoordinator` is defined; `MergeQueue` and test verification gate runners are omitted.
  - Lines 192–216: Only `AgentVote` and `ConsensusArbiter` are defined; `F2PEngine` and `MutationTester` are omitted.
  - Lines 218–239: Only `RepoMapEngine` and `SymbolTag` are defined; `Store`, `VectorIndex`, `ReflexionEngine`, `RuleMatcher`, `SkillCatalog`, and `ThompsonSamplingMatrix` are omitted.
- **Violation**: Key components described in `docs/ARCHITECTURE.md` have missing Go interface definitions.
- **Resolution**: Fully specify all missing structs, interfaces, and methods for each package.

---

### GOSPEC-09: Total Absence of Concurrency Semantics
- **Contract Reference**: Section 1 ("concurrency semantics").
- **Exact Quote from `docs/GO_SPECIFICATION.md`**:
  - Mutexes (`sync.Mutex`) are placed into structs (`Supervisor`, `DurableEngine`, `SagaCoordinator`), but zero concurrency semantics are documented: no lock acquisition order, channel buffering rules, goroutine lifecycles, or context cancellation propagation.
- **Violation**: The document fails to define concurrency guarantees and safety invariants.
- **Resolution**: Document concurrency semantics explicitly for every package (goroutine boundaries, synchronization primitives, lock ordering, channel capacities, thread safety).

---

### GOSPEC-10: Total Absence of Error Handling Contracts
- **Contract Reference**: Section 1 ("error handling contracts").
- **Exact Quote from `docs/GO_SPECIFICATION.md`**:
  - Methods return `error` with no documented sentinel errors, custom error types, retryable classification, or wrapping semantics.
- **Violation**: No error handling contract is established.
- **Resolution**: Define explicit sentinel errors (`var Err... = errors.New(...)`), typed domain errors, wrapping conventions (`%w`), and recovery behaviors for each package.

---

### GOSPEC-11: Missing Internal Go Interfaces from Legacy `API_SPECIFICATION.md`
- **Contract Reference**: Section 1 Note ("`docs/API_SPECIFICATION.md` is considered a legacy/supplementary document. Its internal Go interfaces belong in `docs/GO_SPECIFICATION.md`...").
- **Exact Quote from `docs/API_SPECIFICATION.md`**:
  - Sections 1, 2, and 3 define MCP tools (`dispatch_worker`, `dispatch_speculative`, `verify_f2p`, `refinery_enqueue`, `query_brain_memory`), ACP JSON-RPC envelopes (`initialize`, `agent/sendMessage`), and Gateway proxy handlers.
- **Violation**: These internal Go interfaces were not migrated to `docs/GO_SPECIFICATION.md`.
- **Resolution**: Model concrete Go types, handler interfaces, and dispatch structs for MCP (`internal/mcp`), ACP (`internal/acp`), and Gateway (`internal/gateway`) inside `docs/GO_SPECIFICATION.md`.
