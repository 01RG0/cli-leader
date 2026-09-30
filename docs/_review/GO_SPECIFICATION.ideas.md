# Ideas Migration: GO_SPECIFICATION.md

This ledger tracks all ideas, architecture narratives, configuration templates, or protocol specs moved or redirected from `docs/GO_SPECIFICATION.md` to canonical owner files per `docs/_review/CONTRACT.md`.

---

| Idea | Old Location | New Location |
| :--- | :--- | :--- |
| **Complete Example `cli-leader.yaml` Configuration Block** (Full declarative YAML representation of all gateway, worker, speculative, memory, sandbox, and notification options) | `docs/GO_SPECIFICATION.md` §3 (lines 243–314) | `cli-leader.example.yaml` & `docs/ARCHITECTURE.md` (while Go struct `config.Config` remains in `docs/GO_SPECIFICATION.md`) |
| **Global Exported `pkg/protocol` Package** (Exporting Anthropic API message structs, OpenAI chunk types, and ACP envelopes under public `pkg/`) | `docs/GO_SPECIFICATION.md` §1 (lines 83–87) | `internal/gateway` (Anthropic/OpenAI wire envelopes) and `internal/acp` (ACP JSON-RPC envelopes) per CONTRACT.md §3.2 |
| **Standalone `internal/tui/styles` Subpackage** (Separate package dedicated to Lipgloss themes and styling) | `docs/GO_SPECIFICATION.md` §1 (lines 81–82) | Consolidated into `internal/tui` or `internal/tui/views` per CONTRACT.md §3.2 |
| **External JSON-RPC and SSE Wire Formats** (Low-level HTTP 200 chunk headers, event-stream envelopes, and raw JSON payloads) | `docs/API_SPECIFICATION.md` §1–§3 | `docs/ARCHITECTURE.md` §3 (subsystems 1 & 2) while concrete Go interfaces and handler types reside in `docs/GO_SPECIFICATION.md` |
| **Erlang-Style Worker Fate-Sharing Bundles Narrative** (Conceptual explanation of bundling PTY master + subprocess + auto-reply + ring buffer into a single supervision failure domain) | `docs/GO_SPECIFICATION.md` §2.1 comments | `docs/ARCHITECTURE.md` §3 (Subsystem 3) |
| **Personalized PageRank Mathematics & AST Tag Extraction Concept** (Narrative describing PageRank convergence and Tree-Sitter tag indexing) | `docs/GO_SPECIFICATION.md` §2.5 comments | `docs/ARCHITECTURE.md` §3 (Subsystem 6) and `docs/ECOSYSTEM_INNOVATIONS.md` |
| **Zero-Shot CLI Adapter Discovery via `--help` Introspection** (Dynamic CLI parameter and syntax parser generating Go adapters on the fly) | Implicit in `adapters/generic.go` | `docs/ROADMAP.md` §Phase 7 & `docs/ECOSYSTEM_INNOVATIONS.md` |
| **Distributed Multi-Node Byzantine Quorum Voting** (Consensus mechanisms extending beyond single-host subagent reviews to networked orchestrator swarms) | `docs/GO_SPECIFICATION.md` §2.4 comments | `docs/ECOSYSTEM_INNOVATIONS.md` & `docs/ARCHITECTURE.md` §3 |
