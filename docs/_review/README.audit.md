# README.md Audit Report

This audit evaluates the original `README.md` against the rules, boundaries, and standards set forth in [`docs/_review/CONTRACT.md`](file:///home/rootuser/wt-readme/docs/_review/CONTRACT.md).

## Summary Metrics
- **Drex Marketing Language Score**: 0.937 (Target: <0.05) — Extremely high presence of hype words ("superpowers", "blazing-fast", "rock-solid", "ever-learning").
- **Drex Unmeasured Claims Score**: 0.891 (Target: 0.000) — Pervasive unsubstantiated latency, frame rate, and throughput claims ("sub-millisecond TTFT", "60 FPS", "<10ms startup", "zero duplicate token costs", "zero merge collisions", "Sub-5ms semantic centroid pre-routing").
- **Status Tag Compliance**: 0.0% — Zero features or pillars possessed the mandatory `[v0.1 Core]`, `[Planned]`, or `[Research]` tags required by CONTRACT.md § 2.
- **Contract Ownership Violations**: Missing Quickstart guide, missing configuration example matching `cli-leader.example.yaml`, missing link to `CONTRIBUTING.md`, and leakage of low-level implementation details belonging in `docs/GO_SPECIFICATION.md` or `docs/ARCHITECTURE.md`.

---

## Audit Findings Catalog

| Finding ID | Section / Location | Exact Quote | Contract Violation | Required Resolution |
| :--- | :--- | :--- | :--- | :--- |
| **README-01** | Header (Line 1) | `# cli-leader 🧠⚡` | Marketing hype visual styling (decorative emojis) contrary to crisp technical documentation. | Remove emojis; use clean `# cli-leader` header. |
| **README-02** | Tagline (Line 3) | `> **A high-performance autonomous meta-orchestrator built in Go.**` | Marketing buzzword / unmeasured claim ("high-performance") violating CONTRACT.md § 4.4. | Replace with objective architectural description: "An autonomous meta-orchestrator built in Go." |
| **README-03** | Vision (Line 10) | `"Modern AI coding CLI tools ... each have distinct superpowers: some excel at ultra-long context reasoning, some at fast local AST edits, others at low-cost bulk operations."` | Marketing buzzword ("superpowers") and unquantified marketing statements violating CONTRACT.md § 4. | Rewrite to state technical capabilities and architectural niches objectively. |
| **README-04** | Vision item 2 (Line 14) | `"Universal Multi-Provider Gateway: Native Go reverse proxy ... with sub-millisecond TTFT, SSE streaming state machines, and reasoning block translation."` | Unmeasured latency claim ("sub-millisecond TTFT") violating CONTRACT.md § 4.4; missing feature status tag violating CONTRACT.md § 2. | Tag with `[v0.1 Core]` and reframe latency claim to mechanism description or `Target (unmeasured): sub-millisecond TTFT`. |
| **README-05** | Vision item 3 (Line 15) | `"Ever-Learning Cognitive Memory: Central shared memory store ... that makes every worker CLI smarter over time."` | Anthropomorphic marketing claim ("makes every worker CLI smarter over time"); inconsistent path (`.cursor/rules/*.mdc` instead of `.brain/conventions/`); missing status tag. | Tag with `[Planned]`; use canonical path `.brain/conventions/`; describe concrete verbal reinforcement and rule retrieval mechanics. |
| **README-06** | Vision item 6 (Line 18) | `"Erlang OTP & Temporal Durability: OTP supervision trees ... for crash-proof process bundles, SQLite WAL durable event sourcing for instant zero-token resume, and Sagas..."` | Unmeasured / absolute claims ("crash-proof process bundles", "instant zero-token resume"); missing status tag violating CONTRACT.md § 2. | Tag with `[Planned]`; reframe to supervisor restart policies and SQLite WAL replay mechanics. |
| **README-07** | Vision item 8 (Line 20) | `"Single Static Go Binary: Blazing-fast (<10ms startup), rock-solid concurrency via goroutines, zero-CGO dependencies, and a smooth 60 FPS terminal UI powered by Charm's bubbletea."` | Buzzwords ("Blazing-fast", "rock-solid") and unmeasured performance claims ("<10ms startup", "smooth 60 FPS"); missing status tag. | Convert to `Target (unmeasured): <10ms startup` and `Target (unmeasured): 60 FPS`; describe static compilation and goroutine concurrency objectively. |
| **README-08** | Vision Items 1, 4, 5, 7 (Lines 13-19) | All Vision numbered items (The Brain, Speculative Execution, Mutation Testing, Standards-Compliant) | Missing mandatory feature status tags (`[v0.1 Core]`, `[Planned]`, `[Research]`) violating CONTRACT.md § 2. | Assign appropriate status tag to every vision point. |
| **README-09** | Architecture Diagram (Lines 26-68) | Mermaid diagram node labels: `"cli-leader TUI (Bubbletea & Lipgloss - 60 FPS)"`, `"CognitiveBank"`, `"ControlPlane"`, `"VerificationEngine"` | Unmeasured claim ("60 FPS"); non-canonical subsystem naming violating CONTRACT.md § 3.4 (Brain, Gateway, Supervisor, Worker Manager, Refinery, Durable Engine, Memory Store). | Revise diagram node labels to use canonical component names and replace "60 FPS" with `Target (unmeasured): 60 FPS`. Ensure clean Mermaid syntax. |
| **README-10** | Core Pillars Table (Lines 72-87) | Table headers & contents: `🧠 **Claude Brain, MCP & ACP**`, `🔄 **Multi-Provider Proxy**`, `🏎️ **Speculative Branch Racing**`, etc. | Contains decorative emojis; completely lacks mandatory status tags (`[v0.1 Core]`, `[Planned]`, `[Research]`) violating CONTRACT.md § 2. | Remove emojis; prepend explicit status tag to every pillar title. |
| **README-11** | Core Pillars Claims (Lines 77, 81, 82, 84, 85, 86) | `"zero-allocation SSE streaming"`, `"zero duplicate token costs"`, `"zero merge collisions"`, `"ultra-compact codebase map"`, `"Sub-5ms semantic centroid pre-routing"`, `"smooth rendering"` | Unmeasured performance and absolute claims violating CONTRACT.md § 4.4 without empirical citations. | Rewrite claims to describe concrete mechanisms or mark as `Target (unmeasured): <claim>`. |
| **README-12** | Core Pillars Rule Path (Line 83) | `".cursor/rules/*.mdc files"` | Path inconsistency against CONTRACT.md § 3.1 & § 3.3 (`.brain/conventions/`). | Update path to canonical `.brain/conventions/` directory. |
| **README-13** | Scope Omission: Quickstart | Entire document | Missing Quickstart installation and run guide, violating CONTRACT.md § 1 (Document Ownership Map: `README.md` is canonical owner for Quickstart guide). | Add comprehensive Quickstart section (Prerequisites, Build/Install, Run). |
| **README-14** | Scope Omission: Configuration | Entire document | Missing canonical configuration snippet, violating CONTRACT.md § 1 and § 3.3. | Add YAML configuration snippet matching `cli-leader.example.yaml`. |
| **README-15** | Documentation Links (Lines 90-96) | Missing `CONTRIBUTING.md` link | Omission of `CONTRIBUTING.md` link, violating CONTRACT.md § 1 and execution step requirements. | Add direct relative markdown link to `CONTRIBUTING.md` alongside existing documentation links. |
| **README-16** | Technology Stack Scope Leak (Lines 98-111) | Detailed dependency listings and marketing phrasing (`"High-performance mutation without reflection"`) | Technology Stack section leaks low-level dependency details that canonically belong in `docs/GO_SPECIFICATION.md` § 4. | Streamline Technology Summary Table per CONTRACT.md § 1; log moved details to `docs/_review/README.ideas.md`. |
| **README-17** | Unmapped Architecture Claims | Entire document | Claims in `README.md` lack direct cross-reference mappings to sections in `docs/ARCHITECTURE.md` or `docs/GO_SPECIFICATION.md`. | Add explicit mapping references from README capabilities to corresponding architectural subsystem and specification sections. |

---

## Action Plan for Rewrite
1. **Sanitize Tone & Language**: Replace hype ("superpowers", "blazing-fast", "rock-solid") with objective technical descriptions.
2. **Apply Status Tags**: Tag every capability and pillar with `[v0.1 Core]`, `[Planned]`, or `[Research]`.
3. **Qualify Unmeasured Metrics**: Mark all unverified speed, latency, and throughput claims with `Target (unmeasured): ...`.
4. **Include Quickstart**: Document prerequisites, compilation from source, initial configuration, and launching the daemon.
5. **Include Configuration Example**: Embed a clean, commented YAML configuration excerpt adhering to `cli-leader.example.yaml`.
6. **Harmonize Architecture Diagram**: Rebuild the Mermaid architecture overview with canonical component names and valid syntax.
7. **Ensure Complete Links**: Provide working relative links to `docs/ARCHITECTURE.md`, `docs/GO_SPECIFICATION.md`, `docs/ROADMAP.md`, `docs/ECOSYSTEM_INNOVATIONS.md`, and `CONTRIBUTING.md`.
