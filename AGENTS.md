# AGENTS.md

Guidance for coding agents working in this repository.

## Scope

These instructions apply to the whole repository unless a more specific nested `AGENTS.md` is present.

Follow explicit user instructions first. If user instructions conflict with this file, ask for clarification before making changes that would violate project rules.

## Project

`i5s` is a k9s-inspired terminal UI for Incus, written in Go with Bubble Tea and Lipgloss.

The goal is a polished fullscreen TUI for browsing Incus instances and running common workflows: shell, logs, console logs, lifecycle actions, remote/project switching, and instance config editing.

## Validation Commands

Run these before finishing meaningful code changes:

- Format: `gofmt -w cmd internal`
- Test: `go test ./...`
- Vet: `go vet ./...`
- Build: `go build ./cmd/i5s`

If `go build ./cmd/i5s` creates a root `./i5s` binary, remove it after validation.

For documentation-only changes, full Go validation is optional. For code changes, run all validation commands unless blocked; if blocked, report the blocker and the command output.

## Definition Of Done

- Relevant code is formatted with `gofmt`.
- Behavior changes are covered by tests where practical.
- `go test ./...`, `go vet ./...`, and `go build ./cmd/i5s` pass for code changes.
- Any generated root `./i5s` binary is removed.
- The final response summarizes what changed and which validation commands passed.

## Architecture

- `cmd/i5s`: CLI entrypoint.
- `internal/app`: app wiring and Bubble Tea startup.
- `internal/config`: runtime flag parsing.
- `internal/logging`: debug log setup.
- `internal/incus`: Incus API integration.
- `internal/ui`: Bubble Tea model, key handling, commands, rendering, and UI behavior tests.

Prefer focused files over large catch-all files. Keep responsibilities separated when adding features.

Place new code in the focused file matching its responsibility. Do not grow `service.go` or `model.go` with unrelated behavior when a narrower file exists.

Current Incus package split:

- `service.go`: constructor, service interface, remote/project connection logic.
- `instances.go`: instance listing, enrichment, state, lifecycle, deletion, row conversion.
- `logs.go`: instance logs and console logs.
- `edit.go`: instance config editing.
- `exec_unix.go`: native interactive shell exec on Unix.
- `exec_windows.go`: Windows unsupported shell stub.
- `types.go`: user-facing row types and search text.

Current UI package split:

- `model.go`: root Bubble Tea model and message handling.
- `messages.go`: Bubble Tea message types.
- `commands.go`: asynchronous command builders.
- `keys.go`: keyboard handling.
- `selection.go`: selection helpers.
- `exec_commands.go`: terminal-yielding commands for shell/config edit.
- `views.go`, `render.go`, `table.go`, `styles.go`, `help.go`: rendering and visual language.

## Incus Integration Rules

- Use the Incus Go API, not the `incus` CLI binary, unless explicitly requested.
- Remote/project switching is runtime-only and must not mutate Incus CLI defaults.
- Shell uses native Incus exec behavior equivalent to `incus exec <instance> -- su -l`.
- Config editing uses native Incus API plus `github.com/lxc/incus/v7/shared/cmd.TextEditor`, matching `incus config edit` editor behavior.
- Config editing should use the selected runtime remote/project from `i5s`.
- Delete must only be allowed for stopped instances.
- Best-effort enrichment failures can be logged in debug output, but should not break instance listing.

## Error Handling

- User-triggered operations should surface visible errors in the TUI.
- Non-interactive service calls should use bounded contexts.
- Intentionally interactive editor/shell sessions should not have artificial timeouts unless explicitly requested.
- Best-effort data enrichment may fail without blocking the main list, but failures should remain debuggable through logs.

## UI Behavior

- Binary name is `i5s`.
- Main view layout is `HEADER -> BODY -> FOOTER`.
- The instance table panel fills the remaining terminal height.
- Default table columns are `NAME | STATE | IPV4 | IPV6 | TYPE | SNAPSHOTS`.
- Use Unicode box drawing by default: `╭ ╮ ╰ ╯ ─ │ › …`.
- Preserve the existing fullscreen visual language unless explicitly asked for redesign.
- `ctrl+c` quits globally.
- `enter` opens a shell for the selected running instance.
- `e` edits the selected instance config in the user's editor.
- `l` opens logs; `c` opens console logs.
- `s`, `S`, and `d` stop, start, and delete with existing guards and confirmations.
- Remote/project picker changes are session-local.
- Terminal-yielding workflows such as shell and config edit should use `tea.Exec` so Bubble Tea can suspend and resume the TUI correctly.

## Testing Guidance

- Prefer behavior/input-output tests over implementation-detail tests.
- Keep tests in the same package for now (`package ui`, `package incus`) to avoid exporting internals only for tests.
- Do not add production hooks, exported symbols, or complexity only for tests.
- Use fake services and Bubble Tea messages to test UI behavior.
- Avoid launching real editors or real shells in tests.
- Config-edit YAML helpers are intentionally testable without launching an editor.
- When adding a user-facing key or workflow, test the selected-instance behavior, no-op behavior, failure behavior, and visible discoverability where practical.
- Prefer assertions on meaningful visible output or service calls over internal fields.
- Avoid brittle assertions on exact ANSI/style output unless layout or styling is the behavior being tested.

When changing these areas, preserve or add behavior tests for:

- navigation and filtering
- selection preservation across refreshes
- lifecycle actions and confirmations
- logs and console logs
- remote/project pickers
- rendering and fullscreen layout
- shell behavior
- config edit behavior and YAML parsing

## Coding Style

- Make the smallest correct change.
- Do not add abstractions unless they reduce current duplication or clarify behavior.
- Keep code direct and idiomatic Go.
- Keep comments rare; add them only when behavior is not obvious.
- Document exported Go identifiers with idiomatic doc comments that start with the identifier name.
- Do not comment every unexported helper; add comments only when behavior, constraints, or side effects are not obvious.
- Prefer comments that explain why something exists or what contract it preserves, not comments that restate implementation.
- Avoid backward-compatibility code unless there is a real persisted, shipped, or external compatibility need.
- Prefer bounded contexts for non-interactive service calls.
- Do not impose timeouts on intentionally interactive editor/shell sessions unless explicitly requested.

## Dependency Notes

- `go.yaml.in/yaml/v4` is a direct dependency because config editing parses/dumps Incus YAML.
- Importing Incus `shared/cmd.TextEditor` intentionally brings additional indirect dependencies to match Incus CLI editor behavior.
- Do not replace `TextEditor` with a local editor runner unless there is a concrete reason.

## Agent Workflow

- Build context before editing; inspect relevant files first.
- Preserve unrelated user changes in the worktree.
- Never revert unrelated changes unless explicitly requested.
- Do not create git commits unless explicitly asked.
- Do not use destructive git commands unless explicitly approved.
- After code changes, run the validation commands listed above.
- Report what changed and which validation commands passed.
