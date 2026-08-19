# Terraprism - Claude Code Configuration

## Project Overview

Terraprism is a terminal UI (TUI) tool that visualizes Terraform plan output in a human-friendly, interactive format. It parses `terraform plan` JSON output and presents resource changes with diffs, filtering, and search capabilities.

- **Tech Stack:** Go 1.24, Charm Bubble Tea, Charm Lipgloss, Charm Bubbles
- **Module path:** `github.com/CaptShanks/terraprism`
- **License:** MIT

## Repository Structure

```
cmd/terraprism/       Entry point (main.go)
internal/
  config/             Configuration management
  history/            Plan history tracking
  parser/             Terraform plan JSON parser (core logic)
  tui/                Terminal UI components (Bubble Tea model, views, styles)
  updater/            Self-update mechanism via GitHub releases
testdata/             Sample Terraform plan JSON files for testing
docs/architecture/    Architecture documents for issues (auto-generated)
```

## Development Standards

- Follow standard Go conventions (camelCase unexported, PascalCase exported)
- Wrap errors: `fmt.Errorf("context: %w", err)`
- All packages under `internal/` — no external imports allowed
- Exported functions and types must have godoc comments
- `gofmt` / `goimports` enforced

## Build & Test

```bash
go build -o terraprism ./cmd/terraprism
go test -v -race -cover ./...
golangci-lint run ./...
```

## Testing Conventions

- Go standard `testing` package
- Table-driven tests with subtests (`t.Run`), `t.Parallel()` where safe
- Test data in `testdata/` directory
- Test files colocated with source (`*_test.go`)

## Dependencies

| Package | Purpose |
|---------|---------|
| `charmbracelet/bubbletea` v1.3.10 | TUI framework (Elm architecture) |
| `charmbracelet/bubbles` v0.21.0 | Pre-built TUI components |
| `charmbracelet/lipgloss` v1.1.0 | Terminal styling and layout |
| `muesli/reflow` v0.3.0 | Text wrapping and padding |

No non-TUI dependencies. `blang/semver` and `rhysd/go-github-selfupdate` were
removed with the updater; adding anything that opens a socket is a design change,
not a dependency bump.

## Architecture Decisions

1. **Elm Architecture (TEA):** TUI follows Bubble Tea's Model-Update-View pattern
2. **Parser separation:** `internal/parser` is independent of TUI, produces structured data
3. **No self-update:** Version is injected with `-ldflags`; upgrading means rebuilding from source
4. **Internal packages:** All under `internal/` for free refactoring
5. **Piped input:** Reads from stdin, composable with shell pipelines

## Security Considerations

- Terraform plan output may contain sensitive attributes (passwords, keys, connection strings)
- **The binary makes no network calls.** There is no updater, no telemetry, and no
  code that reads `GITHUB_TOKEN`. Do not reintroduce one.
- Plan/apply output is persisted to `~/.terraprism/` (`0700` dir, `0600` files) and
  deleted once older than `history.MaxHistoryAge` (1 hour), enforced at the start of
  every plan/apply/destroy run
- Binary plan files (`$TMPDIR/terraprism-<pid>.tfplan`) hold unredacted values and are
  removed on exit

## PR Guidelines

- Target `dev` branch for all feature work
- PR title format: `type: description` (e.g., `feat: add resource filtering`)
- Include tests for new functionality
- Run `golangci-lint` before submitting

## Branching

- `main` — production releases only
- `dev` — integration branch for feature work
- `shipwright/issue-N` — feature branches created by automation
