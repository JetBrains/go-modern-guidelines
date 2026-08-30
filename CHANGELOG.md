# Changelog

## Unreleased

### Changed

- Corrected example indentation to use tabs.

### Fixed

- Corrected the Windows PowerShell wrapper invocation to use the call operator (#26) @FenjuFu
- Corrected the availability of the `strings_clone` guideline to Go 1.18 and later (#23) @chiang21fcb

## 1.1.1 - 2026-08-19

### Added

- Guidance for adopting and migrating to `encoding/json/v2` in Go 1.27 and later.

## 1.1.0 - 2026-08-11

### Added

- A CLI tool that detects a project's Go version and dynamically provides only relevant guidelines.
- Dedicated integrations for Junie CLI, Codex, and Cursor.
- Guidelines for Go 1.27 features.

## 1.0.0 - 2026-02-04

### Added

- Modern Go guidelines for Claude Code and Junie, covering features through Go 1.26.
- Installation support via skills.sh (#8) @raeperd

### Changed

- Improved `time.Tick` guidance for Go 1.23 (#2) @endigma
