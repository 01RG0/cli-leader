# ROADMAP Open Questions & Verification Items

As mandated by `CONTRACT.md` Section 4 Rule 3:
> *"If a technical specification, parameter, or name is not established in this CONTRACT, do NOT invent it. Output `TODO(verify): <what>` and record the question in `docs/_review/<FILE>.questions.md`."*

This document tracks all unverified facts, parameters, and scheduling items identified during the `docs/ROADMAP.md` overhaul.

---

## 1. Questions Ledger

### Q-ROAD-01: Milestone Release Dates & Target Schedule
- **Context**: Neither `CONTRACT.md` nor the legacy `ROADMAP.md` specified calendar dates, sprint cycles, or quarter estimates for Phase 0 through Phase 4.
- **In-Text Marker**: `TODO(verify): Target calendar completion dates for Phase 0 through Phase 4`
- **Question**: What are the target calendar release quarters or dates for Phase 0 (v0.1 Core), Phase 1 (Swarm & Gateway), Phase 2 (Refinery & Speculative Racing), Phase 3 (Cognitive Memory & ACP), and Phase 4 (Production Hardening & TUI)?

---

### Q-ROAD-02: Scope of Agent Client Protocol (ACP) Editor Clients
- **Context**: `docs/ARCHITECTURE.md` and `docs/API_SPECIFICATION.md` reference Zed and JetBrains as ACP clients. External communities also experiment with Neovim and VS Code bridges.
- **In-Text Marker**: `TODO(verify): Exact client editor matrix for ACP integration in Phase 3`
- **Question**: Will Phase 3 deliver native ACP client integrations for both Zed and JetBrains simultaneously, or is Zed the primary initial client with JetBrains scheduled as a secondary milestone?

---

### Q-ROAD-03: Sandbox Tier Resolution Precedence
- **Context**: Configuration in `cli-leader.yaml` allows `sandbox.tier: "auto"`.
- **In-Text Marker**: `TODO(verify): Sandbox tier precedence and kernel version fallback logic`
- **Question**: When `sandbox.tier` is set to `auto`, does the engine attempt Bubblewrap (`bwrap`) first and fall back to Landlock LSM, or does it probe kernel capability `/sys/kernel/security/lsm` before choosing?

---

### Q-ROAD-04: Zero-Shot CLI Adapter Generator Backend
- **Context**: Phase 4 introduces a `[Research]` milestone for generating CLI adapters by interrogating `--help` strings.
- **In-Text Marker**: `TODO(verify): Execution backend for zero-shot adapter synthesis`
- **Question**: Does the zero-shot adapter generator invoke the primary cognitive Brain (Claude CLI) via MCP or dispatch prompt synthesis through the local Gateway (`internal/gateway`) using a designated fast utility model?

---

### Q-ROAD-05: Hardware Baseline for Benchmark Validation
- **Context**: Performance claims such as `<100 µs` gateway TTFT and `30Hz / 60 FPS` TUI log batching are marked `Target (unmeasured)`.
- **In-Text Marker**: `TODO(verify): Empirical benchmark runner hardware specifications`
- **Question**: What reference machine configuration (CPU model, core count, RAM, OS kernel version) will be designated as the canonical testbed to generate empirical benchmark citations and validate these targets?
