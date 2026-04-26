# AGENTS.md

This file defines behavioral rules for AI agents (Claude, Copilot, etc.) working on this codebase.
Read it fully before making any changes.

---

## Project Overview

This is a **turn-based tactical TUI game** written in Go.
The game runs entirely in the terminal. There is no web frontend, no GUI, no external server.
Core concerns: game state management, turn logic, entity AI, map/grid rendering, input handling.

---

## Non-Negotiable Rules

### 1. Dependencies — use only what is in `go.mod`

**Never** add a new dependency unless the user explicitly asks for it.

Before writing any import, check `go.mod`:

```bash
cat go.mod
```

If a package is not listed there, do **not** import it.
Do not suggest `go get`, do not add `require` lines, do not vendor new modules.
Solve the problem with the standard library or with what is already listed in `go.mod`.

If you genuinely cannot solve the problem without a new dependency, say so explicitly and
**stop** — do not proceed until the user approves.

### 2. Always verify your work compiles

After every non-trivial change, run:

```bash
go build ./...
```

If it fails, fix the errors before presenting the result.
Never leave the codebase in a broken state.

### 3. Always run tests

After every change, run:

```bash
go test ./...
```

If tests fail, either fix the code **or** fix the test (if the test was wrong).
Explain which you did and why.
Do not skip tests, do not add `t.Skip()` silently.

### 4. Follow CONTRIBUTING.md

Read `CONTRIBUTING.md` before touching anything.
It is the authoritative source on: branch naming, commit message format, code style,
review process, and any project-specific conventions.
If `CONTRIBUTING.md` contradicts something in this file, `CONTRIBUTING.md` wins.

---

## Code Style

- **Formatting**: all code must pass `gofmt`. Run it before finalising any file.
- **Linting**: run `go vet ./...` and fix every warning.
- **Error handling**: never discard errors with `_`. Handle or propagate every error explicitly.
- **No global mutable state** outside of the dedicated game-state package (check `CONTRIBUTING.md`
  for the canonical package name).
- **Naming**: follow standard Go conventions — `PascalCase` for exported symbols,
  `camelCase` for unexported ones, short but descriptive names.
- **Comments**: exported types and functions must have a godoc comment.
  Internal helpers only need a comment when the logic is non-obvious.
- **Panic**: do not use `panic` in game logic. Return errors up the call stack.

---

## Architecture Principles

This is a terminal game. Keep that in mind at every decision.

| Layer | Responsibility |
|---|---|
| `cmd/hera/main.go` | Main file |
| `i18n/` | Translation system and locales|
| `tests` | Test files |
| `core`  | Core files with tiles, grid, rendering and etc | 

- **Game logic must not import TUI packages.** The boundary is strict.
- **TUI packages must not contain game rules.** They only read state and emit commands.
- Prefer small, focused files over large monoliths.
- Keep structs and interfaces in the same package as their primary consumer unless
  sharing across packages is clearly necessary.

---

## TUI-Specific Guidelines

- The TUI library in use is whatever is declared in `go.mod` — inspect it, do not assume.
- Always handle terminal resize events gracefully.
- Do not hard-code terminal dimensions. Query them at runtime.
- Input must be non-blocking where the game loop requires it.
- Render frames only when state changes — avoid unnecessary redraws.
- Leave the terminal in a clean state on exit (restore cursor, disable raw mode, clear
  alternate screen if used).

---

## What Agents Must NOT Do

- Do **not** add new Go modules or packages not in `go.mod` without explicit user approval.
- Do **not** delete or rename existing exported symbols without checking all call sites.
- Do **not** rewrite working subsystems speculatively ("I cleaned up the architecture").
- Do **not** introduce concurrency (goroutines, channels) without a clear, discussed reason.
- Do **not** commit `// TODO` items without a linked issue or a note in the PR description.
- Do **not** ignore compiler warnings or vet output — fix them.
- Do **not** generate placeholder implementations that silently do nothing (fake returns, empty bodies).

---

## Typical Agent Workflow

```text
1. Read CONTRIBUTING.md
2. Read AGENTS.md (this file)
3. Understand the task scope
4. Check go.mod for available dependencies
5. Make the minimal change that solves the problem
6. go build ./...             → fix until clean
7. go test ./...              → fix until green
8. go vet ./...               → fix all warnings
9. golangci-lint run -v ./... → fix all warnings
10. treefmt                   → format
11. Present the diff with a clear explanation
```

---

## Asking for Clarification

If the task is ambiguous — especially regarding game rules, UI behaviour, or whether a
new dependency is acceptable — **ask before coding**. A short clarifying question saves
more time than a large refactor.
