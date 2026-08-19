<p align="center">
  <img src="assets/logo-dark.svg?v=2" alt="Terra-Prism Logo" width="400">
</p>

<p align="center">
  <strong>A beautiful terminal UI for viewing Terraform and OpenTofu plans</strong>
</p>

<p align="center">
  <a href="https://github.com/CaptShanks/terraprism/actions/workflows/ci.yml"><img src="https://github.com/CaptShanks/terraprism/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/CaptShanks/terraprism/releases"><img src="https://img.shields.io/github/v/release/CaptShanks/terraprism?include_prereleases" alt="Release"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/github/go-mod/go-version/CaptShanks/terraprism" alt="Go Version"></a>
  <a href="https://goreportcard.com/report/github.com/CaptShanks/terraprism"><img src="https://goreportcard.com/badge/github.com/CaptShanks/terraprism" alt="Go Report Card"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/CaptShanks/terraprism" alt="License"></a>
</p>

<p align="center">
  Collapsible resources and sub-blocks • Filter & sort • Syntax-highlighted HCL • Vim-style navigation • Auto light/dark mode
</p>

---

# This fork

Everything below this section is upstream's original README, unchanged. This
section describes how this fork differs and supersedes it wherever the two
disagree — specifically the **Installation**, **Upgrading**, **Commands**, and
**Environment Variables** sections still document the update machinery that has
been removed here.

## What changed

Every network call was removed, and plan output is stored more carefully.

| Change | Detail |
|---|---|
| **No phone-home** | The `internal/updater` package is deleted. The TUI no longer runs a background GitHub release check on startup, and `terraprism version` no longer checks for updates. |
| **No token access** | Upstream's update library read `GITHUB_TOKEN` from your environment — falling back to `github.token` in your git config — and sent it to api.github.com. That code path no longer exists. |
| **No self-update** | `terraprism upgrade` is removed; it now exits 1 with a pointer to rebuild. Upstream downloaded and swapped the running binary with **no checksum or signature verification**, despite the release workflow publishing a `checksums.txt` that nothing ever read. |
| **Tighter history** | `~/.terraprism/` is now `0700` and history files `0600` (was `0755`/`0644`, i.e. readable by any other local user). Files are deleted once **older than 1 hour**, checked at the start of every plan/apply/destroy run. |

The binary no longer links `net`, `net/http`, or `crypto/tls` at all — it is
structurally incapable of making a network connection, not merely configured not
to. Verify for yourself:

```bash
go list -deps ./cmd/terraprism | grep -E '^(net|net/http|crypto/tls)$'   # no output
```

Dropping the updater also removed 9 transitive dependencies, including
`go-github`, `oauth2`, and versions of `golang.org/x/net` and `x/crypto` from
2019 and 2020. What remains is the Charm TUI stack only.

**Bug fixed along the way:** history filenames were written in local time but
parsed as UTC, so every timestamp — and therefore the retention cutoff — was
skewed by your zone offset. In UTC+2 a "1 hour" cutoff retained files for three.

**Still worth knowing:** `terraprism apply` writes a binary plan file to
`$TMPDIR/terraprism-<pid>.tfplan`. Unlike the text history, binary plan files
contain *unredacted* values. It is removed on exit, but not if the process is
killed. On Linux, where `/tmp` is shared, treat that path as sensitive.

## Build and install

Build from source. **Do not use the `curl | sh` quick install or the release
binaries documented below** — those are upstream's and still contain everything
listed above.

Requires Go 1.24+:

```bash
make build                          # -> bin/terraprism
```

Or with Docker, if you would rather not install a Go toolchain:

```bash
make docker-build                   # -> dist/terraprism   (Linux, host arch)
make docker-build-mac               # -> dist/terraprism   (macOS, host arch)
```

The container is always Linux, so `docker-build` produces a Linux binary even on
a Mac — hence the second target. Both write to the same path, so the last one you
ran wins; each prints `file dist/terraprism` afterwards so you can tell which you
have. For any other target:

```bash
docker build -o dist --build-arg GOOS=windows --build-arg GOARCH=amd64 .
docker build -o dist --build-arg VERSION=0.12.0-local .
```

Where the binary lands:

| Command | Output |
|---|---|
| `make` / `make build` | `bin/terraprism` |
| `make docker-build` / `make docker-build-mac` | `dist/terraprism` |
| `make release` | `bin/terraprism-<os>-<arch>`, all five platforms |

