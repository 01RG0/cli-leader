# Architecture Open Questions & Verification Log

This document tracks all technical ambiguities, protocol edge cases, or missing specifications identified during the architecture overhaul per `docs/_review/CONTRACT.md`.

---

## Open Questions & Verification Items

1. **TODO(verify): Gemini streaming thinking block envelope specification**
   - *Context*: Subsystem 2 (Multi-Provider Gateway) maps DeepSeek-R1 `reasoning_content` to Anthropic `thinking` blocks. For Gemini 2.5 Pro / Flash via the OpenAI-compatible REST endpoint, how are model thought chunks exposed and does the gateway translate them into `thinking_delta` blocks?
   - *Impact*: Gateway streaming state machine (`internal/gateway/stream.go`).

2. **TODO(verify): ACP agent/sendMessage streaming response structure**
   - *Context*: Subsystem 1 (Brain & Protocols) defines the Agent Client Protocol (ACP) JSON-RPC handshake and `agent/sendMessage` request. What is the canonical JSON-RPC notification schema for streaming incremental assistant text deltas and diff inspections back to Zed and JetBrains IDE buffers?
   - *Impact*: ACP server envelope design (`internal/acp/session.go`).

3. **TODO(verify): Interactive PTY stdin streaming support via MCP**
   - *Context*: Subsystem 3 (Worker Process Engine) uses an automated regex auto-responder (`internal/worker/auto_reply.go`) to answer prompts like `[y/N]`. If an unexpected prompt is encountered, can the Brain stream interactive stdin keystrokes over an MCP tool call, or does the task fail fast to the supervisor?
   - *Impact*: MCP tool interface (`internal/mcp/handlers.go`) and PTY lifecycle.

4. **TODO(verify): Landlock network restriction capabilities across kernel versions**
   - *Context*: Subsystem 7 (Workspace Sandboxing) uses Bubblewrap with Linux Landlock LSM as a fallback. Linux Landlock gained network socket scoping (e.g. `LANDLOCK_ACCESS_NET_CONNECT_TCP`) only in kernel 6.7+. For older kernels supporting filesystem sandboxing but lacking network rules, how is `allow_network: false` enforced in Landlock tier?
   - *Impact*: Sandbox initialization (`internal/sandbox/landlock.go`).

5. **TODO(verify): Thompson Sampling prior distribution parameters**
   - *Context*: Subsystem 6 (Cognitive Memory Bank) and Subsystem 5 (Byzantine Quorum) weight votes and route tasks via Bayesian Thompson Sampling $\text{Beta}(\alpha, \beta)$. What are the canonical default initial prior values (e.g., $\alpha=1.0, \beta=1.0$ uninformative uniform, or $\alpha=2.0, \beta=2.0$) configured in `internal/memory/matrix.go`?
   - *Impact*: Initial worker routing probability and consensus threshold weighting.
