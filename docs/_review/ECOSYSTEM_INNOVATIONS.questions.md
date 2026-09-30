# Open Questions & Verification Items: ECOSYSTEM_INNOVATIONS.md

This document tracks all unverified facts, parameters, library bindings, and external schema specifications identified during the overhaul of `docs/ECOSYSTEM_INNOVATIONS.md` per [`docs/_review/CONTRACT.md`](CONTRACT.md) Section 4 Rule 3 ("Never Invent Facts").

Every item listed here corresponds to a `TODO(verify): <what>` tag inserted into the canonical document text.

---

## Open Questions Table

| Question ID | Section Reference | Item Requiring Verification | Context & Impact |
| :--- | :--- | :--- | :--- |
| **Q-ECO-01** | §1.1 (Aider) | `TODO(verify): Tree-Sitter CGO vs pure-Go bindings for multi-language AST parsing` | README.md specifies a "Single Static Go Binary with zero-CGO dependencies", while ROADMAP.md mentions `smacker/go-tree-sitter` (which traditionally requires CGO for C grammars). Need to confirm whether pure-Go WebAssembly runners (e.g. `wazero`), pure-Go AST parsers, or pre-compiled CGO static archives will be used for multi-language RepoMap extraction. |
| **Q-ECO-02** | §1.1 (Aider) | `TODO(verify): dynamic vs fixed token budgeting formula for RepoMap injection` | `cli-leader.example.yaml` defines `max_repo_map_tokens: 1024`. Need to verify whether the RepoMap engine dynamically scales token allocation when routing to ultra-large context models (e.g. Gemini 1.5/2.0 with 1M+ tokens) or enforces a strict fixed ceiling. |
| **Q-ECO-03** | §1.3 (OpenHands) | `TODO(verify): compatibility layer for OpenHands micro-agent Action/Observation event schemas` | OpenHands standardizes on an event-stream architecture (`Action` and `Observation` events). Need to verify if `cli-leader`'s `internal/worker/adapters` will include a JSON-based event adapter for OpenHands-compatible workers in addition to standard PTY subprocesses. |
| **Q-ECO-04** | §4.1 (ACP) | `TODO(verify): ACP specification schema revision and streaming PTY event support` | Agent Client Protocol (ACP) standardizes agent-editor JSON-RPC interactions. Need to verify the specific draft revision adopted and whether ACP supports raw PTY streaming output or strictly structured file diffs and tool invocation intents. |
| **Q-ECO-05** | §4.2 (Zellij) | `TODO(verify): Zellij IPC socket protocol vs WebAssembly plugin architecture for workspace control` | Zellij supports session manipulation via CLI commands (`zellij action`), an internal Unix IPC socket, or WebAssembly plugins (`zellij-tile`). Need to verify which integration surface `cli-leader`'s worker manager and TUI will target. |
| **Q-ECO-06** | §4.3 (`cargo-mutants`) | `TODO(verify): multi-language AST mutation engine bindings (Go, TypeScript, Python)` | While `cargo-mutants` serves as the architectural reference for Rust, cli-leader targets polyglot codebases. Need to verify whether the adversarial Breaker worker CLI wraps native mutators (`go-mutesting`, `mutmut`, `stryker`) or implements a lightweight Tree-Sitter mutation injector in pure Go. |