Then put it on your `PATH` — `~/.local/bin` needs no sudo:

```bash
install -m 0755 bin/terraprism ~/.local/bin/terraprism
```

Note that `make install` runs `go install`, which writes to `~/go/bin` — a
different directory, and probably not on your `PATH`.

## Upgrading this fork

There is no self-update and no update check. Pull and rebuild:

```bash
git pull && make build
```

## Verify

```bash
make test
terraprism --version    # must return instantly and never mention an update
```

Retention and permissions are covered by `internal/history/history_test.go`.
Both limits are constants in `internal/history/history.go` — `MaxHistoryAge`
(1 hour) and `MaxHistoryFiles` (100) — if you want different retention.

Cleanup is triggered by running terraprism, not by a daemon: if you stop using
it, whatever was written in the last hour stays until the next run. Add
`rm -f ~/.terraprism/*.txt` to a logout hook if you need a harder guarantee.

---

<p align="center">
  <img src="assets/demo.gif" alt="Terra-Prism Demo" width="800">
</p>

## Features

- **Syntax-highlighted HCL** - Full color-coded display of your plan
- **Collapsible resources and sub-blocks** - Expand/collapse resources, large maps, lists, and heredocs
- **Status filter** - Filter resources by action (create, destroy, update, replace, read, etc.)
- **Sort** - Sort by plan order, action, address, or resource type
- **Search** - Find resources by name, type, or address (works with filters)
- **Vim-style navigation** - j/k/gg/G/d/u plus line scrolling for large blocks
- **Auto light/dark mode** - Detects your terminal background
- **Format support** - Works with Terraform 0.11+ and OpenTofu
- **Full-line selection** - Clear visual indicator of selected resource
- **History tracking** - All plans saved with full path, searchable picker
- **Color-coded CLI** - Commands and status colored in history list
- **Environment variables** - Set `TERRAPRISM_TOFU` and `TERRAPRISM_THEME` to avoid passing flags every time
- **State viewer** - Interactive TUI for `state list`, `state show`, and `state rm` with search, sort, multi-select, taint, and untaint
- **Passthrough commands** - Run init, validate, fmt, and other terraform/tofu commands through terraprism

## Installation

### Quick Install (Recommended)

```bash
curl -sSfL https://raw.githubusercontent.com/CaptShanks/terraprism/main/install.sh | sh
```

### Using Go

```bash
go install github.com/CaptShanks/terraprism/cmd/terraprism@latest
```

### From Source

```bash
git clone https://github.com/CaptShanks/terraprism.git
cd terraprism
make build
```

### Manual Download

