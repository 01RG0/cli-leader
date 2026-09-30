# cli-leader 🧠⚡

> **A high-performance autonomous meta-orchestrator built in Go.**  
> Uses **Claude CLI** as the central cognitive brain to direct, supervise, and learn from a swarm of specialized worker AI CLIs across multiple LLM providers.

---

## 🌟 Vision

Modern AI coding CLI tools (Claude Code, Aider, Gemini CLI, Ollama, DeepSeek, etc.) each have distinct superpowers: some excel at ultra-long context reasoning, some at fast local AST edits, others at low-cost bulk operations.

**`cli-leader`** unites them under a single unified intelligence:
1. **The Brain (Claude CLI)**: Acts as the executive leader, breaking down tasks, assigning work to specialized worker CLIs, and arbitrating technical decisions.
2. **Universal Multi-Provider Gateway**: Native Go reverse proxy translating Anthropic Messages API into OpenAI, Google Gemini, DeepSeek, Ollama, and OpenRouter protocols—allowing the Brain to run on any model with full SSE streaming and tool calling.
3. **Ever-Learning Cognitive Memory**: Central shared memory store (Reflexion verbal learning, scoped `.cursor/rules/*.mdc` conventions, and a Bayesian Thompson Sampling capability matrix) that makes every worker CLI smarter over time.
4. **Speculative Execution & Refinery Queue**: Dispatches parallel candidate workers in isolated Git worktrees ("Branch Racing") and verifies merges through a Bors-style Refinery queue.
5. **Adversarial Verification & F2P Testing**: Enforces Fail-to-Pass (F2P) test verification and pairs builders with adversarial breaker CLIs for rigorous audit before merging.
6. **Single Static Go Binary**: Blazing-fast (<10ms startup), rock-solid concurrency via goroutines, zero-CGO dependencies, and a smooth 60 FPS terminal UI powered by Charm's `bubbletea`.

---

## 🏗️ High-Level Architecture

```mermaid
flowchart TD
    User([Developer / Terminal]) --> TUI["cli-leader TUI (Bubbletea & Lipgloss - 60 FPS)"]
    TUI --> Brain["Claude CLI (The Executive Brain)"]

    subgraph Gateway ["Go Multi-Provider Gateway (:8082)"]
        Brain -.->|Anthropic /v1/messages| GatewayServer["Gateway Router & SSE Translator"]
        GatewayServer --> Anthropic["Anthropic (Claude 3.7 / 3.5)"]
        GatewayServer --> OpenAI["OpenAI (GPT-4o / o3-mini)"]
        GatewayServer --> Gemini["Google Gemini (2.5 / 3.0)"]
        GatewayServer --> Ollama["Local Ollama / DeepSeek-R1"]
    end

    subgraph ControlPlane ["Worker Control Plane (MCP & Process Supervisors)"]
        Brain -->|MCP Tool Calls| MCPServer["cli-leader MCP Server"]
        MCPServer --> DualLedger[("Dual Ledger: Task & Progress\n(Loop & Stall Breakers)")]
        MCPServer --> Dispatcher["Worker Dispatcher & PTY Engine"]
        Dispatcher --> W1["Aider (AST Code Refactoring)"]
        Dispatcher --> W2["Gemini CLI (1M+ Token Indexing)"]
        Dispatcher --> W3["Ollama / DeepSeek (Fast Local Tests)"]
        Dispatcher --> W4["Adversarial Breaker (Vulnerability Audit)"]
    end

    subgraph VerificationEngine ["Speculative Execution & Merge Refinery"]
        W1 & W2 & W3 --> Worktrees["Isolated Git Worktrees (Branch Racing)"]
        Worktrees --> F2P["Fail-to-Pass (F2P) TDD Engine"]
        F2P --> Refinery["Refinery Bors-Style Merge Queue"]
        Refinery --> MainGit["Git Main Branch"]
    end

    subgraph CognitiveBank ["Ever-Learning Cognitive Memory"]
        MCPServer <--> EpisodicDB[("SQLite Episodic Ledger\n(modernc.org/sqlite)")]
        MCPServer <--> VectorDB[("Pure-Go Vector Memory\n(chromem-go)")]
        Dispatcher --> Reflexion["Reflexion Engine (Verbal Learning)"]
        Reflexion --> Rules[(".brain/conventions/*.mdc Rules")]
        Reflexion --> Matrix[("Thompson Sampling Capability Matrix")]
    end
```

---

## 📦 Core Pillars

| Pillar | Description |
|---|---|
| 🧠 **Claude Brain & MCP** | Exposes structured MCP tools for Claude CLI to dispatch, monitor, and abort workers, inspect diffs, and manage memory. |
| 🔄 **Multi-Provider Proxy** | Low-latency Anthropic API emulator in Go with SSE chunk framing, reasoning block translation, and `/v1/messages/count_tokens`. |
| 🏎️ **Speculative Branch Racing** | Launches parallel candidate workers in separate Git worktrees and automatically merges the best patch. |
| 🧪 **Fail-to-Pass (F2P) Testing** | Validates that reproduction tests fail on unmodified code before proving the fix passes and prevents regressions. |
| 🚦 **Refinery Merge Queue** | Serializes worktree merges, tests against current HEAD, and squashes cleanly with zero merge collisions. |
| 📚 **Ever-Learning Cognitive Bank** | Extracts verbal reflections from failures into `.cursor/rules/*.mdc` files, indexed via pure-Go vector search (`chromem-go`). |
| 📊 **Thompson Sampling Routing** | Bayesian multi-armed bandit profiling that dynamically routes tasks to the best worker based on cost, speed, and accuracy. |
| ⚡ **Go Native Systems Engine** | `creack/pty` with `Pdeathsig`, 30Hz TUI batching for smooth rendering, and tiered sandboxing (Bubblewrap + Landlock). |

---

## 📖 Documentation

- [System Architecture Deep Dive](docs/ARCHITECTURE.md)
- [Go Technical Specification & Package Layout](docs/GO_SPECIFICATION.md)
- [Development Roadmap & Milestones](docs/ROADMAP.md)

---

## 🚀 Technology Stack

- **Core Language**: Go 1.22+
- **Terminal UI**: [Charm CLI](https://charm.sh/) (`bubbletea`, `lipgloss`, `bubbles`)
- **Protocol**: Model Context Protocol (MCP) Go Server
- **Subprocess & PTY**: `os/exec`, `creack/pty` with Linux `Pdeathsig` process group management
- **Database & Vectors**: Pure Go embedded SQLite (`modernc.org/sqlite`) & in-memory vector store (`github.com/philippgille/chromem-go`)
- **Token Counting**: Pure Go BPE tokenization (`github.com/pkoukk/tiktoken-go`)
- **JSON Optimization**: High-performance mutation without reflection (`github.com/tidwall/gjson`, `github.com/tidwall/sjson`)
- **Sandboxing**: Bubblewrap (`bwrap`) & Linux Landlock LSM (`github.com/landlock-lsm/go-landlock`)
- **VCS Automation**: Native Git Worktree wrapping with mutex lock serialization

---

## 📄 License
MIT License.
