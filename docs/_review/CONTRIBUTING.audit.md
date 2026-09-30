# CONTRIBUTING.md Audit Findings against CONTRACT.md

This document audits the legacy `CONTRIBUTING.md` against the canonical requirements in `docs/_review/CONTRACT.md`.

---

## 1. Audit Summary

| Finding ID | Category | Severity | Contract Section | Short Description |
| :--- | :--- | :--- | :--- | :--- |
| **CONTRIB-01** | Status Tags | High | §2 | Missing canonical feature status tags (`[v0.1 Core]`, `[Planned]`, `[Research]`) across all headings and missing definitions for contributors. |
| **CONTRIB-02** | Document Ownership | High | §1, §4.2 | Misplaced architectural and runtime invariants (Zero-CGO, `Pdeathsig`, TUI batcher, worktree mutexes) belonging in `docs/ARCHITECTURE.md` and `docs/GO_SPECIFICATION.md`. |
| **CONTRIB-03** | Contribution Workflow | High | §1 | Missing contribution workflow, branching guidelines (`feat/*`, `fix/*`, `docs/*`), and pull request submission steps. |
| **CONTRIB-04** | Testing Standards | Critical | §1, §3.1, §3.4 | Vague testing guidelines omitting mandatory PR verification gates (Fail-to-Pass / F2P, AST mutation testing criteria, Bors Refinery merge queue). |
| **CONTRIB-05** | Linting & Static Analysis | Medium | §1, Execution Step 5 | Completely absent code quality and static analysis tooling requirements (`golangci-lint`, formatting). |
| **CONTRIB-06** | Commit Conventions | Medium | §3.2, §3.4 | Commit conventions lack canonical package and component scopes as defined in CONTRACT.md §3.2. |
| **CONTRIB-07** | Specification Cross-References | Medium | §1 | Total absence of relative links to canonical architectural and Go specification documents. |
| **CONTRIB-08** | Sandbox Toolchain | Low | §3.1, §3.3 | Incomplete sandboxing specifications failing to describe Tier 1 (Bubblewrap) vs Tier 2 (Landlock LSM) isolation. |

---

## 2. Detailed Findings

### CONTRIB-01: Missing Canonical Feature Status Tags
- **Contract Reference**: `docs/_review/CONTRACT.md` §2 ("Every feature or capability mentioned anywhere in documentation must have **exactly one** status tag attached to its heading or specification block: `[v0.1 Core]`, `[Planned]`, or `[Research]`.").
- **Exact Quote**:
  ```markdown
  ## 🛠️ Development Setup
  ### Prerequisites
  ### Building from Source
  ### Running Tests
  ## 📐 Architecture & Standards
  ## 📝 Commit Conventions
  ```
- **Analysis**: No section in the original file specifies a status tag. Furthermore, `CONTRIBUTING.md` fails to define the three canonical status tags (`[v0.1 Core]`, `[Planned]`, `[Research]`) that contributors must apply when contributing documentation or proposing capabilities.
- **Remediation**: Add explicit definitions of the three canonical status tags for contributors, and tag all sections within `CONTRIBUTING.md` with `[v0.1 Core]`.

---

### CONTRIB-02: Misplaced Architectural & Low-Level Invariants (Violates Document Ownership Map)
- **Contract Reference**: `docs/_review/CONTRACT.md` §1 ("Each documentation topic has exactly **ONE** canonical owner file. Information must reside in its owner file; cross-references should be relative links rather than duplicated text. System Design & Architecture -> docs/ARCHITECTURE.md; Technical Specification & Go Packages -> docs/GO_SPECIFICATION.md.").
- **Exact Quote**:
  ```markdown
  ## 📐 Architecture & Standards

  Before opening a pull request, ensure your changes adhere to:
  1. **Zero-CGO**: `cli-leader` must compile as a static binary on Linux AMD64/ARM64 and macOS without requiring native C toolchains.
  2. **Process Group Safety**: All spawned processes must respect `Pdeathsig = syscall.SIGTERM` and `-pgid` cleanup.
  3. **Decoupled TUI**: Never send high-frequency messages directly into Bubble Tea's main loop; use the 30Hz ticker batcher.
  4. **Git Serialization**: Worktree lifecycle operations must remain protected by mutex locks.
  ```
- **Analysis**: These four items are deep architectural invariants and implementation-level details. They belong canonically in `docs/ARCHITECTURE.md` and `docs/GO_SPECIFICATION.md`. Duplicating them ad-hoc in `CONTRIBUTING.md` violates document ownership.
- **Remediation**: Log these architectural concepts in `docs/_review/CONTRIBUTING.ideas.md` to ensure they are redirected to their canonical owner files, and replace this inline section with clear relative links to `docs/ARCHITECTURE.md` and `docs/GO_SPECIFICATION.md`.

---

