# Architecture Ideas Routing Log

This document records all ideas, features, sections, or technical details redirected to their canonical owner files per `docs/_review/CONTRACT.md`. **No ideas are deleted.**

---

## Ideas Routing Table

| Idea | Old Location | New Location | Rationale & Owner File Scope |
| :--- | :--- | :--- | :--- |
| Concrete Go package interfaces (`Supervisor`, `DurableEngine`, `SagaCoordinator`, `ConsensusArbiter`, `RepoMapEngine`) | `docs/ARCHITECTURE.md` (Subsystems 3, 4, 5, 6) | `docs/GO_SPECIFICATION.md` | CONTRACT.md Section 1: All concrete Go types, structs, interfaces, method signatures, concurrency semantics, and error handling contracts belong canonically in `docs/GO_SPECIFICATION.md`. |
| Delivery milestones, release numbering (v0.1.0 through v1.0.0), and phase checklists | `docs/ARCHITECTURE.md` (Implicit in subsystem notes) | `docs/ROADMAP.md` | CONTRACT.md Section 1: Delivery phases (Phase 0 through Phase 4+), milestone deliverables, feature status matrices, timelines, and inter-component dependencies belong canonically in `docs/ROADMAP.md`. |
| Extended ecosystem tool comparisons and research analysis (Gastown, Tutti, Hermes Agent, OpenClaw, Maxim Bifrost, RouteLLM, FrugalGPT) | `docs/ARCHITECTURE.md` (Subsystems 1, 2, 6) | `docs/ECOSYSTEM_INNOVATIONS.md` | CONTRACT.md Section 1: Comparative analysis, external surveys, theoretical inspirations, and algorithmic notes belong canonically in `docs/ECOSYSTEM_INNOVATIONS.md`. |
| Full YAML configuration specification schema (`cli-leader.yaml`) | `docs/ARCHITECTURE.md` (Config references) | `docs/GO_SPECIFICATION.md` (Section 3) & `cli-leader.example.yaml` | CONTRACT.md Section 1 & Section 3.3: Configuration schemas and Go struct representations belong in `docs/GO_SPECIFICATION.md`. |
| External HTTP Gateway schemas (`/v1/messages/count_tokens`, `/v1/messages` SSE frames) | `docs/API_SPECIFICATION.md` (Section 1) | `docs/ARCHITECTURE.md` (Subsystem 2 & Protocol Contracts) | CONTRACT.md Section 1 (Note): `docs/API_SPECIFICATION.md` is legacy; its external HTTP protocol contracts belong canonically in `docs/ARCHITECTURE.md`. |
| External Model Context Protocol (MCP) tool schemas (`dispatch_worker`, `dispatch_speculative`, `verify_f2p`, `refinery_enqueue`, `query_brain_memory`) | `docs/API_SPECIFICATION.md` (Section 2) | `docs/ARCHITECTURE.md` (Subsystem 1 & Protocol Contracts) | CONTRACT.md Section 1 (Note): External MCP tool contracts belong canonically in `docs/ARCHITECTURE.md`. |
| External Agent Client Protocol (ACP) JSON-RPC envelopes (`initialize`, `agent/sendMessage`) | `docs/API_SPECIFICATION.md` (Section 3) | `docs/ARCHITECTURE.md` (Subsystem 1 & Protocol Contracts) | CONTRACT.md Section 1 (Note): External ACP protocol contracts belong canonically in `docs/ARCHITECTURE.md`. |
| SQLite WAL Durable Event Sourcing Schema (`workflow_events` DDL) | `docs/API_SPECIFICATION.md` (Section 4) | `docs/GO_SPECIFICATION.md` (`internal/memory/schema.sql` / `internal/state`) | DDL and SQL schema definitions belong canonically in `docs/GO_SPECIFICATION.md` alongside Go state models. |
