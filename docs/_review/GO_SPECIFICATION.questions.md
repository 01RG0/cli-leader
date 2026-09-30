# Questions: GO_SPECIFICATION.md

This document records all open technical questions identified during the specification overhaul where facts or parameters are not explicitly established in `docs/_review/CONTRACT.md`. Each item corresponds to a `TODO(verify): <what>` tag in `docs/GO_SPECIFICATION.md`.

---

### Q-GOSPEC-01: Gateway HTTP Engine Implementation (`internal/gateway`)
- **Question**: CONTRACT.md §3.2 specifies `internal/gateway` as a "FastHTTP reverse proxy translating /v1/messages to upstream providers", while `docs/ROADMAP.md` Phase 1 line 7 states "HTTP listener on :8082 using standard library http.NewServeMux." Which HTTP server engine is canonical?
- **Tagged in Spec**: `TODO(verify): gateway_http_engine`
- **Impact**: Determines whether `internal/gateway.ProxyServer` uses `fasthttp.Server` or `http.Server`. FastHTTP provides zero-allocation advantages but lacks full HTTP/2 and standard `http.Handler` interoperability.

---

### Q-GOSPEC-02: Vector Database Persistence Location (`internal/memory`)
- **Question**: Does `chromem-go` persist its collection directly to `.brain/memory.db` (as serialized BLOBs) or to an independent directory on disk (e.g. `.brain/vectors/`)?
- **Tagged in Spec**: `TODO(verify): vector_db_persistence`
- **Impact**: Affects the backup, initialization, and cleanup semantics of the `internal/memory.VectorIndex` struct.

---

### Q-GOSPEC-03: OTP Supervisor Tree Abstraction (`internal/supervisor`)
- **Question**: CONTRACT.md §3.1 specifies an Erlang OTP supervision tree with `one_for_one` and `one_for_all` strategies, and the dependencies table lists `thejerf/suture/v4`. Does `internal/supervisor.Supervisor` wrap `suture.Supervisor` under the hood or provide an internal native OTP state machine?
- **Tagged in Spec**: `TODO(verify): supervisor_backend_impl`
- **Impact**: Dictates whether child process tokens use `suture.ServiceToken` or a custom identifier type.

---

### Q-GOSPEC-04: Worker Subprocess Subreaper Cleanup (`internal/worker`)
- **Question**: In addition to `SysProcAttr.Setsid = true` and `Pdeathsig = syscall.SIGTERM`, does `cli-leader` register the host daemon as a Linux subreaper via `prctl(PR_SET_CHILD_SUBREAPER, 1)` to prevent orphaned grandchild processes spawned by worker CLI tools?
- **Tagged in Spec**: `TODO(verify): pty_subreaper_support`
- **Impact**: Affects OS-level process cleanup semantics in `internal/worker/pty.go`.

---

### Q-GOSPEC-05: AST Mutation Engine Mechanism (`internal/verification`)
- **Question**: For AST mutation testing in `internal/verification/mutation.go`, does `cli-leader` execute an internal Go AST mutator (using `go/parser`, `go/ast`, and `go/printer`) or invoke an external toolchain runner (e.g., `cargo-mutants` / `go-mutesting`)?
- **Tagged in Spec**: `TODO(verify): mutation_testing_toolchain`
- **Impact**: Determines whether `MutationTester` requires external host binaries or functions as a pure-Go AST mutation pipeline.

---

### Q-GOSPEC-06: ACP Server Transports (`internal/acp`)
- **Question**: Does the Agent Client Protocol (ACP) server operate exclusively over standard I/O (`stdio`) when spawned by Zed/JetBrains, or does it support local UNIX domain sockets (`/tmp/cli-leader-acp.sock`) and TCP for remote instances?
- **Tagged in Spec**: `TODO(verify): acp_transport`
- **Impact**: Defines the transport abstraction in `internal/acp.Server`.