### CONTRIB-03: Missing Contribution Workflow and Git Branch Naming Standards
- **Contract Reference**: `docs/_review/CONTRACT.md` §1 ("Contribution & Developer Guidelines -> CONTRIBUTING.md: Contribution workflow, git branch naming, pull request rules, development environment setup, testing requirements, and code review criteria.").
- **Exact Quote**: Entire original file; no workflow or branch naming guidance exists.
- **Analysis**: New contributors have no instructions on how to branch from `main`, standard branch naming prefixes (`feat/*`, `fix/*`, `docs/*`, `refactor/*`, `test/*`, `chore/*`), how to stage changes, or how to submit pull requests.
- **Remediation**: Provide a comprehensive step-by-step contribution lifecycle detailing branch creation from `main`, naming conventions, local verification, and pull request procedures.

---

### CONTRIB-04: Vague Testing Standards and Omission of PR Verification Gates
- **Contract Reference**: `docs/_review/CONTRACT.md` §1, §3.1 ("Refinery: A serialized Bors-style merge queue that verifies candidate worktree branches with automated test suites and mutation testing gates before squash-merging them into the active branch.", "F2P (Fail-to-Pass): A Test-Driven Development verification protocol..."), and §3.4.
- **Exact Quote**:
  ```markdown
  ### Running Tests
  ```bash
  go test -v -race ./...
  ```
  ```
- **Analysis**: The testing section only lists a single unit test command. It fails to document cli-leader's automated PR verification gates:
  1. Fail-to-Pass (F2P) reproduction test protocol (confirming failure on unmodified `HEAD` and pass on fix).
  2. Race-free unit test completion.
  3. AST mutation testing verification gate (`cargo-mutants` style killing synthetic mutants).
  4. Serialized Bors-style Refinery merge queue requirements.
- **Remediation**: Expand testing into a dedicated PR Verification Gates section detailing F2P, unit tests, AST mutation testing pass criteria, and the Refinery merge queue.

---

### CONTRIB-05: Missing Code Quality & Static Analysis Standards
- **Contract Reference**: Overhaul Execution Step 5 ("linting with golangci-lint"), `docs/_review/CONTRACT.md` §1.
- **Exact Quote**: Entire original file; no mention of linting or code formatters.
- **Analysis**: Contributors are not informed of required linters (`golangci-lint run ./...`), formatting standards (`go fmt`, `goimports`), or static checks required before opening a PR.
- **Remediation**: Document the Go toolchain linting requirements using `golangci-lint`, pre-commit formatting, and add `TODO(verify): Exact golangci-lint version and .golangci.yml linter configuration`.

---

### CONTRIB-06: Incomplete Conventional Commit Scopes & Grounding
- **Contract Reference**: `docs/_review/CONTRACT.md` §3.2 (Canonical Package Names), §3.4.
- **Exact Quote**:
  ```markdown
  ## 📝 Commit Conventions
  We follow the Conventional Commits specification:
  - feat: New capability or user-facing feature
  - fix: Bug fix
  - docs: Documentation changes
  - refactor: Code restructuring without functional changes
  - test: Adding or updating test suites
  ```
- **Analysis**: While Conventional Commits are mentioned, no scopes are provided. Scopes must align with the canonical Go package names (`gateway`, `supervisor`, `worker`, `speculative`, `refinery`, `verification`, `memory`, `state`, `mcp`, `acp`, `sandbox`, `tui`, `notify`, `config`) and components (`brain`, `refinery`, etc.) specified in `CONTRACT.md` §3.2. Breaking change syntax (`!`) is also omitted.
- **Remediation**: Ground commit scopes directly in canonical package names from `CONTRACT.md` §3.2 and illustrate valid scoped commit messages and breaking change notation.

---

### CONTRIB-07: Missing Canonical Cross-References
- **Contract Reference**: `docs/_review/CONTRACT.md` §1 ("Information must reside in its owner file; cross-references should be relative links rather than duplicated text.").
- **Exact Quote**: Zero relative links exist in the original file.
- **Analysis**: Contributors are given no links to `docs/ARCHITECTURE.md`, `docs/GO_SPECIFICATION.md`, or `docs/ROADMAP.md` to learn about package conventions or architectural patterns.
- **Remediation**: Embed clear relative markdown links to `docs/ARCHITECTURE.md`, `docs/GO_SPECIFICATION.md`, and `docs/ROADMAP.md`.

---

### CONTRIB-08: Incomplete Sandboxing Toolchain Specification
- **Contract Reference**: `docs/_review/CONTRACT.md` §3.1 ("Landlock / Bubblewrap: Linux kernel security mechanisms..."), §3.3.
- **Exact Quote**:
  ```markdown
  - *Optional*: Bubblewrap (`bwrap`) for enhanced sandbox isolation
  ```
- **Analysis**: Describes Bubblewrap as an optional standalone tool without explaining the tiered sandbox architecture (Tier 1 Bubblewrap, Tier 2 Linux Landlock LSM fallback) used by cli-leader.
- **Remediation**: Clarify development prerequisites to explain Tier 1 (`bwrap`) and Tier 2 (`Landlock LSM`) sandboxing environments.
