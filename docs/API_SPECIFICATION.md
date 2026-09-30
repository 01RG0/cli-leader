# API & Protocol Specification

This document details the exact protocols, schema payloads, and RPC envelopes used across **`cli-leader`**.

---

## 1. Multi-Provider Gateway API (Anthropic Protocol Emulation)

The gateway listens locally on `http://127.0.0.1:8082`.

### 1.1 `POST /v1/messages/count_tokens`
Evaluates input token counts locally using pure-Go BPE tokenization (`tiktoken-go`).

#### Request Payload:
```json
{
  "model": "claude-3-7-sonnet-20250219",
  "messages": [
    {
      "role": "user",
      "content": "Refactor the authentication middleware."
    }
  ],
  "system": "You are the Executive Brain.",
  "tools": [...]
}
```

#### Response Payload (HTTP 200):
```json
{
  "input_tokens": 142
}
```

---

### 1.2 `POST /v1/messages` (Streaming SSE)
Receives standard Anthropic requests and translates them into OpenAI/Gemini/Ollama upstream requests.

#### Outbound SSE Stream Envelopes:
```http
HTTP/1.1 200 OK
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
X-Accel-Buffering: no

event: message_start
data: {"type":"message_start","message":{"id":"msg_01","type":"message","role":"assistant","content":[],"model":"claude-3-7-sonnet-20250219","stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":142,"output_tokens":0}}}

event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"thinking","thinking":""}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"thinking_delta","thinking":"Analyzing target files..."}}

event: content_block_stop
data: {"type":"content_block_stop","index":0}

event: content_block_start
data: {"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_01","name":"dispatch_worker","input":{}}}

event: content_block_delta
data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"cli_name\":\"aider\",\"task_prompt\":\"Refactor auth.go\"}"}}

event: content_block_stop
data: {"type":"content_block_stop","index":1}

event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"tool_use","stop_sequence":null},"usage":{"output_tokens":48}}

event: message_stop
data: {"type":"message_stop"}
```

---

## 2. Model Context Protocol (MCP) Tool Schemas

Exposed over `stdio` and local SSE for Claude CLI.

### 2.1 `dispatch_worker`
Spawns an isolated worker CLI task in a Git worktree.

```json
{
  "name": "dispatch_worker",
  "description": "Dispatches an isolated worker CLI inside a dedicated Git worktree.",
  "parameters": {
    "type": "object",
    "properties": {
      "cli_name": {
        "type": "string",
        "enum": ["aider", "gemini", "ollama", "generic"],
        "description": "Target CLI tool."
      },
      "task_prompt": {
        "type": "string",
        "description": "Specific, actionable instruction for the worker."
      },
      "target_files": {
        "type": "array",
        "items": { "type": "string" },
        "description": "Target file paths to scope the worker's changes."
      },
      "flags": {
        "type": "array",
        "items": { "type": "string" },
        "description": "Optional CLI flags."
      }
    },
    "required": ["cli_name", "task_prompt"]
  }
}
```

### 2.2 `dispatch_speculative`
Spawns multiple candidate workers in parallel ("Branch Racing").

```json
{
  "name": "dispatch_speculative",
  "description": "Runs parallel candidate workers in separate worktrees and selects the best verified patch.",
  "parameters": {
    "type": "object",
    "properties": {
      "candidates": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "cli_name": { "type": "string" },
            "prompt": { "type": "string" }
          },
          "required": ["cli_name", "prompt"]
        }
      },
      "target_files": {
        "type": "array",
        "items": { "type": "string" }
      }
    },
    "required": ["candidates"]
  }
}
```

### 2.3 `verify_f2p`
Runs Fail-to-Pass (F2P) and mutation testing verification gates.

```json
{
  "name": "verify_f2p",
  "description": "Executes reproduction test and mutation testing gate against a completed worktree.",
  "parameters": {
    "type": "object",
    "properties": {
      "task_id": { "type": "string" },
      "reproduction_test_cmd": { "type": "string" },
      "enable_mutation_gate": { "type": "boolean", "default": true }
    },
    "required": ["task_id", "reproduction_test_cmd"]
  }
}
```

### 2.4 `refinery_enqueue`
Enqueues a verified worktree into the Bors-style Refinery merge queue.

```json
{
  "name": "refinery_enqueue",
  "description": "Enqueues a verified worktree branch for serialized testing and squash-merging into main.",
  "parameters": {
    "type": "object",
    "properties": {
      "task_id": { "type": "string" },
      "commit_message": { "type": "string" }
    },
    "required": ["task_id", "commit_message"]
  }
}
```

### 2.5 `query_brain_memory`
Queries episodic history, `.mdc` rules, and the Tree-Sitter Repo Map.

```json
{
  "name": "query_brain_memory",
  "description": "Searches cognitive memory, project rules, and codebase symbol graph.",
  "parameters": {
    "type": "object",
    "properties": {
      "query": { "type": "string" },
      "category": {
        "type": "string",
        "enum": ["all", "rules", "history", "symbols", "skills"]
      }
    },
    "required": ["query"]
  }
}
```

---

## 3. Agent Client Protocol (ACP) Specification

Standardized JSON-RPC protocol over `stdio` between editors (Zed / JetBrains) and `cli-leader`.

### 3.1 Initial Handshake (`initialize`)
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "1.0",
    "clientInfo": {
      "name": "zed",
      "version": "0.180.0"
    },
    "capabilities": {
      "streaming": true,
      "diffInspection": true
    }
  }
}
```

### 3.2 Thread Message Dispatch (`agent/sendMessage`)
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "agent/sendMessage",
  "params": {
    "threadId": "th_987",
    "message": {
      "role": "user",
      "content": "Implement user password hashing using bcrypt."
    }
  }
}
```

---

## 4. SQLite WAL Durable Event Sourcing Schema

Schema used by the embedded Temporal-grade durable engine (`internal/state/durable_engine.go`):

```sql
CREATE TABLE IF NOT EXISTS workflow_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workflow_id TEXT NOT NULL,
    sequence_num INTEGER NOT NULL,
    event_type TEXT NOT NULL, -- WorkflowStarted, ActivityScheduled, ActivityCompleted, ActivityFailed
    activity_name TEXT,
    payload_json TEXT NOT NULL,
    result_json TEXT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(workflow_id, sequence_num)
);

CREATE INDEX IF NOT EXISTS idx_workflow_events ON workflow_events(workflow_id, sequence_num);
```
