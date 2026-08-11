[![official JetBrains project](http://jb.gg/badges/official.svg)](https://confluence.jetbrains.com/display/ALL/JetBrains+on+GitHub)

# Modern Go Guidelines

This repository contains [guidelines](https://github.com/JetBrains/go-modern-guidelines/blob/main/plugin/skills/use-modern-go/SKILL.md) for code agents that help them write modern Go code.

For example, an agent with these guidelines uses `max(a, b)` instead of an if-else block, `slices.Contains` instead of a manual loop, `cmp.Or(a, b, c)` instead of a chain of nil checks. It also knows about recent additions like `new(42)` to get a pointer to a value and `errors.AsType[T](err)` for type-safe error matching—both from Go 1.26.

The guidelines cover the most useful features from Go 1.0 through Go 1.27, including everything targeted by the `modernize` analyzer. An agent will:

- Detect the project's Go version from `go.mod`
- Use language features and stdlib additions available up to and including that version
- Prefer modern idioms over older patterns

## Motivation

All coding agents tend to generate outdated Go. Two reasons:

1. **Training data lag.** Models don't know about features added after their training cutoff. They can't use `errors.AsType[T]` (Go 1.26) if they've never seen it.

2. **Frequency bias.** Even for features the model knows, it often picks older patterns. There's more `for i := 0; i < n; i++` in the training data than `for i := range n`, so that's what comes out.

These guidelines fix both problems by giving the agent an explicit reference.

This aligns with the Go team's direction. The `modernize` analyzer exists to automatically update existing code to use newer idioms (see [this talk](https://www.youtube.com/watch?v=_VePjjjV9JU) from the Go team). These guidelines serve the same goal for new code: agents write modern Go from the start, so there's less to fix later.

## Requirements

The marketplace integrations run a small CLI that is installed on first use with `go install`. Because of that, the [Go toolchain](https://go.dev/dl/) must be installed and available on your `PATH`.

The CLI is installed into a local cache (for example `~/.cache/go-modern-guidelines`) and never modifies your project. It targets **Go 1.25 or newer**; on an older Go it still works as long as automatic toolchain switching is enabled (`GOTOOLCHAIN=auto`, the default), which lets Go fetch a compatible toolchain on first run.

## Instructions

The guidelines are available for Junie, Claude Code, Codex, and Cursor, and for other agents via skills.sh.

### [Junie](https://junie.jetbrains.com)

#### Junie CLI

Run the following commands inside a Junie CLI session.

1. Add this repository as a marketplace:
```
/extensions marketplace add JetBrains/go-modern-guidelines
```

2. Install the extension:
```
/extensions install modern-go-guidelines
```

Junie invokes the skill automatically when it is relevant to a Go task.

#### Updating

Update the installed extension from inside a Junie CLI session:

```
/extensions update modern-go-guidelines
```

### [Claude Code](https://claude.com/product/claude-code)

#### Installation
Run the following commands inside a Claude Code session.

1. Add this repository as a marketplace:
```
/plugin marketplace add JetBrains/go-modern-guidelines
```

2. Install the plugin:
```
/plugin install modern-go-guidelines@goland-claude-marketplace
```

#### Usage

Claude Code invokes the skill automatically when it is relevant to a Go task.

To invoke it explicitly:

```
/modern-go-guidelines:use-modern-go
```

#### Updating

Claude Code can update the marketplace and installed plugin automatically at startup. Automatic updates are disabled by default for third-party marketplaces, so enable them once:

1. Run `/plugin`.
2. Open **Marketplaces** and select `goland-claude-marketplace`.
3. Select **Enable auto-update**.

When Claude Code reports that the plugin was updated, apply the new version to the current session with:

```
/reload-plugins
```

To update it manually instead, run these commands in a terminal:

```bash
claude plugin marketplace update goland-claude-marketplace
claude plugin update modern-go-guidelines@goland-claude-marketplace
```

### [Codex](https://developers.openai.com/codex/)

#### Installation
Run the following commands in a terminal.

1. Add this repository as a marketplace:
```
codex plugin marketplace add JetBrains/go-modern-guidelines
```

2. Install the plugin:
```
codex plugin add modern-go-guidelines@goland-codex-marketplace
```

#### Updating

Refresh the marketplace and reinstall the plugin so Codex replaces its cached copy:

```bash
codex plugin marketplace upgrade goland-codex-marketplace
codex plugin remove modern-go-guidelines@goland-codex-marketplace
codex plugin add modern-go-guidelines@goland-codex-marketplace
```

### [Cursor](https://cursor.com)

For convenience, the guidelines are distributed as a Cursor plugin.

#### Installation

1. Add this repository as a marketplace by running the following command in a terminal:
```
cursor-agent plugin marketplace add https://github.com/JetBrains/go-modern-guidelines
```

2. Install the plugin with the `/plugins` command inside a Cursor session.

#### Updating

Refresh the marketplace from Git and reopen Cursor so it can pick up the new plugin version:

```bash
cursor-agent plugin marketplace update goland-cursor-marketplace
```

If the installed plugin is still on the previous version, reinstall it with the `/plugins` command. Cursor does not currently provide a non-interactive CLI command for updating an installed plugin.

### Other Agents (via [skills.sh](https://skills.sh))

The same skill package works across other agents such as OpenCode. Install it with:

```bash
npx skills add JetBrains/go-modern-guidelines
```

(`--skill use-modern-go` installs only this skill.)

#### Updating

Update the project-installed skill with:

```bash
npx skills update use-modern-go -p -y
```

For a globally installed skill, replace `-p` with `-g`.

## Local development

To try changes to the CLI in your agent, build this checkout into the tool's cache:

```bash
make dev-install
```

Then set `GO_MODERN_GUIDELINES_DEV=1` in the environment your agent runs in. With it set, any agent using the plugin runs your local build instead of the released version, the same way across Claude Code, Codex, and Cursor. Export it before launching the agent so the agent process inherits it:

```bash
export GO_MODERN_GUIDELINES_DEV=1
```

After editing the CLI, run `make dev-install` again to rebuild; the next call picks it up. To go back to the released version, unset the variable (or run `make dev-uninstall` to remove the build):

```bash
make dev-uninstall
```

This requires the Go toolchain. The dev build is stored in the tool's cache directory (`$XDG_CACHE_HOME/go-modern-guidelines` or `~/.cache/go-modern-guidelines`).

The build is driven by `scripts/dev-install.sh`, which is intentionally separate from the agent-facing wrapper so an agent can never trigger a build. Without `make` (for example on Windows) you can run it directly:

```bash
sh scripts/dev-install.sh install       # or: uninstall
pwsh scripts/dev-install.ps1 install    # PowerShell equivalent
```