Download binaries from the [Releases](https://github.com/CaptShanks/terraprism/releases) page.

## Usage

### Apply Mode (Recommended)

Review and apply in one command:

```bash
# Run plan, review interactively, press 'a' to apply
terraprism apply

# With OpenTofu (set TERRAPRISM_TOFU=1)
TERRAPRISM_TOFU=1 terraprism apply

# Pass arguments to terraform/tofu
terraprism apply -- -target=module.vpc -var="env=prod"
```

### Plan Mode

Run plan and view interactively (no apply):

```bash
terraprism plan
TERRAPRISM_TOFU=1 terraprism plan
```

### State Mode

Interactive TUI for Terraform state with search, sort, show details, remove, taint, and untaint:

```bash
terraprism state list
terraprism state show
terraprism state rm
```

All three subcommands open the same unified TUI. Use **Space** to select items, **Ctrl+Space** to select a range, then **Enter** (show), **t** (taint), **u** (untaint), or **d**/**x** (remove). **Esc** clears selection or dismisses confirmation. Other state subcommands (`state mv`, `state pull`, etc.) pass through to terraform/tofu.

### Pipe Mode

Pipe plan output for viewing:

```bash
terraform plan -no-color | terraprism
tofu plan -no-color | terraprism
```

### Read from file

```bash
terraform plan -no-color > plan.txt
terraprism plan.txt
```

### Print mode (non-interactive)

```bash
terraform plan -no-color | terraprism -p
```

## Keyboard Controls

### Navigation
| Key | Action |
|-----|--------|
| `j` / `↓` | Move to next resource or foldable sub-block |
| `k` / `↑` | Move to previous resource or foldable sub-block |
| `gg` | Jump to first resource |
| `G` | Jump to last resource |
| `d` / `Ctrl+D` | Scroll half page down |
| `u` / `Ctrl+U` | Scroll half page up |
| `Ctrl+E` | Scroll one line down |
| `Ctrl+Y` | Scroll one line up |
| `+` / `=` | Show more unchanged context around diff hunks |
| `-` | Show less unchanged context around diff hunks |

### Expand/Collapse
| Key | Action |
|-----|--------|
| `Enter` / `Space` | Toggle current resource or foldable sub-block |
| `l` / `→` | Expand current resource or foldable sub-block |
| `h` / `←` / `⌫` | Collapse current resource or foldable sub-block |
| `e` | Expand the highlighted item and its foldable sub-blocks (scope) |
| `c` | Collapse the highlighted item and its foldable sub-blocks (scope) |
| `Shift+E` | Expand all visible resources and all nested foldable sub-blocks |
| `Shift+C` | Collapse all visible resources and all nested foldable sub-blocks |

Large maps, lists, and heredocs inside expanded resources become foldable sub-blocks. Large sub-blocks collapse by default; use `l`/`→` or `Enter`/`Space` to expand them, then `Ctrl+E`/`Ctrl+Y` to scroll through the expanded content without moving the selection. When a resource or sub-block is selected, `e` and `c` recursively expand or collapse the foldable content underneath that selection.

Paired remove/add heredocs are shown as one foldable diff section so large values changes can be reviewed as a focused line diff. Use `+`/`=` and `-` to increase or decrease the unchanged context shown around each diff hunk.

### Search
| Key | Action |
|-----|--------|
| `/` | Start search |
| `n` | Next match |
| `N` | Previous match |
| `Esc` | Clear search (or clear filters when filters active) |

### Filter
| Key | Action |
|-----|--------|
| `f` | Open filter picker |
| `Esc` | Clear all filters (from main view or picker) |

In the filter picker: **Space** toggle status, **a** select all, **c** clear all, **Enter** apply, **Esc** clear and close.

### Sort
| Key | Action |
|-----|--------|
| `s` | Open sort picker |

Sort options: default (plan order), by action, by address, by type.

### Apply (in apply mode)
| Key | Action |
|-----|--------|
| `a` | Apply the plan |
| `y` | Confirm apply |

### State Mode (state list/show/rm)
| Key | Action |
|-----|--------|
| `j`/`k`, `↑`/`↓` | Navigate (highlight item) |
| `Space` | Toggle selection of highlighted item |
| `Ctrl+Space` | Select range from anchor to current |
| `Enter` | Show details (opens immediately, no confirmation) |
| `t` | Taint selected (confirmation required) |
| `u` | Untaint selected (confirmation required) |
| `d` or `x` | Remove from state (confirmation required) |
| `/` | Search |
| `f` or `s` | Sort (by address, type, module depth) |
| `Esc` | Clear selection or dismiss confirmation |
| `q` | Quit |

In the detail view: `/` search within content, `e`/`c` expand/collapse all blocks, `j`/`k` scroll, `y` copy full content, `Y` copy line at top (OSC 52).

### Other
| Key | Action |
|-----|--------|
| `q` / `Ctrl+C` | Quit (cancel apply) |

## Color Themes

Terra-Prism automatically detects your terminal background and uses appropriate colors:

### Dark Mode (Catppuccin Mocha)
- Green for resources being created
- Red for resources being destroyed
- Yellow for resources being updated
- Purple for resources being replaced
- Blue for data sources being read

### Light Mode (Catppuccin Latte)
Automatically switches to darker, more visible colors on light backgrounds.

### Force a theme
```bash
TERRAPRISM_THEME=dark terraprism plan.txt   # Force dark mode
TERRAPRISM_THEME=light terraprism plan.txt # Force light mode
```

## Commands

```
terraprism                     # View piped/file input
terraprism plan                # Run terraform plan and view
terraprism apply               # Run plan, view, and apply
terraprism destroy             # Run destroy plan and apply
terraprism state list|show|rm  # Interactive state TUI (search, sort, taint, untaint)
terraprism history             # Manage history files
terraprism version             # Show terraprism and terraform/tofu version
terraprism upgrade             # Upgrade to the latest release
terraprism init|validate|fmt|output|state|import|...  # Pass through to terraform/tofu
```

## Options

```
-h, --help      Show help message
-v, --version   Show version (includes update check and terraform/tofu version)
-p, --print     Print colored output without interactive TUI
```

## Environment Variables

Set these to avoid passing flags every time:

```
TERRAPRISM_TOFU    Set to 1, true, or yes to use OpenTofu instead of Terraform
TERRAPRISM_THEME   Set to "light" or "dark" to force color scheme
TERRAPRISM_SKIP_UPDATE_CHECK   Set to 1, true, or yes to skip update checks
TERRAPRISM_UPDATE_CHECK_INTERVAL  Days between TUI update checks (default: 7)
```

Example: add `export TERRAPRISM_TOFU=1` to your `~/.bashrc` or `~/.zshrc` to always use OpenTofu.

## Upgrading

Upgrade to the latest release:

```bash
terraprism upgrade
```

If self-update fails (e.g. binary in read-only location), re-run the install script:

```bash
curl -sSfL https://raw.githubusercontent.com/CaptShanks/terraprism/main/install.sh | sh
```

When a newer version is available, terraprism shows an update nudge in the TUI footer and when running `terraprism version`.

## Passthrough Commands

Any terraform/tofu command not enhanced by terraprism (init, validate, fmt, output, state mv/pull/push, import, workspace, graph, console, login, logout, providers, force-unlock, show, refresh, taint, untaint) is passed through to the selected engine. Use terraprism as a drop-in replacement:

```bash
terraprism init
terraprism init -upgrade
terraprism validate
terraprism fmt -recursive
terraprism output
terraprism state mv src dst   # state mv, pull, push, etc. pass through
terraprism state list        # state list, show, rm use interactive TUI
```

## History

All plan and apply outputs are automatically saved to `~/.terraprism/` for future reference. History includes the full working directory path and is color-coded for easy reading.

### Listing History

```bash
terraprism history list              # List all history files
terraprism history list --plan       # List only plan files
terraprism history list --apply      # List only apply files
terraprism history list --destroy    # List only destroy files
```

Output shows timestamp, command (colored), status (colored), and full path:
```
  #  TIMESTAMP         COMMAND  STATUS        PATH
--------------------------------------------------------------------------------------
  1  2026-01-14 12:52  plan                   .../infrastructure/aws/prod/eks-cluster
  2  2026-01-14 10:30  apply    [SUCCESS]     .../infrastructure/aws/staging/vpc
```

### Viewing History

```bash
terraprism history view              # Interactive picker with search
terraprism history view 1            # View most recent entry
terraprism history view 3            # View 3rd most recent
terraprism history 1                 # Shorthand for 'view 1'
```

The interactive picker supports:
- **j/k** - Navigate up/down
- **/** - Search (fzf-style, multiple terms)
- **Enter** - Select and view in TUI
- **q/Esc** - Cancel

Search by project, command, status, date, or path:
```
/myproject apply success    # Find successful applies for myproject
/2026-01 destroy            # Find January 2026 destroys
```

### File Naming

Files are named: `YYYY-MM-DD_HH-MM-SS_<project>_<command>[_<status>].txt`

- `plan` - Plan-only commands
- `apply` - Apply commands (status: success, failed, cancelled)
- `destroy` - Destroy commands (status: success, failed, cancelled)

### Cleanup

```bash
terraprism history list --clear      # Delete all history files
```

History is automatically cleaned up when exceeding 100 files (oldest removed first).

## Why Terra-Prism?

Large Terraform plans can be difficult to review:

- Hundreds of resources make it hard to find specific changes
- Long attribute values span multiple lines
- No easy way to focus on specific resources
- Color coding from Terraform can be lost when piping

Terra-Prism solves these problems:

- Collapsible sections and sub-blocks for high-level overview
- Filter by status to focus on creates, destroys, updates, etc.
- Sort by action, address, or type for organized review
- Consistent syntax highlighting
- Search to find specific resources (works with filters)
- Vim-style navigation for resources, sub-blocks, and large expanded content
- Auto-scrolling keeps selection visible

## Inspired By

- [prettyplan](https://prettyplan.chrislewisdev.com/) - Web-based Terraform plan formatter
- [terraform-landscape](https://github.com/coinbase/terraform-landscape) - Ruby-based plan formatter

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development

```bash
# Clone the repo
git clone https://github.com/CaptShanks/terraprism.git
cd terraprism

# Install dependencies
go mod download

# Run tests
make test

# Build
make build

# Run locally
./bin/terraprism
```

## License

GNU AGPL v3.0 - see [LICENSE](LICENSE) for details.

---

Made with Go
