# CONTRIBUTING.md Open Verification Questions

This ledger records all open questions and unverified details encountered during the overhaul of `CONTRIBUTING.md`. In accordance with Global Rule 3 of `docs/_review/CONTRACT.md`, whenever a technical specification, tool version, or parameter is not explicitly defined in `CONTRACT.md`, a `TODO(verify): <what>` marker is inserted into the text and logged below.

---

## 1. Questions Ledger

### Question 1: Linter Configuration & Minimum golangci-lint Version
- **Tag in Text**: `TODO(verify): Exact golangci-lint version and .golangci.yml linter configuration`
- **Question**: What is the pinned or minimum supported version of `golangci-lint` (e.g. v1.57+), and does the repository maintain an authoritative `.golangci.yml` configuration detailing which linters (`govet`, `errcheck`, `staticcheck`, `gosec`, `revive`, etc.) are enforced?
- **Context**: Execution Step 5 mandates documenting linting with `golangci-lint`, but no `.golangci.yml` currently exists in the repository root, nor does `CONTRACT.md` specify the linter ruleset.

### Question 2: Mutation Testing CLI Tooling & Kill Threshold Criteria
- **Tag in Text**: `TODO(verify): Mutation testing CLI command and exact kill threshold percentage for Go code`
- **Question**: What exact CLI command or test harness runs the AST mutation testing gate (e.g. `go-mutesting ./...` or an internal tool invoking `internal/verification/mutation.go`), and what is the numeric pass criteria (e.g. 100% of injected mutants killed in modified ASTs, or a specific mutation score threshold)?
- **Context**: `CONTRACT.md` §3.1 and `docs/ARCHITECTURE.md` Subsystem 5 specify a `cargo-mutants` style AST mutation testing gate in the Refinery merge queue, but do not specify the exact CLI tool invocation or pass threshold percentage for local contributor testing.

### Question 3: Fail-to-Pass (F2P) Local Verification Command
- **Tag in Text**: `TODO(verify): Standard command-line invocation for local Fail-to-Pass (F2P) verification`
- **Question**: Is there a dedicated `cli-leader` command (e.g. `cli-leader verify f2p`) or a standardized Go test flag/tag convention (e.g. `go test -v -tags=f2p ./...`) for contributors to run F2P reproduction verification locally before pushing to the Refinery?
- **Context**: `CONTRACT.md` §3.1 details the two-step F2P verification protocol (step 1 fail, step 2 pass), but does not state the command line invocation for local contributor execution.

### Question 4: Human Reviewer Sign-off alongside Refinery Quorum
- **Tag in Text**: `TODO(verify): Human reviewer approval requirements alongside Refinery Byzantine quorum consensus`
- **Question**: Does merge approval require human code owner approvals (e.g., 1 or 2 maintainer approvals via GitHub pull request reviews) in addition to passing the automated Bors-style Refinery merge queue and Byzantine quorum machine receipts?
- **Context**: `CONTRACT.md` §3.1 and `docs/ARCHITECTURE.md` Subsystem 5 describe automated Byzantine Quorum machine receipts and Bors-style Refinery merge queue, but do not delineate human maintainer code review criteria.
