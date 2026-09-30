# Questions Resolution Ledger

This document tracks the resolution of open questions raised across all worker ledgers during the documentation overhaul.

---

## 1. Questions Decided by CONTRACT.md

The following questions raised in the worker question ledgers have been resolved definitively by `docs/_review/CONTRACT.md`:

| Question ID | Question | Source Ledger | Contract Decision & Rationale |
| :--- | :--- | :--- | :--- |
| **Q-GOSPEC-01** | Gateway HTTP Server Engine (`net/http` vs `fasthttp`) | `GO_SPECIFICATION.questions.md` | **Resolved**: `CONTRACT.md` §3.2 and §3.4 explicitly mandate **FastHTTP** as the canonical HTTP reverse proxy engine for low-latency, zero-allocation SSE streaming. |
| **Q-ROAD-03** | Sandbox Tier Fallback Precedence (`bwrap` vs `landlock`) | `ROADMAP.questions.md` | **Resolved**: `CONTRACT.md` §3.1 & §3.2 establish **Tier 1 (Bubblewrap `bwrap`)** as the primary sandbox for unprivileged filesystem and network containment, with **Tier 2 (Linux Landlock LSM)** as the containerless filesystem restriction fallback. |
| **Q-ECO-02** | RepoMap Token Budgeting Formula | `ECOSYSTEM_INNOVATIONS.questions.md` | **Resolved**: `CONTRACT.md` §3.3 and `cli-leader.example.yaml` mandate a strict token cap: `max_repo_map_tokens: 1024` and `max_injected_rule_tokens: 800`. |
| **Q-ARCH-04** | Landlock Network Isolation | `ARCHITECTURE.questions.md` | **Resolved**: `CONTRACT.md` §3.1 defines network isolation as a Bubblewrap namespace capability (`--unshare-net`). Landlock is strictly designated for filesystem LSM restriction. |
| **Q-CONTRIB-04**| Human Review vs Byzantine Quorum | `CONTRIBUTING.questions.md` | **Resolved**: Byzantine Quorum consensus is an automated machine-receipt review gate within the Refinery merge queue; final production pull requests to `main` still require human maintainer sign-off per `CONTRIBUTING.md`. |

---

## 2. Remaining Open Questions for User Decision

The following technical decisions require project steering or empirical measurement from the user:

### Release Scheduling & Hardware Baselines
1. **Target Calendar Dates for Phases 0–4** (`ROADMAP.questions.md` Q-ROAD-01):
   - What are the target quarters/milestones for Phase 0 (`v0.1 Core`), Phase 1 (Swarm & Gateway), Phase 2 (Refinery & Speculative Racing), Phase 3 (Cognitive Memory & ACP), and Phase 4 (Production TUI)?
2. **Canonical Benchmark Hardware Specification** (`ROADMAP.questions.md` Q-ROAD-05):
   - What reference host machine (CPU architecture, memory, Linux kernel version) will serve as the canonical testbed to validate `Target (unmeasured)` performance claims (e.g. `<100 µs` gateway TTFT, `60 FPS` TUI batching)?

### Core Runtime & Storage Choices
3. **Claude CLI Host Executable Name** (`README.questions.md` Q-README-02):
   - Should `cli-leader` prioritize `claude` (Anthropic's official npm package) or `claude-code`, or require explicit configuration under `workers.claude.binary`?
4. **Configuration Search Hierarchy** (`README.questions.md` Q-README-01):
   - What is the default search order when `--config` is omitted? Proposed: `./cli-leader.yaml` $\to$ `~/.config/cli-leader/config.yaml` $\to$ `/etc/cli-leader/config.yaml`.
5. **Vector Index Persistence Strategy** (`GO_SPECIFICATION.questions.md` Q-GOSPEC-02):
   - Should `chromem-go` vector embeddings persist directly inside `.brain/memory.db` (as binary BLOBs) or in a dedicated on-disk directory (e.g., `.brain/vectors/`)?
6. **AST Mutation Testing Toolchain** (`GO_SPECIFICATION.questions.md` Q-GOSPEC-05):
   - Should AST defect injection be executed by an internal pure-Go walker (`go/ast`) or wrap an external toolchain runner (e.g., `cargo-mutants` / `go-mutesting`)?
7. **Agent Client Protocol (ACP) Transport Scope** (`GO_SPECIFICATION.questions.md` Q-GOSPEC-06 & `ROADMAP.questions.md` Q-ROAD-02):
   - Does ACP communicate exclusively over `stdio` for local editor sub-processes (Zed, JetBrains), or will UNIX domain sockets / TCP listeners be supported in Phase 3?
