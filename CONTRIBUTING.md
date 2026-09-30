# Contributing to cli-leader

Thank you for contributing to **`cli-leader`**!

---

## 🛠️ Development Setup

### Prerequisites
- **Go 1.22+**
- **Git 2.30+** (with `git worktree` support)
- *Optional*: Bubblewrap (`bwrap`) for enhanced sandbox isolation

### Building from Source
```bash
git clone https://github.com/01RG0/cli-leader.git
cd cli-leader
go build -o cli-leader ./cmd/cli-leader
```

### Running Tests
```bash
go test -v -race ./...
```

---

## 📐 Architecture & Standards

Before opening a pull request, ensure your changes adhere to:
1. **Zero-CGO**: `cli-leader` must compile as a static binary on Linux AMD64/ARM64 and macOS without requiring native C toolchains.
2. **Process Group Safety**: All spawned processes must respect `Pdeathsig = syscall.SIGTERM` and `-pgid` cleanup.
3. **Decoupled TUI**: Never send high-frequency messages directly into Bubble Tea's main loop; use the 30Hz ticker batcher.
4. **Git Serialization**: Worktree lifecycle operations must remain protected by mutex locks.

---

## 📝 Commit Conventions
We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:
- `feat:` New capability or user-facing feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `refactor:` Code restructuring without functional changes
- `test:` Adding or updating test suites
