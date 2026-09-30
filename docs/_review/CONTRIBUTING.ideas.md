# CONTRIBUTING.md Migrated & Redirected Ideas

This ledger tracks all ideas, architectural invariants, and specifications moved or redirected from `CONTRIBUTING.md` to maintain strict compliance with `docs/_review/CONTRACT.md`. In accordance with Global Rule 2, **no ideas are deleted**.

---

| Idea | Old Location | New Location |
| :--- | :--- | :--- |
| **Zero-CGO static compilation invariant** (`cli-leader must compile as a static binary on Linux AMD64/ARM64 and macOS without requiring native C toolchains`) | `CONTRIBUTING.md` § Architecture & Standards (L31) | `docs/GO_SPECIFICATION.md` § Standard Project Layout & `docs/ARCHITECTURE.md` § Subsystem 6 |
| **Process Group Safety and Pdeathsig** (`All spawned processes must respect Pdeathsig = syscall.SIGTERM and -pgid cleanup`) | `CONTRIBUTING.md` § Architecture & Standards (L32) | `docs/ARCHITECTURE.md` § Subsystem 3 & `docs/GO_SPECIFICATION.md` § `internal/worker/pty.go` |
| **Decoupled TUI 30Hz ticker batcher** (`Never send high-frequency messages directly into Bubble Tea's main loop; use the 30Hz ticker batcher`) | `CONTRIBUTING.md` § Architecture & Standards (L33) | `docs/ARCHITECTURE.md` § Subsystem 3 & `docs/GO_SPECIFICATION.md` § `internal/tui/batcher.go` |
| **Git Worktree mutex lock serialization** (`Worktree lifecycle operations must remain protected by mutex locks`) | `CONTRIBUTING.md` § Architecture & Standards (L34) | `docs/ARCHITECTURE.md` § Subsystem 4 & `docs/GO_SPECIFICATION.md` § `internal/worker/worktree.go` |
| **Tiered Workspace Sandbox Isolation** (`*Optional*: Bubblewrap (bwrap) for enhanced sandbox isolation`) | `CONTRIBUTING.md` § Prerequisites (L12) | `docs/ARCHITECTURE.md` § Subsystem 7 (Tier 1: Bubblewrap, Tier 2: Linux Landlock LSM fallback) |
